package models

import (
	"encoding/json"
	"testing"
)

func TestSmsEngineConfRoundTripKeepsAllProviders(t *testing.T) {
	original := &SmsEngineConf{
		Channel: SmsChannelAliyun,
		BdSmsEngine: &BdSmsEngine{
			ApiKey: "baidu-key", SecretKey: "baidu-secret", Endpoint: "sms.bj.baidubce.com",
			Template: "smsTpl:test", SignatureId: "sms-sign",
		},
		SmsBaoSmsEngine: &SmsBaoSmsEngine{
			Username: "smsbao-user", Password: "smsbao-secret", Template: "您的验证码是{code}",
		},
		AliyunSmsEngine: &AliyunSmsEngine{
			AccessKeyId: "aliyun-key", AccessKeySecret: "aliyun-secret", RegionId: "cn-hangzhou",
			Endpoint: "dysmsapi.aliyuncs.com", SignName: "JuggleIM", TemplateCode: "SMS_123456",
			TemplateCodeKey: "code",
		},
	}

	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}

	var decoded SmsEngineConf
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Channel != SmsChannelAliyun {
		t.Fatalf("channel = %q, want %q", decoded.Channel, SmsChannelAliyun)
	}
	if decoded.BdSmsEngine == nil || decoded.BdSmsEngine.SecretKey != "baidu-secret" {
		t.Fatalf("Baidu config was not preserved: %#v", decoded.BdSmsEngine)
	}
	if decoded.SmsBaoSmsEngine == nil || decoded.SmsBaoSmsEngine.Password != "smsbao-secret" {
		t.Fatalf("SMSBao config was not preserved: %#v", decoded.SmsBaoSmsEngine)
	}
	if decoded.AliyunSmsEngine == nil || decoded.AliyunSmsEngine.AccessKeySecret != "aliyun-secret" || decoded.AliyunSmsEngine.TemplateCodeKey != "code" {
		t.Fatalf("Aliyun config was not preserved: %#v", decoded.AliyunSmsEngine)
	}
}

func TestSmsEngineConfReadsLegacyAliyunFieldNames(t *testing.T) {
	tests := []struct {
		name    string
		payload string
	}{
		{
			name: "current persisted fields",
			payload: `{
				"channel":"aliyun",
				"aliyun":{
					"access_key":"legacy-key",
					"access_secret":"legacy-secret",
					"region_id":"cn-hangzhou",
					"endpoint":"dysmsapi.aliyuncs.com",
					"sign_name":"JuggleIM",
					"template_code":"SMS_300",
					"template_code_key":"code"
				}
			}`,
		},
		{
			name: "camel case fields",
			payload: `{
				"channel":"aliyun",
				"aliyun":{
					"accessKeyId":"legacy-key",
					"accessKeySecret":"legacy-secret",
					"regionId":"cn-shanghai",
					"endpoint":"dysmsapi.aliyuncs.com",
					"signName":"JuggleIM",
					"templateCode":"SMS_100",
					"templateCodeKey":"code"
				}
			}`,
		},
		{
			name: "ali provider and id credential aliases",
			payload: `{
				"channel":"ali",
				"ali":{
					"access_key_id":"legacy-key",
					"access_key_secret":"legacy-secret",
					"region_id":"cn-hangzhou",
					"endpoint":"dysmsapi.aliyuncs.com",
					"sign_name":"JuggleIM",
					"template_code":"SMS_200",
					"template_code_key":"code"
				}
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var conf SmsEngineConf
			if err := json.Unmarshal([]byte(test.payload), &conf); err != nil {
				t.Fatal(err)
			}
			if conf.Channel != SmsChannelAliyun {
				t.Fatalf("channel = %q, want %q", conf.Channel, SmsChannelAliyun)
			}
			if conf.AliyunSmsEngine == nil {
				t.Fatal("Aliyun config was not loaded")
			}
			if conf.AliyunSmsEngine.AccessKeyId != "legacy-key" {
				t.Fatalf("access key id was not loaded: %#v", conf.AliyunSmsEngine)
			}
			if conf.AliyunSmsEngine.AccessKeySecret != "legacy-secret" {
				t.Fatalf("access key secret was not loaded: %#v", conf.AliyunSmsEngine)
			}
		})
	}
}

func TestSmsEngineConfWritesAliyunFieldsExpectedBySmsEngine(t *testing.T) {
	conf := SmsEngineConf{
		Channel: SmsChannelAliyun,
		AliyunSmsEngine: &AliyunSmsEngine{
			AccessKeyId: "aliyun-key", AccessKeySecret: "aliyun-secret",
		},
	}
	encoded, err := json.Marshal(conf)
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(encoded, &raw); err != nil {
		t.Fatal(err)
	}
	aliyun, ok := raw[SmsChannelAliyun].(map[string]any)
	if !ok {
		t.Fatalf("Aliyun config was not written as an object: %s", encoded)
	}
	if aliyun["access_key"] != "aliyun-key" || aliyun["access_secret"] != "aliyun-secret" {
		t.Fatalf("unexpected Aliyun credential fields: %s", encoded)
	}
	if _, exists := aliyun["access_key_id"]; exists {
		t.Fatalf("legacy alias should not be written: %s", encoded)
	}
}
