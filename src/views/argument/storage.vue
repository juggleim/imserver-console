<script setup>
  import { ref, getCurrentInstance, reactive } from 'vue';
  import FForm from './form.vue';
  import utils from '../../common/utils.js';
  import { Application } from '../../services';
  import { useRouter } from 'vue-router';
  import { RESPONSE } from '../../common/enum';

  let router = useRouter();
  let context = getCurrentInstance();

 

  let {
    currentRoute: {
      _rawValue: {
        params: { app_key },
      },
    },
  } = router;

  let settings = [
    {
      type: 'aws',
      name: 'AWS',
      state: ref({
        access_key: '',
        secret_key: '',
        endpoint: '',
        region: '',
        bucket: '',
      }),
      fields: [
        {
          name: 'access_key',
          label: 'Access Key',
          type: 'input_text',
        },
        {
          name: 'secret_key',
          label: 'Secret Key',
          type: 'input_text',
        },
        {
          name: 'endpoint',
          label: 'Endpoint',
          type: 'input_text',
        },
        {
          name: 'region',
          label: 'Region',
          type: 'input_text',
        },
        {
          name: 'bucket',
          label: 'Buket Name',
          type: 'input_text',
        },
      ],
    },
    {
      type: 'qiniu',
      name: '七牛云',
      state: ref({
        access_key: '',
        secret_key: '',
        bucket: '',
      }),
      fields: [
        {
          name: 'access_key',
          label: 'Access Key',
          type: 'input_text',
        },
        {
          name: 'secret_key',
          label: 'Secret Key',
          type: 'input_text',
        },
        {
          name: 'domain',
          label: 'Domain Name',
          type: 'input_text',
        },
        {
          name: 'bucket',
          label: 'Bucket Name',
          type: 'input_text',
        },
      ],
    },
    {
      type: 'oss',
      name: '阿里云',
      state: ref({
        access_key: '',
        secret_key: '',
        endpoint: '',
        bucket: '',
      }),
      fields: [
        {
          name: 'access_key',
          label: 'Access Key',
          type: 'input_text',
        },
        {
          name: 'secret_key',
          label: 'Secret Key',
          type: 'input_text',
        },
        {
          name: 'endpoint',
          label: 'Endpoint',
          type: 'input_text',
        },
        {
          name: 'bucket',
          label: 'Bucket name',
          type: 'input_text',
        },
      ],
    },
    {
      type: 'minio',
      name: 'MinIO',
      state: ref({
        access_key: '',
        secret_key: '',
        endpoint: '',
        bucket: '',
        use_ssl: false,
      }),
      fields: [
        {
          name: 'access_key',
          label: 'Access Key',
          type: 'input_text',
        },
        {
          name: 'secret_key',
          label: 'Secret Key',
          type: 'input_text',
        },
        {
          name: 'endpoint',
          label: 'Endpoint',
          type: 'input_text',
        },
        {
          name: 'use_ssl',
          label: '启用 HTTPS',
          type: 'radios',
          radios: [
            { name: 'type', value: false, label: '禁用' },
            { name: 'type', value: true, label: '启用' },
          ]
        },
        {
          name: 'bucket',
          label: 'Buket Name',
          type: 'input_text',
        },
      ],
    },
  ];

  let state = reactive({
    checkedValue: { },
    current: '',
    channels: {
      oss: { name: '阿里云', },
      aws: { name: 'AWS', },
      qiniu: { name: '七牛云', },
      minio: { name: 'MinIO', }
    },
    fileConfs: [],
    formTitle: '添加配置',
  });

  const channel = ref(settings[0].type);

  function onTab(setting) {
    channel.value = setting.type;
    search();
  }

  function isModify(){
    let _conf = state.fileConfs.find((conf) => { 
      return conf.channel == channel.value
    });
    return _conf;
  }

  function onSave(item) {
    let _isModify = isModify();
    Application.setStorageConfig({
      app_key,
      channel: channel.value,
      conf: item,
    }).then(() => {
      let value = channel.value;
      let enable = 0;
      let name = state.channels[value].name;
      let conf = { channel, enable, name };
      if(!_isModify){
        state.fileConfs.push(conf);
        return context.proxy.$toast({ icon: 'success', text: '添加成功，请选择存储类型后，保存设置' });;
      }
      context.proxy.$toast({ icon: 'success', text: '修改成功，请选择存储类型后，保存设置' });;
    });
  }

  async function search() {
    const res = await Application.getStorageConfig({ app_key, channel: channel.value });
    utils.forEach(settings, (item) => {
      if (item.type === channel.value) {
        if (res.data && res.data.conf) {
          item.state.value = { ...item.state.value, ...res.data.conf };
        }
      }
    });
    state.formTitle = isModify() ? '修改配置' : '添加配置'
  }
  search();
  function handleChange(e) {
    let value = e.target.value;
    state.checkedValue.value = value;
  }
  function getCheckedChannel() {
    Application.getEnableStorage({app_key}).then((result) => {
      let { code, data = {} } = result;
      if(utils.isEqual(code, RESPONSE.SUCCESS)){
        let { file_confs } = data;
        let useConf = { channel: '' };
        let fileConfs = utils.map(file_confs, (conf) => {
          let channel = state.channels[conf.channel] || { name: '' };
          if(conf.enable){
            useConf = conf;
          }
          return { ...conf, name: channel.name };
        });
        state.current = useConf.channel;
        state.fileConfs = fileConfs;
      }
    })
  }

  function setCheckedChannel() {
    let channel = state.checkedValue.value;
    Application.setEnableStorage({app_key, channel }).then(() => {
      state.current = channel;
      context.proxy.$toast({ icon: 'success', text: '保存成功' });
    })
  }

  getCheckedChannel();
</script>
<template>
  <n-flex vertical>
    <div class="mb-4 app-base cim-cb-box">
      <div class="row cim-cb-row cim-cb-header">
        <div class="cim-cb-form cim-file-form">
         <div class="cim-cb-form-item">
            <label class="col-sm-1 col-form-label cim-form-item-label">
              正在使用存储
            </label>
            <div class="col-sm-7">
              <span class="warn">{{ state.channels[state.current] && state.channels[state.current].name || '未设置' }}</span>
            </div>
         </div>
         <div class="cim-cb-form-item">
            <label class="col-sm-1 col-form-label">
              设置存储类型
            </label>
            <div class="col-sm-4 store-form">
                <div class="form-check-inline store_inline" v-for="setting in state.fileConfs">
                  <input class="form-check-input" type="radio" name="setting.name" :value="setting.channel" :checked="setting.enable" @change="handleChange">
                  <label class="form-check-label">{{ setting.name }}</label>
                </div>
                <div class="form-check-inline store_inline" v-if="state.fileConfs.length == 0">
                  <span class="warn fs-12">请优先添加文件存储配置</span>
                </div>
            </div>
            <div class="col-sm-3 cim-cb-btns">
              <div class="cim-button cim-button-bg warn-bg" v-if="state.fileConfs.length > 0" @click="setCheckedChannel">保存设置</div>
            </div>
         </div>
        </div>
      </div>
    </div>

    <div class="md-6">
      <ul class="nav nav-underline-border" role="tablist">
        <li class="nav-item sw-nav-item" v-for="setting in settings" @click="onTab(setting)">
          <a
            class="nav-link cicon cicon-free"
            :class="{ active: utils.isEqual(channel, setting.type) }"
            > 添加 {{ setting.name }} 配置</a
          >
        </li>
      </ul>
      <div class="tab-content rounded-bottom">
        <div
          class="tab-pane p-3"
          v-for="setting in settings"
          :class="{ active: utils.isEqual(channel, setting.type) }"
        >
        <FForm :fields="setting.fields" :state="setting.state.value" :btn-type="1" :title="state.formTitle" @save="onSave" />
        </div>
      </div>
    </div>
  </n-flex>
</template>
