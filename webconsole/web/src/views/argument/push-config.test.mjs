import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import {
  PUSH_CHANNELS,
  PUSH_SECRET_MASK,
  buildPushTextExtra,
  buildIosPushParams,
  createPushDraft,
  getPushCardValue,
  hasPushErrors,
  isPushFieldActive,
  p8KeyStatus,
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

test('P8 presence depends on raw status rather than a saved filename', () => {
  const ios = PUSH_CHANNELS.find((item) => item.kind === 'ios');
  for (const raw of [false, true]) {
    const draft = createPushDraft(ios, {
      package: 'com.example',
      auth_type: 'p8',
      p8_key_id: 'ABC1234567',
      p8_team_id: 'XYZ1234567',
      p8_key_name: 'saved.p8',
      voip_cert_path: 'retained-voip.p12',
      has_p8_key: raw,
      config_version: 4,
    });
    assert.equal(p8KeyStatus(draft), `appServices.push.status.${raw ? 'keyStored' : 'unset'}`);
    assert.equal(validatePushDraft(ios, draft).p8_key_name, raw ? undefined : 'required');
    draft.p8File = { name: 'new.p8', size: 400 };
    assert.deepEqual(validatePushDraft(ios, draft), {});
    const params = buildIosPushParams('app-1', draft);
    assert.equal(params.config_version, 4);
    assert.equal(params.has_p8_key, undefined);
    draft.p8File = null;
    assert.equal(validatePushDraft(ios, draft).p8_key_name, raw ? undefined : 'required');
    draft.auth_type = 'p12';
    assert.equal(draft._originalAuthType, 'p8');
    assert.deepEqual(validatePushDraft(ios, draft), {});
  }
});

test('P8 card and editor share bilingual saved and unset status', () => {
  for (const locale of ['en-US', 'zh-CN']) {
    const text = readFileSync(
      new URL(`../../locales/${locale}/appServices.json`, import.meta.url),
      'utf8'
    );
    const { push } = JSON.parse(text);
    assert.ok(push.hint.p8Setup);
    assert.ok(push.status.keyStored);
    assert.ok(push.status.unset);
  }
  const view = readFileSync(new URL('./push.vue', import.meta.url), 'utf8');
  const dialog = readFileSync(
    new URL('../../components/push-config-dialog.vue', import.meta.url),
    'utf8'
  );
  assert.match(view, /t\(p8KeyStatus\(item\)\)/);
  assert.match(dialog, /t\(p8KeyStatus\(props.draft\)\)/);
  assert.match(dialog, /!props.draft.has_p8_key/);
});

test('two administrators keep distinct iOS draft versions and reloading clears credential inputs', () => {
  const ios = PUSH_CHANNELS.find((item) => item.kind === 'ios');
  const old = {
    package: 'com.example',
    auth_type: 'p8',
    config_version: 7,
    has_p8_key: true,
    p8_key_id: 'OLD1234567',
    p8_team_id: 'XYZ1234567',
  };
  const a = createPushDraft(ios, old);
  const b = createPushDraft(ios, old);
  b.p8_key_id = 'NEW1234567';
  b.p8File = { name: 'new.p8' };
  assert.equal(buildIosPushParams('app-1', b).config_version, 7);
  a.is_product = 1;
  assert.equal(buildIosPushParams('app-1', a).config_version, 7);
  const reloaded = createPushDraft(ios, { ...old, config_version: 8, p8_key_id: b.p8_key_id });
  assert.equal(buildIosPushParams('app-1', reloaded).config_version, 8);
  assert.equal(reloaded.p8_key_id, 'NEW1234567');
  assert.equal(reloaded.p8File, null);
  assert.equal(reloaded.cert_pwd, '');
  assert.equal(buildIosPushParams('app-1', createPushDraft(ios)).config_version, undefined);
});

test('legacy iOS VoIP-only and unnamed passwordless credentials permit metadata edits', () => {
  const ios = PUSH_CHANNELS.find((item) => item.kind === 'ios');
  for (const metadata of [{ voip_cert_path: 'voip.p12' }, { cert_path: 'app.p12' }, {}]) {
    const draft = createPushDraft(ios, { package: 'com.example', ...metadata, config_version: 2 });
    draft.is_product = 1;
    assert.deepEqual(validatePushDraft(ios, draft), {});
    draft.voipFile = { name: 'new-passwordless.p12' };
    assert.deepEqual(validatePushDraft(ios, draft), {});
    assert.equal(buildIosPushParams('app-1', draft).cert_pwd, '');
    assert.equal(buildIosPushParams('app-1', draft).config_version, 2);
  }
  const newDraft = createPushDraft(ios);
  newDraft.package = 'com.example';
  assert.equal(validatePushDraft(ios, newDraft).cert_path, 'required');
  assert.equal(validatePushDraft(ios, newDraft).cert_pwd, 'required');
});

test('P8 to P12 rollback accepts retained or uploaded VoIP without an ordinary certificate', () => {
  const ios = PUSH_CHANNELS.find((item) => item.kind === 'ios');
  for (const retained of [false, true]) {
    const draft = createPushDraft(ios, {
      package: 'com.example',
      auth_type: 'p8',
      config_version: 7,
      voip_cert_path: retained ? 'retained-voip.p12' : '',
    });
    draft.auth_type = 'p12';
    assert.equal(draft._originalAuthType, 'p8');
    assert.equal(draft.cert_path, '');
    assert.equal(draft.file, null);
    if (!retained) {
      assert.deepEqual(validatePushDraft(ios, draft), { cert_path: 'required' });
      draft.voipFile = { name: 'new-voip.p12' };
    }
    assert.deepEqual(validatePushDraft(ios, draft), {});
    const params = buildIosPushParams('app-1', draft);
    assert.equal(params.auth_type, 'p12');
    assert.equal(params.config_version, 7);
    assert.equal(params.file, null);
    assert.equal(params.cert_pwd, '');
    assert.equal(params.voip_cert_pwd, '');
    assert.equal(params.voipFile, draft.voipFile);
    draft.voipFile = null;
    assert.deepEqual(validatePushDraft(ios, draft), retained ? {} : { cert_path: 'required' });
  }
  const added = createPushDraft(ios);
  added.package = 'com.example';
  added.voipFile = { name: 'new-voip.p12' };
  assert.deepEqual(validatePushDraft(ios, added), {
    cert_path: 'required',
    cert_pwd: 'required',
    voip_cert_pwd: 'required',
  });
});

test('iOS bundles match the shared ASCII 100-character rule', () => {
  const ios = PUSH_CHANNELS.find((item) => item.kind === 'ios');
  const draft = createPushDraft(ios, { package: 'com.original' });
  for (const value of ['com.example app', 'com/example', 'com_foo', 'com.例', 'a'.repeat(101)]) {
    draft.package = value;
    assert.equal(validatePushDraft(ios, draft).package, 'topic');
    assert.equal(draft.original_package, 'com.original');
  }
  draft.package = 'a'.repeat(100);
  assert.deepEqual(validatePushDraft(ios, draft), {});
});

test('multipart carries version and conflict UI offers explicit reload without automatic retry', () => {
  const service = readFileSync(new URL('../../services/application.js', import.meta.url), 'utf8');
  assert.match(service, /form\.append\('config_version', params\.config_version\)/);
  const view = readFileSync(new URL('./push.vue', import.meta.url), 'utf8');
  assert.match(view, /error\.code === 409/);
  assert.match(view, /dialog\.conflict = true/);
  assert.match(view, /@reload="reloadConflictedDraft"/);
  const dialog = readFileSync(
    new URL('../../components/push-config-dialog.vue', import.meta.url),
    'utf8'
  );
  assert.match(dialog, /!props\.saving && !props\.conflict/);
  assert.match(dialog, /emit\('reload'\)/);
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
      oppo: [
        'distribution',
        'channel_id',
        'category',
        'notify_level',
        'badge_operation_type',
        'private_msg_template_id',
        'private_content_parameters',
        'private_title_parameters',
      ],
      vivo: ['distribution', 'category', 'add_badge', 'push_mode'],
      meizu: ['distribution'],
    }
  );
  const badgeOperationType = jpush.jpushTabs
    .find((tab) => tab.type === 'oppo')
    .fields.find((field) => field.payloadName === 'badge_operation_type');
  assert.equal(badgeOperationType.type, 'select');
  assert.equal(badgeOperationType.required, undefined);
  assert.deepEqual(
    badgeOperationType.options.map((option) => option.value),
    [0, 1]
  );
  const oppoFields = jpush.jpushTabs.find((tab) => tab.type === 'oppo').fields;
  assert.equal(
    oppoFields.find((field) => field.payloadName === 'private_content_parameters').jsonObject,
    true
  );
  assert.equal(
    oppoFields.find((field) => field.payloadName === 'private_title_parameters').jsonObject,
    true
  );
  const pushMode = jpush.jpushTabs
    .find((tab) => tab.type === 'vivo')
    .fields.find((field) => field.payloadName === 'push_mode');
  assert.equal(pushMode.labelKey, 'appServices.push.field.pushMode');
  assert.equal(pushMode.type, 'select');
  assert.equal(pushMode.required, undefined);
  assert.deepEqual(
    pushMode.options.map((option) => option.value),
    [0, 1]
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
          oppo: {
            distribution: 'push',
            channel_id: 'oppo-channel',
            category: 'IM',
            notify_level: 2,
            badge_operation_type: 0,
            private_msg_template_id: 'template-id',
            private_content_parameters: { name: 'value' },
            private_title_parameters: { title: 'value' },
          },
          vivo: { distribution: 'push', category: 'IM', push_mode: 0, add_badge: true },
          meizu: { distribution: 'push' },
        },
      },
    },
  };
  const draft = createPushDraft(jpush, item);

  assert.equal(draft.classification, 2);
  assert.equal(draft.jpush_xiaomi_mi_template_id, 'template-id');
  assert.equal(draft.jpush_oppo_distribution, 'push');
  assert.equal(draft.jpush_oppo_notify_level, 2);
  assert.equal(draft.jpush_oppo_badge_operation_type, 0);
  assert.equal(draft.jpush_oppo_private_msg_template_id, 'template-id');
  assert.equal(draft.jpush_oppo_private_content_parameters, '{"name":"value"}');
  assert.equal(draft.jpush_oppo_private_title_parameters, '{"title":"value"}');
  assert.equal(draft.jpush_vivo_push_mode, 0);
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
  draft.jpush_oppo_badge_operation_type = 0;
  assert.equal(
    buildPushTextExtra(jpush, draft).options.third_party_channel.oppo.badge_operation_type,
    0
  );
  draft.jpush_oppo_badge_operation_type = '';
  draft.jpush_vivo_push_mode = 2;
  draft.jpush_oppo_private_content_parameters = '{"name":1}';
  draft.jpush_oppo_private_title_parameters = 'invalid';
  draft.classification = '1.5';
  draft.jpush_oppo_notify_level = 'high';
  assert.deepEqual(validatePushDraft(jpush, draft, []), {
    classification: 'integer',
    jpush_oppo_notify_level: 'integer',
    jpush_vivo_push_mode: 'range',
    jpush_oppo_private_content_parameters: 'stringMap',
    jpush_oppo_private_title_parameters: 'stringMap',
  });

  draft.classification = '';
  draft.jpush_oppo_notify_level = '';
  draft.jpush_vivo_push_mode = '';
  draft.jpush_oppo_private_content_parameters = '{}';
  draft.jpush_oppo_private_title_parameters = '';
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

test('Android secrets refill for editing while iOS passwords never refill', () => {
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
      cert_path: 'app.p12',
      voip_cert_path: 'voip.p12',
      ...(setting.kind === 'ios' ? secrets : { extra: secrets }),
    });
    secretFields.forEach((field) => {
      assert.equal(
        draft[field.name],
        setting.kind === 'ios' ? '' : `plain-${field.name}`,
        `${setting.type}.${field.name}`
      );
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
    voip_cert_path: 'voip.p12',
    is_product: 1,
  });
  assert.deepEqual(validatePushDraft(ios, iosDraft, []), {});
});

test('P8 create and edit enforce key presence, size and Apple IDs without P12 fields', () => {
  const ios = PUSH_CHANNELS.find((item) => item.kind === 'ios');
  const draft = createPushDraft(ios);
  assert.equal(draft.auth_type, 'p12');
  Object.assign(draft, {
    package: 'com.example',
    auth_type: 'p8',
    p8_key_id: 'ABC1234567',
    p8_team_id: 'XYZ1234567',
  });
  assert.deepEqual(validatePushDraft(ios, draft), { p8_key_name: 'required' });
  draft.p8File = { name: 'AuthKey.p8', size: 16384 };
  assert.deepEqual(validatePushDraft(ios, draft), {});
  draft.p8File.size++;
  assert.equal(validatePushDraft(ios, draft).p8_key_name, 'p8Size');
  draft.p8_key_id = 'lowercase0';
  assert.equal(validatePushDraft(ios, draft).p8_key_id, 'appleId');
  const edit = createPushDraft(ios, {
    package: 'com.example',
    auth_type: 'p8',
    has_p8_key: true,
    p8_key_name: 'saved.p8',
    p8_key_id: 'ABC1234567',
    p8_team_id: 'XYZ1234567',
    is_product: 1,
  });
  assert.equal(edit.p8File, null);
  assert.deepEqual(validatePushDraft(ios, edit), {});
  edit.has_p8_key = false;
  assert.equal(validatePushDraft(ios, edit).p8_key_name, 'required');
});

test('iOS authentication switches retain draft files but do not submit inactive credentials', () => {
  const ios = PUSH_CHANNELS.find((item) => item.kind === 'ios');
  const draft = createPushDraft(ios, {
    package: 'com.example',
    cert_path: 'app.p12',
    auth_type: 'p8',
    has_p8_key: true,
    p8_key_id: 'ABC1234567',
    p8_team_id: 'XYZ1234567',
  });
  draft.file = { name: 'replacement.p12' };
  draft.p8File = { name: 'replacement.p8', size: 300 };
  const p8 = buildIosPushParams('app-1', draft);
  assert.equal(p8.file, undefined);
  assert.equal(p8.cert_pwd, undefined);
  assert.equal(p8.p8File, draft.p8File);
  assert.equal(
    ios.fields
      .filter((field) => isPushFieldActive(field, draft))
      .some((field) => field.name === 'cert_pwd'),
    false
  );
  draft.auth_type = 'p12';
  assert.deepEqual(validatePushDraft(ios, draft), {});
  const p12 = buildIosPushParams('app-1', draft);
  assert.equal(p12.p8File, undefined);
  assert.equal(p12.p8_key_id, undefined);
  assert.equal(p12.file, draft.file);
  assert.equal(p12.cert_pwd, '');
  assert.equal(draft.p8File.name, 'replacement.p8');
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
