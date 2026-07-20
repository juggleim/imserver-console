import test from 'node:test';
import assert from 'node:assert/strict';
import {
  PUSH_CHANNELS,
  PUSH_SECRET_MASK,
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

test('edit drafts do not expose secrets and blank secrets preserve existing values', () => {
  const setting = PUSH_CHANNELS.find((item) => item.type === 'Huawei');
  const draft = createPushDraft(setting, {
    package: 'com.example',
    extra: { app_id: '123', app_secret: PUSH_SECRET_MASK },
  });
  assert.equal(draft.app_secret, '');
  assert.equal(draft._secretPresent.app_secret, true);
  assert.deepEqual(validatePushDraft(setting, draft, []), {});
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

  const secretField = huawei.fields.find((field) => field.name === 'app_secret');
  assert.equal(
    getPushCardValue({ extra: { app_secret: 'plaintext' } }, secretField),
    PUSH_SECRET_MASK
  );
});
