<script setup>
import Sidler from './sidler.vue';
import Header from './header.vue';
import { useRouter } from 'vue-router'
import { reactive, watch } from 'vue';
import utils from '../../common/utils';
import menuTools from "./menu-tools";

let state = reactive({
  isCollapse: false,
  isHome: false
});

const router = useRouter()
let { currentRoute: { _rawValue: { fullPath, name } } } = router;
if ( fullPath == '/') {
  menuTools.goHomePage(router);
}
state.isHome = utils.isEqual(name, 'Dashboard');

function onCollapse(isCollapse){
  state.isCollapse = isCollapse;
}
function onMask(){
  onCollapse(false);
}
let useRouterCurrent = reactive(router);
watch(useRouterCurrent, (n) =>{
  let { currentRoute: { name } } = n;
  state.isHome = utils.isEqual(name, 'Dashboard');
});

</script>

<template>
  <div>
    <Header @collapse="onCollapse"></Header>
    <div class="wrapper d-flex min-vh-85 px-3">
      <Sidler :is-collapse="state.isCollapse" :is-home="state.isHome"/>  
      <RouterView v-slot="{ Component, route }" class="cim-container min-vh-85">
        <component :is="Component" :key="route.fullPath"/>
      </RouterView>
      <div :class="{ 'sidebar-backdrop fade show': state.isCollapse && utils.isMobile()}" @click="onMask"></div>
    </div>
  </div>
</template>
