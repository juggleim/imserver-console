package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/services"
	"github.com/juggleim/imserver-console/services/models"
)

type ZegoConf struct {
	AppKey string                `json:"app_key"`
	Conf   *models.ZegoConfigObj `json:"conf"`
}

func SetZegoConf(ctx *gin.Context) {
	var req ZegoConf
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Conf == nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.SetZegoConf(req.AppKey, req.Conf)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

func GetZegoConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, conf := services.GetZegoConf(appkey)
	if code != errs.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
		})
	} else {
		ctxs.SuccessHttpResp(ctx, conf)
	}
}

type AgoraConf struct {
	AppKey string                 `json:"app_key"`
	Conf   *models.AgoraConfigObj `json:"conf"`
}

func SetAgoraConf(ctx *gin.Context) {
	var req AgoraConf
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Conf == nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.SetAgoraConf(req.AppKey, req.Conf)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

func GetAgoraConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, conf := services.GetAgoraConf(appkey)
	if code != errs.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
		})
	} else {
		ctxs.SuccessHttpResp(ctx, conf)
	}
}

type LivekitConf struct {
	AppKey string                   `json:"app_key"`
	Conf   *models.LivekitConfigObj `json:"conf"`
}

func SetLivekitConf(ctx *gin.Context) {
	var req LivekitConf
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Conf == nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.SetLivekitConf(req.AppKey, req.Conf)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

func GetLivekitConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, conf := services.GetLivekitConf(appkey)
	if code != errs.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
		})
	} else {
		ctxs.SuccessHttpResp(ctx, conf)
	}
}
