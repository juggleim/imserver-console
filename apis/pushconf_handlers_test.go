package apis

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/juggleim/imserver-console/commons/dbcommons"
	consolelogs "github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/services/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openPushHandlerMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	t.Helper()
	consolelogs.SetLogger(logrus.New(), logrus.New())
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	restore := dbcommons.SetDbForTesting(gormDB)
	return mock, func() {
		restore()
		_ = sqlDB.Close()
	}
}

func invokePushHandler(t *testing.T, request *http.Request, handler gin.HandlerFunc) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request
	handler(ctx)
	return recorder
}

func multipartPushRequest(t *testing.T, path string, fields map[string]string, files map[string]struct {
	name string
	data string
}) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			t.Fatal(err)
		}
	}
	for field, file := range files {
		part, err := writer.CreateFormFile(field, file.name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := part.Write([]byte(file.data)); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, path, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}

func responseCode(t *testing.T, recorder *httptest.ResponseRecorder) int {
	t.Helper()
	var response struct {
		Code int `json:"code"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("invalid response %q: %v", recorder.Body.String(), err)
	}
	return response.Code
}

func TestListAndroidPushConfsMasksSecrets(t *testing.T) {
	mock, cleanup := openPushHandlerMockDB(t)
	defer cleanup()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `androidpushconfs` WHERE app_key=? and push_channel=? ORDER BY package asc")).
		WithArgs("app-1", "Huawei").
		WillReturnRows(sqlmock.NewRows([]string{"app_key", "push_channel", "package", "push_conf", "push_ext"}).
			AddRow("app-1", "Huawei", "com.example", `{"app_id":"123","app_secret":"plain-secret"}`, nil))

	request := httptest.NewRequest(http.MethodGet, "/apps/androidpushconf/list?app_key=app-1&push_channel=huawei", nil)
	recorder := invokePushHandler(t, request, ListAndroidPushConfs)
	if strings.Contains(recorder.Body.String(), "plain-secret") || !strings.Contains(recorder.Body.String(), models.PushSecretMask) {
		t.Fatalf("unsafe list response: %s", recorder.Body.String())
	}
	if responseCode(t, recorder) != 0 {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUploadFcmPushConfPreservesFileWhenEditingWithoutUpload(t *testing.T) {
	mock, cleanup := openPushHandlerMockDB(t)
	defer cleanup()
	row := sqlmock.NewRows([]string{"app_key", "push_channel", "package", "push_conf", "push_ext"}).
		AddRow("app-1", "FCM", "com.old", "old.json", []byte("old-file"))
	mock.ExpectQuery("SELECT .*androidpushconfs.*package=.*LIMIT").
		WithArgs("app-1", "FCM", "com.old", 1).
		WillReturnRows(row)
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*androidpushconfs.*FOR UPDATE").
		WithArgs("app-1", "FCM", "com.old", 1).
		WillReturnRows(sqlmock.NewRows([]string{"app_key", "push_channel", "package", "push_conf", "push_ext"}).
			AddRow("app-1", "FCM", "com.old", "old.json", []byte("old-file")))
	mock.ExpectExec("UPDATE `androidpushconfs` SET").
		WithArgs("com.new", "old.json", []byte("old-file"), "app-1", "FCM", "com.old").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	request := multipartPushRequest(t, "/apps/fcmpushconf/upload", map[string]string{
		"app_key": "app-1", "package": "com.new", "original_package": "com.old",
	}, nil)
	recorder := invokePushHandler(t, request, UploadFcmPushConf)
	if responseCode(t, recorder) != 0 {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUploadFcmPushConfRollsBackDuplicatePackage(t *testing.T) {
	mock, cleanup := openPushHandlerMockDB(t)
	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `androidpushconfs`").
		WillReturnError(&mysqlDriver.MySQLError{Number: 1062, Message: "duplicate package"})
	mock.ExpectRollback()

	request := multipartPushRequest(t, "/apps/fcmpushconf/upload", map[string]string{
		"app_key": "app-1", "package": "com.duplicate",
	}, map[string]struct {
		name string
		data string
	}{"fcm_conf": {name: "firebase.json", data: "{}"}})
	recorder := invokePushHandler(t, request, UploadFcmPushConf)
	if responseCode(t, recorder) != 1019 {
		t.Fatalf("expected conflict response: %s", recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUploadIosCerReplacesOnlySelectedCertificate(t *testing.T) {
	mock, cleanup := openPushHandlerMockDB(t)
	defer cleanup()
	columns := []string{"app_key", "package", "is_product", "cert_pwd", "voip_cert_pwd", "certificate", "cert_path", "voip_cert", "voip_cert_path"}
	mock.ExpectQuery("SELECT .*ioscertificates.*package=.*LIMIT").
		WithArgs("app-1", "com.example", 1).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("app-1", "com.example", 0, "old-cert-password", "old-voip-password", []byte("old-cert"), "old.p12", []byte("old-voip"), "old-voip.p12"))
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").
		WithArgs("app-1", "com.example", 1).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("app-1", "com.example", 0, "old-cert-password", "old-voip-password", []byte("old-cert"), "old.p12", []byte("old-voip"), "old-voip.p12"))
	mock.ExpectExec("UPDATE `ioscertificates` SET").
		WithArgs("new.p12", "new-cert-password", []byte("new-cert"), 1, "com.example", []byte("old-voip"), "old-voip.p12", "old-voip-password", "app-1", "com.example").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	request := multipartPushRequest(t, "/apps/iospushcer/upload", map[string]string{
		"app_key":          "app-1",
		"package":          "com.example",
		"original_package": "com.example",
		"is_product":       "1",
		"cert_pwd":         "new-cert-password",
	}, map[string]struct {
		name string
		data string
	}{"ioscer": {name: "new.p12", data: "new-cert"}})
	recorder := invokePushHandler(t, request, UploadIosCer)
	if responseCode(t, recorder) != 0 {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSetIosPushConfPreservesFilesAndBlankPasswords(t *testing.T) {
	mock, cleanup := openPushHandlerMockDB(t)
	defer cleanup()
	columns := []string{"app_key", "package", "is_product", "cert_pwd", "voip_cert_pwd", "certificate", "cert_path", "voip_cert", "voip_cert_path"}
	mock.ExpectQuery("SELECT .*ioscertificates.*package=.*LIMIT").
		WithArgs("app-1", "com.example", 1).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("app-1", "com.example", 0, "old-cert-password", "old-voip-password", []byte("old-cert"), "old.p12", []byte("old-voip"), "old-voip.p12"))
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").
		WithArgs("app-1", "com.example", 1).
		WillReturnRows(sqlmock.NewRows(columns).
			AddRow("app-1", "com.example", 0, "old-cert-password", "old-voip-password", []byte("old-cert"), "old.p12", []byte("old-voip"), "old-voip.p12"))
	mock.ExpectExec("UPDATE `ioscertificates` SET").
		WithArgs("old.p12", "old-cert-password", []byte("old-cert"), 1, "com.example", []byte("old-voip"), "old-voip.p12", "old-voip-password", "app-1", "com.example").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	body := `{"app_key":"app-1","package":"com.example","original_package":"com.example","is_product":1,"cert_pwd":"","voip_cert_pwd":""}`
	request := httptest.NewRequest(http.MethodPost, "/apps/iospushcer/set", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := invokePushHandler(t, request, SetIosPushConf)
	if responseCode(t, recorder) != 0 {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
