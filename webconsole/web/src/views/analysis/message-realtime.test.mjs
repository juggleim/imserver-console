import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import {
  ONE_HOUR_MS,
  THREE_DAYS_MS,
  alignRealtimeBucket,
  fillRealtimeSeries,
  getRealtimeBucketMs,
  isRealtimeRangeValid,
} from './message-realtime.mjs';

test('selects the same realtime buckets as the backend', () => {
  assert.equal(getRealtimeBucketMs(ONE_HOUR_MS), 30_000);
  assert.equal(getRealtimeBucketMs(ONE_HOUR_MS + 1), 60_000);
  assert.equal(getRealtimeBucketMs(6 * ONE_HOUR_MS + 1), 300_000);
  assert.equal(getRealtimeBucketMs(12 * ONE_HOUR_MS + 1), 600_000);
  assert.equal(getRealtimeBucketMs(24 * ONE_HOUR_MS + 1), 900_000);
});

test('accepts positive ranges up to three days', () => {
  assert.equal(isRealtimeRangeValid(1_000, 1_001), true);
  assert.equal(isRealtimeRangeValid(1_000, 1_000 + THREE_DAYS_MS), true);
  assert.equal(isRealtimeRangeValid(1_000, 1_000), false);
  assert.equal(isRealtimeRangeValid(1_000, 1_000 + THREE_DAYS_MS + 1), false);
});

test('aligns data and fills missing realtime buckets with zero', () => {
  assert.equal(alignRealtimeBucket(61_000, 30_000), 60_000);
  assert.deepEqual(
    fillRealtimeSeries(
      [
        { time_mark: 61_000, count: 1.25 },
        { time_mark: 121_000, count: 3.5 },
      ],
      60_000,
      150_000,
      30_000
    ),
    [
      [60_000, 1.25],
      [90_000, 0],
      [120_000, 3.5],
      [150_000, 0],
    ]
  );
});

test('registers the realtime route, API and bilingual copy', () => {
  const route = readFileSync(new URL('../../router/index.js', import.meta.url), 'utf8');
  const api = readFileSync(new URL('../../services/api.js', import.meta.url), 'utf8');
  assert.match(route, /AnalysisMessageRealtime/);
  assert.match(api, /apps\/statistic\/msgrealtime/);

  for (const locale of ['en-US', 'zh-CN']) {
    const analysis = JSON.parse(
      readFileSync(new URL(`../../locales/${locale}/analysis.json`, import.meta.url), 'utf8')
    );
    const menu = JSON.parse(
      readFileSync(new URL(`../../locales/${locale}/menu.json`, import.meta.url), 'utf8')
    );
    assert.ok(analysis.realtime.enable);
    assert.ok(analysis.realtime.unitPerSecond);
    assert.ok(menu.analytics.messageRealtime);
  }
});

test('guards stale requests and recreates the chart after toggling realtime statistics', () => {
  const view = readFileSync(new URL('./message-realtime.vue', import.meta.url), 'utf8');
  assert.match(view, /const requestId = \+\+latestRequestId/);
  assert.match(view, /requestId !== latestRequestId \|\| !isRealtimeStatEnabled\.value/);
  assert.match(view, /else if \(value === 0\) \{[\s\S]*disposeChart\(\)/);
  assert.match(view, /if \(value === 1 && isReady\) \{[\s\S]*await nextTick\(\)/);
  assert.match(view, /watch\(locale/);
});
