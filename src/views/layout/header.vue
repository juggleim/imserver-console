<script setup>
import {useRouter} from "vue-router";
import { getCurrentInstance, reactive, watch, nextTick } from 'vue';
import utils from "../../common/utils";
import Storage from "../../common/storage";
import {APP_TYPE, DEPLOY_TYPE, ErrorType, EVENT_NAME, STORAGE} from "../../common/enum";
import menuTools from '../layout/menu-tools';
import appTools from '../../common/app-tools';
import emitter from '../../common/emmit';
import CreateDialog from '../../components/dialog-app.vue';

const context = getCurrentInstance();
const router = useRouter();

let currentUser = Storage.get(STORAGE.USER_TOKEN);
let defaultApp = {app_name: '选择应用', kine: ''};
let defaultNewApp = {name: '', licenseKey: ''}
let state = reactive({
  isPublic: utils.isEqual(currentUser.type, DEPLOY_TYPE.PUBLIC),
  isCollapse: false,
  isShowSetting: false,
  isShowDropdwonMenu: false,
  isShowDropdwonMenuBackdrop: false,
  breadcumbs: [],
  showHeaderApps: false,
  apps: [],
  currentApp: utils.clone(defaultApp),

  newApp: utils.clone(defaultNewApp),
  isShowNewEdit: false,
  newAppNameErrorMsg: '',
  licenses: [
    {value: 1, label: '0-200'},
    {value: 2, label: '200-500'},
    {value: 3, label: '500-1000'},
  ],
});
const props = defineProps(['path']);
const emit = defineEmits(['collapse'])

function onClick() {
  state.isCollapse = !state.isCollapse;
  if (utils.isMobile()) {
    state.isCollapse = true;
  }
  emit('collapse', state.isCollapse);
}

function onSetting() {
  state.isShowSetting = !state.isShowSetting;
}

let onNavigate = (name) => {
  if (utils.isEqual(name, 'dropdown')) {
    return state.isShowDropdwonMenu = !state.isShowDropdwonMenu;
  }
  if (utils.isEqual(name, 'Login')) {
    Storage.remove(STORAGE.USER_TOKEN);
  }
  appTools.setCurrent(utils.clone(defaultApp))
  router.push({name});
}

let useRouterCurrent = reactive(router);

watch(useRouterCurrent, (n) => {
  let {currentRoute: {meta: {titles}, name}} = n;
  let current = appTools.getCurrent();
  if (utils.isEqual(name, 'Dashboard')) {
    current = utils.clone(defaultApp);
  }
  utils.extend(state, {
    breadcumbs: titles,
    currentApp: current
  });
});

let appsOffset = '';
let appHasMore = true;
appTools.fetch({ offset: appsOffset }, ({ apps, offset }) => {
  state.apps = apps;
  appsOffset = offset;
  let {currentRoute: {_rawValue: {params: {app_key}}}} = router;
  let app = apps.filter((app) => {
    return utils.isEqual(app.app_key, app_key);
  })[0] || {};
  let current = utils.isEmpty(app) ? utils.clone(defaultApp) : app;
  let isFirst = true;
  onDropclick(current, isFirst);
})

function onDropclick(app, isFirst) {
  if (!isFirst) {
    menuTools.goBasePage(app);
  }
  appTools.setCurrent(utils.clone(app));
  utils.extend(state.currentApp, appTools.getCurrent())
}

function onHomePage() {
  menuTools.goHomePage(router);
}

function onHideDropmenu() {
  utils.extend(state, {
    isShowDropdwonMenu: false,
    isShowDropdwonMenuBackdrop: false,
    isShowSetting: false
  });
}

function onShowbackdrop() {
  utils.extend(state, {
    isShowDropdwonMenuBackdrop: true
  })
}

function onShowNewEdit(isShow) {
  utils.extend(state, { isShowNewEdit: isShow, isShowDropdwonMenu: false });
}

emitter.$on(EVENT_NAME.CREATE_APP_FINISHED, (app) => {
  let index = utils.find(state.apps, (_app) => {
    return utils.isEqual(_app.app_key, app.app_key);
  });
  if(index == -1){
    state.apps.unshift(app);
  }
  appTools.setCurrent(utils.clone(app));
  utils.extend(state.currentApp, appTools.getCurrent())
  menuTools.goBasePage(app);
});

emitter.$on(EVENT_NAME.CREATE_APP, () => {
  onShowNewEdit(true);
});
emitter.$on(EVENT_NAME.ON_LOGOUT, () => {
  router.push({ name: 'Login' });
});
let canscroll = true;
nextTick(() => {
  let { appList } = context.refs;
  appList.addEventListener("scroll", () => {
    let scrollTop = appList.scrollTop;
    let scrollHeight = appList.scrollHeight;
    let rectHeight = appList.getBoundingClientRect().height;
    let isNeedLoad = scrollHeight - scrollTop - rectHeight < 100;
    if (isNeedLoad && canscroll) {
      if(!appHasMore){
        return;
      }
      canscroll = false;
      appTools.fetch({ offset: appsOffset }, ({ apps, offset, has_more }) => {
        appsOffset = offset;
        appHasMore = has_more;
        utils.forEach(apps, (app) => {
          state.apps.push(app);
        });
        canscroll = true;
      })
    }
  });
});

</script>
<template>
  <header class="header header-sticky p-0">
    <div class="container-fluid border-bottom px-4 cim-header-container">
      <div class="sidebar-header" @click="onHomePage()">
        <div class="sidebar-brand">
          <div class="sildebar-logo">IM 管理后台</div>
        </div>
      </div>
      <ul class="header-nav d-lg-flex">
        <li class="nav-item cim-nav-item">
          <button class="header-toggler cicon cicon-fold cim-header-tg" type="button" @click="onClick"></button>
        </li>
        <li class="nav-item cim-nav-item cim-nav-item-apps cicon cicon-down">
          <div class="dropdown cim-dropdown">
            <div class="row cim-hrow" @click="onNavigate('dropdown')">
              <div class="col-md-6 cim-hdpm-appname">{{ state.currentApp.app_name }}</div>
              <div class="col-md-6 cim-hdpm-appname">{{ state.currentApp.kind }}</div>
            </div>
            <ul class="dropdown-menu cim-dropdown-menu cim-dropdown-menu-app cim-header-app-dropdown" :class="{ 'show': state.isShowDropdwonMenu }"
              @mouseleave="onShowbackdrop()" ref="appList">
              
              <li class="cim-dropdown-item">
                <div class="dropdown-item">
                  <div class="dropdown-item-app">
                    <div class="row cim-hrow">
                      <div class="cicon cicon-add cim-button cim-button-bg" @click="onShowNewEdit(true)">
                        {{ currentUser.is_commercial ? '导入应用' : '创建应用' }}
                      </div>
                    </div>
                  </div>
                </div>
              </li>

              <li class="cim-dropdown-item" v-for="app in state.apps" @click.stop="onDropclick(app)">
                <div class="cim-header-nav-app">
                  <div class="cim-header-nav-appname">{{ app.app_name }}</div>
                  <div class="cim-header-nav-apptype cim-app-public" :class="{'cim-app-private': utils.isEqual(app.app_type, APP_TYPE.PRIVATE)}">{{ app.kind }}</div>
                </div>
              </li>

            </ul>
          </div>
        </li>
      </ul>
      <ul class="header-nav ms-auto">
        <!-- <li class="nav-item cim-nav-item cim-r-nav">
          <div class="nav-link cim-nav-link cicon cicon-finance" @click="onNavigate('FinaDash')">财务管理</div>
        </li>
        <li class="nav-item cim-nav-item cim-r-nav">
          <div class="nav-link cim-nav-link cicon cicon-book" @click="onNavigate('OrderList')">技术支持</div>
        </li> -->
      </ul>
      <ul class="header-nav cim-user-header-nav">
        <li class="nav-item py-1">
          <div class="vr h-100 mx-2 text-body"></div>
        </li>
        <li class="nav-item dropdown cim-dropdown">
          <a class="nav-link py-0 pe-0" href="#" @click="onSetting">
            <div class="avatar cim-avatar cim-header-avatar">
              <img class="avatar-img" src="../../assets/images/res/avatar.jpg">
            </div>
          </a>
          <div class="dropdown-menu cim-dropdown-menu" :class="{ 'show': state.isShowSetting }"
            @mouseleave="onShowbackdrop()">
            <div class="dropdown-item cim-dropdown-item cicon cicon-setting" @click="onNavigate('UserSetting')">账户设置
            </div>
            <div class="dropdown-item cim-dropdown-item cicon cicon-logout" @click="onNavigate('Login')">退出登录</div>
          </div>
        </li>
      </ul>
    </div>
  </header>
  <div class="modal-backdrop cim-o0-modal-backdrop" v-if="state.isShowDropdwonMenuBackdrop" @click="onHideDropmenu()"></div>
  <CreateDialog :show="state.isShowNewEdit" @hide="onShowNewEdit(false)"></CreateDialog>
</template>
