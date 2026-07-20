package apis

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

func GetAndroidPushConf(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	pushChannel, ok := canonicalPushChannel(ctx.Query("push_channel"))
	if appkey == "" || !ok {
		failPushParam(ctx, "app_key and supported push_channel are required")
		return
	}

	ret := &models.AndroidPushConf{AppKey: appkey, PushChannel: pushChannel, Extra: map[string]any{}}
	item, err := (dbs.AndroidPushConfDao{}).Find(appkey, pushChannel)
	if err == nil && item != nil {
		ret.Package = item.Package
		ret.Extra = decodePushExtra(item.PushConf)
	} else if err != nil {
		failPushStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, ret)
}

func ListAndroidPushConfs(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	pushChannel, ok := canonicalPushChannel(ctx.Query("push_channel"))
	if appkey == "" || !ok {
		failPushParam(ctx, "app_key and supported push_channel are required")
		return
	}

	rows, err := (dbs.AndroidPushConfDao{}).List(appkey, pushChannel)
	if err != nil {
		failPushStore(ctx, err)
		return
	}
	items := make([]models.PushConfListItem, 0, len(rows))
	for _, row := range rows {
		item := models.PushConfListItem{
			AppKey:      row.AppKey,
			PushChannel: row.PushChannel,
			Package:     row.Package,
		}
		if pushChannel == string(models.PushChannel_FCM) {
			item.ConfPath = row.PushConf
		} else {
			item.Extra = maskPushExtra(decodePushExtra(row.PushConf))
		}
		items = append(items, item)
	}
	ctxs.SuccessHttpResp(ctx, items)
}

func SetAndroidPushConf(ctx *gin.Context) {
	var req models.AndroidPushConf
	if err := ctx.ShouldBindJSON(&req); err != nil {
		failPushParam(ctx, "param illegal")
		return
	}
	req.AppKey = strings.TrimSpace(req.AppKey)
	req.Package = strings.TrimSpace(req.Package)
	pushChannel, ok := canonicalPushChannel(req.PushChannel)
	if req.AppKey == "" || req.Package == "" || !ok || pushChannel == string(models.PushChannel_FCM) {
		failPushParam(ctx, "app_key, package and supported text push_channel are required")
		return
	}

	extra := req.Extra
	if extra == nil {
		extra = make(map[string]any)
	}
	if strings.TrimSpace(req.OriginalPackage) != "" {
		existing, err := (dbs.AndroidPushConfDao{}).FindByIdentity(req.AppKey, pushChannel, strings.TrimSpace(req.OriginalPackage))
		if err != nil {
			failPushStore(ctx, err)
			return
		}
		extra = mergePushExtra(decodePushExtra(existing.PushConf), extra)
	}
	_, raw, err := normalizeAndValidatePushExtra(pushChannel, extra)
	if err != nil {
		failPushParam(ctx, err.Error())
		return
	}

	err = (dbs.AndroidPushConfDao{}).Save(dbs.AndroidPushConfDao{
		AppKey:      req.AppKey,
		PushChannel: pushChannel,
		Package:     req.Package,
		PushConf:    raw,
	}, req.OriginalPackage)
	if err != nil {
		failPushStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}

func GetFcmPushConf(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.Query("app_key"))
	if appkey == "" {
		failPushParam(ctx, "app_key is required")
		return
	}
	item, err := (dbs.AndroidPushConfDao{}).Find(appkey, string(models.PushChannel_FCM))
	if err != nil {
		failPushStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, item)
}

func ListFcmPushConfs(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	query.Set("push_channel", string(models.PushChannel_FCM))
	ctx.Request.URL.RawQuery = query.Encode()
	ListAndroidPushConfs(ctx)
}

func UploadFcmPushConf(ctx *gin.Context) {
	appkey := strings.TrimSpace(ctx.PostForm("app_key"))
	packageName := strings.TrimSpace(ctx.PostForm("package"))
	originalPackage := strings.TrimSpace(ctx.PostForm("original_package"))
	if appkey == "" || packageName == "" {
		failPushParam(ctx, "app_key and package are required")
		return
	}

	item := dbs.AndroidPushConfDao{
		AppKey:      appkey,
		PushChannel: string(models.PushChannel_FCM),
		Package:     packageName,
	}
	if originalPackage != "" {
		existing, err := (dbs.AndroidPushConfDao{}).FindByIdentity(appkey, string(models.PushChannel_FCM), originalPackage)
		if err != nil {
			failPushStore(ctx, err)
			return
		}
		item.PushConf = existing.PushConf
		item.PushExt = existing.PushExt
	}

	fileHeader, fileErr := ctx.FormFile("fcm_conf")
	if fileErr == nil {
		file, err := fileHeader.Open()
		if err != nil {
			failPushParam(ctx, "unable to open fcm configuration file")
			return
		}
		defer file.Close()
		item.PushExt, err = io.ReadAll(file)
		if err != nil {
			failPushParam(ctx, "unable to read fcm configuration file")
			return
		}
		item.PushConf = fileHeader.Filename
	} else if !errors.Is(fileErr, http.ErrMissingFile) {
		failPushParam(ctx, "invalid fcm configuration file")
		return
	}
	if len(item.PushExt) == 0 || item.PushConf == "" {
		failPushParam(ctx, "fcm configuration file is required")
		return
	}

	if err := (dbs.AndroidPushConfDao{}).Save(item, originalPackage); err != nil {
		failPushStore(ctx, err)
		return
	}
	ctxs.SuccessHttpResp(ctx, nil)
}
