package dbs

import (
	"errors"
	"strings"

	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AndroidPushConfDao struct {
	AppKey      string `gorm:"app_key" json:"app_key"`
	PushChannel string `gorm:"push_channel" json:"push_channel"`
	Package     string `gorm:"package" json:"package"`
	PushConf    string `gorm:"push_conf" json:"conf_path"`
	PushExt     []byte `gorm:"push_ext" json:"-"`
}

func (conf AndroidPushConfDao) TableName() string {
	return "androidpushconfs"
}

func (conf AndroidPushConfDao) Upsert(item AndroidPushConfDao) error {
	err := dbcommons.GetDb().Exec("INSERT INTO androidpushconfs (app_key,push_channel,package,push_conf,push_ext)VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE push_conf=VALUES(push_conf),push_ext=VALUES(push_ext)",
		item.AppKey, item.PushChannel, item.Package, item.PushConf, item.PushExt).Error
	return normalizePushConfError(err)
}

func (conf AndroidPushConfDao) Create(item AndroidPushConfDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (conf AndroidPushConfDao) Find(appkey, pushChannel string) (*AndroidPushConfDao, error) {
	var item AndroidPushConfDao
	err := dbcommons.GetDb().Where("app_key=? and push_channel=?", appkey, pushChannel).Order("package asc").Take(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (conf AndroidPushConfDao) List(appkey, pushChannel string) ([]*AndroidPushConfDao, error) {
	list := make([]*AndroidPushConfDao, 0)
	err := dbcommons.GetDb().Where("app_key=? and push_channel=?", appkey, pushChannel).Order("package asc").Find(&list).Error
	return list, err
}

func (conf AndroidPushConfDao) FindByIdentity(appkey, pushChannel, packageName string) (*AndroidPushConfDao, error) {
	var item AndroidPushConfDao
	err := dbcommons.GetDb().Where("app_key=? and push_channel=? and package=?", appkey, pushChannel, packageName).Take(&item).Error
	if err != nil {
		return nil, normalizePushConfError(err)
	}
	return &item, nil
}

// Save creates a new package-scoped configuration or updates the row identified
// by originalPackage. Package renames and conflict checks are atomic.
func (conf AndroidPushConfDao) Save(item AndroidPushConfDao, originalPackage string) error {
	item.Package = strings.TrimSpace(item.Package)
	originalPackage = strings.TrimSpace(originalPackage)
	return dbcommons.GetDb().Transaction(func(tx *gorm.DB) error {
		if originalPackage == "" {
			return normalizePushConfError(tx.Create(&item).Error)
		}

		var existing AndroidPushConfDao
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("app_key=? and push_channel=? and package=?", item.AppKey, item.PushChannel, originalPackage).
			Take(&existing).Error
		if err != nil {
			return normalizePushConfError(err)
		}

		result := tx.Model(&AndroidPushConfDao{}).
			Where("app_key=? and push_channel=? and package=?", item.AppKey, item.PushChannel, originalPackage).
			Updates(map[string]any{
				"package":   item.Package,
				"push_conf": item.PushConf,
				"push_ext":  item.PushExt,
			})
		if result.Error != nil {
			return normalizePushConfError(result.Error)
		}
		return nil
	})
}

func (conf AndroidPushConfDao) FindByPackage(appkey, packageName string) ([]*AndroidPushConfDao, error) {
	var list []*AndroidPushConfDao
	err := dbcommons.GetDb().Where("app_key=? and package=?", appkey, packageName).Find(&list).Error
	return list, err
}
