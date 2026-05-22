package services

import (
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

func SetZegoConf(appkey string, req *models.ZegoConfigObj) errs.AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(appkey, "zego_config", tools.ToJson(req))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success
}

func GetZegoConf(appkey string) (errs.AdminErrorCode, *models.ZegoConfigObj) {
	zegoConf := &models.ZegoConfigObj{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "zego_config")
	if err == nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), zegoConf)
	} else {
		logs.NewLogEntity().Error(err.Error())
	}

	return errs.AdminErrorCode_Success, zegoConf
}

func SetAgoraConf(appkey string, req *models.AgoraConfigObj) errs.AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(appkey, "agora_config", tools.ToJson(req))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success
}

func GetAgoraConf(appkey string) (errs.AdminErrorCode, *models.AgoraConfigObj) {
	agroaConf := &models.AgoraConfigObj{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "agora_config")
	if err == nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), agroaConf)
	} else {
		logs.NewLogEntity().Error(err.Error())
	}

	return errs.AdminErrorCode_Success, agroaConf
}

func SetLivekitConf(appkey string, req *models.LivekitConfigObj) errs.AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(appkey, "livekit_config", tools.ToJson(req))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success
}

func GetLivekitConf(appkey string) (errs.AdminErrorCode, *models.LivekitConfigObj) {
	livekitConf := &models.LivekitConfigObj{}
	dao := dbs.AppExtDao{}
	conf, err := dao.Find(appkey, "livekit_config")
	if err == nil {
		tools.JsonUnMarshal([]byte(conf.AppItemValue), livekitConf)
	} else {
		logs.NewLogEntity().Error(err.Error())
	}

	return errs.AdminErrorCode_Success, livekitConf
}
