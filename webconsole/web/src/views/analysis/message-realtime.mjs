export const REALTIME_STAT_SWITCH_KEY = 'open_real_time_msg_statistic';
export const ONE_HOUR_MS = 60 * 60 * 1000;
export const THREE_DAYS_MS = 3 * 24 * ONE_HOUR_MS;

const SIX_HOURS_MS = 6 * ONE_HOUR_MS;
const TWELVE_HOURS_MS = 12 * ONE_HOUR_MS;
const ONE_DAY_MS = 24 * ONE_HOUR_MS;

export const REALTIME_RANGES = [
  { titleKey: 'analysis.realtime.range.last15Minutes', durationMs: 15 * 60 * 1000 },
  { titleKey: 'analysis.realtime.range.last1Hour', durationMs: ONE_HOUR_MS },
  { titleKey: 'analysis.realtime.range.last6Hours', durationMs: SIX_HOURS_MS },
  { titleKey: 'analysis.realtime.range.last12Hours', durationMs: TWELVE_HOURS_MS },
  { titleKey: 'analysis.realtime.range.last1Day', durationMs: ONE_DAY_MS },
  { titleKey: 'analysis.realtime.range.last3Days', durationMs: THREE_DAYS_MS },
];

export function getRealtimeBucketMs(rangeMs) {
  if (rangeMs <= ONE_HOUR_MS) return 30 * 1000;
  if (rangeMs <= SIX_HOURS_MS) return 60 * 1000;
  if (rangeMs <= TWELVE_HOURS_MS) return 5 * 60 * 1000;
  if (rangeMs <= ONE_DAY_MS) return 10 * 60 * 1000;
  return 15 * 60 * 1000;
}

export function createRealtimeRange(durationMs = ONE_HOUR_MS, now = Date.now()) {
  return {
    start: new Date(now - durationMs),
    end: new Date(now),
  };
}

export function isRealtimeRangeValid(start, end) {
  const rangeMs = end - start;
  return rangeMs > 0 && rangeMs <= THREE_DAYS_MS;
}

export function alignRealtimeBucket(timeMark, bucketMs) {
  return Math.floor(Number(timeMark) / bucketMs) * bucketMs;
}

export function fillRealtimeSeries(items, rangeStart, rangeEnd, bucketMs) {
  const valueMap = new Map();
  (items || []).forEach((item) => {
    valueMap.set(alignRealtimeBucket(item.time_mark, bucketMs), Number(item.count) || 0);
  });

  const points = [];
  for (
    let timeMark = Math.ceil(rangeStart / bucketMs) * bucketMs;
    timeMark <= rangeEnd;
    timeMark += bucketMs
  ) {
    points.push([timeMark, valueMap.has(timeMark) ? valueMap.get(timeMark) : 0]);
  }
  return points;
}
