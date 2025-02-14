<script setup>
  import { h, ref, reactive } from 'vue';
  import { NButton, NTag, useDialog } from 'naive-ui';
  import { Application } from '@/services';
  import { useRouter } from 'vue-router';
  import { showToast } from '@/common/toast';
  import SensitiveDataEdit from './sensitive_data_edit.vue';
  import ImportWords from './importWords.vue';

  let router = useRouter();
  let {
    currentRoute: {
      _rawValue: {
        params: { app_key },
      },
    },
  } = router;

  let state = reactive({
    list: []
  });
  const tab = ref(0);
  const loading = ref(false);
  const editRef = ref();
  const importRef = ref();
  const queryForm = ref({
    word: '',
  });

  function reloadTable() {
    loadData({ page: 1, size: 10, app_key });
  }

  function addTable() {
    editRef.value.openModal();
  }

  function importWords() {
    importRef.value.openModal();
  }

  function handleSearch() {
    loading.value = true;
    loadData({ page: 1, size: 10, app_key });
  }

  function handleReset() {
    queryForm.value.word = '';
    loadData({ page: 1, size: 10, appKey: app_key });
  }

  const columns = ref([
    {
      title: '词语',
      key: 'word',
    },
    {
      title: '过滤类型',
      key: 'word_type',
      render(row) {
        return h(
          NTag,
          { type: row.word_type === 1 ? 'success' : 'info' },
          {
            default: () => {
              return row.word_type === 1 ? '过滤' : '替换（****）';
            },
          }
        );
      },
    },
    {
      title: '操作',
      key: 'actions',
      render(row) {
        return h(
          NButton,
          {
            strong: true,
            tertiary: true,
            size: 'small',
            type: 'error',
            onClick: () => handleDelete(row),
          },
          { default: () => '删除' }
        );
      },
    },
  ]);
  const data = [];
  const pagination = ref({});
  reloadTable();
  
  function loadData(params) {
    loading.value = true;
    Application.getSensitiveList({ ...params, word: queryForm.value.word || '' })
      .then((res) => {
        state.list = res.data.items || [];
        pagination.value = {
          page: res.data.page,
          itemCount: res.data.total,
          pageSize: 10,
        };
      })
      .finally(() => {
        loading.value = false;
      });
  }

  function handlePageChange(page) {
    loading.value = true;
    Application.getSensitiveList({
      page: page,
      size: 10,
      appKey: app_key,
      word: queryForm.value.word || '',
    })
      .then((res) => {
        data.value = res.data.list;
        pagination.value = {
          page: res.data.page,
          itemCount: res.data.total,
          pageSize: 10,
        };
      })
      .finally(() => {
        loading.value = false;
      });
  }

  const dialog = useDialog();

  function handleDelete(record) {
    dialog.warning({
      title: '警告',
      content: '你确定要删除？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        Application.deleteSensitiveWord({ app_key, word: record.word }).then((res) => {
          if (res.code === 0) {
            showToast({ text: '删除成功' });
          }
          reloadTable();
        });
      },
    });
  }

  // const model = ref({
  //   handler_type: 1,
  //   replace_char: '*',
  // });

  // Application.getSensitiveConf(app_key).then((res) => {
  //   console.log(res);
  //   model.value = res.data;
  // });

  // function handleSave() {
  //   Application.setSensitiveConf(app_key, model.value).then((res) => {
  //     console.log(res);
  //     window.$toast({ icon: 'success', text: '保存成功' });
  //   });
  // }
</script>

<template>
  <div class="tab-content">
    <n-card :bordered="false" class="proCard">
      <!-- <n-form ref="formRef" inline :label-width="80" :model="queryForm">
        <n-form-item label="敏感词" path="model.logType" style="width: 200px">
          <n-input v-model:value="queryForm.word" placeholder="敏感词" />
        </n-form-item>
        <n-form-item>
          <n-button
            attr-type="button"
            type="primary"
            style="margin-right: 10px"
            @click="handleSearch"
          >
            搜索
          </n-button>
          <n-button attr-type="button" type="default" @click="handleReset"> 重置</n-button>
        </n-form-item>
      </n-form> -->
      <div style="margin-bottom: 10px">
        <n-button type="primary" size="small" @click="addTable"> 添加</n-button>
        &nbsp;
        <n-button type="default" size="small" @click="importWords"> 导入词库</n-button>
      </div>
      <n-data-table
        remote
        ref="table"
        :columns="columns"
        :data="state.list"
        :loading="loading"
        :row-key="(row) => row.id"
        @update:page="handlePageChange"
      />
    </n-card>
    <!-- :pagination="pagination" -->

    <SensitiveDataEdit ref="editRef" @reloadTable="reloadTable" />
    <ImportWords ref="importRef" @reloadTable="reloadTable" />
  </div>
</template>