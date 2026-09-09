package dbs

import (
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/juggleim/imserver-console/commons/apnscredentials"
)

func TestIosTwoAdministratorsRejectStaleIDsAfterRotation(t *testing.T) {
	oldKey := generatedP8(t)
	newKey := generatedP8(t)
	mock, cleanup := openPushMockDB(t)
	defer cleanup()
	version := int64(7)
	// B saves a replacement key and ID using the same version A opened.
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WithArgs("app-1", "com.example", 1).
		WillReturnRows(sqlmock.NewRows(iosColumns).AddRow("app-1", "com.example", 0, "", nil, "", "p8", "OLD1234567", "XYZ1234567", oldKey, "old.p8", 7))
	mock.ExpectExec("UPDATE `ioscertificates` SET").WithArgs("p8", "", "", []byte(nil), 8, 0, "NEW1234567", "new.p8", newKey, "XYZ1234567", "com.example", []byte(nil), "", "", "app-1", "com.example").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
	if err := (IosCertificateDao{}).Save(IosCertificateDao{AppKey: "app-1", Package: "com.example", AuthType: "p8", P8KeyID: "NEW1234567", P8KeyName: "new.p8", ExpectedConfigVersion: &version}, "com.example", newKey); err != nil {
		t.Fatal(err)
	}
	// A submits only an intended environment change, but carries the old ID.
	// No UPDATE is permitted after the locked read of B's committed version.
	mock.ExpectBegin()
	mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WithArgs("app-1", "com.example", 1).
		WillReturnRows(sqlmock.NewRows(iosColumns).AddRow("app-1", "com.example", 0, "", nil, "", "p8", "NEW1234567", "XYZ1234567", newKey, "new.p8", 8))
	mock.ExpectRollback()
	err := (IosCertificateDao{}).Save(IosCertificateDao{AppKey: "app-1", Package: "com.example", IsProduct: 1, P8KeyID: "OLD1234567", ExpectedConfigVersion: &version}, "com.example")
	if !errors.Is(err, ErrIosVersionConflict) {
		t.Fatal("stale ID edit was not rejected")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIosEditVersionCompatibility(t *testing.T) {
	for _, test := range []struct {
		name, mode, incoming string
		version              *int64
		conflict             bool
	}{
		{"legacy P12 missing", "p12", "", nil, false},
		{"P8 missing", "p8", "", nil, true},
		{"P8 rollback missing", "p8", "p12", nil, true},
		{"P12 switching missing", "p12", "p8", nil, true},
		{"P12 stale", "p12", "", reviewVersion(3), true},
		{"P12 explicit zero", "p12", "", reviewVersion(0), true},
		{"P12 matching", "p12", "", reviewVersion(4), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			mock, cleanup := openPushMockDB(t)
			defer cleanup()
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WillReturnRows(sqlmock.NewRows(iosColumns).AddRow("app-1", "com.example", 0, "", nil, "", test.mode, "", "", nil, "", 4))
			if test.conflict {
				mock.ExpectRollback()
			} else {
				mock.ExpectExec("UPDATE `ioscertificates` SET").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}
			err := (IosCertificateDao{}).Save(IosCertificateDao{AppKey: "app-1", Package: "com.example", IsProduct: 1, AuthType: test.incoming, ExpectedConfigVersion: test.version}, "com.example")
			if test.conflict && !errors.Is(err, ErrIosVersionConflict) || !test.conflict && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func reviewVersion(v int64) *int64 { return &v }

func TestIosLegacyPasswordlessAndVoipOnlyEdits(t *testing.T) {
	for _, voipOnly := range []bool{false, true} {
		for _, replace := range []bool{false, true} {
			mock, cleanup := openPushMockDB(t)
			var cert, voip []byte
			if voipOnly {
				voip = []byte("voip-retained")
			} else {
				cert = []byte("cert-retained")
			}
			patch := IosCertificateDao{AppKey: "app-1", Package: "com.renamed", IsProduct: 1}
			certPath, voipPath := "", ""
			if replace {
				if voipOnly {
					voip = []byte("voip-new")
					patch.VoipCert, patch.VoipCertPath = voip, "new.p12"
					voipPath = "new.p12"
				} else {
					cert = []byte("cert-new")
					patch.Certificate, patch.CertPath = cert, "new.p12"
					certPath = "new.p12"
				}
			}
			mock.ExpectBegin()
			mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WillReturnRows(sqlmock.NewRows([]string{"app_key", "package", "auth_type", "certificate", "voip_cert", "cert_pwd", "voip_cert_pwd", "config_version"}).AddRow("app-1", "com.example", "p12", cert, voip, "", "", 4))
			mock.ExpectExec("UPDATE `ioscertificates` SET").WithArgs("p12", certPath, "", cert, 5, 1, "", "", []byte(nil), "", "com.renamed", voip, voipPath, "", "app-1", "com.example").WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			if err := (IosCertificateDao{}).Save(patch, "com.example"); err != nil {
				t.Fatal(err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
			cleanup()
		}
	}
}

func TestIosBundleRenameBoundsAndRollback(t *testing.T) {
	for _, name := range []string{"com.example app", "com/example", strings.Repeat("a", 101), strings.Repeat("a", 100)} {
		mock, cleanup := openPushMockDB(t)
		mock.ExpectBegin()
		mock.ExpectQuery("SELECT .*ioscertificates.*FOR UPDATE").WithArgs("app-1", "com.original", 1).WillReturnRows(sqlmock.NewRows(iosColumns).AddRow("app-1", "com.original", 0, "", nil, "", "p12", "", "", nil, "", 2))
		valid := apnscredentials.ValidTopic(name)
		if valid {
			mock.ExpectExec("UPDATE `ioscertificates` SET").WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
		} else {
			mock.ExpectRollback()
		}
		err := (IosCertificateDao{}).Save(IosCertificateDao{AppKey: "app-1", Package: name}, "com.original")
		if valid && err != nil || !valid && !errors.Is(err, ErrIosCredentials) {
			t.Fatalf("unexpected result: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatal(err)
		}
		cleanup()
	}
}
