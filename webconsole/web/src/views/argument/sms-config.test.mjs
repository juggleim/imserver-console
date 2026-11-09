import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';

const enumSource = readFileSync(new URL('../../common/enum.js', import.meta.url), 'utf8');

test('SMS settings expose Aliyun fields and keep SMSBao available', () => {
  assert.match(enumSource, /uid:\s*'smsbao'/);
  assert.match(enumSource, /uid:\s*'aliyun'/);
  for (const key of [
    'access_key',
    'access_secret',
    'region_id',
    'endpoint',
    'sign_name',
    'template_code',
    'template_code_key',
  ]) {
    assert.match(enumSource, new RegExp(`key:\\s*'${key}'`), `missing Aliyun field ${key}`);
  }
  assert.match(enumSource, /defaultValue:\s*'cn-hangzhou'/);
  assert.match(enumSource, /defaultValue:\s*'dysmsapi\.aliyuncs\.com'/);
  assert.match(enumSource, /defaultValue:\s*'code'/);
});

test('Aliyun SMS labels exist in both locales', () => {
  for (const locale of ['zh-CN', 'en-US']) {
    const text = readFileSync(
      new URL(`../../locales/${locale}/smsConfig.json`, import.meta.url),
      'utf8'
    );
    const config = JSON.parse(text);
    assert.ok(config.provider.aliyun);
    for (const key of [
      'access_key',
      'access_secret',
      'region_id',
      'sign_name',
      'template_code',
      'template_code_key',
    ]) {
      assert.ok(config.field[key], `${locale} is missing ${key}`);
    }
  }
});
