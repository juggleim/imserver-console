package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/services/models"
)

type AccessAddress struct {
	Original *OriginalAddress `json:"original"`
	Proxy    *ProxyAddress    `json:"proxy"`
}
type OriginalAddress struct {
	Nav     map[string]string `json:"nav"`
	Api     map[string]string `json:"api"`
	Connect map[string]string `json:"connect"`
}
type ProxyAddress struct {
	Nav     *models.AddressConf `json:"nav"`
	Api     *models.AddressConf `json:"api"`
	Connect *models.AddressConf `json:"connect"`
}

func GetAccessAddress(ctx *gin.Context) {
	ctxs.SuccessHttpResp(ctx, &AccessAddress{
		// Original: &OriginalAddress{
		// 	Nav:     commonservices.GetOriginalNavAddress().NodeConfs,
		// 	Api:     commonservices.GetOriginalApiAddress().NodeConfs,
		// 	Connect: commonservices.GetOriginalConnectAddress().NodeConfs,
		// },
		// Proxy: &ProxyAddress{
		// 	Nav:     commonservices.GetProxyNavAddress(),
		// 	Api:     commonservices.GetProxyApiAddress(),
		// 	Connect: commonservices.GetProxyConnectAddress(),
		// },
	})
}
