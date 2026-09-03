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

const FIELD_DEFS = [
  {
    key: 'alias',
    field: 'alias',
    label: 'appBase.field.alias',
    maxLength: 50,
    fallback: '0',
    requiredMsg: 'appBase.validation.aliasRequired',
    tooLongMsg: 'appBase.validation.aliasTooLong',
    savedMsg: 'appBase.feedback.aliasSaved',
    saveFailedMsg: 'appBase.feedback.aliasSaveFailed',
    requestFailedMsg: 'appBase.feedback.aliasRequestFailed',
    save: ({ app_key, value }) => Application.updateAlias({ app_key, alias: value }),
  },
  {
    key: 'ws',
    field: 'ws_url',
    label: 'appBase.field.wsUrl',
    maxLength: MAX_URL_LENGTH,
    fallback: '',
    requiredMsg: 'appBase.validation.urlRequired',
    tooLongMsg: 'appBase.validation.urlTooLong',
    savedMsg: 'appBase.feedback.urlSaved',
    saveFailedMsg: 'appBase.feedback.urlSaveFailed',
    requestFailedMsg: 'appBase.feedback.urlRequestFailed',
    save: ({ app_key, value }) => Application.updateWsUrl({ app_key, url: value }),
  },
  {
    key: 'api',
    field: 'api_url',
    label: 'appBase.field.apiUrl',
    maxLength: MAX_URL_LENGTH,
    fallback: '',
    requiredMsg: 'appBase.validation.urlRequired',
    tooLongMsg: 'appBase.validation.urlTooLong',
    savedMsg: 'appBase.feedback.urlSaved',
    saveFailedMsg: 'appBase.feedback.urlSaveFailed',
    requestFailedMsg: 'appBase.feedback.urlRequestFailed',
    save: ({ app_key, value }) => Application.updateApiUrl({ app_key, url: value }),
  },
  {
    key: 'app',
    field: 'app_url',
    label: 'appBase.field.appUrl',
    maxLength: MAX_URL_LENGTH,
    fallback: '',
    requiredMsg: 'appBase.validation.urlRequired',
    tooLongMsg: 'appBase.validation.urlTooLong',
    savedMsg: 'appBase.feedback.urlSaved',
    saveFailedMsg: 'appBase.feedback.urlSaveFailed',
    requestFailedMsg: 'appBase.feedback.urlRequestFailed',
    save: ({ app_key, value }) => Application.updateAppUrl({ app_key, url: value }),
  },
];

function makeFieldState(def) {
  return utils.extend({}, [def, {
    mode: 'view',
    draft: '',
    saving: false,
    errorMsg: '',
  }]);
}

let state = reactive({
  appInfo: {
    restricted_fields: {}
  },
  isShowSecret: false,
  aliasField: makeFieldState(FIELD_DEFS[0]),
  urlFields: FIELD_DEFS.slice(1).map(makeFieldState),
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

function onFieldEdit(field) {
  field.draft = state.appInfo[field.field] ?? field.fallback ?? '';
  field.errorMsg = '';
  field.mode = 'edit';
}

function onFieldInput(field) {
  field.errorMsg = '';
}

function onFieldCancel(field) {
  field.mode = 'view';
  field.draft = '';
  field.errorMsg = '';
}

function onFieldSave(field) {
  if (field.saving) {
    return;
  }
  const value = String(field.draft || '').trim();
  const label = t(field.label);
  if (!value) {
    field.errorMsg = t(field.requiredMsg, { field: label });
    return;
  }
  if ([...value].length > field.maxLength) {
    field.errorMsg = t(field.tooLongMsg, { field: label });
    return;
  }

  field.saving = true;
  field.save({ app_key, value }).then(({ code, msg }) => {
    if (utils.isEqual(code, ErrorType.SUCCESS_0.code)) {
      state.appInfo[field.field] = value;
      field.mode = 'view';
      field.draft = '';
      context.proxy.$toast({ icon: 'success', text: t(field.savedMsg, { field: label }) });
      return;
    }
    context.proxy.$toast({
      icon: 'error',
      text: t(field.saveFailedMsg, { field: label, code, msg }),
    });
  }).catch(() => {
    context.proxy.$toast({ icon: 'error', text: t(field.requestFailedMsg, { field: label }) });
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
          <div class="cim-app-base-value cim-app-base-alias" :class="{ 'is-editing': state.aliasField.mode === 'edit' }">
            <template v-if="state.aliasField.mode === 'edit'">
              <input
                class="form-control cim-app-base-alias-input"
                type="text"
                :maxlength="state.aliasField.maxLength"
                :placeholder="t(state.aliasField.label)"
                v-model="state.aliasField.draft"
                @input="onFieldInput(state.aliasField)"
                @keydown.enter="onFieldSave(state.aliasField)"
              >
              <button
                type="button"
                class="cim-button cim-app-base-alias-save"
                :disabled="state.aliasField.saving"
                @click="onFieldSave(state.aliasField)"
              >
                {{ t('common.dialog.save') }}
              </button>
              <button
                type="button"
                class="cim-button cim-app-base-btn-ghost"
                :disabled="state.aliasField.saving"
                @click="onFieldCancel(state.aliasField)"
              >
                {{ t('common.action.cancel') }}
              </button>
              <div class="invalid-feedback feedback cim-app-base-alias-error" v-if="state.aliasField.errorMsg">
                {{ state.aliasField.errorMsg }}
              </div>
            </template>
            <template v-else>
              <span class="cim-app-base-edit-text">{{ state.appInfo.alias || state.aliasField.fallback }}</span>
              <button
                type="button"
                class="cim-button cim-app-base-btn-ghost"
                @click="onFieldEdit(state.aliasField)"
              >
                {{ t('common.action.edit') }}
              </button>
            </template>
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
          <div class="cim-app-base-value cim-app-base-alias" :class="{ 'is-editing': field.mode === 'edit' }">
            <template v-if="field.mode === 'edit'">
              <input
                class="form-control cim-app-base-alias-input"
                type="text"
                :maxlength="field.maxLength"
                :placeholder="t(field.label)"
                v-model="field.draft"
                @input="onFieldInput(field)"
                @keydown.enter="onFieldSave(field)"
              >
              <button
                type="button"
                class="cim-button cim-app-base-alias-save"
                :disabled="field.saving"
                @click="onFieldSave(field)"
              >
                {{ t('common.dialog.save') }}
              </button>
              <button
                type="button"
                class="cim-button cim-app-base-btn-ghost"
                :disabled="field.saving"
                @click="onFieldCancel(field)"
              >
                {{ t('common.action.cancel') }}
              </button>
              <div class="invalid-feedback feedback cim-app-base-alias-error" v-if="field.errorMsg">
                {{ field.errorMsg }}
              </div>
            </template>
            <template v-else>
              <span class="cim-app-base-edit-text">{{ state.appInfo[field.field] || field.fallback || '-' }}</span>
              <button
                type="button"
                class="cim-button cim-app-base-btn-ghost"
                @click="onFieldEdit(field)"
              >
                {{ t('common.action.edit') }}
              </button>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
