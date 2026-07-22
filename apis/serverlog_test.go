package apis

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/errs"
)

func stubVlogServer(t *testing.T, handler http.HandlerFunc) {
	t.Helper()
	server := httptest.NewServer(handler)
	prev := configures.Config.ImAdminDomain
	configures.Config.ImAdminDomain = server.URL
	t.Cleanup(func() {
		configures.Config.ImAdminDomain = prev
		server.Close()
	})
}

func invokeServerLogHandler(t *testing.T, target string, handler gin.HandlerFunc) map[string]interface{} {
	t.Helper()
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, target, nil)
	handler(ctx)

	body := map[string]interface{}{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not json: %s", recorder.Body.String())
	}
	return body
}

func TestServerLogHandlersRejectMissingParams(t *testing.T) {
	stubVlogServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("the im server must not be called for an incomplete request")
	})

	cases := []struct {
		name    string
		target  string
		handler gin.HandlerFunc
	}{
		{"userconnect without user id", "/apps/serverlogs/userconnect?app_key=ak", QryUserConnectLogs},
		{"connect without session", "/apps/serverlogs/connect?app_key=ak&user_id=u1", QryConnectLogs},
		// the node routes the query by user id even though it filters by session
		{"connect without user id", "/apps/serverlogs/connect?app_key=ak&session=s1", QryConnectLogs},
		{"business without user id", "/apps/serverlogs/business?app_key=ak&session=s1&index=2", QryBusinessLogs},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			body := invokeServerLogHandler(t, c.target, c.handler)
			if int(body["code"].(float64)) != int(errs.AdminErrorCode_ParamError) {
				t.Fatalf("code = %v, want param error", body["code"])
			}
		})
	}
}

func TestQryBusinessLogsForwardsSignalIdentity(t *testing.T) {
	var got url.Values
	stubVlogServer(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.Query()
		w.Write([]byte(`{"code":0,"msg":"success","data":{"logs":["{\"service_name\":\"msgdispatcher\",\"method\":\"p_msg\",\"expend\":\"12\",\"parms\":{\"target_id\":\"u2\"}}"]}}`))
	})

	body := invokeServerLogHandler(t, "/apps/serverlogs/business?app_key=ak&session=s1&user_id=u1&index=7&start=1700000000000&count=20", QryBusinessLogs)
	if int(body["code"].(float64)) != int(errs.AdminErrorCode_Success) {
		t.Fatalf("code = %v, want success", body["code"])
	}
	logs := body["data"].(map[string]interface{})["logs"].([]interface{})
	if len(logs) != 1 {
		t.Fatalf("logs len = %d, want 1", len(logs))
	}
	if logs[0].(map[string]interface{})["method"] != "p_msg" {
		t.Errorf("decoded log = %v, want the business log fields", logs[0])
	}

	if got.Get("log_type") != "business" || got.Get("index") != "7" || got.Get("session") != "s1" {
		t.Errorf("forwarded query = %v", got)
	}
}
