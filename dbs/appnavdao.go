package dbs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/juggleim/imserver-console/commons/dbcommons"
	"gorm.io/gorm"
)

const appNavAliasLockName = "imserver-console:appnav-alias"

type AppNavDao struct {
	ID      int64  `gorm:"primary_key"`
	AppKey  string `gorm:"app_key"`
	AliasNo string `gorm:"alias_no"`

	AdminUrl string `gorm:"admin_url"`
	ApiUrl   string `gorm:"api_url"`
	WsUrl    string `gorm:"ws_url"`
	AppUrl   string `gorm:"app_url"`
}

func (app AppNavDao) TableName() string {
	return "appnavs"
}

func (app AppNavDao) FindByAppkey(appkey string) (*AppNavDao, error) {
	var item AppNavDao
	err := dbcommons.GetDb().Where("app_key=?", appkey).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, err
}

func (app AppNavDao) FindByAliasNo(aliasNo string) (*AppNavDao, error) {
	var item AppNavDao
	err := dbcommons.GetDb().Where("alias_no=?", aliasNo).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (app AppNavDao) UpsertAlias(appkey string, aliasNo string) error {
	return dbcommons.GetDb().Exec(
		"INSERT INTO appnavs (app_key,alias_no) VALUES (?,?) ON DUPLICATE KEY UPDATE alias_no=VALUES(alias_no)",
		appkey,
		aliasNo,
	).Error
}

func (app AppNavDao) EnsureNextAlias(appkey string) (string, error) {
	var aliasNo string
	err := dbcommons.GetDb().Connection(func(conn *gorm.DB) (connErr error) {
		var lockResult struct {
			Acquired int `gorm:"column:acquired"`
		}
		if err := conn.Raw("SELECT GET_LOCK(?, 10) AS acquired", appNavAliasLockName).Scan(&lockResult).Error; err != nil {
			return err
		}
		if lockResult.Acquired != 1 {
			return fmt.Errorf("acquire appnav alias lock failed")
		}
		defer func() {
			var releaseResult struct {
				Released int `gorm:"column:released"`
			}
			err := conn.Raw("SELECT RELEASE_LOCK(?) AS released", appNavAliasLockName).Scan(&releaseResult).Error
			if connErr == nil && (err != nil || releaseResult.Released != 1) {
				if err != nil {
					connErr = err
				} else {
					connErr = fmt.Errorf("release appnav alias lock failed")
				}
			}
		}()

		return conn.Transaction(func(tx *gorm.DB) error {
			var item AppNavDao
			findErr := tx.Where("app_key=?", appkey).Take(&item).Error
			if findErr == nil {
				currentAlias := strings.TrimSpace(item.AliasNo)
				if currentAlias != "" && currentAlias != "0" {
					aliasNo = currentAlias
					return nil
				}
			} else if !errors.Is(findErr, gorm.ErrRecordNotFound) {
				return findErr
			}

			var sequence struct {
				MaxAlias int64 `gorm:"column:max_alias"`
			}
			if err := tx.Raw(`
			SELECT COALESCE(MAX(CASE
					WHEN alias_no REGEXP '^[0-9]+$' THEN CAST(alias_no AS UNSIGNED)
					ELSE NULL
				END), 99999) AS max_alias
			FROM appnavs
			`).Scan(&sequence).Error; err != nil {
				return err
			}
			aliasNo = nextAppAliasNo(sequence.MaxAlias)
			if findErr == nil {
				return tx.Model(&AppNavDao{}).Where("app_key=?", appkey).Update("alias_no", aliasNo).Error
			}
			return tx.Create(&AppNavDao{AppKey: appkey, AliasNo: aliasNo}).Error
		})
	})
	if err != nil {
		return "", err
	}
	return aliasNo, nil
}

func nextAppAliasNo(maxAlias int64) string {
	if maxAlias < 99999 {
		maxAlias = 99999
	}
	return strconv.FormatInt(maxAlias+1, 10)
}
