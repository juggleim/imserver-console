<script setup>
import { reactive, getCurrentInstance, nextTick, watch } from 'vue';
import { STORAGE, ErrorType } from "../../../common/enum";
import Storage from "../../../common/storage";
import utils from '../../../common/utils';
import { User } from "../../../services";
import Clipboard from 'clipboard.js';

const context = getCurrentInstance();
const props = defineProps(['conn']);
const emit = defineEmits(['save', 'hide', 'next', 'pre']);
let state = reactive({});

watch(() => props.conn.list, () => {
  nextTick(() => {
    let { signals } = context.refs;
    if (signals) {
      signals.scrollTop = signals.scrollHeight;
    }
  });
})

function onCopy(item){
  Clipboard.copy(item.timestamp, utils.noop, utils.noop);
  context.proxy.$toast({ icon: 'success', text: `【时间】已复制`, duration: 1500 });
}

</script>
<template>
  <div class="cim-signal-contanier">
    <ul class="cim-signals-box" ref="signals">
      <li class="cim-signal fadeinx" v-for="item in props.conn.list">
        <div class="cim-signal-header">
          <div class="cim-signal-time" @click="onCopy(item)">{{ item.timeName }}</div>
          <div class="cim-signal-avatar">
            <span class="cim-signal-icon cicon" :class="['cicon-conn-' + item.type]"></span>
            <span class="cim-signal-name">{{ item.title }}</span>
          </div>
        </div>
        
        <div class="cim-signal-body">
          <ul class="cim-signal-list">
            <li class="cim-signal-item" v-for="info in item.infos" :class="[info.cls ? info.cls : '']">
              {{ info.name }} {{ info.value }}
            </li>
          </ul>
        </div>
      </li>
    </ul>
    <div class="cim-table-footer"  v-if="props.conn.currentCount >= props.conn.count">
      <nav class="cim-navigation">
        <ul class="pagination">
          <li class="page-item">
            <a class="page-link" href="#" aria-label="Next" @click="emit('next', { item: props.conn })">
              <span aria-hidden="true">下一页</span>
            </a>
          </li>
        </ul>
      </nav>
    </div>

    <div class="cim-loading" v-if="props.conn.isLoading">
      <div class="loader-dot"></div>
    </div>
  </div>
</template>
