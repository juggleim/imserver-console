package apis

import (
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

func TestMaskPushExtraNeverReturnsPlaintextSecrets(t *testing.T) {
	masked := maskPushExtra(map[string]any{
		"app_id":        "123",
		"app_secret":    "plain-app-secret",
		"master_secret": "plain-master-secret",
	})
	if masked["app_id"] != "123" {
		t.Fatalf("non-secret value changed: %#v", masked)
	}
	if masked["app_secret"] != models.PushSecretMask || masked["master_secret"] != models.PushSecretMask {
		t.Fatalf("secret was not masked: %#v", masked)
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
