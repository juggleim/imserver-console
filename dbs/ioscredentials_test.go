package dbs

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm/logger"
)

func generatedP8(t *testing.T) []byte {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
}

var iosColumns = []string{"app_key", "package", "is_product", "cert_pwd", "certificate", "cert_path", "auth_type", "p8_key_id", "p8_team_id", "p8_private_key", "p8_key_name", "config_version"}

func TestIosP8CreateStoresExactBytesAndMapsDuplicate(t *testing.T) {
	key := append([]byte("\r\n  "), generatedP8(t)...)
	key = append(key, []byte(" \r\n")...)
	for _, duplicate := range []bool{false, true} {
		t.Run(map[bool]string{false: "create", true: "duplicate"}[duplicate], func(t *testing.T) {
			mock, cleanup := openPushMockDB(t)
			defer cleanup()
			mock.ExpectBegin()
			exec := mock.ExpectExec("INSERT INTO `ioscertificates`").WithArgs("com.example", []byte(nil), "", "app-1", "", 1, []byte(nil), "", "", "p8", "ABC1234567", "XYZ1234567", key, "key.p8", 1)
			if duplicate {
				exec.WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				mock.ExpectRollback()
			} else {
				exec.WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}
			err := (IosCertificateDao{}).Save(IosCertificateDao{AppKey: "app-1", Package: "com.example", IsProduct: 1, AuthType: "p8", P8KeyID: "ABC1234567", P8TeamID: "XYZ1234567", P8KeyName: "key.p8"}, "", key)
			if duplicate && !errors.Is(err, ErrPushConfConflict) || !duplicate && err != nil {
				t.Fatalf("unexpected result: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIosP8LockedMergeRotationAndRollback(t *testing.T) {
	key := generatedP8(t)
	for _, mode := range []string{"", "p8", "p12"} {
		t.Run("mode_"+mode, func(t *testing.T) {
			mock, cleanup := openPushMockDB(t)
			defer cleanup()
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WithArgs("app-1", "com.old", 1).
				WillReturnRows(sqlmock.NewRows(iosColumns).AddRow("app-1", "com.old", 0, "retained-password", []byte("retained-cert"), "old.p12", "p8", "ABC1234567", "XYZ1234567", key, "latest.p8", 9))
			active := mode
			if active == "" {
				active = "p8"
			}
			stored := key
			name := "latest.p8"
			version := int64(9)
			patch := IosCertificateDao{AppKey: "app-1", Package: "com.new", IsProduct: 1, AuthType: mode, ExpectedConfigVersion: &version}
			var replacement []byte
			if mode == "p8" {
				replacement = generatedP8(t)
				stored = replacement
				patch.P8KeyName = "rotated.p8"
				name = patch.P8KeyName
			}
			exec := mock.ExpectExec("UPDATE `ioscertificates` SET")
			exec.WithArgs(active, "old.p12", "retained-password", []byte("retained-cert"), 10, 1, "ABC1234567", name, stored, "XYZ1234567", "com.new", []byte(nil), "", "", "app-1", "com.old")
			exec.WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			if err := (IosCertificateDao{}).Save(patch, "com.old", replacement); err != nil {
				t.Fatal(err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIosP8InvalidMergedConfigurationRollsBack(t *testing.T) {
	key := generatedP8(t)
	for _, reason := range []string{"metadata", "environment", "unknown mode", "tampered key", "invalid upload", "missing key"} {
		t.Run(reason, func(t *testing.T) {
			mock, cleanup := openPushMockDB(t)
			defer cleanup()
			stored := bytes.Clone(key)
			version := int64(3)
			patch := IosCertificateDao{AppKey: "app-1", Package: "com.example", AuthType: "p8", ExpectedConfigVersion: &version}
			var upload []byte
			switch reason {
			case "metadata":
				patch.P8TeamID = "invalid"
			case "environment":
				patch.IsProduct = 2
			case "unknown mode":
				patch.AuthType = "other"
			case "tampered key":
				stored[0] ^= 1
			case "invalid upload":
				upload = []byte("private-invalid")
				patch.P8KeyName = "new.p8"
			case "missing key":
				stored = nil
			}
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WithArgs("app-1", "com.example", 1).
				WillReturnRows(sqlmock.NewRows(iosColumns).AddRow("app-1", "com.example", 0, "retained-password", []byte("retained-cert"), "old.p12", "p8", "ABC1234567", "XYZ1234567", stored, "key.p8", 3))
			mock.ExpectRollback()
			if !errors.Is((IosCertificateDao{}).Save(patch, "com.example", upload), ErrIosCredentials) {
				t.Fatal("invalid config accepted or wrong failure")
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIosP8RawTransitions(t *testing.T) {
	key := generatedP8(t)
	for _, test := range []struct {
		name, mode  string
		raw, upload []byte
		wantErr     error
	}{
		{"missing key rejected", "p8", nil, nil, ErrIosCredentials},
		{"missing key uploaded", "p8", nil, key, nil},
		{"raw retained", "p8", key, nil, nil},
		{"invalid raw never falls back", "p8", []byte("invalid-raw"), nil, ErrIosCredentials},
		{"P12 remains editable without P8", "p12", nil, nil, nil},
		{"P12 preserves inactive raw", "p12", key, nil, nil},
	} {
		t.Run(test.name, func(t *testing.T) {
			mock, cleanup := openPushMockDB(t)
			defer cleanup()
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WithArgs("app-1", "com.old", 1).
				WillReturnRows(sqlmock.NewRows(iosColumns).AddRow("app-1", "com.old", 0, "pwd", []byte("cert"), "cert.p12", test.mode, "ABC1234567", "XYZ1234567", test.raw, "saved.p8", 4))
			patch := IosCertificateDao{AppKey: "app-1", Package: "com.new", IsProduct: 1, ExpectedConfigVersion: reviewVersion(4)}
			name, raw := "saved.p8", test.raw
			if test.upload != nil {
				patch.P8KeyName, name, raw = "new.p8", "new.p8", test.upload
			}
			if test.wantErr != nil {
				mock.ExpectRollback()
			} else {
				query := "UPDATE `ioscertificates` SET `auth_type`=?,`cert_path`=?,`cert_pwd`=?,`certificate`=?,`config_version`=?,`is_product`=?,`p8_key_id`=?,`p8_key_name`=?,`p8_private_key`=?,`p8_team_id`=?,`package`=?,`voip_cert`=?,`voip_cert_path`=?,`voip_cert_pwd`=? WHERE app_key=? and package=?"
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(test.mode, "cert.p12", "pwd", []byte("cert"), 5, 1, "ABC1234567", name, raw, "XYZ1234567", "com.new", []byte(nil), "", "", "app-1", "com.old").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}
			if err := (IosCertificateDao{}).Save(patch, "com.old", test.upload); !errors.Is(err, test.wantErr) {
				t.Fatalf("unexpected result: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIosRawSecretsExcludedFromJSONAndSQLLogs(t *testing.T) {
	key := generatedP8(t)
	row := IosCertificateDao{P8PrivateKey: key, Certificate: []byte("cert-secret"), CertPwd: "password-secret", VoipCert: []byte("voip-secret"), VoipCertPwd: "voip-password"}
	data, err := json.Marshal(row)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"p8_private_key", "cert_pwd", "voip_cert_pwd", `"certificate"`, `"voip_cert"`} {
		if strings.Contains(string(data), forbidden) {
			t.Fatal("DAO leaked credential field")
		}
	}
	mock, cleanup := openPushMockDB(t)
	defer cleanup()
	var output bytes.Buffer
	dbcommons.GetDb().Logger = logger.New(log.New(&output, "", 0), logger.Config{LogLevel: logger.Info})
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `ioscertificates`").WillReturnError(errors.New("database rejected secret parameters"))
	mock.ExpectRollback()
	err = (IosCertificateDao{}).Save(IosCertificateDao{AppKey: "app-1", Package: "com.example", AuthType: "p8", P8KeyID: "ABC1234567", P8TeamID: "XYZ1234567", P8KeyName: "key.p8"}, "", key)
	if err == nil {
		t.Fatal("expected database failure")
	}
	if output.Len() != 0 {
		t.Fatal("credential SQL was logged")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
