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

const integer = (name, label, options = {}) => ({
  name,
  label,
  type: 'input_number',
  integer: true,
  cardVisible: true,
  ...options,
});

const jpushField = (vendor, name, label, options = {}) => ({
  name: `jpush_${vendor}_${name}`,
  payloadName: name,
  label,
  type: 'input_text',
  cardVisible: false,
  path: ['options', 'third_party_channel', vendor, name],
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
      text('badge_class', 'Badge Class', {
        labelKey: 'appServices.push.field.badgeClass',
        omitEmpty: true,
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
      {
        name: 'auth_type',
        labelKey: 'appServices.push.field.authType',
        type: 'radios',
        required: true,
        cardVisible: true,
        defaultValue: 'p12',
        radios: [
          { value: 'p12', label: 'P12' },
          { value: 'p8', label: 'P8' },
        ],
      },
      text('p8_key_id', 'Key ID', {
        labelKey: 'appServices.push.field.p8KeyId',
        required: true,
        authType: 'p8',
      }),
      text('p8_team_id', 'Team ID', {
        labelKey: 'appServices.push.field.p8TeamId',
        required: true,
        authType: 'p8',
      }),
      file('p8_key_name', 'appServices.push.field.p8File', 'p8File', {
        required: true,
        authType: 'p8',
        accept: '.p8',
        maxSize: 16 * 1024,
      }),
      file('cert_path', 'appServices.push.field.certFile', 'file', {
        required: true,
        authType: 'p12',
      }),
      text('cert_pwd', 'Certificate Password', {
        labelKey: 'appServices.push.field.certPassword',
        required: true,
        secret: true,
        authType: 'p12',
      }),
      file('voip_cert_path', 'appServices.push.field.voipCertFile', 'voipFile', {
        authType: 'p12',
      }),
      text('voip_cert_pwd', 'VoIP Certificate Password', {
        labelKey: 'appServices.push.field.voipCertPassword',
        secret: true,
        authType: 'p12',
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
      text('badge_class', 'Badge Class', {
        labelKey: 'appServices.push.field.badgeClass',
        omitEmpty: true,
      }),
      integer('classification', 'Classification', {
        labelKey: 'appServices.push.field.classification',
        path: ['options', 'classification'],
        omitEmpty: true,
        optionsField: true,
      }),
    ],
    jpushTabs: [
      {
        type: 'huawei',
        label: 'Huawei',
        labelKey: 'appServices.push.channel.Huawei',
        fields: [
          jpushField('huawei', 'importance', 'Importance', {
            labelKey: 'appServices.push.field.importance',
          }),
          jpushField('huawei', 'category', 'Category', {
            labelKey: 'appServices.push.field.category',
          }),
        ],
      },
      {
        type: 'xiaomi',
        label: 'Xiaomi',
        labelKey: 'appServices.push.channel.Xiaomi',
        fields: [
          jpushField('xiaomi', 'channel_id', 'Channel ID', {
            labelKey: 'appServices.push.field.channelId',
          }),
          jpushField('xiaomi', 'mi_template_id', 'Mi Template ID', {
            labelKey: 'appServices.push.field.miTemplateId',
          }),
          jpushField('xiaomi', 'mi_template_param', 'Mi Template Param', {
            labelKey: 'appServices.push.field.miTemplateParam',
          }),
        ],
      },
      {
        type: 'honor',
        label: 'Honor',
        labelKey: 'appServices.push.channel.Honor',
        fields: [
          jpushField('honor', 'importance', 'Importance', {
            labelKey: 'appServices.push.field.importance',
          }),
        ],
      },
      {
        type: 'oppo',
        label: 'OPPO',
        labelKey: 'appServices.push.channel.Oppo',
        fields: [
          jpushField('oppo', 'distribution', 'Distribution', {
            labelKey: 'appServices.push.field.distribution',
          }),
          jpushField('oppo', 'channel_id', 'Channel ID', {
            labelKey: 'appServices.push.field.channelId',
          }),
          jpushField('oppo', 'category', 'Category', {
            labelKey: 'appServices.push.field.category',
          }),
          jpushField('oppo', 'notify_level', 'Notify Level', {
            labelKey: 'appServices.push.field.notifyLevel',
            type: 'input_number',
            integer: true,
          }),
          jpushField('oppo', 'badge_operation_type', 'Badge Operation Type', {
            labelKey: 'appServices.push.field.badgeOperationType',
            type: 'select',
            integer: true,
            allowZero: true,
            options: [
              {
                value: 0,
                label: 'Overwrite',
                labelKey: 'appServices.push.option.badgeOverwrite',
              },
              {
                value: 1,
                label: 'Increase',
                labelKey: 'appServices.push.option.badgeIncrease',
              },
            ],
          }),
          jpushField('oppo', 'private_msg_template_id', 'Private Message Template ID', {
            labelKey: 'appServices.push.field.privateMsgTemplateId',
          }),
          jpushField('oppo', 'private_content_parameters', 'Private Content Parameters', {
            labelKey: 'appServices.push.field.privateContentParameters',
            jsonObject: true,
          }),
          jpushField('oppo', 'private_title_parameters', 'Private Title Parameters', {
            labelKey: 'appServices.push.field.privateTitleParameters',
            jsonObject: true,
          }),
        ],
      },
      {
        type: 'vivo',
        label: 'VIVO',
        labelKey: 'appServices.push.channel.Vivo',
        fields: [
          jpushField('vivo', 'distribution', 'Distribution', {
            labelKey: 'appServices.push.field.distribution',
          }),
          jpushField('vivo', 'category', 'Category', {
            labelKey: 'appServices.push.field.category',
          }),
          jpushField('vivo', 'add_badge', 'Add Badge', {
            labelKey: 'appServices.push.field.addBadge',
            type: 'checkbox',
            defaultValue: false,
          }),
          jpushField('vivo', 'push_mode', 'Push Mode', {
            labelKey: 'appServices.push.field.pushMode',
            type: 'select',
            integer: true,
            allowZero: true,
            allowedValues: [0, 1],
            options: [
              {
                value: 0,
                label: 'Production Push',
                labelKey: 'appServices.push.option.pushModeProduction',
              },
              {
                value: 1,
                label: 'Test Push',
                labelKey: 'appServices.push.option.pushModeTest',
              },
            ],
          }),
        ],
      },
      {
        type: 'meizu',
        label: 'Meizu',
        labelKey: 'appServices.push.channel.Meizu',
        fields: [
          jpushField('meizu', 'distribution', 'Distribution', {
            labelKey: 'appServices.push.field.distribution',
          }),
        ],
      },
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
      text('badge_class', 'Badge Class', {
        labelKey: 'appServices.push.field.badgeClass',
        omitEmpty: true,
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
    config_version: item ? item.config_version ?? 1 : undefined,
    _originalAuthType: item ? item.auth_type || 'p12' : '',
    file: null,
    voipFile: null,
    p8File: null,
    has_p8_key: Boolean(item?.has_p8_key),
    _secretPresent: {},
  };
  allDraftFields(setting).forEach((field) => {
    if (field.secret) {
      if (setting.kind === 'ios') {
        draft._secretPresent[field.name] = Boolean(
          source[field.name === 'cert_pwd' ? 'cert_path' : 'voip_cert_path']
        );
        draft[field.name] = '';
        return;
      }
      const storedValue = source[field.name];
      draft._secretPresent[field.name] = Boolean(storedValue);
      draft[field.name] = storedValue !== PUSH_SECRET_MASK ? storedValue || '' : '';
      return;
    }
    const sourceValue = readFieldValue(source, field);
    if (sourceValue !== undefined) {
      draft[field.name] =
        field.jsonObject && typeof sourceValue === 'object'
          ? JSON.stringify(sourceValue)
          : sourceValue;
      return;
    }
    draft[field.name] = field.defaultValue ?? '';
  });
  return draft;
}

export function validatePushDraft(setting, draft, items = []) {
  const errors = {};
  const editing = Boolean(draft.original_package);
  allDraftFields(setting).forEach((field) => {
    if (!isPushFieldActive(field, draft)) return;
    const value = draft[field.name];
    if (field.maxSize && draft[field.model]?.size > field.maxSize) {
      errors[field.name] = 'p8Size';
      return;
    }
    if (field.integer && value !== '' && value !== null && value !== undefined) {
      const numberValue = Number(value);
      if (!Number.isInteger(numberValue)) {
        errors[field.name] = 'integer';
        return;
      }
      if (field.allowedValues && !field.allowedValues.includes(numberValue)) {
        errors[field.name] = 'range';
        return;
      }
    }
    if (field.jsonObject && value !== '' && value !== null && value !== undefined) {
      try {
        const parsed = typeof value === 'string' ? JSON.parse(value) : value;
        if (
          !parsed ||
          Array.isArray(parsed) ||
          typeof parsed !== 'object' ||
          Object.values(parsed).some((item) => typeof item !== 'string')
        ) {
          throw new TypeError('expected an object with string values');
        }
      } catch {
        errors[field.name] = 'stringMap';
        return;
      }
    }
    if (!field.required) {
      return;
    }
    if (setting.kind === 'ios' && editing && field.authType === 'p12') {
      if (field.secret || draft._originalAuthType === 'p12') return;
      if (field.name === 'cert_path' && (draft.voip_cert_path || draft.voipFile?.name)) return;
    }
    if (field.secret && editing && draft._secretPresent[field.name]) {
      return;
    }
    if (field.type === 'file') {
      const saved = field.model === 'p8File' ? draft.has_p8_key : draft[field.name];
      if (!draft[field.model]?.name && !saved) {
        errors[field.name] = 'required';
      }
      return;
    }
    if (value === '' || value === null || value === undefined) {
      errors[field.name] = 'required';
    }
  });

  if (
    draft.auth_type !== 'p8' &&
    !editing &&
    draft.voipFile?.name &&
    !draft.voip_cert_pwd &&
    !draft._secretPresent.voip_cert_pwd
  ) {
    errors.voip_cert_pwd = 'required';
  }

  if (setting.kind === 'ios') {
    if (!['p12', 'p8'].includes(draft.auth_type)) errors.auth_type = 'range';
    if (![0, 1].includes(draft.is_product)) errors.is_product = 'range';
    if (draft.auth_type === 'p8') {
      for (const field of ['p8_key_id', 'p8_team_id']) {
        if (draft[field] && !/^[A-Z0-9]{10}$/.test(draft[field])) errors[field] = 'appleId';
      }
    }
  }

  const packageName = String(draft.package || '').trim();
  if (setting.kind === 'ios' && !/^[A-Za-z0-9.-]{1,100}$/.test(packageName))
    errors.package = 'topic';
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

export function isPushFieldActive(field, source) {
  return !field.authType || field.authType === (source.auth_type || 'p12');
}

export function p8KeyStatus(source) {
  if (source.has_p8_key) return 'appServices.push.status.keyStored';
  return 'appServices.push.status.unset';
}

export function buildIosPushParams(appKey, draft) {
  const params = {
    app_key: appKey,
    package: draft.package,
    original_package: draft.original_package,
    auth_type: draft.auth_type,
    is_product: draft.is_product,
    ...(draft.original_package ? { config_version: draft.config_version } : {}),
  };
  if (draft.auth_type === 'p8') {
    Object.assign(params, {
      p8_key_id: draft.p8_key_id,
      p8_team_id: draft.p8_team_id,
      p8File: draft.p8File,
    });
  } else {
    Object.assign(params, {
      cert_pwd: draft.cert_pwd,
      voip_cert_pwd: draft.voip_cert_pwd,
      file: draft.file,
      voipFile: draft.voipFile,
    });
  }
  return params;
}

function allDraftFields(setting) {
  return [...setting.fields, ...(setting.jpushTabs || []).flatMap((tab) => tab.fields)];
}

function readPath(source, path) {
  return path.reduce((value, key) => value?.[key], source);
}

function readFieldValue(source, field) {
  return field.path ? readPath(source, field.path) : source[field.name];
}

function normalizeDraftValue(field, draft) {
  const value = draft[field.name];
  if (field.type === 'checkbox') {
    return Boolean(value);
  }
  if (field.jsonObject && typeof value === 'string') {
    const trimmed = value.trim();
    return trimmed ? JSON.parse(trimmed) : trimmed;
  }
  if (field.integer && value !== '' && value !== null && value !== undefined) {
    return Number(value);
  }
  return typeof value === 'string' ? value.trim() : value;
}

function isEmptyOptionalValue(field, value) {
  return (
    value === '' ||
    value === null ||
    value === undefined ||
    (field.type === 'checkbox' && !value) ||
    (field.integer && !field.allowZero && value === 0) ||
    (field.jsonObject && typeof value === 'object' && Object.keys(value).length === 0)
  );
}

function buildJpushExtra(setting, draft) {
  const extra = {};
  const options = {};

  setting.fields.forEach((field) => {
    if (field.name === 'package') {
      return;
    }
    const value = normalizeDraftValue(field, draft);
    if (field.path?.[0] === 'options') {
      if (!isEmptyOptionalValue(field, value)) {
        options[field.path[field.path.length - 1]] = value;
      }
      return;
    }
    if (field.omitEmpty && isEmptyOptionalValue(field, value)) {
      return;
    }
    extra[field.name] = value;
  });

  const thirdPartyChannel = {};
  setting.jpushTabs.forEach((tab) => {
    const channel = {};
    tab.fields.forEach((field) => {
      const value = normalizeDraftValue(field, draft);
      if (!isEmptyOptionalValue(field, value)) {
        channel[field.payloadName] = value;
      }
    });
    if (Object.keys(channel).length) {
      thirdPartyChannel[tab.type] = channel;
    }
  });
  if (Object.keys(thirdPartyChannel).length) {
    options.third_party_channel = thirdPartyChannel;
  }
  if (Object.keys(options).length) {
    extra.options = options;
  } else if (draft.original_package) {
    // An explicit null clears previously saved optional settings during editing.
    extra.options = null;
  }
  return extra;
}

export function buildPushTextExtra(setting, draft) {
  if (setting.type === 'Jpush') {
    return buildJpushExtra(setting, draft);
  }
  return setting.fields.reduce((extra, field) => {
    if (field.name === 'package' || field.type !== 'input_text') {
      return extra;
    }
    const value =
      typeof draft[field.name] === 'string' ? draft[field.name].trim() : draft[field.name];
    if (field.omitEmpty && (value === '' || value === null || value === undefined)) {
      return extra;
    }
    extra[field.name] = value;
    return extra;
  }, {});
}

export function getPushCardValue(item, field) {
  const source = { ...item, ...(item.extra || {}) };
  const value = readFieldValue(source, field);
  if (
    (field.name === 'cert_pwd' && source.cert_path) ||
    (field.name === 'voip_cert_pwd' && source.voip_cert_path)
  )
    return PUSH_SECRET_MASK;
  if (field.name === 'auth_type') return value || 'p12';
  if (field.secret && value) {
    return PUSH_SECRET_MASK;
  }
  return value;
}
