package apis

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/dbs"
)

func handlerP8(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return append([]byte("\r\n "), pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})...)
}

func TestUploadP8Handler(t *testing.T) {
	p8 := handlerP8(t)
	mock, cleanup := openPushHandlerMockDB(t)
	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `ioscertificates`").WithArgs("com.example", []byte(nil), "", "app-1", "", 0, []byte(nil), "", "", "p8", "ABC1234567", "XYZ1234567", p8, "key.p8", 1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	request := multipartPushRequest(t, "/apps/iospushcer/upload", map[string]string{"app_key": "app-1", "package": "com.example", "auth_type": "p8", "p8_key_id": "ABC1234567", "p8_team_id": "XYZ1234567"}, map[string]struct {
		name string
		data string
	}{"p8_file": {name: "key.p8", data: string(p8)}})
	response := invokePushHandler(t, request, UploadIosCer)
	if responseCode(t, response) != 0 {
		t.Fatal("valid upload failed")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUploadP8RejectsOversizedBeforeDatabase(t *testing.T) {
	mock, cleanup := openPushHandlerMockDB(t)
	defer cleanup()
	request := multipartPushRequest(t, "/apps/iospushcer/upload", map[string]string{"app_key": "app-1", "package": "com.example", "auth_type": "p8"}, map[string]struct {
		name string
		data string
	}{"p8_file": {name: "key.p8", data: strings.Repeat("x", 16385)}})
	if responseCode(t, invokePushHandler(t, request, UploadIosCer)) != int(errs.AdminErrorCode_ParamError) {
		t.Fatal("oversized P8 accepted")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIosQueriesOnlyReturnSafeDTO(t *testing.T) {
	for _, list := range []bool{false, true} {
		t.Run(map[bool]string{false: "get", true: "list"}[list], func(t *testing.T) {
			mock, cleanup := openPushHandlerMockDB(t)
			defer cleanup()
			columns := []string{"app_key", "package", "cert_pwd", "voip_cert_pwd", "certificate", "voip_cert", "auth_type", "p8_key_id", "p8_team_id", "p8_key_name", "p8_private_key", "config_version"}
			query := mock.ExpectQuery("SELECT .*ioscertificates.*ORDER BY package asc")
			if list {
				query.WithArgs("app-1")
			} else {
				query.WithArgs("app-1", 1)
			}
			query.WillReturnRows(sqlmock.NewRows(columns).AddRow("app-1", "com.example", "hidden-cert-password", "hidden-voip-password", []byte("hidden-cert"), []byte("hidden-voip"), "p8", "ABC1234567", "XYZ1234567", "saved.p8", []byte("hidden-raw-key"), 7))
			handler := GetIosCer
			if list {
				handler = ListIosPushConfs
			}
			response := invokePushHandler(t, httptest.NewRequest(http.MethodGet, "/?app_key=app-1", nil), handler)
			if responseCode(t, response) != 0 {
				t.Fatal("query failed")
			}
			for _, forbidden := range []string{"hidden-", "cert_pwd", "voip_cert_pwd", `"certificate"`, `"voip_cert"`, "p8_private_key", base64.StdEncoding.EncodeToString([]byte("hidden-raw-key"))} {
				if strings.Contains(response.Body.String(), forbidden) {
					t.Fatal("query leaked credentials")
				}
			}
			for _, expected := range []string{`"auth_type":"p8"`, `"has_p8_key":true`, `"config_version":7`, `"p8_key_name":"saved.p8"`} {
				if !strings.Contains(response.Body.String(), expected) {
					t.Fatalf("missing metadata %s", expected)
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIosAppAccessRejectsUnboundBrowserAccount(t *testing.T) {
	for index, handler := range []gin.HandlerFunc{GetIosCer, ListIosPushConfs, SetIosPushConf, UploadIosCer} {
		mock, cleanup := openPushHandlerMockDB(t)
		account := t.Name() + string(rune('A'+index))
		mock.ExpectQuery("SELECT .*accounts.*account=.*LIMIT").WithArgs(account, 1).
			WillReturnRows(sqlmock.NewRows([]string{"account", "state", "role_type"}).AddRow(account, 0, 1))
		mock.ExpectQuery("SELECT .*accountapprels.*account=.*app_key=.*LIMIT").WithArgs(account, "app-other", 1).
			WillReturnRows(sqlmock.NewRows([]string{"account", "app_key"}))
		request := httptest.NewRequest(http.MethodGet, "/?app_key=app-other", nil)
		if index == 2 {
			request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"app_key":"app-other","package":"com.example","original_package":"com.example"}`))
			request.Header.Set("Content-Type", "application/json")
		}
		if index == 3 {
			request = multipartPushRequest(t, "/", map[string]string{"app_key": "app-other", "package": "com.example"}, nil)
		}
		response := invokePushHandler(t, request, func(ctx *gin.Context) { ctx.Set(string(ctxs.CtxKey_Account), account); handler(ctx) })
		if responseCode(t, response) != int(errs.AdminErrorCode_NotPermission) {
			t.Fatal("unbound account accepted")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
		cleanup()
	}
}

func TestIosDatabaseErrorsAreRedacted(t *testing.T) {
	mock, cleanup := openPushHandlerMockDB(t)
	defer cleanup()
	mock.ExpectQuery("SELECT .*ioscertificates").WillReturnError(errors.New("private-key SQL password secret"))
	response := invokePushHandler(t, httptest.NewRequest(http.MethodGet, "/?app_key=app-1", nil), ListIosPushConfs)
	var result map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if responseCode(t, response) == 0 || strings.Contains(response.Body.String(), "secret") {
		t.Fatal("unsafe database failure")
	}
}

func TestIosEditVersionsJSONAndMultipart(t *testing.T) {
	for _, upload := range []bool{false, true} {
		for _, version := range []string{"", "0", "3", "4"} {
			for _, mode := range []string{"p12", "p8"} {
				t.Run(mode+"/"+version+map[bool]string{false: "/json", true: "/multipart"}[upload], func(t *testing.T) {
					mock, cleanup := openPushHandlerMockDB(t)
					defer cleanup()
					mock.ExpectBegin()
					mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WillReturnRows(sqlmock.NewRows([]string{"app_key", "package", "auth_type", "config_version", "voip_cert", "voip_cert_pwd"}).AddRow("app-1", "com.example", mode, 4, []byte("legacy-voip"), ""))
					conflict := version != "4" && (version != "" || mode == "p8")
					if conflict {
						mock.ExpectRollback()
					} else {
						mock.ExpectExec("UPDATE `ioscertificates` SET").WillReturnResult(sqlmock.NewResult(0, 1))
						mock.ExpectCommit()
					}
					// Explicit P12 also proves P8->P12 cannot bypass the version check.
					fields := map[string]string{"app_key": "app-1", "package": "com.example", "original_package": "com.example", "auth_type": "p12"}
					if version != "" {
						fields["config_version"] = version
					}
					var request *http.Request
					handler := SetIosPushConf
					if upload {
						request = multipartPushRequest(t, "/", fields, nil)
						handler = UploadIosCer
					} else {
						body := `{"app_key":"app-1","package":"com.example","original_package":"com.example","auth_type":"p12"`
						if version != "" {
							body += `,"config_version":` + version
						}
						request = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body+"}"))
						request.Header.Set("Content-Type", "application/json")
					}
					response := invokePushHandler(t, request, handler)
					if conflict {
						if response.Code != http.StatusConflict || responseCode(t, response) != 409 {
							t.Fatal("expected safe 409")
						}
						if strings.Contains(response.Body.String(), "legacy-voip") {
							t.Fatal("conflict leaked credentials")
						}
					} else if responseCode(t, response) != 0 {
						t.Fatal("compatible metadata update rejected")
					}
					if err := mock.ExpectationsWereMet(); err != nil {
						t.Fatal(err)
					}
				})
			}
		}
	}
}

func TestIosMultipartRejectsMalformedVersion(t *testing.T) {
	for _, version := range []string{"", "no", "1.5", "9223372036854775808"} {
		mock, cleanup := openPushHandlerMockDB(t)
		request := multipartPushRequest(t, "/", map[string]string{"app_key": "app-1", "package": "com.example", "original_package": "com.example", "config_version": version}, nil)
		if responseCode(t, invokePushHandler(t, request, UploadIosCer)) != int(errs.AdminErrorCode_ParamError) {
			t.Fatal("malformed version accepted")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
		cleanup()
	}
}

func TestIosP8PresenceMetadata(t *testing.T) {
	for _, raw := range []bool{false, true} {
		row := &dbs.IosCertificateDao{P8KeyName: "saved.p8"}
		if raw {
			row.P8PrivateKey = []byte("raw")
		}
		if iosPushItem(row).HasP8Key != raw {
			t.Fatal("key presence must depend on raw bytes, not the filename")
		}
	}
}

func TestIosP8MissingFilePreservesRawOrRejectsMissingKey(t *testing.T) {
	key := handlerP8(t)
	for _, multipart := range []bool{false, true} {
		for _, missing := range []bool{false, true} {
			mock, cleanup := openPushHandlerMockDB(t)
			var raw []byte
			if !missing {
				raw = key
			}
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WithArgs("app-1", "com.old", 1).
				WillReturnRows(sqlmock.NewRows([]string{"app_key", "package", "auth_type", "p8_key_id", "p8_team_id", "p8_private_key", "p8_key_name", "config_version"}).AddRow("app-1", "com.old", "p8", "ABC1234567", "XYZ1234567", raw, "saved.p8", 6))
			if missing {
				mock.ExpectRollback()
			} else {
				mock.ExpectExec("UPDATE `ioscertificates` SET").WithArgs("p8", "", "", []byte(nil), 7, 1, "ABC1234567", "saved.p8", key, "XYZ1234567", "com.new", []byte(nil), "", "", "app-1", "com.old").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"app_key":"app-1","package":"com.new","original_package":"com.old","is_product":1,"config_version":6}`))
			request.Header.Set("Content-Type", "application/json")
			handler := SetIosPushConf
			if multipart {
				request = multipartPushRequest(t, "/", map[string]string{"app_key": "app-1", "package": "com.new", "original_package": "com.old", "is_product": "1", "config_version": "6"}, nil)
				handler = UploadIosCer
			}
			response := invokePushHandler(t, request, handler)
			if missing {
				if responseCode(t, response) != int(errs.AdminErrorCode_ParamError) {
					t.Fatal("missing key was accepted")
				}
			} else if responseCode(t, response) != 0 {
				t.Fatal("raw key was not retained")
			}
			if strings.Contains(response.Body.String(), "PRIVATE KEY") {
				t.Fatal("credentials leaked")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
			cleanup()
		}
	}
}
