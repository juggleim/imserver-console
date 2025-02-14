<script setup>
import { reactive, watch, getCurrentInstance, nextTick } from 'vue';
import utils from '../../common/utils';
import { Misc } from "../../services";
import emitter from '../../common/emmit';
import { ErrorType } from "../../common/enum";
import JsonEditor from 'vue3-ts-jsoneditor';
import { apis } from './apis/api';
import { useRouter } from 'vue-router';
import { RESPONSE } from '../../common/enum';

let router = useRouter();
let context = getCurrentInstance();
let {
    currentRoute: {
      _rawValue: {
        params: { app_key },
      },
    },
  } = router;

let state = reactive({
  apis: apis,
  request: { },
  response: {},
  isJSONContent: false,
  currentAPI: { url: '' }
});
function onShow(item){
  utils.forEach(state.apis, (api) => {
    let _isFold = false;
    let _isActive = false;
    if(utils.isEqual(item.category, api.category)){
      _isFold = !item.isFold;
      _isActive = true;
    }
    // 清理已选中 child 状态
    utils.forEach(api.children, (child) => {
      child.isActive = false;
    });
    api.isFold = _isFold;
    api.isActive = _isActive;
  });
  state.currentAPI = { url: '' };
}
function onAPIClick(item){
  utils.forEach(state.apis, (api) => {
    api.isActive = false;
    utils.forEach(api.children, (child) => {
      child.isActive = utils.isEqual(item.url, child.url);
    });
  });
  state.currentAPI = item;
}
function onInputChange({ text }){
  let content = text.replace(/\n/g, '')
  state.isJSONContent = utils.isJSON(content);
}
function onSend(){
  if(!state.isJSONContent){
    return;
  }
  
  let { request } = state;
  if(utils.isEmpty(request)){
    return;
  }
  if(!utils.isString(request)){
    request = utils.toJSON(request)
  }
  let { currentAPI } = state;
  Misc.ajax({
    path: currentAPI.url,
    app_key: app_key,
    method: currentAPI.method.toUpperCase(),
    body: request
  }).then((result) => {
    state.response = utils.parse(result);
  });
}
watch(() => state.currentAPI, (value) => {
  let item = utils.clone(value);
  state.request = item.body;
  state.response = '';
  if(state.request){
    onInputChange({ text: utils.toJSON(state.request) });
  }
})
</script>
<template>
  <div class="mb-4 cim-api-box">
    <div class="cim-bk-form cim-api-sidebar">
      <ul class="cim-api-sidebar-nav">
        <li class="cim-api-nav-group" v-for="item in state.apis">
          <div class="cim-api-nav-link cicon cim-api-nav-link-toggle" @click="onShow(item)" :class="{'cim-api-nav-link-active': item.isActive}">{{ item.category }}</div>
          <ul class="cim-api-nav-items" :class="{'cim-api-nav-show': item.isFold}">
            <li class="cim-api-nav-item" :class="{'cim-api-nav-item-active': api.isActive}" v-for="api in item.children" @click="onAPIClick(api)">
              <span class="cim-api-nav-method" :class="['cim-api-nav-' + api.method]">{{ api.method }}</span>
              <span class="cim-api-nav-name">{{ api.name }}</span>
            </li>
          </ul>
        </li>
      </ul>
    </div>
    <div class="cim-api-main" v-if="state.currentAPI.url">
      <div class="cim-bk-form cim-api-header">
        <div class="cim-api-header-method" :class="['cim-api-header-'+state.currentAPI.method]">{{ state.currentAPI.method }}</div>
        <div class="cim-api-header-url">
          <input type="text" v-model="state.currentAPI.url">
        </div>
        <div class="cim-api-header-buttons">
          <div class="cim-button cim-button-bg" @click="onSend">发送</div>
        </div>
      </div>
      <div class="cim-bk-form cim-api-main-body">
        <div class="cim-api-request">
          <div class="jug-api-rr-header">
            <span>请求</span>
          </div>
          <JsonEditor
            mode="text"
            :mainMenuBar=false
            :navigationBar=false
            :statusBar=false
            :indentation=2
            @change="onInputChange"
            v-model="state.request"
          />
        </div>
        <div class="cim-api-response">
          <div class="jug-api-rr-header">响应</div>
          <JsonEditor
            mode="text"
            :mainMenuBar=false
            :navigationBar=false
            :statusBar=false
            :readOnly=true
            :indentation=2
            v-model="state.response"
          />
        </div>
      </div>
    </div>
  </div>
</template>
