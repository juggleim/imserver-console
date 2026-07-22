package routers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/apis"
)

func Route(eng *gin.Engine, prefix string) *gin.RouterGroup {
	group := eng.Group("/" + prefix)
	group.Use(corsHandler())
	group.Use(apis.Validate)
	group.OPTIONS("", optionsHandler())
	group.OPTIONS("/*path", optionsHandler())
	group.POST("/imapiagent", apis.ApiAgent)
	group.GET("/common/address", apis.GetAccessAddress)

	group.POST("/login", apis.Login)
	group.POST("/accounts/updpass", apis.UpdPassword)
	group.POST("/accounts/add", apis.AddAccount)
	group.POST("/accounts/delete", apis.DeleteAccounts)
	group.POST("/accounts/disable", apis.DisableAccounts)
	group.POST("/accounts/bindapps", apis.BindApps)
	group.POST("/accounts/unbindapps", apis.UnBindApps)
	group.GET("/accounts/list", apis.QryAccounts)

	group.POST("/apps/active", apis.ActiveApp)
	group.POST("/apps/create", apis.CreateApp)
	group.GET("/apps/list", apis.QryApps)
	group.GET("/apps/info", apis.QryAppInfo)
	group.POST("/apps/alias/set", apis.UpdateAppAlias)
	group.POST("/apps/configs/set", apis.UpdateAppConfigs)
	group.POST("/apps/configs/get", apis.QryAppConfigs)
	group.POST("/apps/eventsubconfig/set", apis.SetEventSubConfig)
	group.GET("/apps/eventsubconfig/get", apis.GetEventSubConfig)
	//translate
	group.POST("/apps/translate/set", apis.SetTranslateConf)
	group.GET("/apps/translate/get", apis.GetTranslateConf)
	//sms
	group.POST("/apps/sms/set", apis.SetSmsConf)
	group.GET("/apps/sms/get", apis.GetSmsConf)
	//rtc
	group.POST("/apps/rtcconf/set", apis.SetRtcConf)
	group.GET("/apps/rtcconf/get", apis.GetRtcConf)
	group.POST("/apps/zegoconf/set", apis.SetZegoConf)
	group.GET("/apps/zegoconf/get", apis.GetZegoConf)
	group.POST("/apps/agoraconf/set", apis.SetAgoraConf)
	group.GET("/apps/agoraconf/get", apis.GetAgoraConf)
	group.POST("/apps/livekitconf/set", apis.SetLivekitConf)
	group.GET("/apps/livekitconf/get", apis.GetLivekitConf)

	group.POST("/apps/iospushcer/set", apis.SetIosPushConf)
	group.POST("/apps/iospushcer/upload", apis.UploadIosCer)
	group.GET("/apps/iospushcer/get", apis.GetIosCer)
	group.GET("/apps/iospushcer/list", apis.ListIosPushConfs)
	group.POST("/apps/fcmpushconf/upload", apis.UploadFcmPushConf)
	group.GET("/apps/fcmpushconf/get", apis.GetFcmPushConf)
	group.GET("/apps/fcmpushconf/list", apis.ListFcmPushConfs)
	group.POST("/apps/androidpushconf/set", apis.SetAndroidPushConf)
	group.GET("/apps/androidpushconf/get", apis.GetAndroidPushConf)
	group.GET("/apps/androidpushconf/list", apis.ListAndroidPushConfs)

	group.POST("/apps/fileconf/set", apis.SetFileConf)
	group.GET("/apps/fileconf/get", apis.GetFileConf)
	group.GET("/apps/fileconf/switch/get", apis.GetFileConfs)
	group.POST("/apps/fileconf/switch/set", apis.SetFileConfSwitch)
	//logs
	group.POST("/apps/clientlogs/notify", apis.ClientLogNtf)
	group.GET("/apps/clientlogs/list", apis.ClientLogList)
	group.GET("/apps/clientlogs/download", apis.ClientLogDownload)
	//connection inspector, proxied to the im server's vlog query
	group.GET("/apps/serverlogs/userconnect", apis.QryUserConnectLogs)
	group.GET("/apps/serverlogs/connect", apis.QryConnectLogs)
	group.GET("/apps/serverlogs/business", apis.QryBusinessLogs)

	//statistic
	group.GET("/apps/statistic/msg", apis.QryMsgStatistic)
	group.GET("/apps/statistic/msgrealtime", apis.QryMsgRealtimeStatistic)
	group.GET("/apps/statistic/useractivity", apis.QryUserActivities)
	group.GET("/apps/statistic/userreg", apis.QryUserRegiste)
	group.GET("/apps/statistic/connectcount", apis.QryConnectCount)
	group.GET("/apps/statistic/maxconnectcount", apis.QryMaxConnectCount)
	group.GET("/apps/statistic/chrmconnectcount", apis.QryChrmConnectCount)
	group.GET("/apps/statistic/maxchrmconnectcount", apis.QryMaxChrmConnectCount)
	group.GET("/apps/statistic/maxchrmconnectcount_v2", apis.QryMaxChrmConnectCountV2)

	group.GET("/apps/monitor/performance/nodes", apis.QryPerformanceNodes)
	group.GET("/apps/monitor/performance/catalog", apis.QryPerformanceCatalog)
	group.GET("/apps/monitor/performance/metrics", apis.QryPerformanceMetric)

	group.GET("/apps/sensitivewords/list", apis.SensitiveWords)
	group.POST("/apps/sensitivewords/import", apis.ImportSensitiveWords)
	group.POST("/apps/sensitivewords/add", apis.AddSensitiveWord)
	group.POST("/apps/sensitivewords/delete", apis.DeleteSensitiveWord)

	// custom interceptors
	group.POST("/apps/interceptors/add", apis.AddInterceptor)
	group.POST("/apps/interceptors/delete", apis.DeleteInterceptor)
	group.POST("/apps/interceptors/update", apis.UpdateInterceptor)
	group.GET("/apps/interceptors/list", apis.ListInterceptors)
	group.POST("/apps/interceptors/conditions/add", apis.AddInterceptorConditions)
	group.POST("/apps/interceptors/conditions/delete", apis.DeleteInterceptorConditions)
	group.POST("/apps/interceptors/conditions/update", apis.UpdateInterceptorConditions)
	group.GET("/apps/interceptors/conditions/list", apis.ListInterceptorConditions)

	group.POST("/apps/file_cred", apis.GetFileCred)
	//users
	group.GET("/apps/users/list", apis.QryUsers)
	group.POST("/apps/users/add")
	group.POST("/apps/users/update")
	group.POST("/apps/users/ban", apis.BanUsers)
	group.POST("/apps/users/unban", apis.UnBanUsers)

	//bots
	group.GET("/apps/bots/list", apis.QryBots)
	group.POST("/apps/bots/add", apis.AddBot)
	group.POST("/apps/bots/update", apis.UpdateBot)

	//groups
	group.GET("/apps/groups/list", apis.QryGroups)
	group.POST("/apps/groups/dissolve", apis.DissolveGroup)

	//convers
	group.GET("/apps/convers/list", apis.QryConversations)
	//history msgs
	group.GET("/apps/historymsgs/list", apis.QryHistoryMsgs)
	group.POST("/apps/historymsgs/recall", apis.RecallHistoryMsg)
	group.POST("/apps/historymsgs/del", apis.DelHistoryMsg)

	//applications
	group.POST("/apps/applications/add", apis.AddApplication)
	group.POST("/apps/applications/update", apis.UpdApplication)
	group.POST("/apps/applications/delete", apis.DelApplications)
	group.GET("/apps/applications/list", apis.QryApplications)

	//email setting
	group.POST("/apps/email/set", apis.SetEmailConf)
	group.GET("/apps/email/get", apis.GetEmailConf)

	return group
}

func corsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		origin := strings.TrimSpace(context.GetHeader("Origin"))
		if origin != "" {
			context.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			context.Writer.Header().Set("Vary", "Origin")
		}
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		context.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders(context.GetHeader("Access-Control-Request-Headers")))
		context.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, X-Token, X-Appid")
		context.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if context.Request.Method == http.MethodOptions {
			context.AbortWithStatus(http.StatusNoContent)
			return
		}
		context.Next()
	}
}

func optionsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Status(http.StatusNoContent)
	}
}

func allowHeaders(requestHeaders string) string {
	requestHeaders = strings.TrimSpace(requestHeaders)
	if requestHeaders != "" {
		return requestHeaders
	}

	return "Content-Type, Content-Length, Authorization, X-Requested-With, X-CSRF-Token, X-Token, X-Appid"
}
