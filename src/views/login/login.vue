<script setup>
import { useRouter } from "vue-router";
import { User } from "../../services";
import { reactive, getCurrentInstance } from 'vue';
import { ErrorType, STORAGE, DEPLOY_TYPE, Errors } from "../../common/enum";
import utils from "../../common/utils";
import Stroage from "../../common/storage";
import menuTools from "../layout/menu-tools";

let context = getCurrentInstance();
const router = useRouter();
let state = reactive({
  account: '',
  password: '',
  pwdErrorMsg: '',
  accountErrorMsg: '',
  loginErrorMsg: '',
  isSending: false,
});

function onLogin(){
  let { account, password } = state;
  if(utils.isEmpty(account)){
    return state.accountErrorMsg = '账号不能为空';
  }
  if(utils.isEmpty(password)){
    return state.pwdErrorMsg = '密码不能为空';
  }
  User.login({ account, password }).then(({ code, data }) => {
    if(utils.isEqual(code, ErrorType.SUCCESS_0.code)){
      data.type = utils.isEqual(data.env, 'private') ? DEPLOY_TYPE.PRIVATE : DEPLOY_TYPE.PUBLIC;
      Stroage.set(STORAGE.USER_TOKEN,  data);
      context.proxy.$toast({
        text: '登录成功',
        icon: 'success'
      });
      menuTools.goHomePage(router);
    }

    let error = Errors.find((item) => { return utils.isEqual(item.code, code)}) || {}
    state.loginErrorMsg = error.msg || '登录失败 ' + code
  });
}

function onInput(name){
  let events = {
    pwd: () => {
      state.pwdErrorMsg = '';
    },
    account: () => {
      state.accountErrorMsg = '';
    }
  };
  state.loginErrorMsg = '';
  events[name]();
}
</script>
<template>
    <div class="bg-body-tertiary min-vh-100 d-flex flex-row align-items-center cim-login-body">
      <div class="cim-login-banner">
        <div class="cim-login-sologin">
          <div class="title">我专属的即时通讯组件</div>
          <div class="grap">
            <p class="content cicon cicon-book">集成简单高效</p>
            <p class="content cicon cicon-product">产品轻量易用</p>
          </div>
          <div class="grap">
            <p class="content cicon cicon-safe">服务可靠稳定</p>
            <p class="content cicon cicon-rokcet">功能灵活多变</p>
          </div>
        </div>
      </div>
      <div class="container">
        <div class="row justify-content-center">
          <div class="col-lg-6">
            <div class="card-group d-block d-md-flex row">
              <div class="card col-md-7 p-4 mb-0 cim-login-card">
           
                <div class="card-body">
                  <h1 class="cim-login-title">登录 IM 后台</h1>
                  <div class="login-input-group mb-3">
                    <input class="form-control" type="text" v-model="state.account" placeholder="账号" @input="onInput('account')">
                    <div class="invalid-feedback feedback" v-if="state.accountErrorMsg">{{ state.accountErrorMsg }}</div>
                  </div>
                  <div class="login-input-group mb-3">
                    <input class="form-control" type="password" v-model="state.password" placeholder="密码"  @input="onInput('pwd')"  @keydown.enter="onLogin()">
                    <div class="invalid-feedback feedback" v-if="state.pwdErrorMsg">{{ state.pwdErrorMsg }}</div>
                  </div>
                  <div class="login-input-group mb-3">
                    <button class="btn btn-primary cim-login-button" type="button" @click="onLogin()">登录</button>
                  </div>
                  <div class="invalid-feedback feedback login-feedback" >{{ state.loginErrorMsg }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
</template>
