package services

import (
	"encoding/json"

	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

type EventSubConfigReq struct {
	AppKey         string                    `json:"app_key"`
	EventSubConfig *models.EventSubConfigObj `json:"event_sub_config"`
	EventSubSwitch *models.EventSubSwitchObj `json:"event_sub_switch"`
}
type EventSubConfigResp struct {
	AppKey         string                    `json:"app_key"`
	EventSubConfig *models.EventSubConfigObj `json:"event_sub_config"`
	EventSubSwitch []*EventSubSwitchModel    `json:"event_sub_switch"`
}
type EventSubSwitchModel struct {
	Name  string               `json:"name"`
	Items []*models.ConfigItem `json:"items"`
}

func SetEventSubConfig(req *EventSubConfigReq) errs.AdminErrorCode {
	dao := dbs.AppExtDao{}
	err := dao.CreateOrUpdate(req.AppKey, "event_sub_config", tools.ToJson(req.EventSubConfig))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	err = dao.CreateOrUpdate(req.AppKey, "event_sub_switch", tools.ToJson(req.EventSubSwitch))
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	return errs.AdminErrorCode_Success
}

func GetEventSubConfig(appkey string) (errs.AdminErrorCode, *EventSubConfigResp) {
	dao := dbs.AppExtDao{}
	ret := &EventSubConfigResp{
		AppKey:         appkey,
		EventSubConfig: &models.EventSubConfigObj{},
		EventSubSwitch: []*EventSubSwitchModel{},
	}
	appExt, err := dao.Find(appkey, "event_sub_config")
	if err == nil && appExt != nil && appExt.AppItemValue != "" {
		json.Unmarshal([]byte(appExt.AppItemValue), ret.EventSubConfig)
	} else if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	appExt, err = dao.Find(appkey, "event_sub_switch")
	switchMap := map[string]int{}
	if err == nil && appExt != nil && appExt.AppItemValue != "" {
		json.Unmarshal([]byte(appExt.AppItemValue), &switchMap)
	} else if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
	//消息订阅
	subSwitchModel := &EventSubSwitchModel{
		Name:  "消息订阅",
		Items: []*models.ConfigItem{},
	}
	subSwitchModel.Items = append(subSwitchModel.Items, &models.ConfigItem{
		Key:   "private_msg_sub_switch",
		Value: getSwitchValue("private_msg_sub_switch", switchMap),
		Name:  "单聊消息订阅",
	})
	subSwitchModel.Items = append(subSwitchModel.Items, &models.ConfigItem{
		Key:   "group_msg_sub_switch",
		Value: getSwitchValue("group_msg_sub_switch", switchMap),
		Name:  "群聊消息订阅",
	})
	subSwitchModel.Items = append(subSwitchModel.Items, &models.ConfigItem{
		Key:   "chatroom_msg_sub_switch",
		Value: getSwitchValue("chatroom_msg_sub_switch", switchMap),
		Name:  "聊天室消息订阅",
	})
	ret.EventSubSwitch = append(ret.EventSubSwitch, subSwitchModel)

	//在线状态订阅
	subSwitchModel = &EventSubSwitchModel{
		Name:  "在线状态",
		Items: []*models.ConfigItem{},
	}
	subSwitchModel.Items = append(subSwitchModel.Items, &models.ConfigItem{
		Key:   "online_sub_switch",
		Value: getSwitchValue("online_sub_switch", switchMap),
		Name:  "上线状态订阅",
	})
	subSwitchModel.Items = append(subSwitchModel.Items, &models.ConfigItem{
		Key:   "offline_sub_switch",
		Value: getSwitchValue("offline_sub_switch", switchMap),
		Name:  "离线状态订阅",
	})
	ret.EventSubSwitch = append(ret.EventSubSwitch, subSwitchModel)

	return errs.AdminErrorCode_Success, ret
}
func getSwitchValue(key string, switchMap map[string]int) int {
	if val, ok := switchMap[key]; ok {
		return val
	} else {
		return 0
	}
}
