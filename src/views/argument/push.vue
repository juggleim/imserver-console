<script setup>
  import { getCurrentInstance, ref } from 'vue';
  import FForm from './form.vue';
  import utils from '../../common/utils.js';
  import { Application } from '../../services';
  import { useRouter } from 'vue-router';
import { RESPONSE } from '../../common/enum';

  let router = useRouter();
  let {
    currentRoute: {
      _rawValue: {
        params: { app_key },
      },
    },
  } = router;

  const context = getCurrentInstance();

  let settings = [
    {
      type: 'Huawei',
      name: '华为',
      state: ref({
        package: '',
        app_id: '',
        app_secret: '',
      }),
      fields: [
        {
          name: 'package',
          label: '包名',
          type: 'input_text',
        },
        {
          name: 'app_id',
          label: 'App Id',
          type: 'input_text',
        },
        {
          name: 'app_secret',
          label: 'App Secret',
          type: 'input_text',
        },
      ],
    },
    {
      type: 'Xiaomi',
      name: '小米',
      state: ref({
        package: '',
        app_secret: '',
      }),
      fields: [
        {
          name: 'package',
          label: '包名',
          type: 'input_text',
        },
        {
          name: 'app_secret',
          label: 'App Secret',
          type: 'input_text',
        },
      ],
    },
    {
      type: 'Oppo',
      name: 'OPPO',
      state: ref({
        package: '',
        app_key: '',
        master_secret: '',
      }),
      fields: [
        {
          name: 'package',
          label: '包名',
          type: 'input_text',
        },
        {
          name: 'app_key',
          label: 'App Key',
          type: 'input_text',
        },
        {
          name: 'master_secret',
          label: 'Master Secret',
          type: 'input_text',
        },
      ],
    },
    {
      type: 'Vivo',
      name: 'vivo',
      state: ref({
        package: '',
        app_id: '',
        app_key: '',
        app_secret: '',
      }),
      fields: [
        {
          name: 'package',
          label: '包名',
          type: 'input_text',
        },
        {
          name: 'app_id',
          label: 'App Id',
          type: 'input_text',
        },
        {
          name: 'app_key',
          label: 'App Key',
          type: 'input_text',
        },
        {
          name: 'app_secret',
          label: 'App Secret',
          type: 'input_text',
        },
      ],
    },
    {
      type: 'ios',
      name: 'iOS',
      state: ref({
        package: '',
        cert_path: '',
        cert_pwd: '',
        voip_ioscer: '',
        voip_cert_path: '',
        voip_cert_pwd: '',
        is_product: 0,
      }),
      fields: [
        {
          name: 'package',
          label: '包名',
          type: 'input_text',
        },
        {
          name: 'cert_path',
          label: '证书文件',
          type: 'file',
        },
        {
          name: 'cert_pwd',
          label: '证书密码',
          type: 'input_text',
        },
        {
          name: 'voip_cert_pwd',
          label: 'VoIP 证书密码',
          type: 'input_text',
        },
        {
          name: 'voip_cert_path',
          label: 'VoIP 证书文件',
          type: 'voipfile',
        },
        {
          name: 'is_product',
          label: '证书环境',
          type: 'radios',
          radios: [
          { name: 'type', value: 0, label: '开发环境' },
          { name: 'type', value: 1, label: '生产环境' },
          ]
        },
      ],
    },
    {
      type: 'fcm',
      name: 'FCM',
      state: ref({
        package: '',
        conf_path: '',
        app_key: '',
        fcm_conf: '',
      }),
      fields: [
        {
          name: 'package',
          label: '包名',
          type: 'input_text',
        },
        {
          name: 'fcm_conf',
          label: '证书文件',
          type: 'file',
        },
      ],
    },
    // {
    //   type: 'Jpush',
    //   name: '极光',
    //   state: ref({
    //     package: '',
    //     app_key: '',
    //     master_secret: '',
    //   }),
    //   fields: [
    //     {
    //       name: 'package',
    //       label: '包名',
    //       type: 'input_text',
    //     },
    //     {
    //       name: 'app_key',
    //       label: 'App Key',
    //       type: 'input_text',
    //     },
    //     {
    //       name: 'master_secret',
    //       label: 'Master Secret',
    //       type: 'input_text',
    //     },
    //   ],
    // },
  ];

  const channel = ref(settings[0].type);

  function onTab(setting) {
    channel.value = setting.type;
    search();
  }

  function onSave(item) {
    if (channel.value.startsWith('ios')) {
      let { file, voipFile } = item;
      if(file.name || voipFile.name){
        return Application.uploadIosPushConfig({
          app_key,
          package: item.package,
          cert_pwd: item.cert_pwd,
          file: item.file,
          is_product: item.is_product,
          voip_cert_pwd: item.voip_cert_pwd,
          voipFile: item.voipFile,
        }).then(() => {
          context.proxy.$toast({ icon: 'success', text: '保存成功' });
        });
      }

      Application.setIosPushConfig({
        app_key,
        package: item.package,
        cert_pwd: item.cert_pwd,
        voip_cert_pwd: item.voip_cert_pwd,
        is_product: item.is_product,
      }).then(() => {
        context.proxy.$toast({ icon: 'success', text: '保存成功' });
      });

    } else if(channel.value.startsWith('fcm')){
      Application.uploadFcmPushConfig({
        app_key,
        package: item.package,
        file: item.file,
      }).then(() => {
        context.proxy.$toast({ icon: 'success', text: '保存成功' });
      });
    } else {
      Application.setAndroidPushConfig({
        app_key,
        push_channel: channel.value,
        package: item.package,
        extra: item,
      }).then(() => {
        context.proxy.$toast({ icon: 'success', text: '保存成功' });
      });
    }
  }

  async function search() {
    if (channel.value.startsWith('ios')) {
      const res = await Application.getIosPushConfig({ app_key, push_channel: channel.value });
      let { data } = res;
      if(!data){
        data = { is_product: 0 };
      }
      utils.forEach(settings, (item) => {
        if (item.type === channel.value) {
          item.state.value = {
            ...item.state.value,
            package: data.package,
            cert_path: data.cert_path,
            cert_pwd: data.cert_pwd,
            voip_cert_pwd: data.voip_cert_pwd,
            voip_cert_path: data.voip_cert_path,
            is_product: data.is_product,
          };
        }
      });
    }  else if(channel.value.startsWith('fcm')){
      const res = await Application.getFcmPushConfig({ app_key, push_channel: channel.value });
      let { data } = res;
      if(!data){
        data = { is_product: 0 };
      }
      utils.forEach(settings, (item) => {
        if (item.type === channel.value) {
          item.state.value = {
            ...item.state.value,
            package: data.package,
            fcm_conf: data.conf_path,
          };
        }
      });
    }else {
      const res = await Application.getAndroidPushConfig({ app_key, push_channel: channel.value });
      let { code,  data = {} } = res;
      if(!utils.isEqual(code, RESPONSE.SUCCESS)){
        return  context.proxy.$toast({ icon: 'error', text: '查询异常，请刷新重试' });;
      }
      utils.forEach(settings, (item) => {
        if (item.type === channel.value) {
          item.state.value = {
            ...item.state.value,
            package: data.package,
            ...data.extra,
          };
        }
      });
    }
  }

  search();
</script>
<template>
  <div class="md-4">
    <ul class="nav nav-underline-border" role="tablist">
      <li class="nav-item sw-nav-item" v-for="setting in settings" @click="onTab(setting)">
        <a
          class="nav-link cicon cicon-free"
          :class="{ active: utils.isEqual(channel, setting.type) }"
          >{{ setting.name }}</a
        >
      </li>
    </ul>
    <div class="tab-content rounded-bottom">
      <div
        class="tab-pane p-3"
        v-for="setting in settings"
        :class="{ active: utils.isEqual(channel, setting.type) }"
      >
        <div class="row cim-sw-row">
          <FForm :fields="setting.fields" :state="setting.state.value" @save="onSave" />
        </div>
      </div>
    </div>
  </div>
</template>
