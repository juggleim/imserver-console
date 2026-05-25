package apis

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

func GetAndroidPushConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	pushChannel := ctx.Query("push_channel")
	if appkey == "" || pushChannel == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_Default,
		})
		return
	}
	ret := &models.AndroidPushConf{
		AppKey:      appkey,
		PushChannel: pushChannel,
	}
	dao := dbs.AndroidPushConfDao{}
	androidDao, err := dao.Find(appkey, pushChannel)
	if err == nil && androidDao != nil {
		ret.Package = androidDao.Package
		_ = json.Unmarshal([]byte(androidDao.PushConf), &ret.Extra)
	} else if err != nil {
		logs.WithContext(ctx).Error(err.Error())
	}
	ctxs.SuccessHttpResp(ctx, ret)
}

func SetAndroidPushConf(ctx *gin.Context) {
	var req models.AndroidPushConf
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	pushConf := ""
	pushChannel := req.PushChannel
	var conf any
	if strings.EqualFold(req.PushChannel, string(models.PushChannel_Huawei)) {
		conf = tools.MapToStruct[models.HuaweiPushConf](req.Extra)
		pushChannel = string(models.PushChannel_Huawei)
	} else if strings.EqualFold(req.PushChannel, string(models.PushChannel_Xiaomi)) {
		conf = tools.MapToStruct[models.XiaomiPushConf](req.Extra)
		pushChannel = string(models.PushChannel_Xiaomi)
	} else if strings.EqualFold(req.PushChannel, string(models.PushChannel_OPPO)) {
		conf = tools.MapToStruct[models.OppoPushConf](req.Extra)
		pushChannel = string(models.PushChannel_OPPO)
	} else if strings.EqualFold(req.PushChannel, string(models.PushChannel_VIVO)) {
		conf = tools.MapToStruct[models.VivoPushConf](req.Extra)
		pushChannel = string(models.PushChannel_VIVO)
	} else if strings.EqualFold(req.PushChannel, string(models.PushChannel_Jpush)) {
		conf = tools.MapToStruct[models.JPushConf](req.Extra)
		pushChannel = string(models.PushChannel_Jpush)
	} else if strings.EqualFold(req.PushChannel, string(models.PushChannel_HONOR)) {
		conf = tools.MapToStruct[models.HonorPushConf](req.Extra)
		pushChannel = string(models.PushChannel_HONOR)
	} else if strings.EqualFold(req.PushChannel, string(models.PushChannel_Getui)) {
		conf = tools.MapToStruct[models.GetuiPushConf](req.Extra)
		pushChannel = string(models.PushChannel_Getui)
	}
	pushConf = tools.ToJson(conf)

	//save to db
	dao := dbs.AndroidPushConfDao{}
	err := dao.Upsert(dbs.AndroidPushConfDao{
		AppKey:      req.AppKey,
		PushChannel: pushChannel,
		Package:     req.Package,
		PushConf:    pushConf,
	})
	if err != nil {
		logs.WithContext(ctx).Error(err.Error())
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func GetFcmPushConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_Default,
		})
		return
	}
	dao := dbs.AndroidPushConfDao{}
	fcmConf, err := dao.Find(appkey, string(models.PushChannel_FCM))
	if err != nil {
		logs.WithContext(ctx).Error(err.Error())
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_Default,
		})
		return
	}
	ctxs.SuccessHttpResp(ctx, fcmConf)
}

func UploadFcmPushConf(ctx *gin.Context) {
	appkey := ctx.PostForm("app_key")
	confPath := ctx.PostForm("conf_path")
	fmcPackage := ctx.PostForm("package")
	fileHeader, err := ctx.FormFile("fcm_conf")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_Default,
		})
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_Default,
		})
		return
	}
	defer file.Close()
	fcmConfBs, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_Default,
		})
		return
	}
	//save to db
	dao := dbs.AndroidPushConfDao{}
	err = dao.Upsert(dbs.AndroidPushConfDao{
		AppKey:      appkey,
		PushChannel: string(models.PushChannel_FCM),
		PushExt:     fcmConfBs,
		PushConf:    confPath,
		Package:     fmcPackage,
	})
	if err != nil {
		logs.WithContext(ctx).Error(err.Error())
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_Default,
		})
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}
