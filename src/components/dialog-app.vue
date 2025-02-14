<script setup>
import ModifyDialog from './dialog.vue';
import { STORAGE, ErrorType, APP_TYPE, DIALOG_TYPE, EVENT_NAME } from "../common/enum";
import { reactive, watch, getCurrentInstance } from 'vue';
import utils from "../common/utils";
import Storage from '../common/storage';
import { Application } from "../services";
import emitter from '../common/emmit';

const props = defineProps(['title', 'text']);
const emit = defineEmits(['hide'])
const context = getCurrentInstance();
let user = Storage.get(STORAGE.USER_TOKEN);

let state = reactive({
  nameError: '',
  licenseError: '',
  isShowNewEdit: false,
  name: '',
  license: '',
  title: user.is_commercial ? '导入应用' : '创建应用',
});

function onShowNewEdit(isShow){
  state.isShowNewEdit = isShow;
}
function onInput(){
  utils.extend(state, { nameError: '', licenseError: '' });
}

function onSave(){
  if(user.is_commercial){
    return importApp();
  }
  return createApp();
}

function createApp(){
  let { name } = state;
  if(utils.isEmpty(name)){
    return state.nameError = '名字不可为空';
  }
  Application.create({ name }).then(onCompleted);
}

function importApp(){
  let { license } = state;
  if(utils.isEmpty(license)){
    return state.licenseError = 'License 不可为空';
  }
  Application.importApp({ license }).then(onCompleted);
}

function onCompleted(result) {
  let { code, msg, data } = result;
  let { is_commercial } = user;
  let tname = is_commercial ? '导入' : '保存';
  let icon = 'error', text = `${tname}失败: ${code} ${msg}`;
  if (utils.isEqual(code, ErrorType.SUCCESS_0.code)) {
    icon = 'success';
    text = `${tname}成功`;
    let { app_key, app_name, app_type } = data;
    let app = { app_key, app_name, app_type, kind: '私有云' };
    emitter.$emit(EVENT_NAME.CREATE_APP_FINISHED, app);
    onHide();
  }
  context.proxy.$toast({icon, text, duration: 4000});
}

function onHide(){
  utils.extend(state, { name: '', license: '' });
  emit('hide', {});
}

watch(() => props.isShow, (value) => {
});

</script>

<template>
  <ModifyDialog :cls="'cim-dialog-createapp'" :title="state.title" @hide="onHide" @save="onSave">
    <div class="row g-2 cim-row"v-if="user.is_commercial">
      <div class="form-floating cimfrom-floating">
        <input class="form-control" placeholder="License" v-model="state.license" @input="onInput">
        <label>License Key</label>
        <div class="invalid-feedback feedback" >{{ state.licenseError }}</div>
      </div>
    </div>
    <div class="row g-2 cim-row" v-else>
      <div class="form-floating cimfrom-floating">
        <input class="form-control" placeholder="应用名称" v-model="state.name" @input="onInput">
        <label>应用名称</label>
        <div class="invalid-feedback feedback" >{{ state.nameError }}</div>
      </div>
    </div>
  </ModifyDialog>
</template>
