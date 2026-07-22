<script setup>
import { reactive, getCurrentInstance, nextTick } from 'vue';
import { STORAGE, ErrorType, RESPONSE, SIGNAL_TYPE, METHOD_MAP, IM_ERRORS } from "../../common/enum";
import Storage from "../../common/storage";
import utils from '../../common/utils';
import { Inspect } from "../../services";
import ConnManger from "./conns/mange.vue";
import ConnSignal from "./conns/signal.vue";
import { useRouter } from "vue-router";
import { t } from '@/i18n';
import PageSection from '@/components/page-section.vue';

let router = useRouter();
let { currentRoute: { _rawValue: { params: { app_key } } } } = router;

let context = getCurrentInstance();
let state = reactive({
  isShowEdit: false,
  current: { session: 'manger' },
  tabs: [
    { connTimeName: t('tools.connection.tab.manager'), session: 'manger', isActive: true, isClose: false, content: '123' },
  ]
});
function onChanged(item, index){
  utils.map(state.tabs, (tab) => {
    tab.isActive = utils.isEqual(item.session, tab.session);
    return tab;
  });
  state.current = item;
}
function onClose(index){
  state.tabs.splice(index, 1);
  let i = index - 1;
  let tab = state.tabs[i];
  onChanged(tab, i);
}
function onCreateCon(item){
  let tab = { ...item };
  let index = utils.find(state.tabs, (_tab) => {
    return utils.isEqual(item.session, _tab.session);
  });
  if(index > -1){
    tab = state.tabs[index];
    return onChanged(tab, index);
  }
  state.tabs.push(tab);
  index = state.tabs.length - 1;
  onChanged(tab, index);
  nextTick(() => {
    let { navs } = context.refs;
    if (navs) {
      navs.scrollLeft = navs.scrollWidth
    }
  });
  getSignals({ index, item: tab, start: 0 }, ({ index: i, list }) => {
    state.tabs[i].list = tab.list.concat(list)
  });
}

function onNext({ item }){
  let index = utils.find(state.tabs, (_tab) => {
    return utils.isEqual(item.session, _tab.session);
  });
  let tab = state.tabs[index];
  tab.isLoading = true;
  let itemIndex = tab.list.length - 1;
  let _item = tab.list[itemIndex];
  getSignals({ index, item: tab, start: _item.timestamp }, ({ index, list }) => {
    tab.list = tab.list.concat(list);
  });
}

function getSignals(params, callback){
  let { item, index, start } = params;

  let { session, user_id, count } = item;
  Inspect.getConn({ start, session, user_id, count, app_key }).then((result) => {
    let { code, msg, data } = result;
    item.isLoading = false;
    if(utils.isEqual(code, RESPONSE.SUCCESS)){
      let { logs } = data;
      item.currentCount = logs.length;
      let _list = formatLogs(logs);
      callback({ list: _list, index });
    }else{
      toastError(code, msg);
    }
  })
}

// a signal carries a seq index; the server-side handling of that same index is
// written to the business log, which is what the drill-down fetches
function onDetail({ item, conn }){
  if(utils.isUndefined(item.seq_index) || item.isDetailLoading){
    return;
  }
  if(item.details){
    item.isDetailOpen = !item.isDetailOpen;
    return;
  }
  item.isDetailLoading = true;
  let start = getLogTime(item);
  Inspect.getSignalLogs({
    app_key,
    session: conn.session,
    user_id: conn.user_id,
    index: item.seq_index,
    // the node only scans the log file covering `start`, and it keeps lines
    // strictly newer than it, so back off one millisecond from the signal
    start: start > 0 ? start - 1 : 0,
    count: 50,
  }).then((result) => {
    let { code, msg, data } = result;
    item.isDetailLoading = false;
    if(utils.isEqual(code, RESPONSE.SUCCESS)){
      item.details = formatDetails(data.logs);
      item.isDetailOpen = true;
    }else{
      toastError(code, msg);
    }
  })
}

function toastError(code, msg){
  let text = utils.isEqual(code, RESPONSE.REQUEST_LIMIT)
    ? t('tools.connection.feedback.rateLimited')
    : t('tools.connection.feedback.requestFailed', { code, msg }, `Error: ${code} ${msg}`);
  context.proxy.$toast({ icon: 'error', text });
}
let actionMap = {
  connect: { type: SIGNAL_TYPE.CONNECTED, method: 'connect' },
  qry: { type: SIGNAL_TYPE.USER },
  qry_ack: { type: SIGNAL_TYPE.REPLY, method: 'qry_ack' },
  u_pub: { type: SIGNAL_TYPE.USER},
  u_pub_ack: { type: SIGNAL_TYPE.REPLY, method: 'u_pub_ack'},
  s_pub: { type: SIGNAL_TYPE.SERVER },
  s_pub_ack: { type: SIGNAL_TYPE.REPLY, method: 's_pub_ack' },
  disconnect: { type: SIGNAL_TYPE.DISCONNECTED, method: 'disconnect' },
};
function formatLogs(logs){
  let _logs = utils.map(logs, (log) => {
    let { action } = log;
    // 移除内置属性
    let _log = utils.clone(log);
    let infos = removeAttrs(_log);

    let actionItem = actionMap[action] || {};
    log = utils.extend(log, { ...actionItem, infos });

    let { method } = log;
    let methodItem = METHOD_MAP[method] || {};
    log = utils.extend(log, methodItem);

    log.timeName = format(getLogTime(log), 'hh:mm:ss.S');
    return log;
  });
  return _logs;
}

// business logs describe how one signal was handled: the fixed fields carry the
// service and its cost, everything else the node collected lands in `parms`
function formatDetails(logs){
  return utils.map(logs || [], (log) => {
    let { service_name, method, expend, parms } = log;
    let infos = [];
    utils.forEach(parms || {}, (v, k) => {
      infos.push({ name: k, value: `: ${v}` });
    });
    return {
      service_name,
      method,
      expend,
      infos,
      timeName: format(getLogTime(log), 'hh:mm:ss.S'),
    };
  });
}

// vlog lines are stamped by the node with `timestamp`; older kv-backed logs
// carried `real_time`
function getLogTime(log){
  return Number(log.real_time || log.timestamp) || 0;
}

function removeAttrs(log){
  let { code } = log;

  let error = { name: t('tools.connection.status.success'), value: '', cls: 'success' }
  // log fields arrive as strings, error codes are compared as numbers, and a
  // logged code of 0 is the server confirming the signal succeeded
  let logCode = Number(code);
  if(!utils.isUndefined(code) && logCode !== 0){
    let index = utils.find(IM_ERRORS, (item) => {
      return utils.isEqual(item.code, logCode);
    });
    let errorItem = IM_ERRORS[index] || { code: logCode, msg: '' }
    error = { name: t('tools.connection.status.failed'), value: `: ${errorItem.code} ${errorItem.msg}`, cls: 'warn' }
  }

  let attrs = ['action', 'app_key', 'appkey', 'method', 'session', 'timestamp', 'code', 'real_time', 'service_name'];
  utils.forEach(attrs, (key) => {
    delete log[key];
  });
  let infos = [];
  utils.forEach(log, (v, k) => {
    infos.push({ name: k, value: `: ${v}`, order: k.charCodeAt(0) });
  });

  infos = utils.sort(infos, (a, b) => {
    return a.order > b.order;
  })
  infos.unshift(error);
  return infos;
}
function format(date, fmt = 'yyyy-MM-dd hh:mm') {
  return utils.formatTime(new Date(date).getTime(), fmt);
}
</script>
<template>
  <PageSection title-key="menu.dev.connectionInspect" body-class="cim-tcon-container">
    <ul class="cim-tcon-headers nav nav-tabs" ref="navs">
      <li class="cim-tcon-nav-item nav-item fadeinx" v-for="(tab, index) in state.tabs" :class="[!tab.isClose ? 'cim-tcon-nav-item-ab' : '']">
        <span v-if="tab.isClose" class="nav-close cicon cicon-close-c" @click="onClose(index)"></span>
        <div v-if="tab.isClose" class="nav-link cim-tcon-nav-item-link" :class="[tab.isActive ? 'active' : '']" href="#" @click="onChanged(tab, index)">
          <span class="cim-tconn-title-userid">{{ tab.user_id }}</span>
          <span class="cim-tconn-title-time">{{ tab.connTimeName }}</span>
        </div>
        <div v-else class="nav-link" :class="[tab.isActive ? 'active' : '']" href="#" @click="onChanged(tab, index)">
          {{ tab.connTimeName }}
        </div>
      </li>
    </ul>

    <ul class="cim-tcon-contents">
      <li class="cim-tcon-content" v-for="(tab, index) in state.tabs" :class="[tab.session == state.current.session ? 'display-flex' : 'display-none']">
        <ConnManger v-if="tab.session == 'manger'" @create="onCreateCon"></ConnManger>
        <ConnSignal v-else :conn="tab" @next="onNext" @detail="onDetail"></ConnSignal>
      </li>
    </ul>
  </PageSection>
</template>
