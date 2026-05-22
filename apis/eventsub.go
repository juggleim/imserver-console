package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/services"
)

func SetEventSubConfig(ctx *gin.Context) {
	var req services.EventSubConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.SetEventSubConfig(&req)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

func GetEventSubConfig(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, config := services.GetEventSubConfig(appkey)
	if code != errs.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
		})
	} else {
		ctxs.SuccessHttpResp(ctx, config)
	}
}
