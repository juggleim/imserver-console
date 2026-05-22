package apis

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
	"gorm.io/gorm"
)

type FileConf struct {
	AppKey  string                 `json:"app_key,omitempty"`
	Channel string                 `json:"channel,omitempty"`
	Enable  int                    `json:"enable"`
	Conf    map[string]interface{} `json:"conf,omitempty"`
}

func GetFileConf(c *gin.Context) {
	appkey := c.Query("app_key")
	channel := c.Query("channel")

	dao := dbs.FileConfDao{}
	fileConf, err := dao.FindFileConf(appkey, channel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctxs.SuccessHttpResp(c, nil)
			return
		}
		logs.NewLogEntity().Error(err.Error())
		c.JSON(http.StatusNotFound, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_Default,
		})
		return
	}

	var resp FileConf
	resp.AppKey = appkey
	resp.Channel = channel

	_ = json.Unmarshal([]byte(fileConf.Conf), &resp.Conf)
	ctxs.SuccessHttpResp(c, resp)

}

func SetFileConf(ctx *gin.Context) {
	var req FileConf
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	pushConf := ""
	var conf any
	if strings.EqualFold(req.Channel, models.ChannelAws) {
		conf = tools.MapToStruct[models.S3Config](req.Conf)
	} else if strings.EqualFold(req.Channel, models.ChannelQiNiu) {
		conf = tools.MapToStruct[models.QiNiuConfig](req.Conf)
	} else if strings.EqualFold(req.Channel, models.ChannelOss) {
		conf = tools.MapToStruct[models.OssConfig](req.Conf)
	} else if strings.EqualFold(req.Channel, models.ChannelMinio) {
		conf = tools.MapToStruct[models.MinioConfig](req.Conf)
	}
	pushConf = tools.ToJson(conf)

	//save to db
	dao := dbs.FileConfDao{}
	err := dao.Upsert(dbs.FileConfDao{
		AppKey:  req.AppKey,
		Channel: req.Channel,
		Conf:    pushConf,
	})
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
	}

	ctxs.SuccessHttpResp(ctx, nil)
}

func GetFileConfCurrSwitch(c *gin.Context) {
	appkey := c.Query("app_key")
	dao := dbs.FileConfDao{}
	fConf, err := dao.FindEnableFileConf(appkey)
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		c.JSON(http.StatusOK, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_NoFileEngine,
		})
		return
	}
	ctxs.SuccessHttpResp(c, gin.H{"channel": fConf.Channel})
}

func GetFileConfs(ctx *gin.Context) {
	appkey := ctx.Query("app_key")
	if appkey == "" {
		ctx.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}
	dao := dbs.FileConfDao{}
	confs, err := dao.FindFileConfs(appkey)

	fileConfs := []*FileConf{}
	if err == nil {
		for _, conf := range confs {
			fileConfs = append(fileConfs, &FileConf{
				Channel: conf.Channel,
				Enable:  conf.Enable,
			})
		}
	} else {
		logs.NewLogEntity().Error(err.Error())
	}
	ctxs.SuccessHttpResp(ctx, map[string]interface{}{
		"file_confs": fileConfs,
	})
}

func SetFileConfSwitch(c *gin.Context) {
	var req struct {
		AppKey  string `json:"app_key"`
		Channel string `json:"channel"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, &ctxs.ApiErrorMsg{
			Code: errs.AdminErrorCode_ParamError,
			Msg:  "param illegal",
		})
		return
	}

	dao := dbs.FileConfDao{}
	err := dao.UpdateEnable(req.AppKey, req.Channel)
	if err != nil {
		logs.NewLogEntity().Error(err.Error())
		ctxs.FailHttpResp(c, errs.AdminErrorCode_Default)
		return
	}
	ctxs.SuccessHttpResp(c, nil)
}
