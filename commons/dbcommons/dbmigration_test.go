package dbcommons

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func migrationMock(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(SetDbForTesting(db))
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	return mock
}

func expectP8Migration(t *testing.T, mock sqlmock.Sqlmock, version string, failAt int, failure error) {
	t.Helper()
	data, err := sqlFs.ReadFile("sqls/" + version + ".sql")
	if err != nil {
		t.Fatal(err)
	}
	index := 0
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		expect := mock.ExpectExec(regexp.QuoteMeta(line))
		if index == failAt {
			expect.WillReturnError(failure)
			return
		}
		expect.WillReturnResult(sqlmock.NewResult(0, 0))
		index++
	}
}

func TestUpgradeVersionRead(t *testing.T) {
	bootstrapStop := errors.New("stop at bootstrap version write")
	for _, test := range []struct {
		name      string
		readErr   error
		bootstrap bool
	}{
		{name: "missing row", bootstrap: true},
		{name: "record not found", readErr: gorm.ErrRecordNotFound, bootstrap: true},
		{name: "missing table", readErr: &mysqlDriver.MySQLError{Number: 1146}, bootstrap: true},
		{name: "wrapped missing table", readErr: fmt.Errorf("lookup: %w", &mysqlDriver.MySQLError{Number: 1146}), bootstrap: true},
		{name: "permission denied", readErr: &mysqlDriver.MySQLError{Number: 1045}},
		{name: "missing database", readErr: &mysqlDriver.MySQLError{Number: 1049}},
		{name: "connection failure", readErr: &mysqlDriver.MySQLError{Number: 2006}},
		{name: "text mentioning 1146", readErr: errors.New("1146 is not a typed MySQL error")},
	} {
		t.Run(test.name, func(t *testing.T) {
			mock := migrationMock(t)
			query := mock.ExpectQuery("SELECT .*globalconfs.*conf_key").WithArgs(JChatDbVersionKey, 1)
			if test.readErr != nil {
				query.WillReturnError(test.readErr)
			} else {
				query.WillReturnRows(sqlmock.NewRows([]string{"conf_value"}))
			}
			want := test.readErr
			if test.bootstrap {
				// The first embedded migration is empty. Reaching its version write
				// proves bootstrap was allowed, without pretending to create a DB.
				mock.ExpectExec("INSERT INTO globalconfs").WithArgs(JChatDbVersionKey, "20260521").WillReturnError(bootstrapStop)
				want = bootstrapStop
			}
			if err := Upgrade(); !errors.Is(err, want) {
				t.Fatalf("unexpected upgrade result: %v", err)
			}
		})
	}
	for _, value := range []string{"", "invalid", "-1", "9223372036854775808", "20260908"} {
		t.Run("version_"+value, func(t *testing.T) {
			mock := migrationMock(t)
			mock.ExpectQuery("SELECT .*globalconfs.*conf_key").WithArgs(JChatDbVersionKey, 1).
				WillReturnRows(sqlmock.NewRows([]string{"conf_value"}).AddRow(value))
			err := Upgrade()
			if value == "20260908" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), "invalid database migration version") {
				t.Fatalf("invalid version did not stop upgrade: %v", err)
			}
		})
	}
}

func TestUpgradeStopsOnSQLOrVersionFailure(t *testing.T) {
	for _, stage := range []string{"sql", "version"} {
		for _, number := range []uint16{1146, 1045} {
			t.Run(fmt.Sprintf("%s_%d", stage, number), func(t *testing.T) {
				mock := migrationMock(t)
				mock.ExpectQuery("SELECT .*globalconfs.*conf_key").WithArgs(JChatDbVersionKey, 1).
					WillReturnRows(sqlmock.NewRows([]string{"conf_value"}).AddRow("20260720"))
				failure := &mysqlDriver.MySQLError{Number: number}
				if stage == "sql" {
					expectP8Migration(t, mock, "20260908", 2, failure)
				} else {
					expectP8Migration(t, mock, "20260908", -1, nil)
					mock.ExpectExec("INSERT INTO globalconfs").WithArgs(JChatDbVersionKey, "20260908").WillReturnError(failure)
				}
				// No remaining statements or version writes are expected after failure.
				if err := Upgrade(); !errors.Is(err, failure) {
					t.Fatalf("migration failure was swallowed: %v", err)
				}
			})
		}
	}
}

func TestUpgradeWritesVersionAfterSuccessfulMigration(t *testing.T) {
	mock := migrationMock(t)
	mock.ExpectQuery("SELECT .*globalconfs.*conf_key").WithArgs(JChatDbVersionKey, 1).
		WillReturnRows(sqlmock.NewRows([]string{"conf_value"}).AddRow("20260720"))
	expectP8Migration(t, mock, "20260908", -1, nil)
	mock.ExpectExec("INSERT INTO globalconfs").WithArgs(JChatDbVersionKey, "20260908").WillReturnResult(sqlmock.NewResult(0, 1))
	if err := Upgrade(); err != nil {
		t.Fatal(err)
	}
}

func TestP8MigrationExecutesEveryStatementAndStopsOnFailure(t *testing.T) {
	filename := "sqls/20260908.sql"
	data, err := sqlFs.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	var statements []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "--") || line == "" {
			continue
		}
		statements = append(statements, line)
	}
	columns := []string{"auth_type", "p8_key_id", "p8_team_id", "p8_private_key", "p8_key_name", "config_version"}
	if strings.Count(string(data), "ALTER TABLE") != len(columns) || !strings.Contains(string(data), "ADD COLUMN p8_private_key BLOB NULL") {
		t.Fatal("P8 migration must only add the six credential columns")
	}
	for _, forbidden := range []string{"DROP ", "UPDATE ", "MODIFY ", "RENAME ", "CHANGE "} {
		if strings.Contains(string(data), forbidden) {
			t.Fatal("raw migration changes historical data or schema")
		}
	}
	if len(statements) != 4*len(columns) {
		t.Fatal("expected independent four-statement column migrations")
	}
	for _, column := range columns {
		if strings.Count(string(data), "COLUMN_NAME = '"+column+"'") != 1 || strings.Count(string(data), "ADD COLUMN "+column+" ") != 1 {
			t.Fatalf("missing idempotent check for %s", column)
		}
	}
	for _, fail := range []bool{false, true} {
		t.Run(map[bool]string{false: "repeat", true: "stop"}[fail], func(t *testing.T) {
			sqlDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer sqlDB.Close()
			db, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB, SkipInitializeWithVersion: true}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
			if err != nil {
				t.Fatal(err)
			}
			restore := SetDbForTesting(db)
			defer restore()
			for range 2 {
				for index, statement := range statements {
					expect := mock.ExpectExec(regexp.QuoteMeta(statement))
					if fail && index == 2 {
						expect.WillReturnError(errors.New("migration rejected"))
						break
					}
					expect.WillReturnResult(sqlmock.NewResult(0, 0))
				}
				err := executeSqlFile(filename)
				if (err != nil) != fail {
					t.Fatal("unexpected migration result")
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
