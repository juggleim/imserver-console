package services

import (
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

type SmsConf struct {
	AppKey string                `json:"app_key"`
	Conf   *models.SmsEngineConf `json:"conf"`
}

func SetSmsConf(appkey string, req *models.SmsEngineConf) errs.AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(appkey, "sms_engine_conf", tools.ToJson(req))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success
}

func GetSmsConf(appkey string) (errs.AdminErrorCode, *models.SmsEngineConf) {
	smsConf := &models.SmsEngineConf{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "sms_engine_conf")
	if err == nil && conf != nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), smsConf)
	} else {
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	return errs.AdminErrorCode_Success, smsConf
}
