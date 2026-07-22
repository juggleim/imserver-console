package services

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/tools"
)

func withVlogServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	prev := configures.Config.ImAdminDomain
	configures.Config.ImAdminDomain = server.URL
	t.Cleanup(func() {
		configures.Config.ImAdminDomain = prev
		server.Close()
	})
	return server
}

func TestQryServerLogsDecodesLinesAndBoundsParams(t *testing.T) {
	var got url.Values
	withVlogServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/console/vlogs/query" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		got = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code":0,"msg":"success","data":{"logs":["{\"action\":\"connect\",\"user_id\":\"u1\",\"timestamp\":1700000000000}","not-json",""]}}`))
	})

	code, logs := QryServerLogs(QryServerLogsReq{
		AppKey:  "appkey1",
		LogType: ServerLogType_Connect,
		UserId:  "u1",
		Session: "sess1",
		Start:   1, // older than the lookback window
		Count:   5000,
	})
	if code != errs.AdminErrorCode_Success {
		t.Fatalf("code = %d, want success", code)
	}
	if len(logs) != 2 {
		t.Fatalf("logs len = %d, want 2 (blank line dropped)", len(logs))
	}
	if logs[0]["action"] != "connect" {
		t.Errorf("first log action = %v, want connect", logs[0]["action"])
	}
	if logs[1]["raw"] != "not-json" {
		t.Errorf("unparsable line = %v, want it kept as raw", logs[1]["raw"])
	}

	if got.Get("log_type") != string(ServerLogType_Connect) {
		t.Errorf("log_type = %s", got.Get("log_type"))
	}
	// the node shards the query by target id and falls back to the user id
	if got.Get("target_id") != "u1" {
		t.Errorf("target_id = %s, want the user id", got.Get("target_id"))
	}
	if got.Get("count") != "1000" {
		t.Errorf("count = %s, want it capped at 1000", got.Get("count"))
	}
	earliest := time.Now().Add(-ServerLogMaxLookback).UnixMilli()
	start, err := tools.String2Int64(got.Get("start"))
	if err != nil {
		t.Fatalf("start is not a number: %s", got.Get("start"))
	}
	if start < earliest-time.Minute.Milliseconds() {
		t.Errorf("start = %d, want it clamped to the 24h lookback (%d)", start, earliest)
	}
}

func TestQryServerLogsRequiresRoutingTarget(t *testing.T) {
	withVlogServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("the im server must not be called without a routing target")
	})

	code, logs := QryServerLogs(QryServerLogsReq{
		AppKey:  "appkey1",
		LogType: ServerLogType_Business,
		Session: "sess1",
	})
	if code != errs.AdminErrorCode_ParamError {
		t.Fatalf("code = %d, want param error", code)
	}
	if logs != nil {
		t.Errorf("logs = %v, want nil", logs)
	}
}

func TestQryServerLogsMapsConcurrencyRejection(t *testing.T) {
	withVlogServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"code":10012,"msg":"too many concurrent requests"}`))
	})

	code, _ := QryServerLogs(QryServerLogsReq{
		AppKey:  "appkey1",
		LogType: ServerLogType_UserConnect,
		UserId:  "u1",
	})
	if code != errs.AdminErrorCode_RequestLimit {
		t.Fatalf("code = %d, want request limit", code)
	}
}
