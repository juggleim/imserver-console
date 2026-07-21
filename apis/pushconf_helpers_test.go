package apis

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/juggleim/imserver-console/services/models"
)

func TestCanonicalPushChannelSupportsConfiguredVendors(t *testing.T) {
	for _, channel := range supportedPushChannels {
		got, ok := canonicalPushChannel(strings.ToLower(string(channel)))
		if !ok || got != string(channel) {
			t.Fatalf("failed to canonicalize %s: %q %v", channel, got, ok)
		}
	}
	if _, ok := canonicalPushChannel("unknown"); ok {
		t.Fatal("unknown channel must be rejected")
	}
}

func TestPreparePushExtraForEditingReturnsAllSecrets(t *testing.T) {
	prepared := preparePushExtraForEditing(map[string]any{
		"app_id":        "123",
		"app_secret":    "plain-app-secret",
		"master_secret": "plain-master-secret",
	})
	if prepared["app_id"] != "123" {
		t.Fatalf("non-secret value changed: %#v", prepared)
	}
	if prepared["app_secret"] != "plain-app-secret" {
		t.Fatalf("app_secret must be returned for editing: %#v", prepared)
	}
	if prepared["master_secret"] != "plain-master-secret" {
		t.Fatalf("master_secret must be returned for editing: %#v", prepared)
	}
}

func TestMergePushExtraPreservesBlankOrMaskedSecret(t *testing.T) {
	existing := map[string]any{"app_id": "old", "app_secret": "real-secret"}
	for _, placeholder := range []string{"", models.PushSecretMask} {
		merged := mergePushExtra(existing, map[string]any{"app_id": "new", "app_secret": placeholder})
		if merged["app_secret"] != "real-secret" || merged["app_id"] != "new" {
			t.Fatalf("unexpected merge for %q: %#v", placeholder, merged)
		}
	}
}

func TestMaskedSecretCannotBeSavedAsNewCredential(t *testing.T) {
	_, _, err := normalizeAndValidatePushExtra(string(models.PushChannel_Huawei), map[string]any{
		"app_id": "123", "app_secret": models.PushSecretMask,
	})
	if err == nil {
		t.Fatal("mask placeholder must not pass credential validation")
	}
}

func TestValidateEveryTextPushChannel(t *testing.T) {
	tests := map[string]map[string]any{
		"Huawei": {"app_id": "id", "app_secret": "secret"},
		"Xiaomi": {"app_secret": "secret"},
		"Oppo":   {"app_key": "key", "master_secret": "secret"},
		"Vivo":   {"app_id": "id", "app_key": "key", "app_secret": "secret"},
		"Jpush":  {"app_key": "key", "master_secret": "secret"},
		"Honor":  {"app_id": "id", "app_key": "key", "app_secret": "secret"},
		"Getui":  {"app_id": "id", "app_key": "key", "master_secret": "secret"},
	}
	for channel, extra := range tests {
		if _, _, err := normalizeAndValidatePushExtra(channel, extra); err != nil {
			t.Errorf("%s should be valid: %v", channel, err)
		}
	}
}

func TestHuaweiAndHonorBadgeClassIsOptional(t *testing.T) {
	tests := []struct {
		channel string
		extra   map[string]any
	}{
		{channel: "Huawei", extra: map[string]any{"app_id": "id", "app_secret": "secret", "badge_class": "   "}},
		{channel: "Honor", extra: map[string]any{"app_id": "id", "app_key": "key", "app_secret": "secret", "badge_class": "   "}},
	}
	for _, test := range tests {
		t.Run(test.channel, func(t *testing.T) {
			withoutBadge, raw, err := normalizeAndValidatePushExtra(test.channel, test.extra)
			if err != nil {
				t.Fatal(err)
			}
			if _, ok := withoutBadge["badge_class"]; ok || strings.Contains(raw, "badge_class") {
				t.Fatalf("empty badge_class must be omitted: %s", raw)
			}

			withBadgeInput := make(map[string]any, len(test.extra)+1)
			for key, value := range test.extra {
				withBadgeInput[key] = value
			}
			withBadgeInput["badge_class"] = "com.example.MainActivity"
			_, raw, err = normalizeAndValidatePushExtra(test.channel, withBadgeInput)
			if err != nil {
				t.Fatal(err)
			}
			var saved map[string]any
			if err := json.Unmarshal([]byte(raw), &saved); err != nil {
				t.Fatal(err)
			}
			if saved["badge_class"] != "com.example.MainActivity" {
				t.Fatalf("badge_class was not preserved: %s", raw)
			}
		})
	}
}

func TestJPushProviderOptionsArePreserved(t *testing.T) {
	extra := map[string]any{
		"app_key":       "key",
		"master_secret": "secret",
		"badge_class":   "com.example.Badge",
		"options": map[string]any{
			"classification": 2,
			"third_party_channel": map[string]any{
				"huawei": map[string]any{"importance": "HIGH", "category": "IM"},
				"xiaomi": map[string]any{
					"channel_id": "xiaomi-channel", "mi_template_id": "template-id", "mi_template_param": "{\"key\":\"value\"}",
				},
				"honor": map[string]any{"importance": "NORMAL"},
				"oppo":  map[string]any{"channel_id": "oppo-channel", "category": "IM", "notify_level": 2},
				"vivo":  map[string]any{"distribution": "push", "category": "IM", "add_badge": true},
				"meizu": map[string]any{"distribution": "push"},
			},
		},
	}

	_, raw, err := normalizeAndValidatePushExtra(string(models.PushChannel_Jpush), extra)
	if err != nil {
		t.Fatal(err)
	}
	var saved models.JPushConf
	if err := json.Unmarshal([]byte(raw), &saved); err != nil {
		t.Fatal(err)
	}
	if saved.BadgeClass != "com.example.Badge" || saved.Options == nil || saved.Options.Classification != 2 {
		t.Fatalf("top-level JPush options were not preserved: %s", raw)
	}
	channels := saved.Options.ThirdPartyChannel
	if channels == nil || channels.Huawei == nil || channels.Huawei.Importance != "HIGH" || channels.Huawei.Category != "IM" {
		t.Fatalf("Huawei channel options were not preserved: %s", raw)
	}
	if channels.Xiaomi == nil || channels.Xiaomi.ChannelId != "xiaomi-channel" || channels.Xiaomi.MiTemplateId != "template-id" || channels.Xiaomi.MiTemplateParam != "{\"key\":\"value\"}" {
		t.Fatalf("Xiaomi channel options were not preserved: %s", raw)
	}
	if channels.Honor == nil || channels.Honor.Importance != "NORMAL" {
		t.Fatalf("Honor channel options were not preserved: %s", raw)
	}
	if channels.Oppo == nil || channels.Oppo.NotifyLevel != 2 {
		t.Fatalf("OPPO channel options were not preserved: %s", raw)
	}
	if channels.Vivo == nil || !channels.Vivo.AddBadge {
		t.Fatalf("VIVO channel options were not preserved: %s", raw)
	}
	if channels.Meizu == nil || channels.Meizu.Distribution != "push" {
		t.Fatalf("Meizu channel options were not preserved: %s", raw)
	}
}

func TestJPushEmptyOptionsAreOmitted(t *testing.T) {
	merged := mergePushExtra(map[string]any{
		"app_key": "key", "master_secret": "secret", "options": map[string]any{"classification": 2},
	}, map[string]any{
		"app_key": "key", "master_secret": "secret", "options": nil,
	})
	_, raw, err := normalizeAndValidatePushExtra(string(models.PushChannel_Jpush), merged)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(raw, "options") {
		t.Fatalf("empty JPush options must be omitted: %s", raw)
	}
}
