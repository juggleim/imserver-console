package services

import (
	"context"
	"time"

	apimodels "github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
)

func AddApplication(ctx context.Context, application *apimodels.Application) (errs.AdminErrorCode, *apimodels.Application) {
	if application.AppKey == "" {
		return errs.AdminErrorCode_AppNotExist, nil
	}
	storage := dbs.ApplicationDao{}
	appId := tools.GenerateUUIDShort11()
	err := storage.Create(dbs.ApplicationDao{
		AppId:    appId,
		AppName:  application.AppName,
		AppIcon:  application.AppIcon,
		AppDesc:  application.AppDesc,
		AppUrl:   application.AppUrl,
		AppOrder: application.AppOrder,
		AppKey:   application.AppKey,
	})
	if err != nil {
		return errs.AdminErrorCode_ServerErr, nil
	}
	now := time.Now().UnixMilli()
	return errs.AdminErrorCode_Success, &apimodels.Application{
		AppId:       appId,
		AppName:     application.AppName,
		AppIcon:     application.AppIcon,
		AppDesc:     application.AppDesc,
		AppUrl:      application.AppUrl,
		AppOrder:    application.AppOrder,
		CreatedTime: now,
		UpdatedTime: now,
		AppKey:      application.AppKey,
	}
}

func UpdApplication(ctx context.Context, application *apimodels.Application) errs.AdminErrorCode {
	storage := dbs.ApplicationDao{}
	err := storage.Update(dbs.ApplicationDao{
		AppId:    application.AppId,
		AppName:  application.AppName,
		AppIcon:  application.AppIcon,
		AppDesc:  application.AppDesc,
		AppUrl:   application.AppUrl,
		AppOrder: application.AppOrder,
		AppKey:   application.AppKey,
	})
	if err != nil {
		return errs.AdminErrorCode_ServerErr
	}
	return errs.AdminErrorCode_Success
}

func DelApplications(ctx context.Context, appIds *apimodels.ApplicationIds) errs.AdminErrorCode {
	storage := dbs.ApplicationDao{}
	if err := storage.BatchDelete(appIds.AppKey, appIds.AppIds); err != nil {
		return errs.AdminErrorCode_ServerErr
	}
	return errs.AdminErrorCode_Success
}

func QryApplications(ctx context.Context, appkey string, page, size int64, isPositive bool) (errs.AdminErrorCode, *apimodels.Applications) {
	ret := &apimodels.Applications{
		Items: []*apimodels.Application{},
		Page:  int(page),
		Size:  int(size),
	}
	storage := dbs.ApplicationDao{}
	items, err := storage.QryApplicationsByPage(appkey, page, size)
	if err != nil {
		return errs.AdminErrorCode_ServerErr, ret
	}
	for _, item := range items {
		ret.Items = append(ret.Items, &apimodels.Application{
			AppId:       item.AppId,
			AppName:     item.AppName,
			AppIcon:     item.AppIcon,
			AppDesc:     item.AppDesc,
			AppUrl:      item.AppUrl,
			AppOrder:    item.AppOrder,
			CreatedTime: item.CreatedTime.UnixMilli(),
			UpdatedTime: item.UpdatedTime.UnixMilli(),
		})
	}
	return errs.AdminErrorCode_Success, ret
}
