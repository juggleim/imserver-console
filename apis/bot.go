package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/imsdk"
	"github.com/juggleim/imserver-console/services"
)

func AddBot(ctx *gin.Context) {
	var req models.BotReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Nickname == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	if imsdk.GetImSdk(req.AppKey) == nil {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "app not found")
		return
	}
	code, bot := services.AddBot(&req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, bot)
}

func UpdateBot(ctx *gin.Context) {
	var req models.BotReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.BotId == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	if imsdk.GetImSdk(req.AppKey) == nil {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "app not found")
		return
	}
	code, bot := services.UpdateBot(&req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, bot)
}
