package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/services"
)

func SetEmailConf(ctx *gin.Context) {
	var req models.EmailConf
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Conf == nil {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.SetEmailConf(ctxs.ToCtx(ctx), &req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func GetEmailConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code, conf := services.GetEmailConf(ctxs.ToCtx(ctx), appkey)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, conf)
}
