<script  setup>
  import { ref } from 'vue';
  import { NModal, useMessage } from 'naive-ui';
  import { useRouter } from 'vue-router';
  import Application from '@/services/application';
  import { showToast } from '@/common/toast';

  const emit = defineEmits(['reloadTable']);
  const message = useMessage();
  const loading = ref(false);
  const showModal = ref(false);
  const formValue = ref(newState({}));
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

  const wordTypes = [
    { label: '过滤', value: 1 },
    { label: '替换（****）', value: 2 },
  ];

  function newState(state = {}) {
    return {
      id: state.id || 0,
      word: state.word || '',
      wordType: state.wordType || 1,
    };
  }

  function openModal() {
    showModal.value = true;
  }

  function confirmForm(e) {
    e.preventDefault();
    formBtnLoading.value = true;
    formRef.value.validate((errors) => {
      if (!errors) {
        Application.addSensitiveWord({ app_key, word: formValue.value.word, word_type: formValue.value.wordType }).then((res) => {
          showToast({ text: '添加成功' });
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

  defineExpose({ openModal });
</script>


<template>
  <div>
    <n-modal
      v-model:show="showModal"
      :mask-closable="false"
      :show-icon="false"
      preset="dialog"
      transform-origin="center"
      :title="'添加敏感词'"
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
                <n-form-item label="敏感词" path="word">
                  <n-input placeholder="请输入敏感词" v-model:value="formValue.word" />
                </n-form-item>
              </n-gi>
              <n-gi span="1">
                <n-form-item label="过滤类型" path="wordType">
                  <n-select placeholder="请选择过滤类型" v-model:value="formValue.wordType" :options="wordTypes" />
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


