<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      :title="'导入词库'"
      :style="{
        width: 840,
      }"
    >
      <n-scrollbar style="max-height: 87vh" class="pr-5">
        <n-spin :show="loading" description="请稍候...">
          <n-form
            ref="formRef"
            :model="formValue"
            :label-placement="'left'"
            :label-width="100"
            class="py-4"
          >
            <n-grid cols="1 s:1 m:1 l:1 xl:1 2xl:1" responsive="screen">
              <n-gi span="1">
                <n-form-item label="词库" path="word">
                  <input type="file" @change="handleFileChange" />
                </n-form-item>
              </n-gi>
            </n-grid>
          </n-form>
        </n-spin>
      </n-scrollbar>
      <template #action>
        <n-space>
          <n-button @click="closeForm"> 取消</n-button>
          <n-button type="info" :loading="formBtnLoading" @click="confirmForm"> 确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup>
  import { ref } from 'vue';
  import { NModal, useMessage } from 'naive-ui';
  import { useRouter } from 'vue-router';
  import { Application } from '@/services';
  import { showToast } from '@/common/toast';
  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const loading = ref(false);
  const showModal = ref(false);
  const formValue = ref({});
  const formRef = ref({});
  const formBtnLoading = ref(false);

  let router = useRouter();
  let {
    currentRoute: {
      _rawValue: {
        params: { app_key },
      },
    },
  } = router;
  let file = null;

  function handleFileChange(e) {
    if (e.target.files.length > 0) {
      file = e.target.files[0];
    }
  }

  function openModal() {
    showModal.value = true;
  }

  function confirmForm(e) {
    e.preventDefault();
    formBtnLoading.value = true;
    formRef.value.validate((errors) => {
      if (!errors) {
        console.log(file);

        Application.importSensitiveWords(app_key, file).then((res) => {
          showToast({ text: '导入成功' });
          emit('reloadTable');
          showModal.value = false;
        }).finally(() => {
          formBtnLoading.value = false;
        });
      } else {
        message.error('请填写完整信息');
      }
    });
  }

  function closeForm() {
    showModal.value = false;
    loading.value = false;
  }

  defineExpose({
    openModal,
  });
</script>

<style scoped></style>
