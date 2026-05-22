package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/services"
)

type ApiBody struct {
	Method string `json:"method"`
	AppKey string `json:"app_key"`
	Path   string `json:"path"`
	Body   string `json:"body"`
}

func ApiAgent(ctx *gin.Context) {
	var req ApiBody
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errs.AdminErrorCode_ParamError)
		return
	}
	appInfo := services.QryApp(req.AppKey)
	if appInfo == nil {
		ctx.JSON(http.StatusForbidden, errs.AdminErrorCode_ParamError)
		return
	}
	httpCode, resp := services.ApiAgent(req.Method, req.Path, req.Body, req.AppKey, appInfo.AppSecret)
	ctx.JSON(httpCode, resp)
}
