package services

import (
	"context"
	"time"

	apimodels "github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/imsdk"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	juggleimsdk "github.com/juggleim/imserver-sdk-go"
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
		Status:      int32(user.Status),
		CreatedTime: user.CreatedTime.UnixMilli(),
	}
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
