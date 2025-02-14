<script setup>
import { reactive, getCurrentInstance, nextTick } from 'vue';
import { APP_TYPE, EVENT_NAME, ErrorType, STORAGE } from "../common/enum";
import utils from '../common/utils';
import { Application } from "../services";
import menuTools from './layout/menu-tools';
import appTools from '../common/app-tools';
import emitter from '../common/emmit';
import Storage from '../common/storage';

const context = getCurrentInstance();
let currentUser = Storage.get(STORAGE.USER_TOKEN);

let state = reactive({
  apps: [ ],
});
function onCopyAppKey(appkey) {
  console.log(appkey);
}
let dashboardApp = {
  offset: '',
  has_more: true
};
function getApps(isFirst, callback = utils.noop){
  if(!dashboardApp.has_more){
    callback();
    return context.proxy.$toast({ icon: 'warn', text: '没有更多啦' });
  }
  Application.getList({ offset: dashboardApp.offset }).then(({ code, data }) => {
    let error = ErrorType.USER_TOKEN_EXPIRE;
    if(utils.isEqual(code, error.code)){
      context.proxy.$toast({ icon: 'error', text: error.msg });
      return menuTools.goLoginPage();
    }
    let { items, offset, has_more } = data;
    utils.extend(dashboardApp, { offset, has_more });

    let apps = items.map((item) => {
      item.created_time = utils.formatTime(item.created_time);
      item.ended_time = item.ended_time == -1 ? '永久有效' : utils.formatTime(item.ended_time);
      item.cur_user_count = utils.numberWithCommas(item.cur_user_count);
      item.max_user_count = utils.numberWithCommas(item.max_user_count);
      item.kind = utils.isEqual(item.app_type, APP_TYPE.PRIVATE) ? '私有云' : '公有云';
      return item;
    });
    if(isFirst){
      state.apps = apps;
    }else{
      utils.forEach(apps, (app) => {
        state.apps.push(app);
      })
    }
    callback();
  });
}
getApps(true);
function onViewDetail(app){
  appTools.setCurrent(app);
  menuTools.goBasePage(app);
}
function onCreateApp(){
  emitter.$emit(EVENT_NAME.CREATE_APP);
}

let canscroll = true;
nextTick(() => {
  let { dashApps } = context.refs;
  dashApps.addEventListener("scroll", () => {
    let scrollTop = dashApps.scrollTop;
    let scrollHeight = dashApps.scrollHeight;
    let rectHeight = dashApps.getBoundingClientRect().height;
    let isNeedLoad = scrollHeight - scrollTop - rectHeight < 50;
    if (isNeedLoad && canscroll) {
      canscroll = false;
      getApps(false, () => {
        canscroll = true;
      });
    }
  });
});


</script>
<template>
  <div class="card-body db-layout">
    <div class="db-contanier">
      <div class="db-header nav-underline-border">
        <strong class="db-header-title">我的应用</strong>
      </div>
      <div class="db-apps" ref="dashApps">
        <div class="db-app">
          <div class="db-app-header bd-bottom">
            <div class="db-app-name">{{ currentUser.is_commercial ? '导入应用' : '创建应用' }}</div>
          </div>
          <div class="db-app-body db-create-body">
            <div class="db-app-create cicon cicon-add" @click="onCreateApp"></div>
          </div>
        </div>
        <div class="db-app" v-for="app in state.apps">
          <div class="db-app-header bd-bottom">
            <div class="db-app-name">{{ app.name }}</div>
          </div>
          <div class="db-app-body">
            <div class="row db-row">
              <div class="col-md-6 db-left">应用名称</div>
              <div class="col-md-6 db-right">{{ app.app_name }}</div>
            </div>
            <div class="row db-row">
              <div class="col-md-6 db-left">App-Key</div>
              <div class="col-md-6 db-right cicon cicon-copy" @click="onCopyAppKey(app.app_key)">{{ app.app_key }}</div>
            </div>
            <div class="row db-row">
              <div class="col-md-6 db-left">创建时间</div>
              <div class="col-md-6 db-right">{{ app.created_time }}</div>
            </div>
            <div class="row db-row">
              <div class="col-md-6 db-left">已注册用户</div>
              <div class="col-md-6 db-right">{{ app.cur_user_count }}</div>
            </div>
            <div class="row db-row">
              <div class="col-md-6 db-left">授权用户总数</div>
              <div class="col-md-6 db-right">{{ app.max_user_count == -1 ? '无限制' : app.max_user_count}}</div>
            </div>
            <div class="row db-row">
              <div class="col-md-6 db-left">部署方式</div>
              <div class="col-md-6 db-right">{{ app.kind }}</div>
            </div>
            <div class="row db-row">
              <div class="col-md-6 db-left">到期时间</div>
              <div class="col-md-6 db-right">{{ app.ended_time }}</div>
            </div>
          </div>
          <div class="db-app-footer bd-top">
            <a class="btn" @click="onViewDetail(app)">查看明细</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
