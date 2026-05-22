package dbs

import (
	"errors"
	"fmt"
	"time"

	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm"
)

type AppStatus int
type AppType int

var (
	AppStatus_Normal AppStatus = 0
	AppStatus_Block  AppStatus = 1
	AppStatus_Expire AppStatus = 2

	AppType_Private AppType = 0
	AppType_Alone   AppType = 1
	AppType_Public  AppType = 2
)

type AppInfoDao struct {
	ID           int64     `gorm:"primary_key"`
	AppName      string    `gorm:"app_name"`
	AppKey       string    `gorm:"app_key"`
	AppSecret    string    `gorm:"app_secret"`
	AppSecureKey string    `gorm:"app_secure_key"`
	AppStatus    int       `gorm:"app_status"`
	AppType      int       `gorm:"app_type"`
	CreatedTime  time.Time `gorm:"created_time"`
	UpdatedTime  time.Time `gorm:"updated_time"`

	// License string `gorm:"license"`
	LicConf string `gorm:"lic_conf"`
}

func (app AppInfoDao) TableName() string {
	return "apps"
}

func (app AppInfoDao) Create(item AppInfoDao) error {
	err := dbcommons.GetDb().Create(&item).Error
	return err
}

func (app AppInfoDao) Upsert(item AppInfoDao) error {
	sql := fmt.Sprintf("INSERT INTO %s (app_name,app_key,app_secret,app_secure_key,app_type,lic_conf)VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE lic_conf=?", app.TableName())
	return dbcommons.GetDb().Exec(sql, item.AppName, item.AppKey, item.AppSecret, item.AppSecureKey, item.AppType, item.LicConf, item.LicConf).Error
}

func (app AppInfoDao) FindByAppkey(appkey string) (*AppInfoDao, error) {
	var appItem AppInfoDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Take(&appItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &appItem, nil
}

func (app AppInfoDao) FindById(id int64) *AppInfoDao {
	var appItem AppInfoDao
	err := dbcommons.GetDb().Where("id=?", id).Take(&appItem).Error
	if err != nil {
		return nil
	}
	return &appItem
}

func (app AppInfoDao) QryApps(limit int64, offset int64) ([]*AppInfoDao, error) {
	var list []*AppInfoDao
	err := dbcommons.GetDb().Where("id < ?", offset).Order("id desc").Limit(int(limit)).Find(&list).Error
	return list, err
}

func (app AppInfoDao) UpdateLicConf(appkey, licConf string) error {
	return dbcommons.GetDb().Model(&AppInfoDao{}).Where("app_key=?", appkey).Update("lic_conf", licConf).Error
}

/*
upd := map[string]interface{}{}
	if nickname != "" {
		upd["nickname"] = nickname
	}
	if userPortrait != "" {
		upd["user_portrait"] = userPortrait
	}
	if len(upd) > 0 {
		upd["updated_time"] = time.Now()
	} else {
		return fmt.Errorf("do nothing")
	}
	err := dbcommons.GetDb().Model(&UserDao{}).Where("app_key=? and user_id=?", appkey, userId).Update(upd).Error
	return err
*/
