<script setup>
  import { ref, watch } from 'vue';
  import Dialog from './dialog.vue';
  import { t } from '@/i18n';

  const props = defineProps({
    show: Boolean,
    title: String,
    setting: Object,
    draft: Object,
    errors: Object,
    saving: Boolean,
  });
  const emit = defineEmits(['save', 'hide']);
  const fileInputs = ref({});
  const visibleSecrets = ref({});
  const activeJpushTab = ref(props.setting?.jpushTabs?.[0]?.type || '');
  const jpushOptionsExpanded = ref(false);

  watch(
    () => props.draft,
    () => {
      visibleSecrets.value = {};
      activeJpushTab.value = props.setting?.jpushTabs?.[0]?.type || '';
      jpushOptionsExpanded.value = false;
    }
  );

  watch(
    () => props.errors,
    (errors) => {
      if (optionFields().some((field) => errors?.[field.name])) {
        jpushOptionsExpanded.value = true;
      }
    }
  );

  function fieldLabel(field) {
    return field.labelKey ? t(field.labelKey, {}, field.label) : field.label;
  }

  function fieldPlaceholder(field) {
    if (field.secret && props.draft.original_package) {
      return t('appServices.push.hint.secretUnchanged');
    }
    if (field.integer) {
      return t('appServices.push.placeholder.integerOptional');
    }
    if (field.jsonObject) {
      return t('appServices.push.placeholder.jsonObjectOptional');
    }
    return t('appServices.push.placeholder.input', { name: fieldLabel(field) });
  }

  function errorText(field) {
    const error = props.errors[field.name];
    if (error === 'duplicate') {
      return t('appServices.push.validation.packageDuplicate');
    }
    if (error === 'integer') {
      return t('appServices.push.validation.integer', { name: fieldLabel(field) });
    }
    if (error === 'stringMap') {
      return t('appServices.push.validation.stringMap', { name: fieldLabel(field) });
    }
    return error ? t('appServices.push.validation.required', { name: fieldLabel(field) }) : '';
  }

  function activeJpushFields() {
    return props.setting?.jpushTabs?.find((tab) => tab.type === activeJpushTab.value)?.fields || [];
  }

  function baseFields() {
    return props.setting?.fields?.filter((field) => !field.optionsField) || [];
  }

  function optionFields() {
    return props.setting?.fields?.filter((field) => field.optionsField) || [];
  }

  function dialogClass() {
    if (props.setting?.type !== 'Jpush') {
      return 'cim-push-dialog';
    }
    return `cim-push-dialog cim-jpush-dialog${
      jpushOptionsExpanded.value ? ' is-options-expanded' : ''
    }`;
  }

  function setFileInput(field, element) {
    if (element) {
      fileInputs.value[field.model] = element;
    } else {
      delete fileInputs.value[field.model];
    }
  }

  function selectFile(field) {
    fileInputs.value[field.model]?.click();
  }

  function onFile(field, event) {
    props.draft[field.model] = event.target.files?.[0] || null;
  }

  function removeSelectedFile(field) {
    props.draft[field.model] = null;
    const input = fileInputs.value[field.model];
    if (input) {
      input.value = '';
    }
  }

  function fileName(field) {
    return (
      props.draft[field.model]?.name ||
      props.draft[field.name] ||
      t('appServices.push.status.noFile')
    );
  }

  function onSave() {
    if (!props.saving) {
      emit('save');
    }
  }

  function isSecretVisible(field) {
    return Boolean(visibleSecrets.value[field.name]);
  }

  function toggleSecretVisibility(field) {
    visibleSecrets.value[field.name] = !isSecretVisible(field);
  }

  function secretToggleLabel(field) {
    const key = isSecretVisible(field)
      ? 'appServices.push.action.hideSecret'
      : 'appServices.push.action.showSecret';
    return t(key, { name: fieldLabel(field) });
  }
</script>

<template>
  <Dialog
    :title="props.title"
    :show="props.show"
    :cls="dialogClass()"
    @hide="emit('hide')"
    @save="onSave"
  >
    <div class="cim-push-dialog-fields" v-if="props.setting && props.draft">
      <div class="cim-push-dialog-field" v-for="field in baseFields()" :key="field.name">
        <label class="cim-push-dialog-label">
          {{ fieldLabel(field) }}
          <span class="cim-push-required" v-if="field.required">*</span>
        </label>

        <div
          class="cim-push-dialog-input-wrap"
          v-if="field.type === 'input_text' || field.type === 'input_number'"
        >
          <input
            v-model="props.draft[field.name]"
            class="form-control cim-push-dialog-input"
            :class="{
              'is-invalid': props.errors[field.name],
              'has-secret-toggle': field.secret,
            }"
            :type="
              field.secret && !isSecretVisible(field)
                ? 'password'
                : field.type === 'input_number'
                  ? 'number'
                  : 'text'
            "
            :step="field.integer ? 1 : undefined"
            :inputmode="field.integer ? 'numeric' : undefined"
            :placeholder="fieldPlaceholder(field)"
            autocomplete="new-password"
          />
          <button
            v-if="field.secret"
            type="button"
            class="cim-push-secret-toggle"
            :class="{ 'is-hidden': !isSecretVisible(field) }"
            :aria-label="secretToggleLabel(field)"
            :title="secretToggleLabel(field)"
            @click="toggleSecretVisibility(field)"
          >
            <svg viewBox="0 0 24 24" aria-hidden="true">
              <path d="M2 12s3.5-6 10-6 10 6 10 6-3.5 6-10 6S2 12 2 12Z" />
              <circle cx="12" cy="12" r="2.75" />
            </svg>
          </button>
        </div>

        <div
          v-else-if="field.type === 'file'"
          class="cim-push-dialog-upload"
          :class="{ 'is-invalid': props.errors[field.name] }"
        >
          <input
            :ref="(element) => setFileInput(field, element)"
            class="cim-push-upload-input"
            type="file"
            @change="onFile(field, $event)"
          />
          <span class="cim-push-dialog-file-name">{{ fileName(field) }}</span>
          <button
            v-if="props.draft[field.model]?.name"
            type="button"
            class="cim-push-upload-button"
            @click="removeSelectedFile(field)"
          >
            {{ t('common.action.remove') }}
          </button>
          <button type="button" class="cim-push-upload-button" @click="selectFile(field)">
            {{ t('common.action.upload') }}
          </button>
        </div>

        <div v-else-if="field.type === 'radios'" class="cim-push-dialog-radios">
          <label class="cim-push-dialog-radio" v-for="option in field.radios" :key="option.value">
            <input
              v-model="props.draft[field.name]"
              class="form-check-input"
              type="radio"
              :name="`${props.setting.type}-${field.name}`"
              :value="option.value"
            />
            <span>{{ option.labelKey ? t(option.labelKey, {}, option.label) : option.label }}</span>
          </label>
        </div>

        <div class="cim-push-field-error" v-if="errorText(field)">{{ errorText(field) }}</div>
        <div class="cim-push-field-hint" v-else-if="field.secret && props.draft.original_package">
          {{ t('appServices.push.hint.secretUnchanged') }}
        </div>
      </div>

      <section
        class="cim-jpush-options"
        v-if="optionFields().length || props.setting.jpushTabs?.length"
      >
        <button
          type="button"
          class="cim-jpush-options-title"
          :aria-expanded="jpushOptionsExpanded"
          @click="jpushOptionsExpanded = !jpushOptionsExpanded"
        >
          <span class="cim-jpush-options-label">{{ t('appServices.push.section.options') }}</span>
          <span class="cim-jpush-options-line" aria-hidden="true"></span>
          <span
            class="cim-jpush-options-chevron"
            :class="{ 'is-expanded': jpushOptionsExpanded }"
            aria-hidden="true"
          ></span>
        </button>

        <div class="cim-jpush-options-content" v-if="jpushOptionsExpanded">
          <div class="cim-jpush-options-fields" v-if="optionFields().length">
            <div class="cim-push-dialog-field" v-for="field in optionFields()" :key="field.name">
              <label class="cim-push-dialog-label" :for="field.name">{{ fieldLabel(field) }}</label>
              <input
                :id="field.name"
                v-model="props.draft[field.name]"
                class="form-control cim-push-dialog-input"
                :class="{ 'is-invalid': props.errors[field.name] }"
                :type="field.type === 'input_number' ? 'number' : 'text'"
                :step="field.integer ? 1 : undefined"
                :inputmode="field.integer ? 'numeric' : undefined"
                :placeholder="fieldPlaceholder(field)"
              />
              <div class="cim-push-field-error" v-if="errorText(field)">{{ errorText(field) }}</div>
            </div>
          </div>

          <div class="cim-jpush-channel-title">{{
            t('appServices.push.section.channelOptions')
          }}</div>
          <div class="cim-jpush-tabs" role="tablist">
            <button
              v-for="tab in props.setting.jpushTabs"
              :key="tab.type"
              type="button"
              role="tab"
              class="cim-jpush-tab"
              :class="{ 'is-active': activeJpushTab === tab.type }"
              :aria-selected="activeJpushTab === tab.type"
              @click="activeJpushTab = tab.type"
            >
              {{ tab.labelKey ? t(tab.labelKey, {}, tab.label) : tab.label }}
            </button>
          </div>

          <div class="cim-jpush-tab-panel" role="tabpanel">
            <div
              class="cim-push-dialog-field"
              v-for="field in activeJpushFields()"
              :key="field.name"
            >
              <label class="cim-push-dialog-label" :for="field.name">{{ fieldLabel(field) }}</label>

              <label class="cim-jpush-checkbox" v-if="field.type === 'checkbox'">
                <input :id="field.name" v-model="props.draft[field.name]" type="checkbox" />
                <span>{{ t('appServices.push.option.enabled') }}</span>
              </label>

              <select
                v-else-if="field.type === 'select'"
                :id="field.name"
                v-model="props.draft[field.name]"
                class="form-select cim-push-dialog-input"
              >
                <option value="">{{ t('appServices.push.option.notSet') }}</option>
                <option v-for="option in field.options" :key="option.value" :value="option.value">
                  {{ option.labelKey ? t(option.labelKey, {}, option.label) : option.label }}
                </option>
              </select>

              <input
                v-else
                :id="field.name"
                v-model="props.draft[field.name]"
                class="form-control cim-push-dialog-input"
                :class="{ 'is-invalid': props.errors[field.name] }"
                :type="field.type === 'input_number' ? 'number' : 'text'"
                :step="field.integer ? 1 : undefined"
                :inputmode="field.integer ? 'numeric' : undefined"
                :placeholder="fieldPlaceholder(field)"
              />
              <div class="cim-push-field-error" v-if="errorText(field)">{{ errorText(field) }}</div>
            </div>
          </div>
        </div>
      </section>

      <div class="cim-push-saving" v-if="props.saving">{{
        t('appServices.push.status.saving')
      }}</div>
    </div>
  </Dialog>
</template>
