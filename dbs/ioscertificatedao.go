package dbs

import (
	"errors"
	"strings"

	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IosCertificateDao struct {
	Package     string `gorm:"package" json:"package"`
	Certificate []byte `gorm:"certificate" json:"certificate"`
	CertPath    string `gorm:"cert_path" json:"cert_path"`
	AppKey      string `gorm:"app_key" json:"app_key"`
	CertPwd     string `gorm:"cert_pwd" json:"cert_pwd"`
	IsProduct   int    `gorm:"is_product" json:"is_product"`

	VoipCert     []byte `gorm:"voip_cert" json:"voip_cert"`
	VoipCertPwd  string `gorm:"voip_cert_pwd" json:"voip_cert_pwd"`
	VoipCertPath string `gorm:"voip_cert_path" json:"voip_cert_path"`
	// CreatedTime time.Time `gorm:"created_time"`
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
	err := dbcommons.GetDb().Exec("INSERT INTO ioscertificates (app_key,package,is_product,cert_pwd,voip_cert_pwd,certificate,cert_path,voip_cert,voip_cert_path) VALUES (?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE is_product=VALUES(is_product),cert_pwd=VALUES(cert_pwd),voip_cert_pwd=VALUES(voip_cert_pwd),certificate=VALUES(certificate),cert_path=VALUES(cert_path),voip_cert=VALUES(voip_cert),voip_cert_path=VALUES(voip_cert_path)",
		item.AppKey, item.Package, item.IsProduct, item.CertPwd, item.VoipCertPwd, item.Certificate, item.CertPath, item.VoipCert, item.VoipCertPath).Error
	return normalizePushConfError(err)
}

func (cer IosCertificateDao) List(appkey string) ([]*IosCertificateDao, error) {
	list := make([]*IosCertificateDao, 0)
	err := dbcommons.GetDb().Where("app_key=?", appkey).Order("package asc").Find(&list).Error
	return list, err
}

// Save creates a new package-scoped certificate or updates the row identified
// by originalPackage. Callers must merge file bytes they intend to preserve.
func (cer IosCertificateDao) Save(item IosCertificateDao, originalPackage string) error {
	item.Package = strings.TrimSpace(item.Package)
	originalPackage = strings.TrimSpace(originalPackage)
	return dbcommons.GetDb().Transaction(func(tx *gorm.DB) error {
		if originalPackage == "" {
			return normalizePushConfError(tx.Create(&item).Error)
		}

		var existing IosCertificateDao
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("app_key=? and package=?", item.AppKey, originalPackage).
			Take(&existing).Error
		if err != nil {
			return normalizePushConfError(err)
		}

		result := tx.Model(&IosCertificateDao{}).
			Where("app_key=? and package=?", item.AppKey, originalPackage).
			Updates(map[string]any{
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
	err := dbcommons.GetDb().Create(&item).Error
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
