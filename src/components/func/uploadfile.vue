<template>
  <div class="form-control mb-4" style="height: 100%;">
    <p>{{ value }}</p>
    <input type="file" style="display: none;" @change="onSelect" />
    <div class="cim-form-file cicon cicon-add cim-form-add cim-upload" @click="onAdd"></div>
  </div>
</template>

<script setup lang="ts">
  import { Application } from '../../services';
  import { defineEmits, defineProps, computed, getCurrentInstance } from 'vue';

  let context = getCurrentInstance();

  const props = defineProps<{value: string}>();
  const emit = defineEmits(['update:value']);

  const value = computed(() => props.value);

  function onAdd(e){
    e.target.parentElement.childNodes[1].click();
  }
  function onSelect(event: any) {
    const file = event.target.files[0];

    console.log(file);
    Application.uploadFile(file).then((res) => {
      value.value = res.data.path;
      emit('update:value', res.data.path);
      context.proxy.$toast({ icon: 'success', text: 'Upload success', duration: 4000 });
    });
  }
</script>
