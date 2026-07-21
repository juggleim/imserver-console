import test from 'node:test';
import assert from 'node:assert/strict';
import {
  PUSH_CHANNELS,
  PUSH_SECRET_MASK,
  buildPushTextExtra,
  createPushDraft,
  getPushCardValue,
  hasPushErrors,
  validatePushDraft,
} from './push-config.mjs';

test('defines all nine push channel tabs with package fields', () => {
  assert.equal(PUSH_CHANNELS.length, 9);
  assert.deepEqual(
    PUSH_CHANNELS.map((setting) => setting.type),
    ['Huawei', 'Xiaomi', 'Oppo', 'Vivo', 'ios', 'fcm', 'Jpush', 'Honor', 'Getui']
  );
  PUSH_CHANNELS.forEach((setting) => {
    assert.equal(setting.fields[0].name, 'package');
    assert.equal(setting.fields[0].required, true);
  });
});

test('every channel accepts its declared required fields and files', () => {
  PUSH_CHANNELS.forEach((setting) => {
    const draft = createPushDraft(setting);
    setting.fields.forEach((field) => {
      if (field.name === 'package') {
        draft.package = `com.example.${setting.type.toLowerCase()}`;
      } else if (field.type === 'input_text' && field.required) {
        draft[field.name] = 'credential';
      } else if (field.type === 'file' && field.required) {
        draft[field.model] = { name: `${setting.type}.config` };
      }
    });
    const errors = validatePushDraft(setting, draft, []);
    assert.equal(hasPushErrors(errors), false, `${setting.type}: ${JSON.stringify(errors)}`);
  });
});

test('Huawei add validates package, App ID and App Secret', () => {
  const setting = PUSH_CHANNELS.find((item) => item.type === 'Huawei');
  const draft = createPushDraft(setting);
  assert.deepEqual(validatePushDraft(setting, draft, []), {
    package: 'required',
    app_id: 'required',
    app_secret: 'required',
  });

  Object.assign(draft, { package: 'com.example', app_id: '123', app_secret: 'secret' });
  assert.deepEqual(validatePushDraft(setting, draft, []), {});
});

test('Huawei and Honor omit blank badge_class and include a trimmed value', () => {
  for (const channel of ['Huawei', 'Honor']) {
    const setting = PUSH_CHANNELS.find((item) => item.type === channel);
    const draft = createPushDraft(setting);
    setting.fields.forEach((field) => {
      if (field.type === 'input_text' && field.required) {
        draft[field.name] = 'credential';
      }
    });

    assert.equal(Object.hasOwn(buildPushTextExtra(setting, draft), 'badge_class'), false);
    draft.badge_class = '  com.example.MainActivity  ';
    assert.equal(buildPushTextExtra(setting, draft).badge_class, 'com.example.MainActivity');
  }
});

test('JPush exposes only its own provider option fields in six channel tabs', () => {
  const jpush = PUSH_CHANNELS.find((item) => item.type === 'Jpush');
  assert.deepEqual(
    jpush.fields.slice(-2).map((field) => field.name),
    ['badge_class', 'classification']
  );
  assert.equal(jpush.fields.find((field) => field.name === 'classification').integer, true);
  assert.equal(jpush.fields.find((field) => field.name === 'classification').optionsField, true);
  assert.equal(jpush.fields.find((field) => field.name === 'badge_class').required, undefined);
  assert.deepEqual(
    jpush.jpushTabs.map((tab) => tab.type),
    ['huawei', 'xiaomi', 'honor', 'oppo', 'vivo', 'meizu']
  );
  assert.deepEqual(
    Object.fromEntries(
      jpush.jpushTabs.map((tab) => [tab.type, tab.fields.map((field) => field.payloadName)])
    ),
    {
      huawei: ['importance', 'category'],
      xiaomi: ['channel_id', 'mi_template_id', 'mi_template_param'],
      honor: ['importance'],
      oppo: ['channel_id', 'category', 'notify_level'],
      vivo: ['distribution', 'category', 'add_badge'],
      meizu: ['distribution'],
    }
  );
  PUSH_CHANNELS.filter((item) => item.type !== 'Jpush').forEach((setting) => {
    assert.equal(setting.jpushTabs, undefined, setting.type);
  });
});

test('JPush provider options round-trip with nested API field names and numeric types', () => {
  const jpush = PUSH_CHANNELS.find((item) => item.type === 'Jpush');
  const item = {
    package: 'com.example.jpush',
    extra: {
      app_key: 'key',
      master_secret: 'secret',
      badge_class: 'com.example.Badge',
      options: {
        classification: 2,
        third_party_channel: {
          huawei: { importance: 'HIGH', category: 'IM' },
          xiaomi: {
            channel_id: 'xiaomi-channel',
            mi_template_id: 'template-id',
            mi_template_param: '{"key":"value"}',
          },
          honor: { importance: 'NORMAL' },
          oppo: { channel_id: 'oppo-channel', category: 'IM', notify_level: 2 },
          vivo: { distribution: 'push', category: 'IM', add_badge: true },
          meizu: { distribution: 'push' },
        },
      },
    },
  };
  const draft = createPushDraft(jpush, item);

  assert.equal(draft.classification, 2);
  assert.equal(draft.jpush_xiaomi_mi_template_id, 'template-id');
  assert.equal(draft.jpush_oppo_notify_level, 2);
  assert.equal(draft.jpush_vivo_add_badge, true);
  assert.deepEqual(buildPushTextExtra(jpush, draft), item.extra);
});

test('JPush omits empty options and validates optional integer values when present', () => {
  const jpush = PUSH_CHANNELS.find((item) => item.type === 'Jpush');
  const draft = createPushDraft(jpush);
  Object.assign(draft, {
    package: 'com.example.jpush',
    app_key: 'key',
    master_secret: 'secret',
  });

  assert.deepEqual(buildPushTextExtra(jpush, draft), {
    app_key: 'key',
    master_secret: 'secret',
  });
  draft.classification = '1.5';
  draft.jpush_oppo_notify_level = 'high';
  assert.deepEqual(validatePushDraft(jpush, draft, []), {
    classification: 'integer',
    jpush_oppo_notify_level: 'integer',
  });

  draft.classification = '';
  draft.jpush_oppo_notify_level = '';
  draft.original_package = draft.package;
  assert.equal(buildPushTextExtra(jpush, draft).options, null);
});

test('masked secret placeholders stay out of edit drafts and preserve existing values', () => {
  const setting = PUSH_CHANNELS.find((item) => item.type === 'Huawei');
  const draft = createPushDraft(setting, {
    package: 'com.example',
    extra: { app_id: '123', app_secret: PUSH_SECRET_MASK },
  });
  assert.equal(draft.app_secret, '');
  assert.equal(draft._secretPresent.app_secret, true);
  assert.deepEqual(validatePushDraft(setting, draft, []), {});
});

test('all secret fields refill plaintext values for editing', () => {
  PUSH_CHANNELS.forEach((setting) => {
    const secretFields = setting.fields.filter((field) => field.secret);
    if (!secretFields.length) {
      return;
    }
    const secrets = Object.fromEntries(
      secretFields.map((field) => [field.name, `plain-${field.name}`])
    );
    const draft = createPushDraft(setting, {
      package: `com.example.${setting.type.toLowerCase()}`,
      ...(setting.kind === 'ios' ? secrets : { extra: secrets }),
    });
    secretFields.forEach((field) => {
      assert.equal(draft[field.name], `plain-${field.name}`, `${setting.type}.${field.name}`);
      assert.equal(draft._secretPresent[field.name], true, `${setting.type}.${field.name}`);
    });
  });
});

test('duplicate packages are scoped to the current channel list', () => {
  const setting = PUSH_CHANNELS.find((item) => item.type === 'Oppo');
  const draft = createPushDraft(setting);
  Object.assign(draft, {
    package: 'com.example',
    app_key: 'key',
    master_secret: 'secret',
  });
  assert.equal(
    validatePushDraft(setting, draft, [{ package: 'com.example' }]).package,
    'duplicate'
  );
  assert.equal(validatePushDraft(setting, draft, [{ package: 'com.other' }]).package, undefined);
});

test('FCM and iOS edit drafts keep existing file names without new File objects', () => {
  const fcm = PUSH_CHANNELS.find((item) => item.type === 'fcm');
  const fcmDraft = createPushDraft(fcm, { package: 'com.fcm', conf_path: 'firebase.json' });
  assert.deepEqual(validatePushDraft(fcm, fcmDraft, []), {});

  const ios = PUSH_CHANNELS.find((item) => item.type === 'ios');
  const iosDraft = createPushDraft(ios, {
    package: 'com.ios',
    cert_path: 'app.p12',
    cert_pwd: PUSH_SECRET_MASK,
    voip_cert_path: 'voip.p12',
    voip_cert_pwd: PUSH_SECRET_MASK,
    is_product: 1,
  });
  assert.deepEqual(validatePushDraft(ios, iosDraft, []), {});
});

test('drafts are isolated across tabs and card helpers always mask secret values', () => {
  const huawei = PUSH_CHANNELS.find((item) => item.type === 'Huawei');
  const xiaomi = PUSH_CHANNELS.find((item) => item.type === 'Xiaomi');
  const first = createPushDraft(huawei);
  const second = createPushDraft(xiaomi);
  first.package = 'com.changed';
  assert.equal(second.package, '');

  PUSH_CHANNELS.forEach((setting) => {
    setting.fields
      .filter((field) => field.secret)
      .forEach((field) => {
        const item =
          setting.kind === 'ios'
            ? { [field.name]: 'plaintext' }
            : { extra: { [field.name]: 'plaintext' } };
        assert.equal(
          getPushCardValue(item, field),
          PUSH_SECRET_MASK,
          `${setting.type}.${field.name}`
        );
      });
  });
});
