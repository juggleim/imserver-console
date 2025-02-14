<script setup>
const props = defineProps(['item']);
const emit = defineEmits(['save'])

import { reactive, watch } from 'vue';
import utils from '../../common/utils';

let state = reactive({
  value: props.item.options[0].key
});

function onSave(){
  let { id } = props.item;
  let { value } = state;
  emit('save', { id, value });
}
watch(() => props.item.value, (value) => {
  utils.extend(state, { value });
})
</script>

<template>
  <div class="cim-sw-form">
    <div class="cim-form-check form-switch">
      <label class="form-check-label">{{ props.item.name }}</label>
      <select class="form-select" v-model="state.value">
        <option :value="op.key" v-for="op in props.item.options">{{ op.value }}</option>
      </select>
      <div class="cim-button" @click="onSave">保存</div>
    </div>
  </div>
</template>
