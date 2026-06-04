package services

import (
	"context"
	"encoding/json"
	"time"

	apimodels "github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/imsdk"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

const (
	userExtKeyBotConf     = "bot_conf"
	userExtKeyBotSettings = "bot_settings"
)

func QryUsers(ctx context.Context, appkey, userId, name, offset string, limit int64, isPositive bool) (errs.AdminErrorCode, *apimodels.Users) {
	var startId int64
	if offset != "" {
		startId, _ = tools.DecodeInt(offset)
	}
	ret := &apimodels.Users{
		Items: []*apimodels.User{},
	}
	storage := dbs.UserDao{}
	if userId != "" {
		user, err := storage.FindByUserId(appkey, userId)
		if err == nil && user != nil {
			ret.Items = append(ret.Items, userDaoToAPI(user))
		}
	} else {
		users, err := storage.QryUsers(appkey, name, startId, limit, isPositive)
		if err == nil {
			for _, user := range users {
				ret.Offset, _ = tools.EncodeInt(user.ID)
				ret.Items = append(ret.Items, userDaoToAPI(user))
			}
		}
	}
	return errs.AdminErrorCode_Success, ret
}

func QryBots(ctx context.Context, appkey, userId, name, offset string, limit int64, isPositive bool) (errs.AdminErrorCode, *apimodels.Bots) {
	var startId int64
	if offset != "" {
		startId, _ = tools.DecodeInt(offset)
	}
	ret := &apimodels.Bots{
		Items: []*apimodels.Bot{},
	}
	storage := dbs.UserDao{}
	if userId != "" {
		user, err := storage.FindByUserId(appkey, userId)
		if err == nil && user != nil && user.UserType == int(dbs.UserType_Bot) {
			ret.Items = append(ret.Items, userDaoToBot(appkey, user))
		}
	} else {
		bots, err := storage.QryBots(appkey, name, startId, limit, isPositive)
		if err == nil {
			for _, bot := range bots {
				ret.Offset, _ = tools.EncodeInt(bot.ID)
				ret.Items = append(ret.Items, userDaoToBot(appkey, bot))
			}
		}
	}
	return errs.AdminErrorCode_Success, ret
}

func AddBot(req *apimodels.BotReq) (errs.AdminErrorCode, *apimodels.Bot) {
	sdk := imsdk.GetImSdk(req.AppKey)
	if sdk == nil {
		return errs.AdminErrorCode_AppNotExist, nil
	}
	botInfo := botReqToSdk(req)
	resp, code, _, err := sdk.RegisterBot(botInfo)
	if ec := botSdkErrCode(code, err); ec != errs.AdminErrorCode_Success {
		return ec, nil
	}
	botId := req.BotId
	if resp != nil && resp.UserId != "" {
		botId = resp.UserId
	}
	syncBotLocal(req.AppKey, botId, req)
	return errs.AdminErrorCode_Success, loadBot(req.AppKey, botId, req)
}

func UpdateBot(req *apimodels.BotReq) (errs.AdminErrorCode, *apimodels.Bot) {
	sdk := imsdk.GetImSdk(req.AppKey)
	if sdk == nil {
		return errs.AdminErrorCode_AppNotExist, nil
	}
	storage := dbs.UserDao{}
	user, err := storage.FindByUserId(req.AppKey, req.BotId)
	if err != nil || user == nil || user.UserType != int(dbs.UserType_Bot) {
		return errs.AdminErrorCode_ParamError, nil
	}
	botInfo := botReqToSdk(req)
	code, _, err := sdk.UpdateBot(botInfo)
	if ec := botSdkErrCode(code, err); ec != errs.AdminErrorCode_Success {
		return ec, nil
	}
	syncBotLocal(req.AppKey, req.BotId, req)
	return errs.AdminErrorCode_Success, loadBot(req.AppKey, req.BotId, req)
}

func QryUserInfo(appkey, userId string) *apimodels.User {
	storage := dbs.UserDao{}
	user, err := storage.FindByUserId(appkey, userId)
	if err != nil || user == nil {
		return &apimodels.User{UserId: userId}
	}
	return userDaoToAPI(user)
}

func userDaoToAPI(user *dbs.UserDao) *apimodels.User {
	return &apimodels.User{
		UserId:      user.UserId,
		Nickname:    user.Nickname,
		Avatar:      user.UserPortrait,
		Pinyin:      user.Pinyin,
		UserType:    user.UserType,
		Phone:       user.Phone,
		Email:       user.Email,
		Status:      int32(user.Status),
		Account:     user.LoginAccount,
		CreatedTime: user.CreatedTime.UnixMilli(),
	}
}

func userDaoToBot(appkey string, user *dbs.UserDao) *apimodels.Bot {
	bot := &apimodels.Bot{
		BotId:       user.UserId,
		Nickname:    user.Nickname,
		Avatar:      user.UserPortrait,
		Pinyin:      user.Pinyin,
		UserType:    user.UserType,
		CreatedTime: user.CreatedTime.UnixMilli(),
	}
	extDao := dbs.UserExtDao{}
	extFields, err := extDao.QryExtFieldsByItemKeys(appkey, user.UserId, []string{userExtKeyBotConf, userExtKeyBotSettings})
	if err != nil {
		return bot
	}
	if item, ok := extFields[userExtKeyBotConf]; ok && item.ItemValue != "" {
		var conf apimodels.BotConf
		if json.Unmarshal([]byte(item.ItemValue), &conf) == nil {
			bot.BotConf = &conf
		}
	}
	if item, ok := extFields[userExtKeyBotSettings]; ok && item.ItemValue != "" {
		var settings apimodels.BotSettings
		if json.Unmarshal([]byte(item.ItemValue), &settings) == nil {
			bot.BotSettings = &settings
		}
	}
	return bot
}

func botReqToSdk(req *apimodels.BotReq) juggleimsdk.BotInfo {
	info := juggleimsdk.BotInfo{
		BotId:    req.BotId,
		Nickname: req.Nickname,
		Portrait: req.Avatar,
	}
	if req.BotConf != nil {
		info.BotConf = &juggleimsdk.BotConf{
			BotId:    req.BotConf.BotId,
			Url:      req.BotConf.Url,
			ApiKey:   req.BotConf.ApiKey,
			IsStream: req.BotConf.IsStream,
		}
	}
	if req.BotSettings != nil {
		info.BotSettings = &juggleimsdk.BotSettings{
			OnlyMentioned: req.BotSettings.OnlyMentioned,
		}
	}
	if len(req.ExtFields) > 0 {
		info.ExtFields = req.ExtFields
	} else if req.Pinyin != "" {
		info.ExtFields = map[string]string{"pinyin": req.Pinyin}
	}
	return info
}

func syncBotLocal(appkey, botId string, req *apimodels.BotReq) {
	if botId == "" {
		return
	}
	userDao := dbs.UserDao{}
	if req.Nickname != "" || req.Avatar != "" || req.Pinyin != "" {
		if err := userDao.UpdateBotProfile(appkey, botId, req.Nickname, req.Avatar, req.Pinyin); err != nil {
			logs.NewLogEntity().Error(err.Error())
		}
	}
	upsertBotExts(appkey, botId, req.BotConf, req.BotSettings)
}

func upsertBotExts(appkey, userId string, conf *apimodels.BotConf, settings *apimodels.BotSettings) {
	var items []dbs.UserExtDao
	if conf != nil {
		if conf.BotId == "" {
			conf.BotId = userId
		}
		bs, err := json.Marshal(conf)
		if err == nil {
			items = append(items, dbs.UserExtDao{
				AppKey: appkey, UserId: userId, ItemKey: userExtKeyBotConf, ItemValue: string(bs), ItemType: 0,
			})
		}
	}
	if settings != nil {
		bs, err := json.Marshal(settings)
		if err == nil {
			items = append(items, dbs.UserExtDao{
				AppKey: appkey, UserId: userId, ItemKey: userExtKeyBotSettings, ItemValue: string(bs), ItemType: 0,
			})
		}
	}
	if len(items) == 0 {
		return
	}
	extDao := dbs.UserExtDao{}
	if err := extDao.BatchUpsert(items); err != nil {
		logs.NewLogEntity().Error(err.Error())
	}
}

func loadBot(appkey, botId string, req *apimodels.BotReq) *apimodels.Bot {
	storage := dbs.UserDao{}
	user, err := storage.FindByUserId(appkey, botId)
	if err == nil && user != nil && user.UserType == int(dbs.UserType_Bot) {
		return userDaoToBot(appkey, user)
	}
	bot := &apimodels.Bot{
		BotId:       botId,
		Nickname:    req.Nickname,
		Avatar:      req.Avatar,
		Pinyin:      req.Pinyin,
		UserType:    int(dbs.UserType_Bot),
		BotConf:     req.BotConf,
		BotSettings: req.BotSettings,
	}
	if bot.BotConf != nil && bot.BotConf.BotId == "" {
		bot.BotConf.BotId = botId
	}
	return bot
}

func botSdkErrCode(code juggleimsdk.ApiCode, err error) errs.AdminErrorCode {
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		return errs.AdminErrorCode_ServerErr
	}
	if code == juggleimsdk.ApiCode_Success {
		return errs.AdminErrorCode_Success
	}
	return errs.AdminErrorCode(code)
}

func BanUsers(ctx context.Context, req *apimodels.BanUsersReq) errs.AdminErrorCode {
	userStorage := dbs.UserDao{}
	appkey := req.AppKey
	banUsers := &juggleimsdk.BanUsers{
		Items: []*juggleimsdk.BanUser{},
	}
	for _, user := range req.Items {
		endTime := user.EndTime
		if endTime == 0 && user.EndTimeOffset > 0 {
			endTime = time.Now().UnixMilli() + user.EndTimeOffset
		}
		banUsers.Items = append(banUsers.Items, &juggleimsdk.BanUser{
			UserId:  user.UserId,
			EndTime: endTime,
		})
		_ = userStorage.UpdateStatus(appkey, user.UserId, dbs.UserStatus_Ban)
	}
	if sdk := imsdk.GetImSdk(appkey); sdk != nil {
		_, _, _ = sdk.BanUsers(banUsers)
	}
	return errs.AdminErrorCode_Success
}

func UnBanUsers(ctx context.Context, req *apimodels.BanUsersReq) errs.AdminErrorCode {
	userStorage := dbs.UserDao{}
	banUsers := &juggleimsdk.BanUsers{
		Items: []*juggleimsdk.BanUser{},
	}
	appkey := req.AppKey
	for _, user := range req.Items {
		banUsers.Items = append(banUsers.Items, &juggleimsdk.BanUser{
			UserId: user.UserId,
		})
		_ = userStorage.UpdateStatus(appkey, user.UserId, dbs.UserStatus_Normal)
	}
	if sdk := imsdk.GetImSdk(appkey); sdk != nil {
		_, _, _ = sdk.UnBanUsers(banUsers)
	}
	return errs.AdminErrorCode_Success
}
