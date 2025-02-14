<script setup>
import { reactive } from 'vue';
import utils from '../../common/utils';
let state = reactive({
  files: [
    { name: 'a.png', url: '', isImage: true },
    { name: 'b.png', url: '', isImage: true },
    { name: 'c.pdf', url: '', isImage: false },
    { name: 'd.xls', url: '', isImage: false },
  ]
});
function onRemove(index){
  state.files.splice(index, 1);
}
function onAdd(){
  state.files.push({ name: `${Date.now()}.png`, isImage: true })
}
</script>
<template>
  <div class="mb-4">
    <div class="card-body">
      <div class="tab-content rounded-bottom">
        <div class="tab-pane p-3 active preview">
          <div class="mb-3">
            <label class="form-label" >工单标题</label>
            <input class="form-control">
          </div>
          <div class="mb-3">
            <label class="form-label" >问题描述</label>
            <textarea class="form-control" rows="5"></textarea>
          </div>
          <div class="mb-3">
            <label class="form-label" >添加附件</label>
            <div class="form-control cim-form-files">
              <div class="cim-form-file cicon cicon-add cim-form-add" @click="onAdd()"></div>
              <div class="cim-form-file" v-for="(file, index)  in state.files" :class="{'cim-form-image':file.isImage}">
                <div class="cicon cicon-close cim-from-file-close" @click="onRemove(index)"></div>
                <div class="cim-form-filename">{{ file.name }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
