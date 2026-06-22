package apis

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"
)

func QryPerformanceNodes(ctx *gin.Context) {
	items := services.QryPerformanceNodes()
	ctxs.SuccessHttpResp(ctx, items)
}

func QryPerformanceCatalog(ctx *gin.Context) {
	items := services.QryPerformanceCatalog()
	ctxs.SuccessHttpResp(ctx, items)
}

func QryPerformanceMetric(ctx *gin.Context) {
	nodeName := ctx.Query("node_name")
	metricType := ctx.Query("metric_type")
	start, end, ok := parsePerformanceTimeRange(ctx)
	if !ok {
		return
	}
	result, err := services.QryPerformanceMetric(nodeName, metricType, start, end)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  performanceMetricErrorMsg(err),
		})
		return
	}
	ctxs.SuccessHttpResp(ctx, result)
}

func parsePerformanceTimeRange(ctx *gin.Context) (start int64, end int64, ok bool) {
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")
	if startStr == "" || endStr == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return 0, 0, false
	}
	startVal, err := tools.String2Int64(startStr)
	if err != nil || startVal <= 0 {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return 0, 0, false
	}
	endVal, err := tools.String2Int64(endStr)
	if err != nil || endVal <= 0 {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return 0, 0, false
	}
	return startVal, endVal, true
}

func performanceMetricErrorMsg(err error) string {
	switch {
	case errors.Is(err, services.ErrPerformanceMetricTypeInvalid):
		return "metric_type illegal"
	case errors.Is(err, services.ErrPerformanceRangeTooLarge):
		return "time range exceeds 24 hours"
	default:
		return "param illegal"
	}
}
