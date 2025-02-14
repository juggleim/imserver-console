<script setup lang="ts">
  import { RouterView, useRoute } from 'vue-router';
  import { watch } from 'vue';
  import {
    dateZhCN,
    GlobalThemeOverrides,
    NConfigProvider,
    NDialogProvider,
    NMessageProvider,
    zhCN,
  } from 'naive-ui';
  import Storage from '@/common/storage';
  import { STORAGE } from '@/common/enum';

  const themeOverrides: GlobalThemeOverrides = {
    common: {
      primaryColor: '#2d4af2',
    },
    Button: {
      textColor: '#2d4af2',
      hoverColor: '#2d4af2',
    },
    Select: {
      peers: {
        InternalSelection: {
          textColor: '#2d4af2',
        },
      },
    },
    // ...
  };
  const route = useRoute();
  Storage.set(STORAGE.APP_KEY,  route.params.app_key);
  watch(
    () => route.params.app_key,
    (newId) => {
      Storage.set(STORAGE.APP_KEY,  newId);
    }
  );
</script>

<template>
  <n-config-provider :locale="zhCN" :date-locale="dateZhCN" :theme-overrides="themeOverrides">
    <n-dialog-provider>
      <n-message-provider>
        <RouterView />
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>
