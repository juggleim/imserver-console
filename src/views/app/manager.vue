<script setup>
import { reactive, getCurrentInstance } from 'vue';
import ModifyDialog from '../../components/dialog.vue';
import { ErrorType } from "../../common/enum";
import utils from '../../common/utils';
import { Application } from "../../services";

let context = getCurrentInstance();

let defaultApp = {
  name: '',
  count: 0
}

let state = reactive({
  apps: [],
  isShowEdit: false,
  licenses: [
    { value: 1, label: '0-200' },
    { value: 2, label: '200-500' },
    { value: 3, label: '500-1000' },
  ],
  app: utils.clone(defaultApp)
});

function onShowEdit(isShow){
  state.isShowEdit = isShow;
}

function onSave(){
  let app = { name: '', count: '' };
  Application.create(app).then(({ code, msg }) => {
    let icon = 'error', text = msg;
    if(utils.isEqual(code, ErrorType.SUCCESS_0.code)){
      icon = 'success';
      text = '创建成功';
      state.apps.push(app);
      state.app = utils.clone(defaultApp);
    }
    context.proxy.$toast({ icon, text, duration: 4000 });
  });
}

function getApps(){
  Application.getList().then(({ data: { items } }) => {
    let apps = items.map((item) => {
      item.created_time = utils.formatTime(item.created_time);
      item.ended_time = utils.formatTime(item.ended_time);
      item.user_count = utils.numberWithCommas(item.user_count)
      return item;
    });
    state.apps = apps;
  });
}
getApps();
</script>
<template>
  <div class="mb-4">
    <div class="header cim-header">
      <div class="cim-title">应用列表</div>
      <div class="cicon cicon-add cim-button cim-button-bg" @click="onShowEdit(true)" @save="onSave()">创建应用</div>
    </div>
    <table class="table cim-table">
      <thead>
        <tr>
          <th scope="col">应用名称</th>
          <th scope="col">授权个数</th>
          <th scope="col">到期时间</th>
          <th scope="col">创建时间</th>
          <th scope="col">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="app in state.apps">
          <td>{{ app.app_name }}</td>
          <td>{{ app.user_count }}</td>
          <td>{{ app.ended_time }}</td>
          <td>{{ app.created_time }}</td>
          <td>
            <!-- <a class="btn-link cim-btn-link" type="button" @click="onShowEdit(true)">修改</a> -->
            <a class="btn-link cim-btn-link" type="button" @click="">查看</a>
          </td>
        </tr>
      </tbody>
    </table>
    <ModifyDialog :show="state.isShowEdit" :title="'创建应用'" @hide="onShowEdit(false)">
      <div class="row g-2 cim-row">
          <div class="form-floating">
            <input class="form-control" placeholder="应用名称">
            <label>应用名称</label>
          </div>
          <div class="form-floating">
            <select class="form-select">
              <option :value="license.value" v-for="license in state.licenses">{{ license.label }}</option>
            </select>
            <label>授权个数</label>
          </div>
      </div>
    </ModifyDialog>
  </div>
</template>