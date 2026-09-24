import { computed, reactive } from 'vue';
import { dateEnUS, dateZhCN, enUS, zhCN } from 'naive-ui';
import Storage from '@/common/storage';
import { STORAGE } from '@/common/enum';
import { resources } from './resources';
import {
  DEFAULT_LOCALE,
  FALLBACK_LOCALE,
  normalizeSupportedLocale,
  resolveLocalePreference,
} from './locale-preference.mjs';

export { DEFAULT_LOCALE, FALLBACK_LOCALE };

export const LOCALE_OPTIONS = [
  { value: 'en-US', labelKey: 'common.language.english' },
  { value: 'zh-CN', labelKey: 'common.language.simplifiedChinese' },
];

const naiveLocaleMap = {
  'zh-CN': {
    locale: zhCN,
    dateLocale: dateZhCN,
  },
  'en-US': {
    locale: enUS,
    dateLocale: dateEnUS,
  },
};

const state = reactive({
  locale: DEFAULT_LOCALE,
  fallbackLocale: FALLBACK_LOCALE,
  isReady: false,
  source: 'default',
  lastError: '',
});

function isSupportedLocale(locale) {
  return Object.prototype.hasOwnProperty.call(resources, locale);
}

function getMessage(locale, key) {
  return key.split('.').reduce((result, current) => {
    if (result && typeof result === 'object' && current in result) {
      return result[current];
    }
    return undefined;
  }, resources[locale]);
}

function formatMessage(message, params = {}) {
  if (typeof message !== 'string') {
    return '';
  }
  return message.replace(/\{(\w+)\}/g, (_, key) => {
    if (key in params) {
      return `${params[key]}`;
    }
    return `{${key}}`;
  });
}

function storeLocale(locale) {
  Storage.set(STORAGE.LOCALE, locale);
}

export function initI18n() {
  const browserLocales =
    typeof navigator === 'undefined'
      ? []
      : navigator.languages?.length
        ? navigator.languages
        : [navigator.language || navigator.userLanguage];
  const preference = resolveLocalePreference(Storage.get(STORAGE.LOCALE), browserLocales);
  state.locale = preference.locale;
  state.source = preference.source;
  state.isReady = true;
  updateDocumentTitle();
  return state.locale;
}

export function setLocale(nextLocale) {
  const locale = normalizeSupportedLocale(nextLocale);
  if (!locale || !isSupportedLocale(locale)) {
    return false;
  }
  state.locale = locale;
  state.source = 'storage';
  state.lastError = '';
  state.isReady = true;
  storeLocale(locale);
  updateDocumentTitle();
  return true;
}

function updateDocumentTitle() {
  if (typeof document !== 'undefined') {
    document.title = t('common.documentTitle');
  }
}

export function getNaiveLocaleConfig(locale = state.locale) {
  return naiveLocaleMap[locale] || naiveLocaleMap[DEFAULT_LOCALE];
}

export function t(key, params = {}, fallback = '') {
  const currentMessage = getMessage(state.locale, key);
  if (typeof currentMessage === 'string') {
    return formatMessage(currentMessage, params);
  }

  const fallbackMessage = getMessage(state.fallbackLocale, key);
  if (typeof fallbackMessage === 'string') {
    return formatMessage(fallbackMessage, params);
  }

  if (fallback) {
    return formatMessage(fallback, params);
  }

  return getMessage(FALLBACK_LOCALE, 'common.feedback.missing') || '';
}

export function useI18n() {
  return {
    locale: computed(() => state.locale),
    localeOptions: LOCALE_OPTIONS,
    isReady: computed(() => state.isReady),
    t,
    setLocale,
  };
}

export function installI18n(app) {
  app.config.globalProperties.$t = t;
  app.config.globalProperties.$setLocale = setLocale;
}

export function getCurrentLocale() {
  return state.locale;
}

export function getLocaleState() {
  return state;
}
