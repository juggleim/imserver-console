<script setup>
import { reactive, getCurrentInstance } from 'vue';
import utils from '../../common/utils';
import PullLogDialog from '../../components/dialog-pull-log.vue';
import { useRouter } from "vue-router";
import { Log } from "../../services";
import { ErrorType, STORAGE, RESPONSE, LOG_PULL_STATUS, PLATFORMAS } from '../../common/enum';
import Storage from "../../common/storage";

const context = getCurrentInstance();
let router = useRouter();
let { currentRoute: { _rawValue: { params: { app_key } } } } = router;
let defaultParams = {
  offset: 1,
  limit: 50,
  app_key: app_key,
  user_id: '',
  description: ''
};
let state = reactive({
  isShowPullDialog: false,
  params: utils.clone(defaultParams),
  list: [ ],
  range: {
    start: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000),
    end: new Date(Date.now() + 1 * 24 * 60 * 60 * 1000),
  }
});
function format(date, fmt = 'yyyy-MM-dd') {
  return utils.formatTime(new Date(date).getTime(), fmt);
}
function onShowPullDialog(isShow){
  state.isShowPullDialog = isShow;
}
function onCreatePull(params){
  let data = { ...params, app_key };
  Log.create(data).then((result) => {
    let { code, msg } = result;
    if(utils.isEqual(code, RESPONSE.SUCCESS)){
      utils.extend(state.params, utils.clone(defaultParams))
      onShowPullDialog(false);
      search();
    }else{
      context.proxy.$toast({ icon: 'error', text: `Error: ${code} ${msg}` });
    }
  });
}
function onSearch(){
  search();
}
function onNext(){
  let { offset } = state.params;
  offset += 1;
  utils.extend(state.params, { offset });
  search();
}
function onPre(){
  let { offset } = state.params;
  offset -= 1;
  utils.extend(state.params, { offset });
  search();
}
function onDwonload(item){
  Log.download({ id: item.id, app_key }).then((blob) => {
    downloadLog(blob, `imlog-${format(Date.now(), 'yyyyMMddhhmmss')}.txt`);
  });

  function downloadLog(blob, filename) {
    const a = document.createElement('a')
    a.download = filename
    const blobUrl = URL.createObjectURL(blob)
    a.href = blobUrl
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(blobUrl)
  }
}
function search(){
  let { start, end } = state.range;
  let params = { ...state.params, start: start.getTime(), end: end.getTime() };
  Log.getList(params).then((result) => {
    let { code, data, msg = '' } = result;
    if(utils.isEqual(code, RESPONSE.SUCCESS)){
      let { items } = data;
      items = utils.map(items, (item) => {
        let { created_time, start, end, state: status } = item;
        // let index = utils.find(PLATFORMAS, (item) => { return item.value == platform});
        // let plat = PLATFORMAS[index] || { name: '' }
        // let platformName = plat.name;
        let statusName = LOG_PULL_STATUS[status];
        let startName = format(start, 'yyyy-MM-dd hh:mm:ss');
        let endName = format(end, 'yyyy-MM-dd hh:mm:ss');
        let createName = format(created_time, 'yyyy-MM-dd hh:mm:ss');
        utils.extend(item, {  statusName, startName, endName, createName });
        return item;
      });
      state.list = items.reverse();
    }else{
      context.proxy.$toast({ icon: 'error', text: `Error: ${code} ${msg}` });
    }
  });
}
search()

</script>
<template>
   <div class="mb-4 cim-log-contanier">
    <div class="cim-log-header">
      <ul class="cim-log-header-lf-box">
        <li class="cim-log-lf-item">
          <VDatePicker v-model.range="state.range" >
            <template #default="{ togglePopover }">
              <div class="form-control cim-as-date-content" @click="togglePopover">{{ format(state.range.start) }} 至 {{ format(state.range.end) }}</div>
            </template>
          </VDatePicker>
        </li>
        <li class="cim-log-lf-item">
          <input class="form-control" type="text" v-model="state.params.user_id" placeholder="用户 ID" autocomplete="off" @keydown.enter="onSearch">
        </li>
        <li class="cim-log-lf-item">
          <input class="form-control" type="text" v-model="state.params.description" placeholder="备注信息" autocomplete="off" @keydown.enter="onSearch">
        </li>
        <li class="cim-log-lf-item">
          <div class="cim-button cim-button-bg" @click="onSearch">查询</div>
        </li>
      </ul>
      <ul class="cim-log-header-rg-box">
        <li class="cim-log-lf-item">
          <div class="cim-button cim-button-bg cim-button-green" @click="onShowPullDialog(true)">拉取日志</div>
        </li>
      </ul>
    </div>
    <div class="cim-log-body">
      <table class="table cim-table">
        <thead>
          <tr>
            <th class="cim-td-c">用户 ID</th>
            <!-- <th class="cim-td-c">指令 ID</th> -->
            <th class="cim-td-c">平台</th>
            <th class="cim-td-c">收集时间</th>
            <th class="cim-td-c">拉取时间</th>
            <th class="cim-td-c">日志状态</th>
            <th class="cim-td-c">日志备注</th>
            <!-- <th class="cim-td-c">异常信息</th> -->
            <th class="cim-td-c">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in state.list">
            <td class="cim-td-c">{{ item.user_id }}</td>
            <!-- <td class="cim-td-c">{{ item.msg_id }}</td> -->
            <td class="cim-td-c">{{ item.platform }}</td>
            <td class="cim-td-c">{{ item.startName }} 至 {{ item.endName }}</td>
            <td class="cim-td-c">{{ item.createName }}</td>
            <td class="cim-td-c">
              <span class="cicon cim-log-status" :class="['cicon-status-' + item.state]">{{ item.statusName }}</span>
            </td>
            <td class="cim-td-c">{{ item.description }}</td>
            <!-- <td class="cim-td-c">{{ item.fail_reason }}</td> -->
            <td class="cim-td-c">
              <ul class="cim-table-tools">
                <li class="cim-table-tool" v-if="utils.isEqual(item.state, LOG_PULL_STATUS.COMPLETE)">
                  <a class="btn-link" :href="item.log_url" v-if="item.log_url">下载</a>
                  <a class="btn-link" href="#" v-else  @click="onDwonload(item)">下载</a>
                </li>
              </ul>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="cim-log-footer">
      <nav class="cim-navigation">
        <ul class="pagination">
          <li class="page-item" v-if="state.params.offset > 1" >
            <a class="page-link" href="#" aria-label="Previous"@click="onPre">
              <span aria-hidden="true">上一页</span>
            </a>
          </li>
          <li class="page-item">
            <a class="page-link" href="#" v-if="state.list.length >= state.params.limit"  aria-label="Next" @click="onNext">
              <span aria-hidden="true">下一页</span>
            </a>
          </li>
        </ul>
      </nav>
    </div>
    <PullLogDialog :show="state.isShowPullDialog" :title="'拉取日志'" :text="'确定'" :type="state.dialogType" @hide="onShowPullDialog(false)" @save="onCreatePull"></PullLogDialog>
  </div>
</template>
