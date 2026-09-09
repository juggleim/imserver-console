package dbs

import (
	"errors"
	"strings"

	"github.com/juggleim/imserver-console/commons/apnscredentials"
	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type IosCertificateDao struct {
	Package     string `gorm:"package" json:"package"`
	Certificate []byte `gorm:"certificate" json:"-"`
	CertPath    string `gorm:"cert_path" json:"cert_path"`
	AppKey      string `gorm:"app_key" json:"app_key"`
	CertPwd     string `gorm:"cert_pwd" json:"-"`
	IsProduct   int    `gorm:"is_product" json:"is_product"`

	VoipCert              []byte `gorm:"voip_cert" json:"-"`
	VoipCertPwd           string `gorm:"voip_cert_pwd" json:"-"`
	VoipCertPath          string `gorm:"voip_cert_path" json:"voip_cert_path"`
	AuthType              string `gorm:"column:auth_type" json:"auth_type"`
	P8KeyID               string `gorm:"column:p8_key_id" json:"p8_key_id"`
	P8TeamID              string `gorm:"column:p8_team_id" json:"p8_team_id"`
	P8PrivateKey          []byte `gorm:"column:p8_private_key" json:"-"`
	P8KeyName             string `gorm:"column:p8_key_name" json:"p8_key_name"`
	ConfigVersion         int64  `gorm:"column:config_version" json:"config_version"`
	ExpectedConfigVersion *int64 `gorm:"-" json:"-"`
}

func (cer IosCertificateDao) TableName() string {
	return "ioscertificates"
}

func (cer IosCertificateDao) FindByPackage(appkey, packageName string) (*IosCertificateDao, error) {
	var item IosCertificateDao
	err := dbcommons.GetDb().Where("app_key=? and package=?", appkey, packageName).Take(&item).Error
	if err != nil {
		return nil, normalizePushConfError(err)
	}
	return &item, nil
}

func (cer IosCertificateDao) Upsert(item IosCertificateDao) error {
	err := dbcommons.GetDb().Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Exec("INSERT INTO ioscertificates (app_key,package,is_product,cert_pwd,voip_cert_pwd,certificate,cert_path,voip_cert,voip_cert_path) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE is_product=VALUES(is_product),cert_pwd=VALUES(cert_pwd),voip_cert_pwd=VALUES(voip_cert_pwd),certificate=VALUES(certificate),cert_path=VALUES(cert_path),voip_cert=VALUES(voip_cert),voip_cert_path=VALUES(voip_cert_path),config_version=config_version+1",
		item.AppKey, item.Package, item.IsProduct, item.CertPwd, item.VoipCertPwd, item.Certificate, item.CertPath, item.VoipCert, item.VoipCertPath).Error
	return normalizePushConfError(err)
}

func (cer IosCertificateDao) List(appkey string) ([]*IosCertificateDao, error) {
	list := make([]*IosCertificateDao, 0)
	err := dbcommons.GetDb().Where("app_key=?", appkey).Order("package asc").Find(&list).Error
	return list, err
}

var ErrIosCredentials = errors.New("invalid iOS push credentials or metadata")
var ErrIosVersionConflict = errors.New("iOS configuration changed; reload before saving")

// Save merges a patch only after locking the current row. P8 uploads are stored
// verbatim after validation.
func (cer IosCertificateDao) Save(item IosCertificateDao, originalPackage string, privateKey ...[]byte) error {
	item.Package = strings.TrimSpace(item.Package)
	originalPackage = strings.TrimSpace(originalPackage)
	// Disable SQL parameter logging even when database debug logging is enabled.
	return dbcommons.GetDb().Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Transaction(func(tx *gorm.DB) error {
		var existing IosCertificateDao
		if originalPackage != "" {
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("app_key=? and package=?", item.AppKey, originalPackage).
				Take(&existing).Error
			if err != nil {
				return normalizePushConfError(err)
			}
			version := existing.ConfigVersion
			if version < 1 {
				version = 1
			}
			p8Edit := existing.AuthType == "p8" || item.AuthType == "p8" || item.P8KeyID != "" || item.P8TeamID != "" || len(privateKey) > 0 && privateKey[0] != nil
			if item.ExpectedConfigVersion != nil && *item.ExpectedConfigVersion != version || item.ExpectedConfigVersion == nil && p8Edit {
				return ErrIosVersionConflict
			}
		}
		previousAuthType := existing.AuthType
		existing.AppKey, existing.Package, existing.IsProduct = item.AppKey, item.Package, item.IsProduct
		if item.AuthType != "" {
			existing.AuthType = item.AuthType
		}
		if existing.AuthType == "" {
			existing.AuthType = "p12"
		}
		if item.CertPwd != "" && item.CertPwd != "********" {
			existing.CertPwd = item.CertPwd
		}
		if item.VoipCertPwd != "" && item.VoipCertPwd != "********" {
			existing.VoipCertPwd = item.VoipCertPwd
		}
		if item.Certificate != nil {
			existing.Certificate, existing.CertPath = item.Certificate, item.CertPath
		}
		if item.VoipCert != nil {
			existing.VoipCert, existing.VoipCertPath = item.VoipCert, item.VoipCertPath
		}
		if item.P8KeyID != "" {
			existing.P8KeyID = item.P8KeyID
		}
		if item.P8TeamID != "" {
			existing.P8TeamID = item.P8TeamID
		}
		if existing.AppKey == "" || !apnscredentials.ValidTopic(existing.Package) || (existing.IsProduct != 0 && existing.IsProduct != 1) {
			return ErrIosCredentials
		}
		var p8 []byte
		if len(privateKey) > 0 {
			p8 = privateKey[0]
		}
		if p8 != nil {
			if len(item.P8KeyName) == 0 || len(item.P8KeyName) > 255 || apnscredentials.Validate(p8, existing.P8KeyID, existing.P8TeamID) != nil {
				return ErrIosCredentials
			}
			existing.P8PrivateKey, existing.P8KeyName = p8, item.P8KeyName
		}
		switch existing.AuthType {
		case "p12":
			// Preserve legacy metadata-only edits, including VoIP-only records and
			// passwordless certificates. Empty passwords are not validity evidence.
			if originalPackage == "" && (len(existing.Certificate) == 0 || existing.CertPath == "" || existing.CertPwd == "" || len(existing.VoipCert) > 0 && existing.VoipCertPwd == "") {
				return ErrIosCredentials
			}
			if item.Certificate != nil && (len(item.Certificate) == 0 || item.CertPath == "") || item.VoipCert != nil && (len(item.VoipCert) == 0 || item.VoipCertPath == "") {
				return ErrIosCredentials
			}
			if previousAuthType == "p8" && len(existing.Certificate) == 0 && len(existing.VoipCert) == 0 {
				return ErrIosCredentials
			}
		case "p8":
		default:
			return ErrIosCredentials
		}
		if existing.AuthType == "p8" || item.P8KeyID != "" || item.P8TeamID != "" || p8 != nil {
			if apnscredentials.Validate(existing.P8PrivateKey, existing.P8KeyID, existing.P8TeamID) != nil {
				return ErrIosCredentials
			}
		}
		if originalPackage == "" {
			existing.ConfigVersion = 1
			return normalizePushConfError(tx.Create(&existing).Error)
		}
		if existing.ConfigVersion < 1 {
			existing.ConfigVersion = 1
		}
		existing.ConfigVersion++
		item = existing
		result := tx.Model(&IosCertificateDao{}).
			Where("app_key=? and package=?", item.AppKey, originalPackage).
			Updates(map[string]any{
				"auth_type":      item.AuthType,
				"p8_key_id":      item.P8KeyID,
				"p8_team_id":     item.P8TeamID,
				"p8_private_key": item.P8PrivateKey,
				"p8_key_name":    item.P8KeyName,
				"config_version": item.ConfigVersion,
				"package":        item.Package,
				"is_product":     item.IsProduct,
				"cert_pwd":       item.CertPwd,
				"voip_cert_pwd":  item.VoipCertPwd,
				"certificate":    item.Certificate,
				"cert_path":      item.CertPath,
				"voip_cert":      item.VoipCert,
				"voip_cert_path": item.VoipCertPath,
			})
		if result.Error != nil {
			return normalizePushConfError(result.Error)
		}
		return nil
	})
}

func (cer IosCertificateDao) Create(item IosCertificateDao) error {
	if item.AuthType == "" {
		item.AuthType = "p12"
	}
	if item.ConfigVersion < 1 {
		item.ConfigVersion = 1
	}
	err := dbcommons.GetDb().Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Create(&item).Error
	return err
}

func (cer IosCertificateDao) Find(appkey string) (*IosCertificateDao, error) {
	var item IosCertificateDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Order("package asc").Take(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
