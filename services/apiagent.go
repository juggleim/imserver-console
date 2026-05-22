package services

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/juggleim/imserver-console/commons/configures"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func ApiAgent(method, path, body, appkey, secret string) (int, string) {
	ImApiUrl := fmt.Sprintf("%s/apigateway", configures.Config.ImApiDomain)
	respBs, code, err := tools.HttpDoBytes(method, fmt.Sprintf("%s%s", ImApiUrl, path), getSignatureHeaders(appkey, secret), body)
	if err == nil {
		return code, string(respBs)
	} else {
		logs.NewLogEntity().Error(err.Error())
		return http.StatusBadRequest, ""
	}
}

func getSignatureHeaders(appkey, secret string) map[string]string {
	nonce := fmt.Sprintf("%d", rand.Int31n(10000))
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signature := tools.SHA1(fmt.Sprintf("%s%s%s", secret, nonce, timestamp))
	headers := map[string]string{
		"Content-Type": "application/json",
		"appkey":       appkey,
		"nonce":        nonce,
		"timestamp":    timestamp,
		"signature":    signature,
	}
	return headers
}
