<script setup>
import { reactive, getCurrentInstance } from 'vue';
import ModifyDialog from '../../components/dialog.vue';
import { STORAGE, ErrorType } from "../../common/enum";
import Storage from "../../common/storage";
import utils from '../../common/utils';
import { User } from "../../services";

let context = getCurrentInstance();
let state = reactive({
  isShowEdit: false,
  user: Storage.get(STORAGE.USER_TOKEN),
  pwd: '',
  newPwd: '',
  pwdErrorMsg: '',
  newPwdErrorMsg: ''
});

function onShowEdit(isShow){
  state.isShowEdit = isShow;
}
function onInput(name){
  let events = {
    pwd: () => {
      state.pwdErrorMsg = '';
    },
    newPwd: () => {
      state.newPwdErrorMsg = '';
    }
  };
  events[name]();
}
function onSave(){
  let { newPwd, pwd, user } = state;
  if(utils.isEmpty(pwd)){
    return state.pwdErrorMsg = '原密码不能为空';
  }
  if(utils.isEmpty(newPwd)){
    return state.newPwdErrorMsg = '新密码不能为空';
  }
  User.updatePwd({
    account: user.account,
    password: pwd,
    new_password: newPwd
  }).then(({ code, msg }) => {
    if(utils.isEqual(code, ErrorType.SUCCESS_0.code)){
      onShowEdit(false);
      context.proxy.$toast({
        text: '修改成功',
        icon: 'success'
      });
    }else if(utils.isEqual(code, ErrorType.USER_OLDPWD_WRONG.code)){
      state.newPwdErrorMsg = ErrorType.USER_OLDPWD_WRONG.msg;
    }else{
      context.proxy.$toast({
        text: `${code}: ${msg}`,
        icon: 'error'
      });
    }
  });
}
</script>
<template>
  <div class="md-4 app-base">
    <ul class="nav nav-underline-border ab-underline-border">
      <li class="nav-item"><a class="nav-link active cicon cicon-product">基本信息</a></li>
    </ul>
    <div class="cim-us-table">
      <div class="row cim-ab-row cim-us-row">
        <label class="col-sm-1 col-form-label">账户名称</label>
        <div class="col-sm-2">
          <input class="form-control-plaintext" type="text" readonly :value="state.user.account">
        </div>
        <div class="col-sm-4">
          <!-- <a class="btn-link cim-btn-link" type="button" @click="onShowEdit(true)">修改</a> -->
        </div>
      </div>

      <div class="row cim-ab-row cim-us-row">
        <label class="col-sm-1 col-form-label">账户密码</label>
        <div class="col-sm-2">
          <input class="form-control-plaintext" type="text" readonly value="**** ****">
        </div>
        <div class="col-sm-4">
          <a class="btn-link cim-btn-link" type="button" @click="onShowEdit(true)">修改</a>
        </div>
      </div>

      <!-- <div class="row cim-ab-row cim-us-row">
        <label class="col-sm-1 col-form-label">注册时间</label>
        <div class="col-sm-2">
          <input class="form-control-plaintext" type="text" readonly value="2024-01-01 18:09:20">
        </div>
        <div class="col-sm-4">
        </div>
      </div> -->
    </div>

    <ModifyDialog :show="state.isShowEdit" :title="'修改密码'" @hide="onShowEdit(false)" @save="onSave()">
      <div class="row g-2 cim-row">
          <div class="form-floating">
            <input class="form-control" placeholder="账户名称" :value="state.user.account" disabled>
            <label>账户名称</label>
          </div>
          <div class="form-floating">
            <input class="form-control" placeholder="原密码" v-model="state.pwd" @input="onInput('pwd')">
            <label>原密码</label>
            <div class="invalid-feedback feedback" v-if="state.pwdErrorMsg">{{ state.pwdErrorMsg }}</div>
          </div>
          <div class="form-floating">
            <input class="form-control" placeholder="新密码" v-model="state.newPwd" @input="onInput('newPwd')">
            <label>新密码</label>
            <div class="invalid-feedback feedback" v-if="state.newPwdErrorMsg">{{ state.newPwdErrorMsg }}</div>
          </div>
      </div>
    </ModifyDialog>
  </div>
</template>
