package models

import "encoding/json"

const (
	SmsChannelBaidu  = "baidu"
	SmsChannelSmsBao = "smsbao"
	SmsChannelAliyun = "aliyun"
)

type SmsEngineConf struct {
	Channel         string           `json:"channel,omitempty"`
	BdSmsEngine     *BdSmsEngine     `json:"baidu,omitempty"`
	SmsBaoSmsEngine *SmsBaoSmsEngine `json:"smsbao,omitempty"`
	AliyunSmsEngine *AliyunSmsEngine `json:"aliyun,omitempty"`
}

type BdSmsEngine struct {
	ApiKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	Endpoint    string `json:"endpoint"`
	Template    string `json:"template"`
	SignatureId string `json:"signature_id"`
}

type SmsBaoSmsEngine struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Template string `json:"template"`
}

type AliyunSmsEngine struct {
	AccessKeyId     string `json:"access_key"`
	AccessKeySecret string `json:"access_secret"`
	RegionId        string `json:"region_id"`
	Endpoint        string `json:"endpoint"`
	SignName        string `json:"sign_name"`
	TemplateCode    string `json:"template_code"`
	TemplateCodeKey string `json:"template_code_key"`
}

func (conf *SmsEngineConf) UnmarshalJSON(data []byte) error {
	type smsEngineConfJSON struct {
		Channel         string           `json:"channel"`
		BdSmsEngine     *BdSmsEngine     `json:"baidu"`
		SmsBaoSmsEngine *SmsBaoSmsEngine `json:"smsbao"`
		AliyunSmsEngine *AliyunSmsEngine `json:"aliyun"`
		AliSmsEngine    *AliyunSmsEngine `json:"ali"`
		AliyunSmsLegacy *AliyunSmsEngine `json:"aliyun_sms"`
		AliyunSmsFlat   *AliyunSmsEngine `json:"aliyunsms"`
	}
	var decoded smsEngineConfJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	conf.Channel = normalizeSmsChannel(decoded.Channel)
	conf.BdSmsEngine = decoded.BdSmsEngine
	conf.SmsBaoSmsEngine = decoded.SmsBaoSmsEngine
	conf.AliyunSmsEngine = firstAliyunEngine(
		decoded.AliyunSmsEngine,
		decoded.AliSmsEngine,
		decoded.AliyunSmsLegacy,
		decoded.AliyunSmsFlat,
	)
	return nil
}

func (engine *AliyunSmsEngine) UnmarshalJSON(data []byte) error {
	type aliyunSmsEngineJSON struct {
		AccessKeyId      string `json:"access_key"`
		AccessKeyIdAlias string `json:"access_key_id"`
		AccessKeyIdCamel string `json:"accessKeyId"`
		AccessKeyIDCamel string `json:"accessKeyID"`
		AccessKeyIdFlat  string `json:"accesskeyid"`

		AccessKeySecret      string `json:"access_secret"`
		AccessKeySecretAlias string `json:"access_key_secret"`
		AccessKeySecretCamel string `json:"accessKeySecret"`
		AccessKeySecretOld   string `json:"secret_key"`
		AccessKeySecretFlat  string `json:"accesskeysecret"`

		RegionId      string `json:"region_id"`
		RegionIdCamel string `json:"regionId"`
		Endpoint      string `json:"endpoint"`
		SignName      string `json:"sign_name"`
		SignNameCamel string `json:"signName"`

		TemplateCode         string `json:"template_code"`
		TemplateCodeCamel    string `json:"templateCode"`
		TemplateCodeKey      string `json:"template_code_key"`
		TemplateCodeKeyCamel string `json:"templateCodeKey"`
	}
	var decoded aliyunSmsEngineJSON
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}
	engine.AccessKeyId = firstNonEmpty(
		decoded.AccessKeyId,
		decoded.AccessKeyIdAlias,
		decoded.AccessKeyIdCamel,
		decoded.AccessKeyIDCamel,
		decoded.AccessKeyIdFlat,
	)
	engine.AccessKeySecret = firstNonEmpty(
		decoded.AccessKeySecret,
		decoded.AccessKeySecretAlias,
		decoded.AccessKeySecretCamel,
		decoded.AccessKeySecretOld,
		decoded.AccessKeySecretFlat,
	)
	engine.RegionId = firstNonEmpty(decoded.RegionId, decoded.RegionIdCamel)
	engine.Endpoint = decoded.Endpoint
	engine.SignName = firstNonEmpty(decoded.SignName, decoded.SignNameCamel)
	engine.TemplateCode = firstNonEmpty(decoded.TemplateCode, decoded.TemplateCodeCamel)
	engine.TemplateCodeKey = firstNonEmpty(decoded.TemplateCodeKey, decoded.TemplateCodeKeyCamel)
	return nil
}

func normalizeSmsChannel(channel string) string {
	switch channel {
	case "ali", "aliyun_sms", "aliyunsms":
		return SmsChannelAliyun
	default:
		return channel
	}
}

func firstAliyunEngine(engines ...*AliyunSmsEngine) *AliyunSmsEngine {
	for _, engine := range engines {
		if engine != nil {
			return engine
		}
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
