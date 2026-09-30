<script setup>
  import { computed, getCurrentInstance, nextTick, onBeforeUnmount, reactive, watch } from 'vue';
  import { useRoute } from 'vue-router';
  import { t, useI18n } from '@/i18n';
  import utils from '../../common/utils';
  import { CONVERSATION_TYPE, RESPONSE } from '../../common/enum';
  import { Analysis, Application } from '../../services';
  import {
    ONE_HOUR_MS,
    REALTIME_RANGES,
    REALTIME_STAT_SWITCH_KEY,
    createRealtimeRange,
    fillRealtimeSeries,
    getRealtimeBucketMs,
    isRealtimeRangeValid,
  } from './message-realtime.mjs';

  const context = getCurrentInstance();
  const route = useRoute();
  const { locale } = useI18n();
  const appKey = computed(() => route.params.app_key);

  const state = reactive({
    activeDurationMs: ONE_HOUR_MS,
    buttons: REALTIME_RANGES,
    channelType: CONVERSATION_TYPE.PRIVATE,
    range: createRealtimeRange(),
    switchValue: 0,
    switchSaving: false,
  });

  const isRealtimeStatEnabled = computed(() => Number(state.switchValue) === 1);
  const channelTypes = computed(() => [
    {
      label: t('analysis.realtime.channel.private'),
      value: CONVERSATION_TYPE.PRIVATE,
    },
    {
      label: t('analysis.realtime.channel.group'),
      value: CONVERSATION_TYPE.GROUP,
    },
    {
      label: t('analysis.realtime.channel.chatroom'),
      value: CONVERSATION_TYPE.CHATROOM,
    },
  ]);

  let chart = null;
  let isReady = false;
  let latestRequestId = 0;

  function notify(icon, text) {
    context.proxy.$toast({ icon, text });
  }

  function formatDateTime(date) {
    return utils.formatTime(new Date(date).getTime(), 'yyyy-MM-dd hh:mm');
  }

  function formatAvgPerSec(value) {
    return Number(value).toLocaleString(undefined, {
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    });
  }

  function syncActiveDurationFromRange() {
    const rangeMs = new Date(state.range.end).getTime() - new Date(state.range.start).getTime();
    const matched = REALTIME_RANGES.find((item) => Math.abs(item.durationMs - rangeMs) < 3000);
    state.activeDurationMs = matched ? matched.durationMs : null;
  }

  function getAxisLabelFormatter(bucketMs) {
    return (value) => {
      if (bucketMs <= 30 * 1000) return utils.formatTime(value, 'hh:mm:ss');
      if (bucketMs <= 60 * 1000) return utils.formatTime(value, 'hh:mm');
      return utils.formatTime(value, 'MM-dd hh:mm');
    };
  }

  function drawChart(result, rangeStart, rangeEnd) {
    if (!chart) {
      chart = context.proxy.$echat.init(context.refs.realtimeChart);
    }

    function disposeChart() {
      chart?.dispose();
      chart = null;
    }

    const bucketMs = getRealtimeBucketMs(rangeEnd - rangeStart);
    const upData = fillRealtimeSeries(result.upMsgs, rangeStart, rangeEnd, bucketMs);
    const downData = fillRealtimeSeries(result.downMsgs, rangeStart, rangeEnd, bucketMs);
    const dispatchData = fillRealtimeSeries(result.dispatchMsgs, rangeStart, rangeEnd, bucketMs);
    const pointCount = Math.max(upData.length, downData.length, dispatchData.length);
    const legend = {
      up: t('analysis.realtime.legendUp'),
      down: t('analysis.realtime.legendDown'),
      dispatch: t('analysis.realtime.legendDispatch'),
    };

    chart.setOption(
      {
        legend: { data: [legend.up, legend.down, legend.dispatch] },
        tooltip: {
          trigger: 'axis',
          axisPointer: { type: 'cross' },
          formatter(params) {
            if (!params || !params.length) return '';
            const timeFormat = bucketMs <= 30 * 1000 ? 'yyyy-MM-dd hh:mm:ss' : 'yyyy-MM-dd hh:mm';
            const lines = [utils.formatTime(params[0].value[0], timeFormat)];
            params.forEach((item) => {
              lines.push(
                `${item.marker}${item.seriesName}: ${formatAvgPerSec(item.value[1])} ${t(
                  'analysis.realtime.unitPerSecond'
                )}`
              );
            });
            return lines.join('<br/>');
          },
        },
        grid: { left: '3%', right: '4%', bottom: 24, containLabel: true },
        xAxis: {
          type: 'time',
          min: rangeStart,
          max: rangeEnd,
          boundaryGap: false,
          axisLabel: { formatter: getAxisLabelFormatter(bucketMs), hideOverlap: true },
        },
        yAxis: {
          type: 'value',
          name: t('analysis.realtime.unitPerSecond'),
          axisLabel: { formatter: (value) => formatAvgPerSec(value) },
        },
        series: [
          {
            name: legend.up,
            type: 'line',
            smooth: true,
            showSymbol: pointCount <= 120,
            data: upData,
            lineStyle: { color: '#5470C6' },
          },
          {
            name: legend.down,
            type: 'line',
            smooth: true,
            showSymbol: pointCount <= 120,
            data: downData,
            lineStyle: { color: '#008000' },
          },
          {
            name: legend.dispatch,
            type: 'line',
            smooth: true,
            showSymbol: pointCount <= 120,
            data: dispatchData,
            lineStyle: { color: '#EE6666' },
          },
        ],
      },
      true
    );
  }

  async function queryRealtimeStatistics(params) {
    try {
      const { data, code } = await Analysis.getRealtimeMessageChat(params);
      if (!utils.isEqual(code, RESPONSE.SUCCESS)) {
        notify('error', t('analysis.realtime.fetchFailed'));
        return null;
      }
      return {
        dispatchMsgs: data?.msg_dispatch?.items || [],
        downMsgs: data?.msg_down?.items || [],
        upMsgs: data?.msg_up?.items || [],
      };
    } catch (_error) {
      notify('error', t('analysis.realtime.fetchFailed'));
      return null;
    }
  }

  async function loadChart() {
    if (!isRealtimeStatEnabled.value) return;
    const start = new Date(state.range.start).getTime();
    const end = new Date(state.range.end).getTime();
    if (!isRealtimeRangeValid(start, end)) {
      notify('error', t('analysis.realtime.invalidRange'));
      return;
    }

    const requestId = ++latestRequestId;
    const result = await queryRealtimeStatistics({
      app_key: appKey.value,
      start,
      end,
      channel_type: state.channelType,
    });
    if (requestId !== latestRequestId || !isRealtimeStatEnabled.value) return;
    if (result) drawChart(result, start, end);
  }

  function onShortcutClick(item) {
    state.activeDurationMs = item.durationMs;
    state.range = createRealtimeRange(item.durationMs);
  }

  async function loadSetting() {
    if (utils.isEmpty(appKey.value)) return;
    try {
      const { data, code } = await Application.getSetting({
        app_key: appKey.value,
        config_keys: [REALTIME_STAT_SWITCH_KEY],
      });
      if (!utils.isEqual(code, RESPONSE.SUCCESS)) {
        notify('error', t('analysis.realtime.settingFetchFailed'));
        return;
      }
      const value = data?.configs?.[REALTIME_STAT_SWITCH_KEY];
      state.switchValue = utils.isEmpty(value) ? 0 : Number(value);
    } catch (_error) {
      notify('error', t('analysis.realtime.settingFetchFailed'));
    }
  }

  async function onSwitchToggle(event) {
    const input = event.target;
    const value = input.checked ? 1 : 0;
    state.switchSaving = true;
    try {
      const { code } = await Application.updateSetting({
        app_key: appKey.value,
        id: REALTIME_STAT_SWITCH_KEY,
        value,
      });
      if (!utils.isEqual(code, RESPONSE.SUCCESS)) {
        input.checked = isRealtimeStatEnabled.value;
        notify('error', t('analysis.realtime.settingSaveFailed'));
        return;
      }
      state.switchValue = value;
      notify('success', t('analysis.realtime.settingSaveSuccess'));
      if (value === 1 && isReady) {
        await nextTick();
        loadChart();
      } else if (value === 0) {
        latestRequestId += 1;
        disposeChart();
      }
    } catch (_error) {
      input.checked = isRealtimeStatEnabled.value;
      notify('error', t('analysis.realtime.settingSaveFailed'));
    } finally {
      state.switchSaving = false;
    }
  }

  function resizeChart() {
    chart?.resize();
  }

  watch(
    () => [state.range.start, state.range.end, state.channelType],
    () => {
      if (!isReady || !isRealtimeStatEnabled.value) return;
      syncActiveDurationFromRange();
      loadChart();
    }
  );

  watch(locale, () => {
    if (chart && isRealtimeStatEnabled.value) loadChart();
  });

  watch(appKey, async () => {
    latestRequestId += 1;
    disposeChart();
    isReady = false;
    state.switchValue = 0;
    await loadSetting();
    state.range = createRealtimeRange();
    state.activeDurationMs = ONE_HOUR_MS;
    isReady = true;
    if (isRealtimeStatEnabled.value) {
      await nextTick();
      loadChart();
    }
  });

  loadSetting().then(() => {
    nextTick(() => {
      state.activeDurationMs = ONE_HOUR_MS;
      isReady = true;
      if (isRealtimeStatEnabled.value) loadChart();
    });
  });

  window.addEventListener('resize', resizeChart);
  onBeforeUnmount(() => {
    latestRequestId += 1;
    window.removeEventListener('resize', resizeChart);
    disposeChart();
  });
</script>

<template>
  <div class="mb-4 cim-as-box realtime-statistics-page">
    <div class="realtime-statistics-switch">
      <div>
        <div class="realtime-statistics-title">{{ t('analysis.realtime.enable') }}</div>
        <div class="realtime-statistics-description">
          {{ t('analysis.realtime.description') }}
        </div>
      </div>
      <div class="form-check form-switch realtime-statistics-toggle">
        <input
          class="form-check-input"
          type="checkbox"
          :checked="isRealtimeStatEnabled"
          :disabled="state.switchSaving"
          @change="onSwitchToggle"
        />
      </div>
    </div>

    <template v-if="isRealtimeStatEnabled">
      <div class="cim-as-tools">
        <div class="cim-as-tool realtime-statistics-tools">
          <div
            v-for="item in state.buttons"
            :key="item.durationMs"
            class="cim-as-button"
            :class="{ 'cim-as-button-active': state.activeDurationMs === item.durationMs }"
            @click="onShortcutClick(item)"
          >
            {{ t(item.titleKey) }}
          </div>
          <div class="cim-as-date cicon cicon-date">
            <VDatePicker
              v-model.range="state.range"
              mode="dateTime"
              is24hr
              class="cim-as-date-picker"
            >
              <template #default="{ togglePopover }">
                <div class="cim-as-date-content" @click="togglePopover">
                  {{ formatDateTime(state.range.start) }} {{ t('common.word.to') }}
                  {{ formatDateTime(state.range.end) }}
                </div>
              </template>
            </VDatePicker>
          </div>
          <div class="realtime-statistics-channel">
            <select v-model="state.channelType" class="form-select">
              <option v-for="item in channelTypes" :key="item.value" :value="item.value">
                {{ item.label }}
              </option>
            </select>
          </div>
        </div>
      </div>
      <div class="row cim-as-body realtime-statistics-chart-row">
        <div ref="realtimeChart" class="cim-bk-form"></div>
      </div>
    </template>
  </div>
</template>

<style scoped>
  .realtime-statistics-page {
    padding-top: 0;
  }

  .realtime-statistics-switch {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 24px;
    min-height: 82px;
    padding: 18px 24px;
    border-bottom: 1px solid #e8edf3;
  }

  .realtime-statistics-title {
    color: #1f2937;
    font-size: 16px;
    font-weight: 600;
  }

  .realtime-statistics-description {
    margin-top: 4px;
    color: #7a8596;
    font-size: 13px;
  }

  .realtime-statistics-toggle {
    flex: 0 0 auto;
    margin: 0;
  }

  .realtime-statistics-toggle .form-check-input {
    width: 44px;
    height: 24px;
    margin: 0;
    cursor: pointer;
  }

  .realtime-statistics-tools {
    flex-wrap: wrap;
    gap: 12px;
  }

  .realtime-statistics-channel {
    width: 150px;
  }

  .realtime-statistics-chart-row {
    min-height: 420px;
  }

  @media (max-width: 768px) {
    .realtime-statistics-switch {
      align-items: flex-start;
      padding: 16px;
    }

    .realtime-statistics-channel {
      width: 100%;
    }
  }
</style>
