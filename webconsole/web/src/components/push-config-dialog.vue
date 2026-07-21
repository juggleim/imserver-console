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

  watch(
    () => props.draft,
    () => {
      visibleSecrets.value = {};
    }
  );

  function fieldLabel(field) {
    return field.labelKey ? t(field.labelKey, {}, field.label) : field.label;
  }

  function fieldPlaceholder(field) {
    if (field.secret && props.draft.original_package) {
      return t('appServices.push.hint.secretUnchanged');
    }
    return t('appServices.push.placeholder.input', { name: fieldLabel(field) });
  }

  function errorText(field) {
    const error = props.errors[field.name];
    if (error === 'duplicate') {
      return t('appServices.push.validation.packageDuplicate');
    }
    return error ? t('appServices.push.validation.required', { name: fieldLabel(field) }) : '';
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
    cls="cim-push-dialog"
    @hide="emit('hide')"
    @save="onSave"
  >
    <div class="cim-push-dialog-fields" v-if="props.setting && props.draft">
      <div class="cim-push-dialog-field" v-for="field in props.setting.fields" :key="field.name">
        <label class="cim-push-dialog-label">
          {{ fieldLabel(field) }}
          <span class="cim-push-required" v-if="field.required">*</span>
        </label>

        <div class="cim-push-dialog-input-wrap" v-if="field.type === 'input_text'">
          <input
            v-model="props.draft[field.name]"
            class="form-control cim-push-dialog-input"
            :class="{
              'is-invalid': props.errors[field.name],
              'has-secret-toggle': field.secret,
            }"
            :type="field.secret && !isSecretVisible(field) ? 'password' : 'text'"
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
      <div class="cim-push-saving" v-if="props.saving">{{
        t('appServices.push.status.saving')
      }}</div>
    </div>
  </Dialog>
</template>
