package dbs

import (
	"errors"
	"time"

	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm"
)

type UserDao struct {
	ID           int64     `gorm:"primary_key"`
	UserType     int       `gorm:"user_type"`
	UserId       string    `gorm:"user_id"`
	Nickname     string    `gorm:"nickname"`
	UserPortrait string    `gorm:"user_portrait"`
	Pinyin       string    `gorm:"pinyin"`
	Phone        string    `gorm:"phone"`
	Email        string    `gorm:"email"`
	LoginAccount string    `gorm:"login_account"`
	LoginPass    string    `gorm:"login_pass"`
	CreatedTime  time.Time `gorm:"created_time"`
	UpdatedTime  time.Time `gorm:"updated_time"`
	AppKey       string    `gorm:"app_key"`
}

func (user UserDao) TableName() string {
	return "users"
}

func (user UserDao) FindByUserId(appkey, userId string) (*UserDao, error) {
	var item UserDao
	err := dbcommons.GetDb().Where("app_key=? and user_id=?", appkey, userId).Take(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (user UserDao) FindByPhone(appkey, phone string) (*UserDao, error) {
	var item UserDao
	err := dbcommons.GetDb().Where("app_key=? and phone=?", appkey, phone).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (user UserDao) FindByEmail(appkey, email string) (*UserDao, error) {
	var item UserDao
	err := dbcommons.GetDb().Where("app_key=? and email=?", appkey, email).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (user UserDao) Count(appkey string) int64 {
	var count int64
	err := dbcommons.GetDb().Model(&UserDao{}).Where("app_key=?", appkey).Count(&count).Error
	if err != nil {
		return 0
	}
	return count
}

func (user UserDao) CountByTime(appkey string, start, end int64) int64 {
	var count int64
	err := dbcommons.GetDb().Model(&UserDao{}).Where("app_key=? and created_time>=? and created_time<=?", appkey, time.UnixMilli(start), time.UnixMilli(end)).Count(&count).Error
	if err != nil {
		return count
	}
	return count
}
