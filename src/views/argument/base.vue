<script setup>
import { reactive, getCurrentInstance } from 'vue';
import utils from '../../common/utils';
import { Application } from "../../services";
import { useRouter } from "vue-router";
import { ErrorType, STORAGE, APP_TYPE, APP_STATUS } from '../../common/enum';
import Storage from "../../common/storage";

const context = getCurrentInstance();
let router = useRouter();
let { currentRoute: { _rawValue: { params: { app_key } } } } = router;
let currentUser = Storage.get(STORAGE.USER_TOKEN);

let state = reactive({
  appInfo: {
    restricted_fields: {}
  },
  isShowSecret: false,
  currentAppStatus: APP_STATUS.ONLINE,
});

function fetchApp() {
  Application.getOne({ app_key }).then(({ data }) => {
    let { cur_user_count, max_user_count } = data;
    utils.formatProps(data, { count: 'num', time: 'date' })
    data.use_percent = Math.floor(cur_user_count / max_user_count) * 100;
    data.expired_time = data.expired_time == -1 ? '永久有效' : utils.formatTime(data.expired_time);
    data.n_app_secret = '********************';
    utils.extend(state.appInfo, data);
  });
}
fetchApp();

function onShowSecret() {
  let isShowSecret = !state.isShowSecret;
  utils.extend(state, { isShowSecret })
}
function onInput() {
  state.licenseErrorMsg = '';
}
function onHideDialog() {
  utils.extend(state, { isShowModifyEdit: false });
}
function onCopyLicnse(){
  context.proxy.$toast({ icon: 'success', text: '复制成功' });
}

</script>
<template>
  <div class="mb-4 app-base">
    <ul class="nav nav-underline-border ab-underline-border">
          <li class="nav-item"><a class="nav-link active cicon cicon-product">基本信息</a></li>
        </ul>
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">App 名称</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">{{ state.appInfo.app_name }}</div>
          </div>
          <div class="col-sm-4"></div>
        </div>
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">App Key</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">{{ state.appInfo.app_key }}</div>
          </div>
          <div class="col-sm-4"></div>
        </div>
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">App Secret</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext cim-app-secret" :class="{ 'redfont': state.isShowSecret }">
              <span class="cim-secret-text">{{ state.isShowSecret ? state.appInfo.app_secret : state.appInfo.n_app_secret
                }}</span>
              <span class="cicon cicon-hide cim-secret-btn" @click="onShowSecret"></span>
            </div>
          </div>
          <div class="col-sm-4"></div>
        </div>
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">App 到期时间</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">{{ state.appInfo.expired_time }}</div>
          </div>
          <div class="col-sm-4"></div>
        </div>

        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">授权数量</label>
          <div class="col-sm-2">
            {{ state.appInfo.n_max_user_count == -1 ? '无限制' : state.appInfo.n_max_user_count }}
          </div>
          <div class="col-sm-4"></div>
        </div>
      
        <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">App 状态</label>
          <div class="col-sm-4">
            <div class="form-control-plaintext">已上线</div>
          </div>
          <div class="col-sm-4"></div>
        </div>

        <!-- <div class="row cim-ab-row">
          <label class="col-sm-1 col-form-label">部署进度</label>
          <div class="col-sm-4 cim-app-deploy-box">
            <div class="cim-app-deploy-hr"></div>
            <div class="cim-app-deploy-list">
              <div class="cim-app-deply-item cicon cicon-circel-dui">准备安装</div>
              <div class="cim-app-deply-item cicon cicon-circel">正在部署</div>
              <div class="cim-app-deply-item cicon cicon-circel">部署完成</div>
            </div>
          </div>
        </div> -->

  </div>
</template>
