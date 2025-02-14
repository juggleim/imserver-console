<script setup>
import { reactive } from 'vue';
import ModifyDialog from '../../components/dialog.vue';

let state = reactive({
  users: [
    { name: '张晓曦', time: '2023-10-10 23:03', account: 'zhangxiaoxin@163.com', status: '已启用' },
    { name: '张晓曦', time: '2023-10-10 23:03', account: 'zhangxiaoxin@163.com', status: '已启用' },
    { name: '张晓曦', time: '2023-10-10 23:03', account: 'zhangxiaoxin@163.com', status: '已启用' },
  ],
  radios: [
    { name: 'type', value: 1, label: '启用' },
    { name: 'type', value: 2, label: '禁用' },
  ],
  isOpen: 1,
  isShowEdit: false
});

function onShowEdit(isShow) {
  state.isShowEdit = isShow;
}

</script>
<template>
  <div class="mb-4">
    <div class="header cim-header">
      <div class="cim-title">角色管理</div>
      <div class="cicon cicon-add cim-button cim-button-bg" @click="onShowEdit(true)">添加用户</div>
    </div>
    <table class="table cim-table">
      <thead>
        <tr>
          <th>用户名称</th>
          <th>用户账号</th>
          <th>是否启用</th>
          <th>创建时间</th>
          <th class="cim-td-c">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in state.users">
          <td>{{ user.name }}</td>
          <td>{{ user.account }}</td>
          <td>{{ user.status }}</td>
          <td>{{ user.time }}</td>
          <td class="cim-td-c">
            <a class="btn-link cim-btn-link" type="button" @click="onShowEdit(true)">修改</a>
            <a class="btn-link cim-btn-link" type="button">禁用</a>
          </td>
        </tr>
      </tbody>
    </table>
    <ModifyDialog :show="state.isShowEdit" :title="'添加用户'" @hide="onShowEdit(false)">
      <div class="row g-2 cim-row">
        <div class="form-floating">
          <input class="form-control" placeholder="用户名称">
          <label>用户名称</label>
        </div>
        <div class="form-floating">
          <input class="form-control" placeholder="用户密码">
          <label>用户密码</label>
        </div>
        <div class="form-floating">
          <div class="form-control">
            <div class="form-check form-check-inline" v-for="radio in state.radios">
              <input class="form-check-input" type="radio" name="radio.name" :value="radio.value"
                v-model="state.isOpen" @change="onRadieChanged(radio.value)">
              <label class="form-check-label">{{ radio.label }}</label>
            </div>
          </div>
          <label>是否启用</label>
        </div>
      </div>
    </ModifyDialog>
  </div>
</template>