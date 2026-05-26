package dbs

import (
	"errors"
	"fmt"
	"time"

	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm"
)

type ApplicationDao struct {
	ID          int64     `gorm:"primary_key"`
	AppId       string    `gorm:"app_id"`
	AppName     string    `gorm:"app_name"`
	AppIcon     string    `gorm:"app_icon"`
	AppDesc     string    `gorm:"app_desc"`
	AppUrl      string    `gorm:"app_url"`
	AppOrder    int       `gorm:"app_order"`
	CreatedTime time.Time `gorm:"created_time"`
	UpdatedTime time.Time `gorm:"updated_time"`
	AppKey      string    `gorm:"app_key"`
}

func (app ApplicationDao) TableName() string {
	return "applications"
}

func (app ApplicationDao) Create(item ApplicationDao) error {
	return dbcommons.GetDb().Exec(fmt.Sprintf("INSERT INTO %s (app_id,app_name,app_icon,app_desc,app_url,app_order,app_key)VALUES(?,?,?,?,?,?,?)", app.TableName()), item.AppId, item.AppName, item.AppIcon, item.AppDesc, item.AppUrl, item.AppOrder, item.AppKey).Error
}

func (app ApplicationDao) Update(item ApplicationDao) error {
	upd := map[string]interface{}{
		"app_name":  item.AppName,
		"app_icon":  item.AppIcon,
		"app_desc":  item.AppDesc,
		"app_url":   item.AppUrl,
		"app_order": item.AppOrder,
	}
	return dbcommons.GetDb().Model(&ApplicationDao{}).Where("app_key=? and app_id=?", item.AppKey, item.AppId).Updates(upd).Error
}

func (app ApplicationDao) BatchDelete(appkey string, appIds []string) error {
	return dbcommons.GetDb().Where("app_key=? and app_id in (?)", appkey, appIds).Delete(&ApplicationDao{}).Error
}

func (app ApplicationDao) QryApplicationsByPage(appkey string, page, size int64) ([]*ApplicationDao, error) {
	var items []*ApplicationDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Order("app_order asc").Offset(int((page - 1) * size)).Limit(int(size)).Find(&items).Error
	return items, err
}

func (app ApplicationDao) FindByAppId(appkey, appId string) (*ApplicationDao, error) {
	var item ApplicationDao
	err := dbcommons.GetDb().Where("app_key=? and app_id=?", appkey, appId).Take(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
