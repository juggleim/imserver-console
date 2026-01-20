<script setup>
import { reactive, getCurrentInstance } from 'vue';
import { useRouter } from "vue-router";
import utils from '../../common/utils';
import { Application } from "../../services";
import common from '../../common/common';

let router = useRouter();
let { currentRoute: { _rawValue: { params: { app_key } } } } = router;
const context = getCurrentInstance();

Application.getEventHook({ app_key }).then(({ data = {} }) => {
  let { event_sub_config, event_sub_switch } = data;
  utils.extend(state, { 
    config: event_sub_config,
    hooks: event_sub_switch
  });
});
let state = reactive({
  hooks: [],
  config: {},
});
function onHookChanged(e){
  utils.forEach(state.hooks, (hook) => {
    utils.forEach(hook.items, (item) => {
      if(utils.isEqual(item.key, e.target.id)){
        item.value = Number(e.target.checked);
      }
    })
  });
  onSaveConfig();
}
function onSaveConfig(){
  let { config, hooks } = state;
  let { event_sub_url } = config;
  if(!common.isValidateUrl(event_sub_url) && event_sub_url !== ''){
    return context.proxy.$toast({ icon: 'error', text: '回调地址格式不正确，检查是否包含协议头' });
  }
  Application.setEventHook({ app_key,  config, hooks }).then(() => {
    context.proxy.$toast({ icon: 'success', text: '保存成功' });
  });
}
</script>
<template>
  <div class="mb-4 app-base cim-cb-box">
    <div class="row cim-cb-row cim-cb-header">
      <div class="cim-cb-form">
        <label class="col-sm-1 col-form-label">
          <div class="cim-cb-input-item">设置事件地址</div>
          <div class="cim-cb-input-item">设置鉴权凭证</div>
        </label>
        <div class="col-sm-7">
          <input class="form-control cim-cb-input-item" type="text" v-model="state.config.event_sub_url" placeholder="请输入事件地址,示例 http[s]://example.com/submsg">
          <input class="form-control cim-cb-input-item" type="text" v-model="state.config.event_sub_auth" placeholder="请输入鉴权凭证，凭证会与事件一起回调给业务服务器">
        </div>
        <div class="col-sm-3 cim-cb-btns">
          <div class="cim-button cim-button-bg" @click="onSaveConfig">保存</div>
        </div>
      </div>
    </div>

    <div class="row cim-cb-row cim-cb-body">
      <div class="cim-cb-group" v-for="hook in state.hooks">
        <div class="cim-cb-group-title">{{ hook.name }}</div>
        <div class="cim-cb-group-content">
          <div class="cim-form-check form-switch" v-for="item in hook.items">
            <label class="form-check-label">{{ item.name }}</label>
            <input class="form-check-input" type="checkbox" :checked="item.value" :id="item.key" @change="onHookChanged">
          </div>
        </div>
      </div>
    </div>
  </div>
</template>



