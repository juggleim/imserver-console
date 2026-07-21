export const DEFAULT_LOCALE = 'en-US';
export const FALLBACK_LOCALE = 'en-US';

export function normalizeSupportedLocale(input) {
  if (typeof input !== 'string') {
    return null;
  }
  const locale = input.trim().toLowerCase();
  if (locale.startsWith('en')) {
    return 'en-US';
  }
  if (locale.startsWith('zh')) {
    return 'zh-CN';
  }
  return null;
}

function resolveBrowserLocale(browserLocales) {
  const candidates = Array.isArray(browserLocales) ? browserLocales : [browserLocales];
  for (const candidate of candidates) {
    const locale = normalizeSupportedLocale(candidate);
    if (locale) {
      return locale;
    }
  }
  return null;
}

export function resolveLocalePreference(storedLocale, browserLocales) {
  const locale = normalizeSupportedLocale(storedLocale);
  if (locale) {
    return { locale, source: 'storage' };
  }
  const browserLocale = resolveBrowserLocale(browserLocales);
  if (browserLocale) {
    return { locale: browserLocale, source: 'browser' };
  }
  return { locale: DEFAULT_LOCALE, source: 'default' };
}
