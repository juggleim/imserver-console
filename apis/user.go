package apis

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"
)

func QryUsers(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	userId := ctx.Query("user_id")
	name := ctx.Query("name")
	offset := ctx.Query("offset")
	var count int64 = 20
	countStr := ctx.Query("count")
	if countStr != "" {
		if val, err := tools.String2Int64(countStr); err == nil {
			count = val
		}
	}
	isPositiveOrder := false
	if orderStr := ctx.Query("order"); orderStr != "" {
		if order, err := strconv.Atoi(orderStr); err == nil && order > 0 {
			isPositiveOrder = true
		}
	}
	code, users := services.QryUsers(ctxs.ToCtx(ctx), appkey, userId, name, offset, count, isPositiveOrder)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, users)
}

func QryBots(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	userId := ctx.Query("user_id")
	name := ctx.Query("name")
	offset := ctx.Query("offset")
	var count int64 = 20
	countStr := ctx.Query("count")
	if countStr != "" {
		if val, err := tools.String2Int64(countStr); err == nil {
			count = val
		}
	}
	isPositiveOrder := false
	if orderStr := ctx.Query("order"); orderStr != "" {
		if order, err := strconv.Atoi(orderStr); err == nil && order > 0 {
			isPositiveOrder = true
		}
	}
	code, bots := services.QryBots(ctxs.ToCtx(ctx), appkey, userId, name, offset, count, isPositiveOrder)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, bots)
}

func BanUsers(ctx *gin.Context) {
	var req models.BanUsersReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.BanUsers(ctxs.ToCtx(ctx), &req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func UnBanUsers(ctx *gin.Context) {
	var req models.BanUsersReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.UnBanUsers(ctxs.ToCtx(ctx), &req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}
