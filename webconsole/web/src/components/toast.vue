<script setup>
const props = defineProps(['icon', 'title', 'text']);
</script>

<template>
  <div class="toaster">
      <div class="toast fade show cim-toast-box">
        <div class="cim-toast-icon cicon" :class="'cicon-'+props.icon"></div>
        <div class="cim-toast-content">
          <div class="cim-toast-title">{{ props.title }}</div>
          <div class="cim-toast-text">{{ props.text }}</div>
        </div>
      </div>
    </div>
</template>
