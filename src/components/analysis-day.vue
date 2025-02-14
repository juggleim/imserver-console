<script setup>
const props = defineProps(['title', 'item']);
</script>

<template>
  <div class="col-sm-3 cim-asr-col">
    <span class="cim-ars-memo">{{ props.title }}</span>
    <div class="cim-ars-num">{{ item.count }}</div>
    <div class="cim-ars-percent">
      较前一日<span class="cicon cim-ars-direction" :class="[item.isUp ? 'cicon-ac-up' : 'cicon-ac-down']">{{ item.percent }}%</span>
    </div>
  </div>
</template>
