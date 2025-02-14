<script setup>
import DialogDomain from '../../components/dialog-modify-domain.vue';

const props = defineProps(['item']);
const emit = defineEmits(['save'])
import { reactive, watch } from 'vue';
import utils from '../../common/utils';
let state = reactive({
  isShowModal: false,
  domains: []
});

function onSave({ domains }){
  let { id } = props.item;
  state.domains = domains;
  emit('save', { id, value: utils.toJSON({ domains }) });
  onShowModal(false);
}
function onShowModal(isShow){
  state.isShowModal = isShow;
}
watch(() => props.item.value, (value) => {
  value = value || '{ "domains": [] }';
  let val = utils.parse(value);
  utils.extend(state, { domains: val.domains });
})
</script>

<template>
   <div class="cim-sw-form">
    <div class="cim-form-check form-switch">
      <label class="form-check-label">{{ props.item.name }}</label>
      <div class="cim-secrity-dm-content">
        <ul>
          <li v-if="state.domains.length == 0">未设置</li>
          <li v-else>{{ state.domains.join(';') }}</li>
        </ul>
      </div>
      <div class="cim-button" @click="onShowModal(true)">修改</div>
    </div>
  </div>
  <DialogDomain :show="state.isShowModal" :domains="state.domains" @hide="onShowModal(false)" @save="onSave"></DialogDomain>
</template>
