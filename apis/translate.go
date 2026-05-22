package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/services"
)

func SetTranslateConf(ctx *gin.Context) {
	var req services.TranslateConf
	if err := ctx.ShouldBindJSON(&req); err != nil || req.AppKey == "" || req.Conf == nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.SetTranslateConf(req.AppKey, req.Conf)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

func GetTranslateConf(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, conf := services.GetTranslateConf(appkey)
	if code != errs.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
		})
	} else {
		ctxs.SuccessHttpResp(ctx, conf)
	}
}
