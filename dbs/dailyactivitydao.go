package dbs

import (
	"github.com/juggleim/imserver-console/commons/dbcommons"
)

type DailyActivityDao struct {
	ID       int64  `gorm:"primary_key"`
	TimeMark int64  `gorm:"time_mark"`
	Count    int64  `gorm:"count"`
	AppKey   string `gorm:"app_key"`
}

func (stat DailyActivityDao) TableName() string {
	return "dailyactivities"
}

func (stat DailyActivityDao) QryStats(appkey string, start, end int64) []*DailyActivityDao {
	var items []*DailyActivityDao
	err := dbcommons.GetDb().Where("app_key=? and time_mark>=? and time_mark<=?", appkey, start, end).Limit(1000).Find(&items).Error
	if err == nil {
		return items
	}
	return []*DailyActivityDao{}
}
