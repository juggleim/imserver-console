<script setup>
import Dialog from './dialog.vue';
import { reactive, getCurrentInstance } from 'vue';
import utils from '../common/utils';

const context = getCurrentInstance();
const props = defineProps(['title','text']);
const emit = defineEmits(['save', 'hide'])

let state = reactive({
  params: {
    user_id: '',
    description: '',
    platform: 'Android',
  },
  range: {
    start: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000),
    end: new Date(),
  },
  times: [{value: '1'}],
  categories: [
    { name: '业务日志', value: 1 }
  ],
  platforms: [
    { name: 'Android', value: 'Android' },
    { name: 'iOS', value: 'iOS' },
    { name: 'Web', value: 'Web' },
    { name: 'PC', value: 'PC' },
  ]
});

function onSave(){
  let { user_id, description } = state.params;
  if(utils.isEqual(user_id.length, 0)){
    return context.proxy.$toast({ icon: 'error', text: `用户 ID 不可为空` });
  }
  if(utils.isEqual(description.length, 0)){
    return context.proxy.$toast({ icon: 'error', text: `备注 不可为空` });
  }
  let params = utils.clone(state.params);
  let { start, end } = state.range;
  emit('save', { ...params, start: start.getTime(), end: end.getTime() });
}
function onHide(e){
  emit('hide', {});
}
function onInput(){
  state.licenseErrorMsg = '';
}
function format(date) {
  return utils.formatTime(new Date(date).getTime(), 'yyyy-MM-dd hh:mm');
}
</script>

<template>
 <Dialog :title="props.title" :btn-title="props.text" :cls="'cim-pull-dialog'" @hide="onHide" @save="onSave" >
    <div class="row g-2 cim-row cim-dialog-pull-row">
        <!-- <div class="form-floating">
          <input class="form-control" :disabled="true" placeholder="AppKey" @input="onInput">
          <label>AppKey</label>
        </div> -->

        <!-- <div class="form-floating">
          <select class="form-select" >
            <option :value="category.value" v-for="category in state.categories" >{{ category.name }}</option>
          </select>
          <label>日志类型</label>
        </div> -->
        
        <div class="form-floating">
          <select class="form-select" v-model="state.params.platform">
            <option :value="item.value" v-for="item in state.platforms" >{{ item.name }}</option>
          </select>
          <label>平台</label>
        </div>
        <div class="form-floating cim-dialog-pull-datetime">
          <VDatePicker v-model.range="state.range" mode="dateTime" is24hr>
            <template #default="{ togglePopover }" >
              <div class="form-control cim-as-date-content" @click="togglePopover">{{ format(state.range.start) }} 至 {{ format(state.range.end) }}</div>
            </template>
          </VDatePicker>
          <label>时间范围</label>
        </div>
        <div class="form-floating cim-from-must cicon cicon-must">
          <input class="form-control" placeholder="用户 ID" v-model="state.params.user_id">
          <label>用户 ID </label>
        </div>
        <div class="form-floating cim-from-must cicon cicon-must">
          <input class="form-control" placeholder="备注说明" v-model="state.params.description">
          <label>备注说明</label>
        </div>
    </div>
 </Dialog>
</template>
