package dbs

import (
	"fmt"

	"github.com/juggleim/imserver-console/commons/dbcommons"
)

type UserActivityDao struct {
	ID       int64  `gorm:"primary_key"`
	UserId   string `gorm:"user_id"`
	TimeMark int64  `gorm:"time_mark"`
	Count    int64  `gorm:"count"`
	AppKey   string `gorm:"app_key"`
}

func (stat UserActivityDao) TableName() string {
	return "useractivities"
}

func (stat UserActivityDao) IncrByStep(appkey, userId string, timeMark, step int64) error {
	sql := fmt.Sprintf("insert into %s (user_id,time_mark,count,app_key)values(?,?,?,?) ON DUPLICATE KEY UPDATE count=count+?", stat.TableName())
	return dbcommons.GetDb().Exec(sql, userId, timeMark, step, appkey, step).Error
}

func (stat UserActivityDao) CountUserActivities(appkey string, timeMark int64) int64 {
	var count int64
	err := dbcommons.GetDb().Model(&UserActivityDao{}).Where("app_key=? and time_mark=?", appkey, timeMark).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}
