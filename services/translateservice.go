package services

import (
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

type TranslateConf struct {
	AppKey string                  `json:"app_key"`
	Conf   *models.TransEngineConf `json:"conf"`
}

func SetTranslateConf(appkey string, req *models.TransEngineConf) errs.AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(appkey, "trans_engine_conf", tools.ToJson(req))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success
}

func GetTranslateConf(appkey string) (errs.AdminErrorCode, *models.TransEngineConf) {
	transConf := &models.TransEngineConf{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "trans_engine_conf")
	if err == nil && conf != nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), transConf)
	} else {
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}

	return errs.AdminErrorCode_Success, transConf
}
