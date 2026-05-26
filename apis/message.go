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

func QryHistoryMsgs(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	channelType := ctx.Query("channel_type")
	fromId := ctx.Query("from_id")
	targetId := ctx.Query("target_id")

	var start int64
	if startTimeStr := ctx.Query("start"); startTimeStr != "" {
		if val, err := tools.String2Int64(startTimeStr); err == nil && val > 0 {
			start = val
		}
	}
	var count int64
	if countStr := ctx.Query("count"); countStr != "" {
		if val, err := tools.String2Int64(countStr); err == nil && val > 0 {
			count = val
		}
	}
	isPositive := ctx.Query("order") == "1"

	ret := &models.HisMsgs{
		Msgs: []*models.HisMsg{},
	}
	if sdk := imsdk.GetImSdk(appkey); sdk != nil {
		cType := juggleimsdk.ChannelType_Private
		if channelType == "2" {
			cType = juggleimsdk.ChannelType_Group
		}
		resp, code, _, err := sdk.QryHisMsgs(fromId, targetId, cType, start, int(count), isPositive)
		if err == nil && code == juggleimsdk.ApiCode_Success && resp != nil {
			for _, msg := range resp.Msgs {
				ret.Msgs = append(ret.Msgs, &models.HisMsg{
					Sender:     services.QryUserInfo(appkey, msg.SenderId),
					MsgId:      msg.MsgId,
					MsgTime:    msg.MsgTime,
					MsgType:    msg.MsgType,
					MsgContent: msg.MsgContent,
				})
			}
		}
	}
	ctxs.SuccessHttpResp(ctx, ret)
}

func RecallHistoryMsg(ctx *gin.Context) {
	var req models.RecallHisMsgReq
	if err := ctx.BindJSON(&req); err != nil || req.AppKey == "" || req.FromId == "" || req.TargetId == "" || req.MsgId == "" || req.ChannelType == 0 {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	if sdk := imsdk.GetImSdk(req.AppKey); sdk != nil {
		cType := juggleimsdk.ChannelType_Private
		if req.ChannelType == 2 {
			cType = juggleimsdk.ChannelType_Group
		}
		code, _, err := sdk.RecallMsg(&juggleimsdk.RecallMsgReq{
			FromId:      req.FromId,
			TargetId:    req.TargetId,
			ChannelType: int32(cType),
			MsgId:       req.MsgId,
			MsgTime:     req.MsgTime,
			Exts:        req.Exts,
		})
		if err != nil {
			ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ServerErr)
			return
		}
		if code != juggleimsdk.ApiCode_Success {
			ctxs.FailHttpResp(ctx, errs.AdminErrorCode(code))
			return
		}
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func DelHistoryMsg(ctx *gin.Context) {
	var req models.DelHisMsgsReq
	if err := ctx.BindJSON(&req); err != nil || req.AppKey == "" || req.FromId == "" || req.TargetId == "" || req.ChannelType == 0 || len(req.Msgs) <= 0 {
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError)
		return
	}
	if sdk := imsdk.GetImSdk(req.AppKey); sdk != nil {
		cType := juggleimsdk.ChannelType_Private
		if req.ChannelType == 2 {
			cType = juggleimsdk.ChannelType_Group
		}
		code, _, err := sdk.DelMsgs(&juggleimsdk.DelMsgsReq{
			FromId:      req.FromId,
			TargetId:    req.TargetId,
			ChannelType: int32(cType),
			DelScope:    1,
			Msgs:        req.Msgs,
		})
		if err != nil {
			ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ServerErr)
			return
		}
		if code != juggleimsdk.ApiCode_Success {
			ctxs.FailHttpResp(ctx, errs.AdminErrorCode(code))
			return
		}
	}
	ctxs.SuccessHttpResp(ctx, nil)
}
