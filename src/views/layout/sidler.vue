<script setup>
import { useRouter } from "vue-router";
import { reactive, watch } from 'vue';
import utils from "../../common/utils";
import Stroage from "../../common/storage";
import { STORAGE } from "../../common/enum";
import menuTools from "./menu-tools";
import appTools from '../../common/app-tools';

const props = defineProps(["isCollapse", "isHome"]);
const router = useRouter();
let useRouterCurrent = reactive(router);

let state = reactive({
  menus: []
});
menuTools.setMenuConfig({ state, router });

function updateActive(active){
  state.menus.map((_menu) => {
    // _menu.isUnfold = utils.isInclude(active.name, _menu.tag);
    if(_menu.children){
      _menu.children = _menu.children.map((child) => {
        child.isActive = utils.isEqual(active.name, child.name);
        return child;
      });
    }
    return _menu;
  });
}
let onFold = (menu) => {
  state.menus.map((_menu) => {
    if(utils.isEqual(menu.title, _menu.title)){
      _menu.isUnfold = !_menu.isUnfold;
    }
    return _menu;
  });
};

let onNavigate = (menu) => {
  if(menu.children){
    return onFold(menu);
  }
  let currentApp = appTools.getCurrent();
  router.push({ name: menu.name, params: { app_key: currentApp.app_key } });
}

watch(useRouterCurrent, (n) =>{
  let { currentRoute: { name } } = n;
  if(!menuTools.isSameGroup(state.menus, name)){
    menuTools.showMenus(router);
  }
  if(utils.isEqual(name, 'Dashboard')){
    state.menus = [];
  }
  updateActive({ name })
});

function onLogout(){
  Stroage.remove(STORAGE.USER_TOKEN);
  router.push({ name: 'Login' });
}
let { currentRoute: { _rawValue: { name } } } = router;

//非首页，根据 URL 展示菜单
if(!utils.isEqual(name, 'Dashboard')){
  menuTools.showMenus(router);
}
// 展示菜单后更新选中状态，和展示菜单先后顺序不能改变
updateActive({ name });
</script>
<template>
  <div class="sidebar" :class="{ 'hide': props.isCollapse && utils.isMobile(), 'show':   props.isCollapse && utils.isMobile(), 'hide': props.isHome}" id="sidebar">
 
    <ul class="sidebar-nav" data-coreui="navigation" data-simplebar>
      <li class="nav-group" v-for="menu in state.menus" :class="{ 'show': menu.isUnfold, 'hideMenu': menu.isHidden }">
        <div class="nav-link cicon" :class="{ 'nav-group-toggle': menu.children, 'active': menu.isActive }" href="#" @click="onNavigate(menu)">
          <span class="cim-nav-icon cicon" :class="[menu.icon]"></span>
          {{ menu.title }}
        </div>
        <ul class="nav-group-items compact">
          <li class="nav-item" v-for="child in menu.children" :class="{'hideMenu': child.isHidden }">
            <div class="nav-link cicon" :class="{ 'active': child.isActive}" href="#" @click="onNavigate(child)">
              <span class="nav-icon">
                <!-- <span class="nav-icon-bullet"></span> -->
              </span> {{ child.title }}
            </div>
          </li>
        </ul>
      </li>
    </ul>
    <div class="sidebar-footer">     
      <div class="cicon cicon-logout cim-button cim-sider-logout" @click="onLogout">退出登录</div>
    </div>
  </div>
</template>
