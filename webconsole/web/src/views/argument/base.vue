<script setup>
import { getCurrentInstance, reactive } from 'vue';
import utils from '../../common/utils';
import { Application } from "../../services";
import { useRouter } from "vue-router";
import { t } from '@/i18n';
import { ErrorType } from '../../common/enum';

const context = getCurrentInstance();
let router = useRouter();
let { currentRoute: { _rawValue: { params: { app_key } } } } = router;

const MAX_URL_LENGTH = 255;

const URL_FIELD_DEFS = [
  {
    key: 'ws',
    field: 'ws_url',
    label: 'appBase.field.wsUrl',
    save: ({ app_key, url }) => Application.updateWsUrl({ app_key, url }),
  },
  {
    key: 'api',
    field: 'api_url',
    label: 'appBase.field.apiUrl',
    save: ({ app_key, url }) => Application.updateApiUrl({ app_key, url }),
  },
  {
    key: 'app',
    field: 'app_url',
    label: 'appBase.field.appUrl',
    save: ({ app_key, url }) => Application.updateAppUrl({ app_key, url }),
  },
];

let state = reactive({
  appInfo: {
    restricted_fields: {}
  },
  isShowSecret: false,
  isSavingAlias: false,
  editingAlias: '0',
  aliasErrorMsg: '',
  urlFields: URL_FIELD_DEFS.map((def) => utils.extend({}, [def, {
    editing: '',
    saving: false,
    errorMsg: '',
  }])),
});

function fetchApp() {
  Application.getOne({ app_key }).then(({ data }) => {
    let { cur_user_count, max_user_count } = data;
    utils.formatProps(data, { count: 'num', time: 'date' })
    data.use_percent = Math.floor(cur_user_count / max_user_count) * 100;
    data.raw_expired_time = data.expired_time;
    data.expired_time = utils.isPermanentExpireTime(data.expired_time) ? '' : utils.formatTime(data.expired_time);
    data.n_app_secret = '********************';
    utils.extend(state.appInfo, data);
    state.editingAlias = data.alias || '0';
    utils.forEach(state.urlFields, (field) => {
      field.editing = data[field.field] || '';
    });
  });
}
fetchApp();

function onShowSecret() {
  let isShowSecret = !state.isShowSecret;
  utils.extend(state, { isShowSecret })
}

function getExpireTimeLabel() {
  return utils.isPermanentExpireTime(state.appInfo.raw_expired_time)
    ? t('common.label.permanentValidity')
    : state.appInfo.expired_time;
}

function onAliasInput() {
  state.aliasErrorMsg = '';
}

function onSaveAlias() {
  if (state.isSavingAlias) {
    return;
  }
  const alias = String(state.editingAlias || '').trim();
  if (!alias) {
    state.aliasErrorMsg = t('appBase.validation.aliasRequired');
    return;
  }
  if ([...alias].length > 50) {
    state.aliasErrorMsg = t('appBase.validation.aliasTooLong');
    return;
  }

  state.isSavingAlias = true;
  Application.updateAlias({ app_key, alias }).then(({ code, msg }) => {
    if (utils.isEqual(code, ErrorType.SUCCESS_0.code)) {
      state.appInfo.alias = alias;
      state.editingAlias = alias;
      context.proxy.$toast({ icon: 'success', text: t('appBase.feedback.aliasSaved') });
      return;
    }
    context.proxy.$toast({
      icon: 'error',
      text: t('appBase.feedback.aliasSaveFailed', { code, msg }),
    });
  }).catch(() => {
    context.proxy.$toast({ icon: 'error', text: t('appBase.feedback.aliasRequestFailed') });
  }).finally(() => {
    state.isSavingAlias = false;
  });
}

function onUrlInput(field) {
  field.errorMsg = '';
}

function onSaveUrl(field) {
  if (field.saving) {
    return;
  }
  const url = String(field.editing || '').trim();
  const label = t(field.label);
  if (!url) {
    field.errorMsg = t('appBase.validation.urlRequired', { field: label });
    return;
  }
  if ([...url].length > MAX_URL_LENGTH) {
    field.errorMsg = t('appBase.validation.urlTooLong', { field: label });
    return;
  }

  field.saving = true;
  field.save({ app_key, url }).then(({ code, msg }) => {
    if (utils.isEqual(code, ErrorType.SUCCESS_0.code)) {
      state.appInfo[field.field] = url;
      field.editing = url;
      context.proxy.$toast({ icon: 'success', text: t('appBase.feedback.urlSaved', { field: label }) });
      return;
    }
    context.proxy.$toast({
      icon: 'error',
      text: t('appBase.feedback.urlSaveFailed', { field: label, code, msg }),
    });
  }).catch(() => {
    context.proxy.$toast({ icon: 'error', text: t('appBase.feedback.urlRequestFailed', { field: label }) });
  }).finally(() => {
    field.saving = false;
  });
}

</script>
<template>
  <div class="mb-4 app-base cim-app-base-page">
    <div class="cim-app-base-card">
      <div class="cim-app-base-head">
        <h2 class="cim-app-base-title">{{ t('appBase.sectionTitle') }}</h2>
      </div>
      <div class="cim-app-base-divider"></div>
      <div class="cim-app-base-list">
        <div class="cim-app-base-item">
          <div class="cim-app-base-label">{{ t('appBase.field.appName') }}</div>
          <div class="cim-app-base-value">{{ state.appInfo.app_name }}</div>
        </div>
        <div class="cim-app-base-item">
          <div class="cim-app-base-label">{{ t('appBase.field.appKey') }}</div>
          <div class="cim-app-base-value">{{ state.appInfo.app_key }}</div>
        </div>
        <div class="cim-app-base-item">
          <div class="cim-app-base-label">{{ t('appBase.field.alias') }}</div>
          <div class="cim-app-base-value cim-app-base-alias">
            <input
              class="form-control cim-app-base-alias-input"
              type="text"
              maxlength="50"
              :placeholder="t('appBase.field.alias')"
              v-model="state.editingAlias"
              @input="onAliasInput"
              @keydown.enter="onSaveAlias"
            >
            <button
              type="button"
              class="cim-button cim-app-base-alias-save"
              :disabled="state.isSavingAlias"
              @click="onSaveAlias"
            >
              {{ t('common.dialog.save') }}
            </button>
            <div class="invalid-feedback feedback cim-app-base-alias-error" v-if="state.aliasErrorMsg">
              {{ state.aliasErrorMsg }}
            </div>
          </div>
        </div>
        <div class="cim-app-base-item">
          <div class="cim-app-base-label">{{ t('appBase.field.appSecret') }}</div>
          <div class="cim-app-base-value cim-app-base-secret">
            <span class="cim-secret-text">{{ state.isShowSecret ? state.appInfo.app_secret : state.appInfo.n_app_secret }}</span>
            <button type="button" class="cim-secret-btn" @click="onShowSecret">
              <span class="cicon cicon-hide"></span>
            </button>
          </div>
        </div>
        <div class="cim-app-base-item">
          <div class="cim-app-base-label">{{ t('appBase.field.expireTime') }}</div>
          <div class="cim-app-base-value">{{ getExpireTimeLabel() }}</div>
        </div>
        <div class="cim-app-base-item">
          <div class="cim-app-base-label">{{ t('appBase.field.licenseCount') }}</div>
          <div class="cim-app-base-value">
            {{ state.appInfo.n_max_user_count == -1 ? t('common.label.unlimited') : state.appInfo.n_max_user_count }}
          </div>
        </div>
        <div class="cim-app-base-item">
          <div class="cim-app-base-label">{{ t('appBase.field.status') }}</div>
          <div class="cim-app-base-value">{{ t('appBase.status.online') }}</div>
        </div>
        <div class="cim-app-base-item" v-for="field in state.urlFields" :key="field.key">
          <div class="cim-app-base-label">{{ t(field.label) }}</div>
          <div class="cim-app-base-value cim-app-base-alias">
            <input
              class="form-control cim-app-base-alias-input"
              type="text"
              :maxlength="MAX_URL_LENGTH"
              :placeholder="t(field.label)"
              v-model="field.editing"
              @input="onUrlInput(field)"
              @keydown.enter="onSaveUrl(field)"
            >
            <button
              type="button"
              class="cim-button cim-app-base-alias-save"
              :disabled="field.saving"
              @click="onSaveUrl(field)"
            >
              {{ t('common.dialog.save') }}
            </button>
            <div class="invalid-feedback feedback cim-app-base-alias-error" v-if="field.errorMsg">
              {{ field.errorMsg }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
