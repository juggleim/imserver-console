<script setup>
import { reactive, getCurrentInstance } from 'vue';
import { useRouter } from "vue-router";
import utils from '../../common/utils';
import { R3d } from "../../services";
import { RESPONSE, RTC_CHANNELS, R3D_USE_TYPE } from '../../common/enum';
import TranslateDialog from "../../components/dialog-3rd.vue";

let router = useRouter();
let { currentRoute: { _rawValue: { params: { app_key } } } } = router;
const context = getCurrentInstance();

let channels = utils.clone(RTC_CHANNELS);
let globalData = {};
let state = reactive({
  list: [],
  current: channels[0],
  isShowDialog: false,
});

function onShowDialog(isShow, item){
  state.isShowDialog = isShow;
  if(isShow){
    state.current = item;
  }
}

function onSave(params){
  let { uid, fields } = params;
  let _params = { };
  let conf = {};
  utils.forEach(fields, (field) => {
    let { key, value } = field;
    conf[key] = value;
  });
  _params[uid] = conf;
  if(utils.isEqual(R3D_USE_TYPE.ENABLE, conf.is_used)){
    _params.channel = uid;
    utils.forEach(state.list, (_item) => {
      _item.isUsed = utils.isEqual(_item.uid, uid);
    });
  }
  delete _params[uid].is_used;
  utils.extend(globalData, _params);
  create(globalData, _params);
}
function create(conf, _params){
  R3d.setRTC({ app_key, conf }).then((result) => {
    let { code, data } = result;
    if(!utils.isEqual(code, RESPONSE.SUCCESS)){
      return context.proxy.$toast({ icon: 'error', text: `添加失败: ${code}` });
    }
    context.proxy.$toast({ icon: 'success', text: `操作成功` });
    updateList(_params);
    onShowDialog(false);
  });
}

function updateConf(conf){
  R3d.setRTC({ app_key, conf }).then((result) => {
    let { code, data } = result;
    if(!utils.isEqual(code, RESPONSE.SUCCESS)){
      return context.proxy.$toast({ icon: 'error', text: `修改失败: ${code}` });
    }
    context.proxy.$toast({ icon: 'success', text: `操作成功` });
    onShowDialog(false);
  });
}

R3d.getRTC({ app_key }).then((result) => {
  let { code, data } = result;
  if(!utils.isEqual(code, RESPONSE.SUCCESS)){
    return context.proxy.$toast({ icon: 'error', text: `获取配置失败: ${code}` });
  }
  globalData = utils.clone(data);
  let isAdd = true;
  updateList(data, isAdd);
});

function updateList(data, isAdd){
  let { channel: currentChannel } = data;
  utils.forEach(channels, (channel) => {
    let isUsed = utils.isEqual(channel.uid, currentChannel);
    let { uid, children } = channel;
    let _channel = data[uid] || {};

    if(utils.isEmpty(_channel) && !isAdd){
      return;
    }

    channel.children = utils.map(children, (child) => {
      let { key } = child;
      child.value = _channel[key] || '';
      return child;
    });
    if(isUsed){
      state.current = channel;
    }
    utils.extend(channel, { isUsed });
    if(isAdd){
      state.list.push(channel)
    }
  });
}

function onDisable(item, index){
  item.isUsed = false;
  utils.extend(globalData, { channel: '' });
  updateConf(globalData);
}

function onEnable(item, index){
  utils.forEach(state.list, (_item) => {
    _item.isUsed = false;
  });
  item.isUsed = true;
  utils.extend(globalData, { channel: item.uid });
  updateConf(globalData);
}

</script>
<template>
  <div class="mb-4 app-base cim-r3-container">
    <ul class="cim-r3-body">
      <li class="cim-r3-item" v-for="(item, index) in state.list">
        <div class="cim-r3-item-header">
          <div class="cim-r3-header-info">
            <div class="cim-r3-avatar" :class="[item.icon + '-avatar']"></div>
            <div class="cim-r3-item-name">
              {{ item.name }}
            </div>
          </div>
          <div class="cim-r3-header-status">
            <span class="life-status" v-if="item.children[0].value">【已启用】</span>
          </div>
        </div>
        <ul class="cim-r3-item-contents">
          <li class="cim-rtc-item-content" v-for="child in item.children">
            <div class="title" v-if="child.type == 'text' || child.type == 'number'">{{ child.name }}:</div>
            <div class="value"  v-if="(child.type == 'text' || child.type == 'number') && child.value">{{ child.secretValue || child.value }}</div>
            <div class="value unset"  v-if="(child.type == 'text' || child.type == 'number') && !child.value">未设置</div>
          </li>
        </ul>
        <ul class="cim-r3-item-tools">
          <li class="cim-r3-item-tool cicon wr-update" @click="onShowDialog(true, item)">设置</li>
          <!-- <li class="cim-r3-item-tool cicon wr-disable" v-if="item.isUsed" @click="onDisable(item, index)">停用</li>
          <li class="cim-r3-item-tool cicon wr-enable" v-if="!item.isUsed && item.children[0].value"  @click="onEnable(item, index)">启用</li> -->
        </ul>
      </li>
    </ul>
    <TranslateDialog :show="state.isShowDialog" :title="'音视频设置'" :custom="'音视频通道'" :channel="state.current" :list="state.list" @save="onSave" @hide="onShowDialog(false)"></TranslateDialog>
  </div>
</template>