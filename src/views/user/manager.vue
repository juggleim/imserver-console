<script setup>
import { reactive, getCurrentInstance } from 'vue';
import ModifyDialog from '../../components/dialog.vue';
import { USER_STATE, ErrorType, USER_ROLE, ROLES } from "../../common/enum";
import utils from '../../common/utils';
import { User } from "../../services";

let context = getCurrentInstance();
let defaltUser = {
  account: '',
  password: '',
  confirmPasswrod: '',
  state: USER_STATE.ENABLE,
  role: USER_ROLE.ADMIN
}
let state = reactive({
  users: [],
  roles: utils.clone(ROLES),
  radios: [
    { name: 'type', value: USER_STATE.ENABLE, label: '启用' },
    { name: 'type', value: USER_STATE.DISABLE, label: '禁用' },
  ],
  isShowEdit: false,
  user: utils.clone(defaltUser),

  accountErrorMsg: '',
  pwdErrorMsg: '',
  conPwdErrorMsg: '',
});

User.getUsers().then(({ data }) => {
  let { items } = data;
  state.users = items.map((item) => {
    item.time = utils.formatTime(item.created_time)
    return item;
  })
});
function onShowEdit(isShow, user) {
  state.isShowEdit = isShow;
  if(user){
    state.user = user;
  }
  if(!isShow){
    state.user = utils.clone(defaltUser);
  }
}
function onRadieChanged(type){
  state.user.state = type;
}
function onOperate(user){
  let userState;
  if(utils.isEqual(user.state, USER_STATE.ENABLE)){
    userState = USER_STATE.DISABLE;
  }
  if(utils.isEqual(user.state, USER_STATE.DISABLE)){
    userState = USER_STATE.ENABLE;
  }
  User.disable({ 
    accounts: [user.account],
    is_disable: userState
  }).then(() => {
    state.users.map((_user) => {
      if(utils.isEqual(user.account, _user.account)){
        utils.extend(_user, { state: userState });
      }
      return _user;
    });
    context.proxy.$toast({
      icon: 'success',
      text: '操作成功',
      duration: 4000
    })
  })
  
}
function onDelete(index){
  let user = state.users[index];
  User.remove({ accounts: [user.account] }).then(() => {
    state.users.splice(index, 1);
    context.proxy.$toast({
      icon: 'success',
      text: '删除成功',
      duration: 4000
    });
  });
}
function onSave(){
  let { user } = state;
  let { account, password, confirmPasswrod, role } = user;
  if(utils.isEmpty(account)){
    return state.accountErrorMsg = '账户名称不能为空';
  }
  if(utils.isEmpty(password)){
    return state.pwdErrorMsg = '密码不能为空';
  }
  if(!utils.isEqual(confirmPasswrod, password)){
    return state.conPwdErrorMsg = '两次密码输入不一致';
  }
  if(user.created_time){
    return User.updatePwd({
      account: account,
      password: password,
      new_password: confirmPasswrod,
      role_id: role,
    }).then(({ code, msg }) => {
      if(utils.isEqual(code, ErrorType.SUCCESS_0.code)){
        context.proxy.$toast({
          icon: 'success',
          text: '修改密码成功',
          duration: 3000
        });
      }else{
        context.proxy.$toast({
          icon: 'error',
          text: `${code}:${msg}`,
          duration: 3000
        });
      }
    });
  }
  User.add({ account, password, state: user.state, role_id: role }).then(({ code }) => {
    let icon = 'success', text = '保存成功';
    if(utils.isEqual(code, ErrorType.SUCCESS_0.code)){
      user.time = utils.formatTime(Date.now());
      state.users.push(user);
      onShowEdit(false);
    }else if(utils.isEqual(code, ErrorType.USER_EXISTS.code)){
      icon = 'error';
      text = ErrorType.USER_EXISTS.msg;
    }else{
      icon = 'error';
      text = '保存失败，请重试';
    }
    context.proxy.$toast({
      icon,
      text,
      duration: 4000
    });
  });
}
function onInput(name){
  let events = {
    account: () => {
      state.accountErrorMsg = '';
    },
    pwd: () => {
      state.pwdErrorMsg = '';
    },
    confirm: () => {
      state.conPwdErrorMsg = '';
    }
  };
  events[name]();
}
</script>
<template>
  <div class="mb-4">
    <div class="header cim-header">
      <div class="cim-title">用户管理</div>
      <div class="cicon cicon-add cim-button cim-button-bg" @click="onShowEdit(true)">添加用户</div>
    </div>
    <table class="table cim-table">
      <thead>
        <tr>
          <th>用户名称</th>
          <th>用户账号</th>
          <th>用户密码</th>
          <th>是否启用</th>
          <th>创建时间</th>
          <th class="cim-td-c">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(user, index) in state.users">
          <td>{{ user.account }}</td>
          <td>{{ user.account }}</td>
          <td>**** ****</td>
          <td class="cicon" :class="{ 'cicon-success': user.state == USER_STATE.ENABLE, 'cicon-error': user.state == USER_STATE.DISABLE }">
            {{ user.state == USER_STATE.ENABLE ? '已启用' : '已禁用' }}
          </td>
          <td>{{ user.time }}</td>
          <td class="cim-td-c cim-td-operate">
            <!-- <a class="btn-link cim-btn-link" type="button" @click="onShowEdit(true, user)">修改</a> -->
            <a class="btn-link cim-btn-link" type="button" @click="onOperate(user)">{{ user.state == USER_STATE.DISABLE ? '启用' : '禁用'}}</a>
            <a class="btn-link cim-btn-link" type="button" @click="onDelete(index)">删除</a>
          </td>
        </tr>
      </tbody>
    </table>
    <ModifyDialog :show="state.isShowEdit" :title="'用户变更'" @hide="onShowEdit(false)" @save="onSave()">
      <div class="row g-2 cim-row">
        <div class="form-floating">
          <input class="form-control" :disabled="state.user.created_time" v-model="state.user.account" type="text" placeholder="用户名称" autocomplete="off"  @input="onInput('account')">
          <label>用户名称</label>
          <div class="invalid-feedback feedback" v-if="state.accountErrorMsg">{{ state.accountErrorMsg }}</div>
        </div>
        <div class="form-floating">
          <input class="form-control" v-model="state.user.password" type="text" placeholder="用户密码" @input="onInput('pwd')">
          <label>用户密码</label>
          <div class="invalid-feedback feedback" v-if="state.pwdErrorMsg">{{ state.pwdErrorMsg }}</div>
        </div>
        <div class="form-floating">
          <input class="form-control" v-model="state.user.confirmPasswrod" type="text" placeholder="确认密码" @input="onInput('confirm')">
          <label>确认密码</label>
          <div class="invalid-feedback feedback" v-if="state.conPwdErrorMsg">{{ state.conPwdErrorMsg }}</div>
        </div>
        <div class="form-floating">
          <select class="form-select" v-model="state.user.role">
            <option :value="item.value" v-for="item in state.roles" >{{ item.name }}</option>
          </select>
          <label>用户角色</label>
        </div>
        <!-- <div class="form-floating">
          <div class="form-control">
            <div class="form-check form-check-inline" v-for="radio in state.radios">
              <input class="form-check-input" type="radio" name="radio.name" :value="radio.value"
                v-model="state.user.state" @change="onRadieChanged(radio.value)">
              <label class="form-check-label">{{ radio.label }}</label>
            </div>
          </div>
          <label>是否启用</label>
        </div> -->
        <input type="password" autocomplete="new-password" style="display: none;" />
      </div>
    </ModifyDialog>
  </div>
</template>