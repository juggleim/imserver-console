package services

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
)

type ServerLogType string

const (
	ServerLogType_UserConnect ServerLogType = "userconnect"
	ServerLogType_Connect     ServerLogType = "connect"
	ServerLogType_Business    ServerLogType = "business"
)

// A vlog query greps the log files of an im node, so the node fences it:
// at most a 24h lookback, 1000 lines per query and 2 concurrent queries.
// Mirror the bounds here so an out-of-range console request is corrected
// instead of being rejected on the other side.
const (
	ServerLogMaxLookback  = 24 * time.Hour
	ServerLogDefaultCount = 50
	ServerLogMaxCount     = 1000

	// the node's own scan is capped at 4s inside a 5s rpc, so the default 5s
	// http timeout would give up right when a slow query is about to answer
	serverLogQryTimeout = 10 * time.Second
)

type QryServerLogsReq struct {
	AppKey  string
	LogType ServerLogType
	UserId  string
	Session string
	// the node routes the query by target_id, falling back to user_id
	TargetId string
	Index    int64
	Start    int64
	Count    int64
}

type vlogsResp struct {
	Code errs.AdminErrorCode `json:"code"`
	Msg  string              `json:"msg"`
	Data struct {
		Logs []string `json:"logs"`
	} `json:"data"`
}

// QryServerLogs proxies to the im server's /console/vlogs/query. Each returned
// log line is already a json object built by the node, so it is decoded here
// and handed to the console as structured data.
func QryServerLogs(req QryServerLogsReq) (errs.AdminErrorCode, []map[string]interface{}) {
	if req.AppKey == "" || req.LogType == "" {
		return errs.AdminErrorCode_ParamError, nil
	}
	if req.TargetId == "" {
		req.TargetId = req.UserId
	}
	if req.TargetId == "" {
		return errs.AdminErrorCode_ParamError, nil
	}
	if req.Count <= 0 {
		req.Count = ServerLogDefaultCount
	} else if req.Count > ServerLogMaxCount {
		req.Count = ServerLogMaxCount
	}
	earliest := time.Now().Add(-ServerLogMaxLookback).UnixMilli()
	if req.Start < earliest {
		req.Start = earliest
	}

	params := url.Values{}
	params.Set("app_key", req.AppKey)
	params.Set("log_type", string(req.LogType))
	params.Set("user_id", req.UserId)
	params.Set("session", req.Session)
	params.Set("target_id", req.TargetId)
	params.Set("index", tools.Int2String(req.Index))
	params.Set("start", tools.Int2String(req.Start))
	params.Set("count", tools.Int2String(req.Count))

	qryUrl := fmt.Sprintf("%s/console/vlogs/query?%s", configures.Config.ImAdminDomain, params.Encode())
	respBs, httpCode, err := tools.HttpDoBytesWithTimeout(http.MethodGet, qryUrl, GetImConsoleHeaders(), "", serverLogQryTimeout)
	if err != nil {
		logs.NewLogEntity().Error(fmt.Sprintf("qry server logs failed. log_type:%s err:%s", req.LogType, err.Error()))
		return errs.AdminErrorCode_QryServerLogsFail, nil
	}
	if httpCode != http.StatusOK {
		// the node rejects a query when both of its scan slots are busy
		if httpCode == http.StatusTooManyRequests {
			return errs.AdminErrorCode_RequestLimit, nil
		}
		if httpCode == http.StatusBadRequest {
			return errs.AdminErrorCode_ParamError, nil
		}
		logs.NewLogEntity().Error(fmt.Sprintf("qry server logs failed. log_type:%s http_code:%d resp:%s", req.LogType, httpCode, string(respBs)))
		return errs.AdminErrorCode_QryServerLogsFail, nil
	}

	resp := &vlogsResp{}
	if len(respBs) > 0 {
		if err := tools.JsonUnMarshal(respBs, resp); err != nil {
			return errs.AdminErrorCode_QryServerLogsFail, nil
		}
	}
	if resp.Code != errs.AdminErrorCode_Success {
		return errs.AdminErrorCode_QryServerLogsFail, nil
	}

	items := []map[string]interface{}{}
	for _, log := range resp.Data.Logs {
		if log == "" {
			continue
		}
		item := map[string]interface{}{}
		if err := tools.JsonUnMarshal([]byte(log), &item); err != nil {
			// keep an unparsable line visible rather than silently dropping it
			item = map[string]interface{}{"raw": log}
		}
		items = append(items, item)
	}
	return errs.AdminErrorCode_Success, items
}
