package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/juggleim/imserver-console/commons/ctxs"
	"github.com/juggleim/imserver-console/commons/errs"
	"github.com/juggleim/imserver-console/commons/logs"
	"github.com/juggleim/imserver-console/commons/tools"
	"github.com/juggleim/imserver-console/dbs"
	"github.com/juggleim/imserver-console/services/models"
)

var supportedPushChannels = []models.PushChannel{
	models.PushChannel_Huawei,
	models.PushChannel_Xiaomi,
	models.PushChannel_OPPO,
	models.PushChannel_VIVO,
	models.PushChannel_Jpush,
	models.PushChannel_HONOR,
	models.PushChannel_Getui,
	models.PushChannel_FCM,
}

var secretPushFields = map[string]struct{}{
	"app_secret":    {},
	"master_secret": {},
	"cert_pwd":      {},
	"voip_cert_pwd": {},
}

func canonicalPushChannel(value string) (string, bool) {
	value = strings.TrimSpace(value)
	for _, channel := range supportedPushChannels {
		if strings.EqualFold(value, string(channel)) {
			return string(channel), true
		}
	}
	return "", false
}

func decodePushExtra(raw string) map[string]any {
	extra := make(map[string]any)
	if raw != "" {
		_ = json.Unmarshal([]byte(raw), &extra)
	}
	return extra
}

func mergePushExtra(existing, incoming map[string]any) map[string]any {
	merged := make(map[string]any, len(existing)+len(incoming))
	for key, value := range existing {
		merged[key] = value
	}
	for key, value := range incoming {
		if _, secret := secretPushFields[key]; secret {
			text := strings.TrimSpace(fmt.Sprint(value))
			if text == "" || text == models.PushSecretMask {
				continue
			}
		}
		merged[key] = value
	}
	return merged
}

func preparePushExtraForEditing(extra map[string]any) map[string]any {
	prepared := make(map[string]any, len(extra))
	for key, value := range extra {
		prepared[key] = value
	}
	return prepared
}

func normalizeAndValidatePushExtra(channel string, extra map[string]any) (map[string]any, string, error) {
	for key, value := range extra {
		if _, secret := secretPushFields[key]; secret && strings.TrimSpace(fmt.Sprint(value)) == models.PushSecretMask {
			return nil, "", fmt.Errorf("%s must be provided as a real value", key)
		}
	}
	var conf any
	switch channel {
	case string(models.PushChannel_Huawei):
		value := tools.MapToStruct[models.HuaweiPushConf](extra)
		value.BadgeClass = strings.TrimSpace(value.BadgeClass)
		if !value.Valid() {
			return nil, "", errors.New("app_id and app_secret are required")
		}
		conf = value
	case string(models.PushChannel_Xiaomi):
		value := tools.MapToStruct[models.XiaomiPushConf](extra)
		if !value.Valid() {
			return nil, "", errors.New("app_secret is required")
		}
		conf = value
	case string(models.PushChannel_OPPO):
		value := tools.MapToStruct[models.OppoPushConf](extra)
		if !value.Valid() {
			return nil, "", errors.New("app_key and master_secret are required")
		}
		conf = value
	case string(models.PushChannel_VIVO):
		value := tools.MapToStruct[models.VivoPushConf](extra)
		if !value.Valid() {
			return nil, "", errors.New("app_id, app_key and app_secret are required")
		}
		conf = value
	case string(models.PushChannel_Jpush):
		value := tools.MapToStruct[models.JPushConf](extra)
		if !value.Valid() {
			return nil, "", errors.New("app_key and master_secret are required")
		}
		if value.Options != nil && value.Options.ThirdPartyChannel != nil {
			vivo := value.Options.ThirdPartyChannel.Vivo
			if vivo != nil && vivo.PushMode != nil && (*vivo.PushMode < 0 || *vivo.PushMode > 1) {
				return nil, "", errors.New("vivo push_mode must be 0 or 1")
			}
		}
		conf = value
	case string(models.PushChannel_HONOR):
		value := tools.MapToStruct[models.HonorPushConf](extra)
		value.BadgeClass = strings.TrimSpace(value.BadgeClass)
		if !value.Valid() {
			return nil, "", errors.New("app_id, app_key and app_secret are required")
		}
		conf = value
	case string(models.PushChannel_Getui):
		value := tools.MapToStruct[models.GetuiPushConf](extra)
		if !value.Valid() {
			return nil, "", errors.New("app_id, app_key and master_secret are required")
		}
		conf = value
	default:
		return nil, "", errors.New("unsupported push channel")
	}

	raw := tools.ToJson(conf)
	normalized := decodePushExtra(raw)
	return normalized, raw, nil
}

func failPushParam(ctx *gin.Context, message string) {
	ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ParamError, message)
}

func failPushStore(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, dbs.ErrPushConfConflict):
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_PushConfExisted, "package already configured")
	case errors.Is(err, dbs.ErrPushConfNotFound):
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_PushConfNotExist, "push configuration not found")
	default:
		logs.WithContext(ctx).Error(err.Error())
		ctxs.FailHttpResp(ctx, errs.AdminErrorCode_ServerErr, "push configuration operation failed")
	}
}
