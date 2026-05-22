package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"
	"github.com/juggleim/imserver-console/services/models"
	juggleimsdk "github.com/juggleim/imserver-sdk-go"
)

func QryAppInfo(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	appinfo := services.QryApp(appkey)
	ctxs.SuccessHttpResp(ctx, appinfo)
}

func ActiveApp(ctx *gin.Context) {
	var req CreateAppReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.License == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	appinfo, code, _, err := juggleimsdk.ActiveApp(configures.Config.ImApiDomain, req.License)
	if handleSdkErr(ctx, code, err) {
		return
	}
	ctxs.SuccessHttpResp(ctx, appinfo)
}

func CreateApp(ctx *gin.Context) {
	var req models.AppInfo
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, appinfo := services.CreateApp(req)
	if code != errs.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
			Msg:  "",
		})
	} else {
		ctxs.SuccessHttpResp(ctx, appinfo)
	}
}

type CreateAppReq struct {
	License string `json:"license"`
}

func QryApps(ctx *gin.Context) {
	offsetStr := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 50
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	account := ctx.Query("account")
	code, apps := services.QryApps(ctxs.ToCtx(ctx), account, limit, offsetStr)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, apps)
}

func UpdateAppConfigs(ctx *gin.Context) {
	var req services.AppConfigs
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.UpdateAppConfigs(req.AppKey, req.Configs)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

func QryAppConfigs(ctx *gin.Context) {
	var req QryConfigsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, resp := services.QryAppConfigs(req.AppKey, req.ConfigKeys)
	if code != errs.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
		})
	} else {
		ctxs.SuccessHttpResp(ctx, resp)
	}
}

type QryConfigsReq struct {
	AppKey     string   `json:"app_key"`
	ConfigKeys []string `json:"config_keys"`
}
