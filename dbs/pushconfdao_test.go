package dbs

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func openPushMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	t.Helper()
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

func TestAndroidPushConfListReturnsEveryPackageInOrder(t *testing.T) {
	mock, cleanup := openPushMockDB(t)
	defer cleanup()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `androidpushconfs` WHERE app_key=? and push_channel=? ORDER BY package asc")).
		WithArgs("app-1", "Huawei").
		WillReturnRows(sqlmock.NewRows([]string{"app_key", "push_channel", "package", "push_conf", "push_ext"}).
			AddRow("app-1", "Huawei", "com.example.a", `{"app_id":"a"}`, nil).
			AddRow("app-1", "Huawei", "com.example.b", `{"app_id":"b"}`, nil))

	items, err := (AndroidPushConfDao{}).List("app-1", "Huawei")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].Package != "com.example.a" || items[1].Package != "com.example.b" {
		t.Fatalf("unexpected list: %#v", items)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAndroidPushConfSaveMapsConcurrentDuplicate(t *testing.T) {
	mock, cleanup := openPushMockDB(t)
	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `androidpushconfs`").
		WillReturnError(&mysqlDriver.MySQLError{Number: 1062, Message: "duplicate business key"})
	mock.ExpectRollback()

	err := (AndroidPushConfDao{}).Save(AndroidPushConfDao{
		AppKey: "app-1", PushChannel: "Huawei", Package: "com.example", PushConf: `{}`,
	}, "")
	if !errors.Is(err, ErrPushConfConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAndroidPushConfSaveRenamesOnlyOriginalPackage(t *testing.T) {
	mock, cleanup := openPushMockDB(t)
	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*androidpushconfs.*FOR UPDATE").
		WithArgs("app-1", "Huawei", "com.old", 1).
		WillReturnRows(sqlmock.NewRows([]string{"app_key", "push_channel", "package", "push_conf", "push_ext"}).
			AddRow("app-1", "Huawei", "com.old", `{}`, nil))
	mock.ExpectExec("UPDATE `androidpushconfs` SET").
		WithArgs("com.new", `{"app_id":"new"}`, sqlmock.AnyArg(), "app-1", "Huawei", "com.old").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := (AndroidPushConfDao{}).Save(AndroidPushConfDao{
		AppKey: "app-1", PushChannel: "Huawei", Package: "com.new", PushConf: `{"app_id":"new"}`,
	}, "com.old")
	if err != nil {
		t.Fatal(err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIosCertificateListDoesNotCollapseLegacyRow(t *testing.T) {
	mock, cleanup := openPushMockDB(t)
	defer cleanup()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `ioscertificates` WHERE app_key=? ORDER BY package asc")).
		WithArgs("app-1").
		WillReturnRows(sqlmock.NewRows([]string{"app_key", "package", "cert_path", "cert_pwd", "certificate"}).
			AddRow("app-1", "com.legacy", "legacy.p12", "secret", []byte("certificate")))

	items, err := (IosCertificateDao{}).List("app-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].Package != "com.legacy" {
		t.Fatalf("unexpected legacy list: %#v", items)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
