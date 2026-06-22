package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"
)

func QryMsgStatistic(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	statTypeStrArr := ctx.QueryArray("stat_type")
	statTypes := []services.StatType{}
	for _, statTypeStr := range statTypeStrArr {
		intVal, err := tools.String2Int64(statTypeStr)
		if err == nil && intVal > 0 {
			statTypes = append(statTypes, services.StatType(intVal))
		}
	}
	if len(statTypes) <= 0 {
		statTypes = append(statTypes, services.StatType_Up)
		statTypes = append(statTypes, services.StatType_Down)
		statTypes = append(statTypes, services.StatType_Dispatch)
	}

	channelTypeStr := ctx.Query("channel_type")
	var channelType int64 = 0
	if channelTypeStr != "" {
		intVal, err := tools.String2Int64(channelTypeStr)
		if err == nil && intVal > 0 {
			channelType = intVal
		}
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	items := services.QryMsgStatistic(appkey, statTypes, int(channelType), start, end)
	ctxs.SuccessHttpResp(ctx, items)
}

func QryMsgRealtimeStatistic(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	statTypeStrArr := ctx.QueryArray("stat_type")
	statTypes := []services.StatType{}
	for _, statTypeStr := range statTypeStrArr {
		intVal, err := tools.String2Int64(statTypeStr)
		if err == nil && intVal > 0 {
			statTypes = append(statTypes, services.StatType(intVal))
		}
	}
	if len(statTypes) <= 0 {
		statTypes = append(statTypes, services.StatType_Up)
		statTypes = append(statTypes, services.StatType_Down)
		statTypes = append(statTypes, services.StatType_Dispatch)
	}

	channelTypeStr := ctx.Query("channel_type")
	var channelType int64 = 0
	if channelTypeStr != "" {
		intVal, err := tools.String2Int64(channelTypeStr)
		if err == nil && intVal > 0 {
			channelType = intVal
		}
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	if start > 0 && end > 0 && end-start > services.ThreeDaysMs() {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	items := services.QryMsgRealtimeStatistic(appkey, statTypes, int(channelType), start, end)
	ctxs.SuccessHttpResp(ctx, items)
}

func QryUserActivities(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	items := services.QryUserActivities(appkey, start, end)
	ctxs.SuccessHttpResp(ctx, items)
}

func QryConnectCount(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	ret := services.QryConnect(appkey, start, end)
	ctxs.SuccessHttpResp(ctx, ret)
}

func QryChrmConnectCount(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	ret := services.QryChrmConnect(appkey, start, end)
	ctxs.SuccessHttpResp(ctx, ret)
}

func QryUserRegiste(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	ret := services.QryUserRegiste(appkey, start, end)
	ctxs.SuccessHttpResp(ctx, ret)
}

func QryMaxConnectCount(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	ret := services.QryMaxConnect(appkey, start, end)
	ctxs.SuccessHttpResp(ctx, ret)
}

func QryMaxChrmConnectCount(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	ret := services.QryMaxChrmConnect(appkey, start, end)
	ctxs.SuccessHttpResp(ctx, ret)
}

func QryMaxChrmConnectCountV2(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	startStr := ctx.Query("start")
	var start int64 = 0
	if startStr != "" {
		intVal, err := tools.String2Int64(startStr)
		if err == nil && intVal > 0 {
			start = intVal
		}
	}
	endStr := ctx.Query("end")
	var end int64 = 0
	if endStr != "" {
		intVal, err := tools.String2Int64(endStr)
		if err == nil && intVal > 0 {
			end = intVal
		}
	}
	ret := services.QryMaxChrmConnectV2(appkey, start, end)
	ctxs.SuccessHttpResp(ctx, ret)
}
