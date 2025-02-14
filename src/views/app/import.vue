<script setup>
import { reactive, getCurrentInstance } from 'vue';
import ModifyDialog from '../../components/dialog.vue';
import { ErrorType } from "../../common/enum";
import utils from '../../common/utils';
import { Application } from "../../services";
import { useRouter } from "vue-router";
import appTools from '../../common/app-tools';

let context = getCurrentInstance();
let router = useRouter();

let defaltApp = {
  name: '',
  licenseKey: ''
}

let state = reactive({
  apps: [],
  isShowEdit: false,
  nameErrorMsg: '',
  licenseErrorMsg: '',
  app: utils.clone(defaltApp)
});

function onShowEdit(isShow){
  state.isShowEdit = isShow;
}

function onSave(){
  let { app } = state;
  if(utils.isEmpty(app.name)){
    return state.nameErrorMsg = '名称不能为空';
  }
  if(utils.isEmpty(app.licenseKey)){
    return state.licenseErrorMsg = '授权不能为空';
  }

  Application.create(app).then(({ code, msg }) => {
    let icon = 'error', text = `导入失败 ${msg}`;
    if(utils.isEqual(code, ErrorType.SUCCESS_0.code)){
      icon = 'success';
      text = '创建成功';
      state.apps.push(app);
      state.app = utils.clone(defaltApp);
    }
    context.proxy.$toast({ icon, text, duration: 4000 });
  });
}

function getApps(){
  Application.getList().then(({ data: { items } }) => {
    let apps = items.map((item) => {
      item.created_time = utils.formatTime(item.created_time);
      item.ended_time = utils.formatTime(item.ended_time);
      item.cur_user_count = utils.numberWithCommas(item.cur_user_count)
      item.max_user_count = utils.numberWithCommas(item.max_user_count)
      return item;
    });
    state.apps = apps;
  });
}
getApps();
function onInput(){
  state.nameErrorMsg = '';
  state.licenseErrorMsg = '';
}
function onViewDetal(app){
  appTools.setCurrent(app);
  router.push({ name: 'ArguBase', params: { app_key: app.app_key } });
}
</script>
<template>
  <div class="mb-4">
    <div class="header cim-header">
      <div class="cim-title">应用列表</div>
      <div class="cicon cicon-add cim-button cim-button-bg" @click="onShowEdit(true)" @save="onSave()">导入应用</div>
    </div>
    <table class="table cim-table">
      <thead>
        <tr>
          <th>应用名称</th>
          <th class="jd-td-c">已使用授权数量</th>
          <th class="jd-td-c">未使用授权数量</th>
          <th class="jd-td-c">到期时间</th>
          <th class="jd-td-c">创建时间</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="app in state.apps">
          <td>{{ app.app_name }}</td>
          <td class="jd-td-c">{{ app.cur_user_count }}</td>
          <td class="jd-td-c">{{ app.max_user_count }}</td>
          <td class="jd-td-c">{{ app.ended_time }}</td>
          <td class="jd-td-c">{{ app.created_time }}</td>
          <td>
            <!-- <a class="btn-link cim-btn-link" type="button" @click="onShowEdit(true)">修改</a> -->
            <a class="btn-link cim-btn-link" type="button" @click="onViewDetal(app)">查看</a>
          </td>
        </tr>
      </tbody>
    </table>
    <ModifyDialog :show="state.isShowEdit" :title="'导入应用'" @hide="onShowEdit(false)" @save="onSave()">
      <div class="row g-2 cim-row">
          <div class="form-floating">
            <input class="form-control" v-model="state.app.name" placeholder="应用名称"  @input="onInput()">
            <label>应用名称</label>
            <div class="invalid-feedback feedback" v-if="state.nameErrorMsg">{{ state.nameErrorMsg }}</div>
          </div>
          <div class="form-floating">
            <input class="form-control" v-model="state.app.licenseKey" placeholder="License Key"  @input="onInput()">
            <label>License Key</label>
            <div class="invalid-feedback feedback" v-if="state.licenseErrorMsg">{{ state.licenseErrorMsg }}</div>
          </div>
      </div>
    </ModifyDialog>
  </div>
</template>