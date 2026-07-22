package apis

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"
)

type ServerLogsResp struct {
	Logs []map[string]interface{} `json:"logs"`
}

// QryUserConnectLogs lists the connect attempts of one user, newest window
// first. It is the entry point of the connection inspector.
func QryUserConnectLogs(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	userId := strings.TrimSpace(ctx.Query("user_id"))
	if appkey == "" || userId == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "app_key and user_id are required")
		return
	}
	qryServerLogs(ctx, services.QryServerLogsReq{
		AppKey:  appkey,
		LogType: services.ServerLogType_UserConnect,
		UserId:  userId,
		Start:   qryInt64(ctx, "start"),
		Count:   qryInt64(ctx, "count"),
	})
}

// QryConnectLogs returns the signal trace of a single connection. The node
// filters by session but routes the query by user id, so both are required.
func QryConnectLogs(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	session := strings.TrimSpace(ctx.Query("session"))
	userId := strings.TrimSpace(ctx.Query("user_id"))
	if appkey == "" || session == "" || userId == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "app_key, session and user_id are required")
		return
	}
	qryServerLogs(ctx, services.QryServerLogsReq{
		AppKey:  appkey,
		LogType: services.ServerLogType_Connect,
		UserId:  userId,
		Session: session,
		Start:   qryInt64(ctx, "start"),
		Count:   qryInt64(ctx, "count"),
	})
}

// QryBusinessLogs returns how the server processed one signal of a connection,
// identified by its seq index. The node only scans the log file covering
// `start`, so pass the timestamp of the signal itself.
func QryBusinessLogs(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	session := strings.TrimSpace(ctx.Query("session"))
	userId := strings.TrimSpace(ctx.Query("user_id"))
	if appkey == "" || session == "" || userId == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, "app_key, session and user_id are required")
		return
	}
	qryServerLogs(ctx, services.QryServerLogsReq{
		AppKey:  appkey,
		LogType: services.ServerLogType_Business,
		UserId:  userId,
		Session: session,
		Index:   qryInt64(ctx, "index"),
		Start:   qryInt64(ctx, "start"),
		Count:   qryInt64(ctx, "count"),
	})
}

func qryServerLogs(ctx *gin.Context, req services.QryServerLogsReq) {
	ctx.Set(string(ctxs.CtxKey_AppKey), req.AppKey)
	code, logs := services.QryServerLogs(req)
	if code != errs.AdminErrorCode_Success {
		ctxs.FailHttpResp(ctx, code)
		return
	}
	ctxs.SuccessHttpResp(ctx, &ServerLogsResp{Logs: logs})
}

func qryInt64(ctx *gin.Context, key string) int64 {
	val, err := tools.String2Int64(strings.TrimSpace(ctx.Query(key)))
	if err != nil {
		return 0
	}
	return val
}
