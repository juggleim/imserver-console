<script setup>
import { reactive, getCurrentInstance } from 'vue';
import FInput from "../../components/func/input.vue";
import FSwitch from "../../components/func/switch.vue";
import FInputModal from "../../components/func/input-modal.vue";
import FSelect from "../../components/func/select.vue";
import { FUNC_TYPE } from "../../common/enum";
import utils from '../../common/utils';
import { Application } from "../../services";
import { useRouter } from "vue-router";

let router = useRouter();
let { currentRoute: { _rawValue: { params: { app_key } } } } = router;

const context = getCurrentInstance();

let settings = [
{
    type: 'app', 
    name: '应用相关', 
    list: [ 
      { id: 'token_effective_minute', type: 'input', name: 'Token 有效时长（小时）', value: 0 },
      { id: 'kick_mode', type: 'switch', name: '允许同设备多端登录', value: 0 },
      { id: 'security_domains', type: FUNC_TYPE.INPUT_MODAL, name: '安全域名', value: '{ "domains": [] }' },
    ] 
  },
  {
    type: 'message', 
    name: '消息相关', 
    list: [ 
      { id: 'is_hide_msg_before_join_group', type: 'switch', name: '入群后获取之前的历史消息', value: 0 },
      { id: 'not_check_grp_member', type: 'switch', name: '不在群组是否可以获取群消息', value: 0 },
      { id: 'his_msg_save_day', type: 'select', name: '历史消息存储时长 (天)', value: '7', options: [{ key: '7', value: '7 天' }, { key: '360', value: '1 年' }] },
    ] 
  },
  {
    type: 'group', 
    name: '群组相关', 
    list: [ 
      { id: 'max_grp_member_count', type: 'input', name: '群人数上限', value: 1000 },
    ] 
  },
  {
    type: 'chatroom', 
    name: '聊天室相关', 
    list: [ 
      { id: 'chrm_msg_cache_max_count', type: 'input', name: '单个聊天室消息桶大小', value: 50 },
      { id: 'chrm_att_max_count', type: 'input', name: '单个聊天室属性数量', value: 100 },
      { id: 'chrm_event_ntf', type: 'switch', name: '是否开启聊天室事件通知', value: false },
      { id: 'chrm_event_cache_max_count', type: 'input', name: '单个聊天室事件桶大小', value: 50 },
    ] 
  },
];
let state = reactive({
  settings: settings,
  current: settings[0].type
});

function onTab(setting){
  utils.extend(state, {
    current: setting.type
  });
}

function onSave(item){
  Application.updateSetting({...item, app_key}).then(() => {
    context.proxy.$toast({ icon: 'success', text: '保存成功' });
  });
}
function iterate(list, callback){
  utils.forEach(list, (item) => {
    utils.forEach(item.list, (i) => {
      callback(i);
    });
  });
}
function search(){
  let config_keys = [];
  iterate(settings, (item) => {
    config_keys.push(item.id);
  });
  Application.getSetting({ app_key, config_keys }).then(({ data}) => {
    let { configs } = data;
    iterate(state.settings, (item) => {
      utils.forEach(configs, (v, k) => {
        if(utils.isEqual(item.id, k)){
          item.value = v;
        }
      });
    });
  });
}
search();
</script>
<template>
   <div class="md-4">
    <ul class="nav nav-underline-border" role="tablist">
      <li class="nav-item sw-nav-item" v-for="setting in state.settings" @click="onTab(setting)">
        <a class="nav-link cicon cicon-free" :class="{'active': utils.isEqual(state.current, setting.type)}">{{ setting.name }}</a>
      </li>
    </ul>
    <div class="tab-content rounded-bottom">
      <div class="tab-pane p-3" v-for="setting in state.settings" :class="{'active': utils.isEqual(state.current, setting.type)}">
        <div class="row cim-sw-row">
          <div class="col-sm-4 cim-sw-col" v-for="item in setting.list">
              <FInput  v-if="utils.isEqual(item.type, FUNC_TYPE.INPUT)" :item="item" @save="onSave"></FInput>
              <FSelect v-if="utils.isEqual(item.type, FUNC_TYPE.SELECT)" :item="item" @save="onSave"></FSelect>
              <FSwitch v-if="utils.isEqual(item.type, FUNC_TYPE.SWITCH)" :item="item" @save="onSave"></FSwitch>
              <FInputModal v-if="utils.isEqual(item.type, FUNC_TYPE.INPUT_MODAL)" :item="item" @save="onSave"></FInputModal>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
