<script setup>
  import { computed, getCurrentInstance, reactive, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { Application } from '../../services';
  import { RESPONSE } from '../../common/enum';
  import { t } from '@/i18n';
  import PageSection from '@/components/page-section.vue';
  import PushConfigDialog from '@/components/push-config-dialog.vue';
  import {
    PUSH_CHANNELS,
    buildPushTextExtra,
    buildIosPushParams,
    createPushDraft,
    getPushCardValue,
    hasPushErrors,
    isPushFieldActive,
    p8KeyStatus,
    validatePushDraft,
  } from './push-config.mjs';

  const router = useRouter();
  const appKey = router.currentRoute.value.params.app_key;
  const context = getCurrentInstance();
  const settings = reactive(
    PUSH_CHANNELS.map((setting) => ({
      ...setting,
      items: [],
      loading: false,
      failed: false,
    }))
  );
  const channel = ref(settings[0].type);
  const dialog = reactive({
    show: false,
    mode: 'add',
    draft: null,
    errors: {},
    saving: false,
    conflict: false,
  });

  const currentSetting = computed(
    () => settings.find((setting) => setting.type === channel.value) || settings[0]
  );
  const dialogTitle = computed(() => {
    const key =
      dialog.mode === 'edit'
        ? 'appServices.push.dialog.editTitle'
        : 'appServices.push.dialog.addTitle';
    return t(key, { channel: getSettingName(currentSetting.value) });
  });

  function getSettingName(setting) {
    return setting.nameKey ? t(setting.nameKey, {}, setting.name) : setting.name;
  }

  function getFieldLabel(field) {
    return field.labelKey ? t(field.labelKey, {}, field.label) : field.label;
  }

  function toast(icon, text) {
    context.proxy.$toast({ icon, text });
  }

  function ensureSuccess(result) {
    if (result?.code === RESPONSE.SUCCESS) {
      return result;
    }
    const error = new Error(result?.msg || 'request failed');
    error.code = result?.code;
    throw error;
  }

  async function loadSetting(setting) {
    setting.loading = true;
    setting.failed = false;
    try {
      let result;
      if (setting.kind === 'ios') {
        result = await Application.getIosPushConfigList({ app_key: appKey });
      } else if (setting.kind === 'fcm') {
        result = await Application.getFcmPushConfigList({ app_key: appKey });
      } else {
        result = await Application.getAndroidPushConfigList({
          app_key: appKey,
          push_channel: setting.type,
        });
      }
      ensureSuccess(result);
      setting.items = Array.isArray(result.data) ? result.data : [];
    } catch (error) {
      setting.failed = true;
      setting.items = [];
      toast('error', t('appServices.push.feedback.queryFailed'));
    } finally {
      setting.loading = false;
    }
  }

  function closeDialog() {
    if (dialog.saving) {
      return;
    }
    dialog.show = false;
    dialog.draft = null;
    dialog.errors = {};
  }

  function onTab(setting) {
    closeDialog();
    channel.value = setting.type;
    loadSetting(setting);
  }

  function openAdd() {
    dialog.conflict = false;
    dialog.mode = 'add';
    dialog.draft = createPushDraft(currentSetting.value);
    dialog.errors = {};
    dialog.show = true;
  }

  function openEdit(item) {
    dialog.conflict = false;
    dialog.mode = 'edit';
    dialog.draft = createPushDraft(currentSetting.value, item);
    dialog.errors = {};
    dialog.show = true;
  }

  async function saveDraft() {
    if (dialog.conflict || dialog.saving) return;
    const setting = currentSetting.value;
    dialog.errors = validatePushDraft(setting, dialog.draft, setting.items);
    if (hasPushErrors(dialog.errors)) {
      return;
    }

    const draft = dialog.draft;
    draft.package = String(draft.package).trim();
    dialog.saving = true;
    try {
      let result;
      if (setting.kind === 'ios') {
        const params = buildIosPushParams(appKey, draft);
        if (
          !draft.original_package ||
          params.file?.name ||
          params.voipFile?.name ||
          params.p8File?.name
        ) {
          result = await Application.uploadIosPushConfig(params);
        } else {
          result = await Application.setIosPushConfig(params);
        }
      } else if (setting.kind === 'fcm') {
        result = await Application.uploadFcmPushConfig({
          app_key: appKey,
          package: draft.package,
          original_package: draft.original_package,
          file: draft.file,
        });
      } else {
        result = await Application.setAndroidPushConfig({
          app_key: appKey,
          push_channel: setting.type,
          package: draft.package,
          original_package: draft.original_package,
          extra: buildPushTextExtra(setting, draft),
        });
      }
      ensureSuccess(result);
      toast('success', t('appServices.push.feedback.saveSuccess'));
      dialog.show = false;
      dialog.draft = null;
      dialog.errors = {};
      await loadSetting(setting);
    } catch (error) {
      if (setting.kind === 'ios' && error.code === 409) {
        dialog.conflict = true;
        return;
      }
      if (error.code === RESPONSE.PUSH_CONF_EXISTED) {
        dialog.errors = { ...dialog.errors, package: 'duplicate' };
      }
      toast(
        'error',
        t('appServices.push.feedback.saveFailed', {
          code: error.code || '',
          msg: error.message || '',
        })
      );
    } finally {
      dialog.saving = false;
    }
  }

  async function reloadConflictedDraft() {
    if (dialog.saving) return;
    dialog.saving = true;
    const originalPackage = dialog.draft.original_package;
    const setting = currentSetting.value;
    try {
      await loadSetting(setting);
      if (setting.failed) return;
      const latest = setting.items.find((item) => item.package === originalPackage);
      if (latest) openEdit(latest);
      else {
        dialog.show = false;
        dialog.draft = null;
        dialog.conflict = false;
        toast('error', t('appServices.push.feedback.configMissing'));
      }
    } finally {
      dialog.saving = false;
    }
  }

  function getCardValue(item, field) {
    const value = getPushCardValue(item, field);
    if (field.type === 'radios') {
      const option = field.radios.find((candidate) => candidate.value === value);
      return option ? (option.labelKey ? t(option.labelKey, {}, option.label) : option.label) : '';
    }
    return value;
  }

  function cardFields(setting, item) {
    return setting.fields.filter((field) => field.cardVisible && isPushFieldActive(field, item));
  }

  loadSetting(settings[0]);
</script>

<template>
  <PageSection
    title-key="menu.app.pushSettings"
    shell-class="cim-push-page"
    body-class="cim-push-content"
  >
    <ul class="nav nav-underline-border cim-push-tabs" role="tablist">
      <li class="nav-item sw-nav-item" v-for="setting in settings" :key="setting.type">
        <button
          type="button"
          class="nav-link cicon cicon-free"
          :class="{ active: channel === setting.type }"
          :aria-selected="channel === setting.type"
          @click="onTab(setting)"
        >
          {{ getSettingName(setting) }}
        </button>
      </li>
    </ul>

    <div class="cim-push-state" v-if="currentSetting.loading">
      {{ t('appServices.push.status.loading') }}
    </div>
    <div class="cim-push-state cim-push-state-error" v-else-if="currentSetting.failed">
      <span>{{ t('appServices.push.feedback.queryFailed') }}</span>
      <button type="button" class="cim-button" @click="loadSetting(currentSetting)">
        {{ t('appServices.push.action.retry') }}
      </button>
    </div>
    <div class="cim-push-card-grid" v-else>
      <article class="cim-push-card" v-for="item in currentSetting.items" :key="item.package">
        <header class="cim-push-card-header">
          <h3 class="cim-push-card-title">{{ item.package }}</h3>
        </header>
        <div class="cim-push-card-body">
          <div
            class="cim-push-card-row"
            v-for="field in cardFields(currentSetting, item)"
            :key="field.name"
          >
            <span class="cim-push-card-label">{{ getFieldLabel(field) }}</span>
            <span class="cim-push-card-value" :class="{ 'is-unset': !getCardValue(item, field) }">
              {{ getCardValue(item, field) || t('appServices.push.status.unset') }}
            </span>
          </div>
          <div v-if="item.auth_type === 'p8'" class="cim-push-card-row">
            <span class="cim-push-card-label">{{ t('appServices.push.field.p8Status') }}</span>
            <span class="cim-push-card-value">{{ t(p8KeyStatus(item)) }}</span>
          </div>
        </div>
        <footer class="cim-push-card-footer">
          <button type="button" class="cim-push-settings-button" @click="openEdit(item)">
            {{ t('appServices.push.action.settings') }}
          </button>
        </footer>
      </article>

      <button
        type="button"
        class="cim-push-add-card"
        @click="openAdd"
        :aria-label="t('appServices.push.action.add')"
      >
        <span class="cim-push-add-icon" aria-hidden="true">+</span>
        <span>{{ t('appServices.push.action.add') }}</span>
      </button>
    </div>

    <PushConfigDialog
      v-if="dialog.draft"
      :show="dialog.show"
      :title="dialogTitle"
      :setting="currentSetting"
      :draft="dialog.draft"
      :errors="dialog.errors"
      :saving="dialog.saving"
      :conflict="dialog.conflict"
      @reload="reloadConflictedDraft"
      @hide="closeDialog"
      @save="saveDraft"
    />
  </PageSection>
</template>
