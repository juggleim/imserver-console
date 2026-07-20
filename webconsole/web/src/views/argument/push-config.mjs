export const PUSH_SECRET_MASK = '********';

const text = (name, label, options = {}) => ({
  name,
  label,
  type: 'input_text',
  cardVisible: true,
  ...options,
});

const file = (name, labelKey, model, options = {}) => ({
  name,
  labelKey,
  model,
  type: 'file',
  cardVisible: true,
  ...options,
});

export const PUSH_CHANNELS = [
  {
    type: 'Huawei',
    name: 'Huawei',
    nameKey: 'appServices.push.channel.Huawei',
    kind: 'text',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      text('app_id', 'App ID', { labelKey: 'appServices.push.field.appId', required: true }),
      text('app_secret', 'App Secret', {
        labelKey: 'appServices.push.field.appSecret',
        required: true,
        secret: true,
      }),
    ],
  },
  {
    type: 'Xiaomi',
    name: 'Xiaomi',
    nameKey: 'appServices.push.channel.Xiaomi',
    kind: 'text',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      text('app_secret', 'App Secret', {
        labelKey: 'appServices.push.field.appSecret',
        required: true,
        secret: true,
      }),
      text('channel_id', 'Channel ID', { labelKey: 'appServices.push.field.channelId' }),
    ],
  },
  {
    type: 'Oppo',
    name: 'OPPO',
    nameKey: 'appServices.push.channel.Oppo',
    kind: 'text',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      text('app_key', 'App Key', { labelKey: 'appServices.push.field.appKey', required: true }),
      text('master_secret', 'Master Secret', {
        labelKey: 'appServices.push.field.masterSecret',
        required: true,
        secret: true,
      }),
      text('channel_id', 'Channel ID', { labelKey: 'appServices.push.field.channelId' }),
    ],
  },
  {
    type: 'Vivo',
    name: 'VIVO',
    nameKey: 'appServices.push.channel.Vivo',
    kind: 'text',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      text('app_id', 'App ID', { labelKey: 'appServices.push.field.appId', required: true }),
      text('app_key', 'App Key', { labelKey: 'appServices.push.field.appKey', required: true }),
      text('app_secret', 'App Secret', {
        labelKey: 'appServices.push.field.appSecret',
        required: true,
        secret: true,
      }),
    ],
  },
  {
    type: 'ios',
    name: 'iOS',
    nameKey: 'appServices.push.channel.ios',
    kind: 'ios',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      file('cert_path', 'appServices.push.field.certFile', 'file', { required: true }),
      text('cert_pwd', 'Certificate Password', {
        labelKey: 'appServices.push.field.certPassword',
        required: true,
        secret: true,
      }),
      file('voip_cert_path', 'appServices.push.field.voipCertFile', 'voipFile'),
      text('voip_cert_pwd', 'VoIP Certificate Password', {
        labelKey: 'appServices.push.field.voipCertPassword',
        secret: true,
      }),
      {
        name: 'is_product',
        label: 'Certificate Environment',
        labelKey: 'appServices.push.field.certEnv',
        type: 'radios',
        required: true,
        cardVisible: true,
        defaultValue: 0,
        radios: [
          { value: 0, label: 'Development', labelKey: 'appServices.push.option.dev' },
          { value: 1, label: 'Production', labelKey: 'appServices.push.option.prod' },
        ],
      },
    ],
  },
  {
    type: 'fcm',
    name: 'FCM',
    nameKey: 'appServices.push.channel.fcm',
    kind: 'fcm',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      file('conf_path', 'appServices.push.field.configFile', 'file', { required: true }),
    ],
  },
  {
    type: 'Jpush',
    name: 'JPush',
    nameKey: 'appServices.push.channel.Jpush',
    kind: 'text',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      text('app_key', 'App Key', { labelKey: 'appServices.push.field.appKey', required: true }),
      text('master_secret', 'Master Secret', {
        labelKey: 'appServices.push.field.masterSecret',
        required: true,
        secret: true,
      }),
    ],
  },
  {
    type: 'Honor',
    name: 'Honor',
    nameKey: 'appServices.push.channel.Honor',
    kind: 'text',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      text('app_id', 'App ID', { labelKey: 'appServices.push.field.appId', required: true }),
      text('app_key', 'App Key', { labelKey: 'appServices.push.field.appKey', required: true }),
      text('app_secret', 'App Secret', {
        labelKey: 'appServices.push.field.appSecret',
        required: true,
        secret: true,
      }),
    ],
  },
  {
    type: 'Getui',
    name: 'Getui',
    nameKey: 'appServices.push.channel.Getui',
    kind: 'text',
    fields: [
      text('package', 'Package Name', {
        labelKey: 'appServices.push.field.package',
        required: true,
        cardVisible: false,
      }),
      text('app_id', 'App ID', { labelKey: 'appServices.push.field.appId', required: true }),
      text('app_key', 'App Key', { labelKey: 'appServices.push.field.appKey', required: true }),
      text('master_secret', 'Master Secret', {
        labelKey: 'appServices.push.field.masterSecret',
        required: true,
        secret: true,
      }),
    ],
  },
];

export function createPushDraft(setting, item = null) {
  const source = item ? { ...item, ...(item.extra || {}) } : {};
  const draft = {
    original_package: item?.package || '',
    file: null,
    voipFile: null,
    _secretPresent: {},
  };
  setting.fields.forEach((field) => {
    if (field.secret) {
      draft._secretPresent[field.name] = Boolean(source[field.name]);
      draft[field.name] = '';
      return;
    }
    if (Object.prototype.hasOwnProperty.call(source, field.name)) {
      draft[field.name] = source[field.name];
      return;
    }
    draft[field.name] = field.defaultValue ?? '';
  });
  return draft;
}

export function validatePushDraft(setting, draft, items = []) {
  const errors = {};
  const editing = Boolean(draft.original_package);
  setting.fields.forEach((field) => {
    if (!field.required) {
      return;
    }
    if (field.secret && editing && draft._secretPresent[field.name]) {
      return;
    }
    if (field.type === 'file') {
      if (!draft[field.model]?.name && !draft[field.name]) {
        errors[field.name] = 'required';
      }
      return;
    }
    const value = draft[field.name];
    if (value === '' || value === null || value === undefined) {
      errors[field.name] = 'required';
    }
  });

  if (draft.voipFile?.name && !draft.voip_cert_pwd && !draft._secretPresent.voip_cert_pwd) {
    errors.voip_cert_pwd = 'required';
  }

  const packageName = String(draft.package || '').trim();
  const duplicate = items.some(
    (item) => item.package === packageName && item.package !== draft.original_package
  );
  if (packageName && duplicate) {
    errors.package = 'duplicate';
  }
  return errors;
}

export function hasPushErrors(errors) {
  return Object.keys(errors).length > 0;
}

export function getPushCardValue(item, field) {
  const source = { ...item, ...(item.extra || {}) };
  const value = source[field.name];
  if (field.secret && value) {
    return PUSH_SECRET_MASK;
  }
  return value;
}
