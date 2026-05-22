package apis

import (
	"net/http"

	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"

	"github.com/gin-gonic/gin"
)

func Login(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code, account := services.CheckLogin(req.Account, req.Password)
	if code == errs.AdminErrorCode_Success {
		authStr, err := generateAuthorization(req.Account)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &ctxs.ApiErrorMsg{
				Code: errs.AdminErrorCode_Default,
				Msg:  "auth fail",
			})
			return
		}
		ctxs.SuccessHttpResp(ctx, &LoginResp{
			Account:       req.Account,
			Authorization: authStr,
			Env:           "private", //public
			// RoleId:        account.RoleId,
			RoleType:     account.RoleType,
			IsCommercial: configures.Config.IsCommercial,
		})
	} else {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
			Msg:  "login failed",
		})
	}
}

type AccountReq struct {
	Account     string `json:"account"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
	// RoleId      int    `json:"role_id"`
	RoleType int `json:"role_type"`
}

type LoginResp struct {
	Account       string `json:"account"`
	Authorization string `json:"authorization"`
	Env           string `json:"env"`
	// RoleId        int    `json:"role_id"`
	RoleType     int  `json:"role_type"`
	IsCommercial bool `json:"is_commercial"`
}

func AddAccount(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.AddAccount(ctxs.ToCtx(ctx), req.Account, req.Password, req.RoleType)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

func UpdPassword(ctx *gin.Context) {
	var req AccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.UpdPassword(req.Account, req.Password, req.NewPassword)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

func DisableAccounts(ctx *gin.Context) {
	var req AccountsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.Accounts) <= 0 {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.DisableAccounts(ctxs.ToCtx(ctx), req.Accounts, req.IsDisable)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

type BindAppsReq struct {
	Account string   `json:"account"`
	AppKeys []string `json:"app_keys"`
}

func BindApps(ctx *gin.Context) {
	var req BindAppsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Account == "" || len(req.AppKeys) <= 0 {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.BindApps(ctxs.ToCtx(ctx), req.Account, req.AppKeys)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func UnBindApps(ctx *gin.Context) {
	var req BindAppsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Account == "" || len(req.AppKeys) <= 0 {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	code := services.UnBindApps(ctxs.ToCtx(ctx), req.Account, req.AppKeys)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func DeleteAccounts(ctx *gin.Context) {
	var req AccountsReq
	if err := ctx.ShouldBindJSON(&req); err != nil || len(req.Accounts) <= 0 {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	code := services.DeleteAccounts(ctxs.ToCtx(ctx), req.Accounts)
	ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
		Code: code,
	})
}

type AccountsReq struct {
	Accounts  []string `json:"accounts"`
	IsDisable int      `json:"is_disable"`
}

func QryAccounts(ctx *gin.Context) {
	offsetStr := ctx.Query("offset")
	limitStr := ctx.Query("limit")
	var limit int64 = 50
	if limitStr != "" {
		intVal, err := tools.String2Int64(limitStr)
		if err == nil && intVal > 0 && intVal <= 100 {
			limit = intVal
		}
	}
	code, accounts := services.QryAccounts(ctxs.ToCtx(ctx), limit, offsetStr)
	if code != errs.AdminErrorCode_Success {
		ctx.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: code,
		})
		return
	}
	ctxs.SuccessHttpResp(ctx, accounts)
}
