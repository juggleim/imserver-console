package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/services"
)

func GetFileCred(ctx *gin.Context) {
	req := models.QryFileCredReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code, resp := services.GetFileCred(ctxs.ToCtx(ctx), &req)
	switch code {
	case services.FileCredSuccess:
		ctxs.SuccessHttpResp(ctx, resp)
	case services.FileCredSignErr:
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ServerErr, "file sign error")
	default:
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_NoFileEngine)
	}
}
