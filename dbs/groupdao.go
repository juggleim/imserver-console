package dbs

import (
	"errors"
	"time"

	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm"
)

type GroupDao struct {
	ID            int64     `gorm:"primary_key"`
	GroupId       string    `gorm:"group_id"`
	GroupName     string    `gorm:"group_name"`
	GroupPortrait string    `gorm:"group_portrait"`
	CreatorId     string    `gorm:"creator_id"`
	CreatedTime   time.Time `gorm:"created_time"`
	UpdatedTime   time.Time `gorm:"updated_time"`
	AppKey        string    `gorm:"app_key"`
	IsMute        int       `gorm:"is_mute"`
}

func (group GroupDao) TableName() string {
	return "groupinfos"
}

func (group GroupDao) FindById(appkey, groupId string) (*GroupDao, error) {
	var item GroupDao
	err := dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Take(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (group GroupDao) Delete(appkey, groupId string) error {
	return dbcommons.GetDb().Where("app_key=? and group_id=?", appkey, groupId).Delete(&GroupDao{}).Error
}

func (group GroupDao) QryGroups(appkey, name string, startId, limit int64, isPositive bool) ([]*GroupDao, error) {
	var items []*GroupDao
	whereStr := "app_key=?"
	params := []interface{}{appkey}
	orderBy := "id desc"
	if isPositive {
		orderBy = "id asc"
		whereStr = whereStr + " and id>?"
		params = append(params, startId)
	} else if startId > 0 {
		whereStr = whereStr + " and id<?"
		params = append(params, startId)
	}
	if name != "" {
		whereStr = whereStr + " and group_name like ?"
		params = append(params, "%"+name+"%")
	}
	err := dbcommons.GetDb().Where(whereStr, params...).Order(orderBy).Limit(int(limit)).Find(&items).Error
	return items, err
}
