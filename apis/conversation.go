package apis

import (
	"github.com/gin-gonic/gin"
	juggleimsdk "github.com/juggleim/imserver-sdk-go"
	"github.com/juggleim/imserver-console/apis/models"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/imsdk"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/services"
)

func QryConversations(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	var start int64
	if startStr := ctx.Query("start"); startStr != "" {
		if val, err := tools.String2Int64(startStr); err == nil && val > 0 {
			start = val
		}
	}
	var count int64 = 20
	if countStr := ctx.Query("count"); countStr != "" {
		if val, err := tools.String2Int64(countStr); err == nil && val > 0 {
			count = val
		}
	}
	var targetId *string
	if targetIdStr := ctx.Query("target_id"); targetIdStr != "" {
		targetId = &targetIdStr
	}
	var channelType *int32
	if channelTypeStr := ctx.Query("channel_type"); channelTypeStr != "" {
		if val, err := tools.String2Int64(channelTypeStr); err == nil && val > 0 {
			c := int32(val)
			channelType = &c
		}
	}
	ret := &models.GlobalConversations{
		Items: []*models.GlobalConversation{},
	}
	if sdk := imsdk.GetImSdk(appkey); sdk != nil {
		resp, code, _, err := sdk.QryGlobalConvers(start, int(count), targetId, channelType)
		if err == nil && code == juggleimsdk.ApiCode_Success && resp != nil {
			for _, item := range resp.Items {
				conver := &models.GlobalConversation{
					ChannelType: item.ChannelType,
					Sender:      services.QryUserInfo(appkey, item.UserId),
					Time:        item.Time,
				}
				if item.ChannelType == 1 {
					conver.Receiver = services.QryUserInfo(appkey, item.TargetId)
				} else if item.ChannelType == 2 {
					conver.Group = services.QryGroupInfo(appkey, item.TargetId)
				}
				ret.Items = append(ret.Items, conver)
			}
		}
	}
	ctxs.SuccessHttpResp(ctx, ret)
}
