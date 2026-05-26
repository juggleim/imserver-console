package apis

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"
)

func AddApplication(ctx *gin.Context) {
	var req models.Application
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, resp := services.AddApplication(ctxs.ToCtx(ctx), &req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, resp)
}

func UpdApplication(ctx *gin.Context) {
	var req models.Application
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.AppId == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.UpdApplication(ctxs.ToCtx(ctx), &req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func DelApplications(ctx *gin.Context) {
	var req models.ApplicationIds
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || len(req.AppIds) <= 0 {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.DelApplications(ctxs.ToCtx(ctx), &req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func QryApplications(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	var page int64 = 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if val, err := tools.String2Int64(pageStr); err == nil {
			page = val
		}
	}
	var size int64 = 20
	if sizeStr := ctx.Query("size"); sizeStr != "" {
		if val, err := tools.String2Int64(sizeStr); err == nil {
			size = val
		}
	}
	isPositiveOrder := false
	if orderStr := ctx.Query("order"); orderStr != "" {
		if order, err := strconv.Atoi(orderStr); err == nil && order > 0 {
			isPositiveOrder = true
		}
	}
	code, resp := services.QryApplications(ctxs.ToCtx(ctx), appkey, page, size, isPositiveOrder)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, resp)
}
