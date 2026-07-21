import test from 'node:test';
import assert from 'node:assert/strict';
import {
  DEFAULT_LOCALE,
  FALLBACK_LOCALE,
  normalizeSupportedLocale,
  resolveLocalePreference,
} from './locale-preference.mjs';

test('supported browser language is used when no language was saved', () => {
  assert.equal(DEFAULT_LOCALE, 'en-US');
  assert.equal(FALLBACK_LOCALE, 'en-US');
  assert.deepEqual(resolveLocalePreference(undefined, 'zh-CN'), {
    locale: 'zh-CN',
    source: 'browser',
  });
  assert.deepEqual(resolveLocalePreference(undefined, ['fr-FR', 'en-GB']), {
    locale: 'en-US',
    source: 'browser',
  });
});

test('unsupported browser language falls back to English', () => {
  assert.deepEqual(resolveLocalePreference(undefined, 'fr-FR'), {
    locale: 'en-US',
    source: 'default',
  });
});

test('a saved language is restored before considering the browser language', () => {
  assert.deepEqual(resolveLocalePreference('zh-CN', 'en-US'), {
    locale: 'zh-CN',
    source: 'storage',
  });
  assert.deepEqual(resolveLocalePreference('en-US', 'zh-CN'), {
    locale: 'en-US',
    source: 'storage',
  });
});

test('invalid saved values fall back to English', () => {
  assert.equal(normalizeSupportedLocale('unsupported'), null);
  assert.deepEqual(resolveLocalePreference('unsupported', 'unsupported'), {
    locale: 'en-US',
    source: 'default',
  });
});
