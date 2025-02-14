<script setup>
import { reactive } from 'vue';
let state = reactive({
  files: [
    { name: 'a.png', url: '', isImage: true },
  ]
});
</script>
<template>
   <div class="md-4 app-base ent-body">
    <ul class="nav nav-underline-border ab-underline-border">
      <li class="nav-item"><a class="nav-link active cicon cicon-product">认证信息</a></li>
    </ul>
    <div class="cim-us-table">
      <div class="row cim-ab-row">
        <label class="col-sm-2 col-form-label cim-td-r cicon cicon-require">公司名称</label>
        <div class="col-sm-9">
          <input class="form-control" type="text"  value="北京追风未来科技有限公司">
        </div>
      </div>

      <div class="row cim-ab-row">
        <label class="col-sm-2 col-form-label cim-td-r cicon cicon-require">注册地址</label>
        <div class="col-sm-9">
          <input class="form-control" type="text"  value="北京市延庆区">
        </div>
      </div>

      <div class="row cim-ab-row">
        <label class="col-sm-2 col-form-label cim-td-r cicon cicon-require">统一社会信任代码</label>
        <div class="col-sm-9">
          <input class="form-control" type="text"  value="">
        </div>
      </div>

      <div class="row cim-ab-row">
        <label class="col-sm-2 col-form-label cim-td-r cicon cicon-require">公司法人</label>
        <div class="col-sm-9">
          <input class="form-control" type="text"  value="">
        </div>
      </div>
      
      <div class="row cim-ab-row">
        <label class="col-sm-2 col-form-label cim-td-r cicon cicon-require">公司电话</label>
        <div class="col-sm-9">
          <input class="form-control" type="text"  value="">
        </div>
      </div>

      <div class="row cim-ab-row">
        <label class="col-sm-2 col-form-label cim-td-r cicon cicon-require ent-website">公司网站</label>
        <div class="col-sm-9">
          <input class="form-control" type="text"  value="">
        </div>
      </div>

      <div class="row cim-ab-row">
        <label class="col-sm-2 col-form-label cim-td-r cicon cicon-require">营业执照</label>
        <div class="col-sm-9">
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
</template>
