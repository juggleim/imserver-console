package services

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

var appFieldsMap map[string]bool

func init() {
	appFieldsMap = make(map[string]bool)
	appFieldsMap["is_hide_msg_before_join_group"] = true
	appFieldsMap["file_config"] = true
	appFieldsMap["event_sub_config"] = true
	appFieldsMap["event_sub_switch"] = true
	appFieldsMap["his_msg_save_day"] = true
	appFieldsMap["kick_mode"] = true
}

func CreateApp(appInfo models.AppInfo) (errs.AdminErrorCode, *models.AppInfo) {
	dao := dbs.AppInfoDao{}
	if appInfo.AppKey == "" {
		appInfo.AppKey = tools.RandStr(16)
	}
	dbAppInfo := dao.FindByAppkey(appInfo.AppKey)
	if dbAppInfo != nil && dbAppInfo.AppKey == appInfo.AppKey {
		return errs.AdminErrorCode_AppHasExisted, &models.AppInfo{
			AppName:     dbAppInfo.AppName,
			AppKey:      dbAppInfo.AppKey,
			AppSecret:   dbAppInfo.AppSecret,
			CreatedTime: dbAppInfo.CreatedTime.UnixMilli(),
		}
	}
	if len(appInfo.AppSecret) != 16 {
		appInfo.AppSecret = tools.RandStr(16)
	}
	newApp := dbs.AppInfoDao{
		AppName:      appInfo.AppName,
		AppKey:       appInfo.AppKey,
		AppSecret:    appInfo.AppSecret,
		AppSecureKey: tools.RandStr(16),
		AppType:      appInfo.AppType,
		CreatedTime:  time.Now(),
		UpdatedTime:  time.Now(),
	}
	err := dao.Upsert(newApp)
	if err != nil {
		return errs.AdminErrorCode_AddAppFail, nil
	}
	return errs.AdminErrorCode_Success, &models.AppInfo{
		AppType:   newApp.AppType,
		AppName:   newApp.AppName,
		AppKey:    newApp.AppKey,
		AppSecret: newApp.AppSecret,
	}
}

func QryApps(ctx context.Context, account string, limit int64, offset string) (errs.AdminErrorCode, *models.Apps) {
	curAccount, exist := GetAccountInfo(ctxs.GetAccountFromCtx(ctx))
	if !exist || curAccount == nil {
		return errs.AdminErrorCode_AccountNotExist, nil
	}
	apps := &models.Apps{
		Items:   []*models.SimpleApp{},
		HasMore: false,
		Offset:  "",
	}
	if curAccount.RoleType == RoleType_SuperAdmin && account == "" {
		dao := dbs.AppInfoDao{}
		offsetInt, err := tools.DecodeInt(offset)
		if err != nil {
			offsetInt = math.MaxInt64
		}
		dbApps, err := dao.QryApps(limit+1, offsetInt)
		if err == nil {
			if len(dbApps) > int(limit) {
				dbApps = dbApps[:len(dbApps)-1]
				apps.HasMore = true
			}
			var id int64 = math.MaxInt64
			for _, dbApp := range dbApps {
				app := &models.SimpleApp{
					AppKey:       dbApp.AppKey,
					AppName:      dbApp.AppName,
					CreatedTime:  dbApp.CreatedTime.UnixMilli(),
					AppType:      dbApp.AppType,
					MaxUserCount: 100,
				}
				userDdao := dbs.UserDao{}
				app.CurUserCount = userDdao.Count(dbApp.AppKey)
				fillSimpleAppQuotaFromApi(app)
				apps.Items = append(apps.Items, app)
				if dbApp.ID < id {
					id = dbApp.ID
				}
			}
			if id > 0 {
				offset, _ := tools.EncodeInt(id)
				apps.Offset = offset
			}
		} else {
			logs.NewLogEntity().Error(err.Error())
		}
	} else {
		acc := curAccount.Account
		if curAccount.RoleType == RoleType_SuperAdmin && account != "" {
			acc = account
		}
		dao := dbs.AccountAppRelDao{}
		offsetInt, err := tools.DecodeInt(offset)
		if err != nil {
			offsetInt = math.MaxInt64
		}
		dbApps, err := dao.QryApps(acc, limit+1, offsetInt)
		if err == nil {
			if len(dbApps) > int(limit) {
				dbApps = dbApps[:len(dbApps)-1]
				apps.HasMore = true
			}
			var id int64 = math.MaxInt64
			for _, dbApp := range dbApps {
				app := &models.SimpleApp{
					AppKey:       dbApp.AppKey,
					AppName:      dbApp.AppName,
					CreatedTime:  dbApp.CreatedTime.UnixMilli(),
					AppType:      dbApp.AppType,
					MaxUserCount: 100,
				}
				userDdao := dbs.UserDao{}
				app.CurUserCount = userDdao.Count(dbApp.AppKey)
				fillSimpleAppQuotaFromApi(app)
				apps.Items = append(apps.Items, app)
				if dbApp.ID < id {
					id = dbApp.ID
				}
			}
			if id > 0 {
				offset, _ := tools.EncodeInt(id)
				apps.Offset = offset
			}
		} else {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	return errs.AdminErrorCode_Success, apps
}

func QryApp(appkey string) *models.AppInfo {
	dao := dbs.AppInfoDao{}
	dbApp := dao.FindByAppkey(appkey)
	if dbApp == nil {
		return nil
	}
	appInfo := &models.AppInfo{
		AppType:      dbApp.AppType,
		AppName:      dbApp.AppName,
		AppKey:       dbApp.AppKey,
		AppSecret:    dbApp.AppSecret,
		CreatedTime:  dbApp.CreatedTime.UnixMilli(),
		UpdateTime:   dbApp.UpdatedTime.UnixMilli(),
		AppStatus:    dbApp.AppStatus,
		ConfigFields: make(map[string]string),
		MaxUserCount: 100,
	}
	userDao := dbs.UserDao{}
	appInfo.CurUserCount = userDao.Count(dbApp.AppKey)
	fillAppInfoQuotaFromApi(appInfo)
	//appext
	extDao := dbs.AppExtDao{}
	dbExts := extDao.FindListByAppkey(appkey)
	for _, dbExt := range dbExts {
		appInfo.ConfigFields[dbExt.AppItemKey] = dbExt.AppItemValue
	}

	return appInfo
}

func UpdateAppConfigs(appkey string, configFields map[string]interface{}) errs.AdminErrorCode {
	//check fields
	// for fieldKey, _ := range configFields {
	// 	if _, exist := appFieldsMap[fieldKey]; !exist {
	// 		return AdminErrorCode_NotSupportField
	// 	}
	// }
	dao := dbs.AppExtDao{}
	for fieldKey, fieldValue := range configFields {
		err := dao.CreateOrUpdate(appkey, fieldKey, fmt.Sprintf("%s", fieldValue))
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	return errs.AdminErrorCode_Success
}

func QryAppConfigs(appkey string, keys []string) (errs.AdminErrorCode, *AppConfigs) {
	ret := &AppConfigs{
		AppKey:  appkey,
		Configs: map[string]interface{}{},
	}
	dao := dbs.AppExtDao{}
	extList, err := dao.FindByItemKeys(appkey, keys)
	extMap := map[string]string{}
	if err == nil {
		for _, ext := range extList {
			extMap[ext.AppItemKey] = ext.AppItemValue
		}
	} else {
		logs.NewLogEntity().Error(err.Error())
	}
	for _, key := range keys {
		if val, ok := extMap[key]; ok {
			ret.Configs[key] = val
		} else {
			ret.Configs[key] = ""
		}
	}
	return errs.AdminErrorCode_Success, ret
}

type AppConfigs struct {
	AppKey  string                 `json:"app_key"`
	Configs map[string]interface{} `json:"configs"`
}

func fillSimpleAppQuotaFromApi(app *models.SimpleApp) {
	if app == nil || app.AppKey == "" || !configures.Config.IsCommercial {
		return
	}
	info := qryRemoteAppInfo(app.AppKey)
	if info == nil {
		return
	}
	if info.ExpiredTime > 0 {
		app.EndedTime = info.ExpiredTime
	}
	if info.MaxUserCount > 0 {
		app.MaxUserCount = info.MaxUserCount
	}
}

func fillAppInfoQuotaFromApi(appInfo *models.AppInfo) {
	if appInfo == nil || appInfo.AppKey == "" {
		return
	}
	info := qryRemoteAppInfo(appInfo.AppKey)
	if info == nil {
		return
	}
	if info.ExpiredTime > 0 {
		appInfo.ExpiredTime = info.ExpiredTime
	}
	if info.MaxUserCount > 0 {
		appInfo.MaxUserCount = info.MaxUserCount
	}
}

func qryRemoteAppInfo(appKey string) *models.AppInfo {
	url := fmt.Sprintf("%s/console/apps/info?app_key=%s", configures.Config.ImAdminDomain, appKey)
	respBs, code, err := tools.HttpDoBytes(http.MethodGet, url, GetImConsoleHeaders(), "")
	if err != nil || code != 200 {
		return nil
	}
	resp := &models.AppInfoResp{}
	if len(respBs) > 0 {
		err = tools.JsonUnMarshal(respBs, resp)
		if err != nil {
			return nil
		}
	}
	return resp.Data
}

type ActiveAppReq struct {
	License string `json:"license"`
}

func ActiveApp(req models.ActiveAppReq) (errs.AdminErrorCode, *models.AppInfo) {
	url := fmt.Sprintf("%s/console/apps/active", configures.Config.ImAdminDomain)
	respBs, code, err := tools.HttpDoBytes(http.MethodPost, url, GetImConsoleHeaders(), tools.ToJson(req))
	if err != nil || code != 200 {
		return errs.AdminErrorCode_AddAppFail, nil
	}
	resp := &models.AppInfoResp{}
	if len(respBs) > 0 {
		err = tools.JsonUnMarshal(respBs, resp)
		if err != nil {
			return errs.AdminErrorCode_AddAppFail, nil
		}
	}
	return errs.AdminErrorCode_Success, resp.Data
}

func GetImConsoleHeaders() map[string]string {
	headers := map[string]string{}
	timestamp := tools.Int2String(time.Now().Unix())
	nonce := tools.RandStr(5)
	headers["timestamp"] = timestamp
	headers["nonce"] = nonce
	secret := configures.Config.AdminSecret
	str := fmt.Sprintf("%s%s%s", secret, nonce, timestamp)
	sig := tools.SHA1(str)
	headers["signature"] = sig
	return headers
}
