package services

import (
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

type RtcConfReq struct {
	AppKey string   `json:"app_key"`
	Conf   *RtcConf `json:"conf"`
}
type RtcConf struct {
	ZegoConf    *models.ZegoConfigObj    `json:"zego_conf,omitempty"`
	AgoraConf   *models.AgoraConfigObj   `json:"agora_conf,omitempty"`
	LivekitConf *models.LivekitConfigObj `json:"livekit_conf,omitempty"`
}

func SetRtcConf(appkey string, req *RtcConf) errs.AdminErrorCode {
	dao := dbs.AppExtDao{}
	if req.ZegoConf != nil {
		err := dao.CreateOrUpdate(appkey, "zego_config", tools.ToJson(req.ZegoConf))
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	if req.AgoraConf != nil {
		err := dao.CreateOrUpdate(appkey, "agora_config", tools.ToJson(req.AgoraConf))
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	if req.LivekitConf != nil {
		err := dao.CreateOrUpdate(appkey, "livekit_config", tools.ToJson(req.LivekitConf))
		if err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	return errs.AdminErrorCode_Success
}

func GetRtcConf(appkey string) (errs.AdminErrorCode, *RtcConf) {
	ret := &RtcConf{}
	dao := dbs.AppExtDao{}
	exts, err := dao.FindByItemKeys(appkey, []string{"zego_config", "agora_config", "livekit_config"})
	if err == nil {
		for _, ext := range exts {
			if ext.AppItemKey == "zego_config" {
				zegoConf := &models.ZegoConfigObj{}
				err = tools.JsonUnMarshal([]byte(ext.AppItemValue), zegoConf)
				if err == nil {
					ret.ZegoConf = zegoConf
				} else {
					logs.NewLogEntity().Error(err.Error())
				}
			} else if ext.AppItemKey == "agora_config" {
				agoraConf := &models.AgoraConfigObj{}
				err = tools.JsonUnMarshal([]byte(ext.AppItemValue), agoraConf)
				if err == nil {
					ret.AgoraConf = agoraConf
				} else {
					logs.NewLogEntity().Error(err.Error())
				}
			} else if ext.AppItemKey == "livekit_config" {
				livekitConf := &models.LivekitConfigObj{}
				err = tools.JsonUnMarshal([]byte(ext.AppItemValue), livekitConf)
				if err == nil {
					ret.LivekitConf = livekitConf
				} else {
					logs.NewLogEntity().Error(err.Error())
				}
			}
		}
	} else {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success, ret
}
