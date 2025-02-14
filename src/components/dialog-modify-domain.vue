<script setup>
import Dialog from './dialog.vue';
import { reactive, getCurrentInstance, watch } from 'vue';
import utils from '../common/utils';

const context = getCurrentInstance();
const props = defineProps(['domains']);
const emit = defineEmits(['save', 'hide'])

let state = reactive({
  domains: props.domains,
  text: props.domains.join(';\n')
});

function onSave(){
  let { text } = state;
  let domains = text.split(';')
  domains = utils.filter(domains, (domain) => {
    return !utils.isEmpty(domain);
  });
  domains = utils.map(domains, (domain) => {
    return domain.replace('\n', '');
  });
  for(let i = 0; i < domains.length; i++){
    let domain = domains[i];
    
    if(!utils.isURL(domain)){
      return context.proxy.$toast({ icon: 'error', text: `${domain} 格式不正确` });
    }
  }
  emit('save', { domains });
}
function onHide(e){
  emit('hide', {});
}
watch(() => props.domains, () => {
  utils.extend(state, { 
    domains: props.domains,
    text: props.domains.join(';\n')
  });
})
</script>

<template>
 <Dialog :title="'安全域名'" :btn-title="'保存'" :class="['cim-security-dialog']" @hide="onHide" @save="onSave" >
    <div class="form-floating">
      <textarea class="form-control" style="width: 36rem;height: 10rem;resize: none;" placeholder=" " v-model="state.text"></textarea>
      <label>请输入安全域名, 多个域名请用【英文分号】分割，取消设置安全域名，清空保存即可</label>
    </div>
    <div class="form-example">
      <div class="title">域名示例：</div>
      <ul>
        <li>https://example.com;http://im.fake.com;http://127.0.0.1:8305</li>
        <li class="warn">配置安全域名后，Web 只能在已设置的域名下连接 IM Server，生产环境请谨慎使用</li>
      </ul>
    </div>
 </Dialog>
</template>
