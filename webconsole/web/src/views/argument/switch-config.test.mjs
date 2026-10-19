import test from 'node:test';
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';

const view = readFileSync(new URL('./switch.vue', import.meta.url), 'utf8');
const switchComponent = readFileSync(
  new URL('../../components/func/switch.vue', import.meta.url),
  'utf8'
);

test('应用相关配置包含默认关闭的 open_e2ee 开关', () => {
  const definition = view.match(/let settings = ([\s\S]*?);\s*let state/)[1];
  const settings = new Function('FUNC_TYPE', `return (${definition});`)({
    INPUT_MODAL: 'input_modal',
  });
  const item = settings
    .find((section) => section.type === 'app')
    .list.find((entry) => entry.id === 'open_e2ee');
  assert.equal(item.type, 'switch');
  assert.equal(item.value, false);
  assert.equal(item.labelKey, 'switchConfig.item.open_e2ee');
  assert.equal(item.checkedValue, 'true');
  assert.equal(item.uncheckedValue, 'false');
});

test('开关组件支持 true/false 配置值，同时兼容原有 0/1 配置值', () => {
  const script = switchComponent
    .match(/<script setup>([\s\S]*?)<\/script>/)[1]
    .replace(/^\s*import .*;$/gm, '');
  for (const item of [
    { id: 'open_e2ee', value: false, checkedValue: 'true', uncheckedValue: 'false' },
    { id: 'kick_mode', value: 0 },
  ]) {
    const saves = [];
    let update;
    const component = new Function(
      'defineProps',
      'defineEmits',
      'reactive',
      'watch',
      'utils',
      `${script}; return { state, onSave };`
    )(
      () => ({ item }),
      () => (event, payload) => saves.push(payload),
      (value) => value,
      (source, callback) => {
        update = callback;
      },
      { extend: Object.assign }
    );
    assert.equal(component.state.value, false);
    for (const value of ['', undefined, null, false, 'false', '0']) {
      update(value);
      assert.equal(component.state.value, false);
    }
    for (const value of item.id === 'open_e2ee' ? [true, 'true'] : [true, 1, '1']) {
      update(value);
      assert.equal(component.state.value, true);
    }
    component.onSave({ target: { checked: true } });
    component.onSave({ target: { checked: false } });
    assert.deepEqual(saves, [
      { id: item.id, value: item.checkedValue ?? '1' },
      { id: item.id, value: item.uncheckedValue ?? '0' },
    ]);
  }
});

test('open_e2ee 提供中英文文案', () => {
  for (const locale of ['zh-CN', 'en-US']) {
    const messages = JSON.parse(
      readFileSync(new URL(`../../locales/${locale}/switchConfig.json`, import.meta.url), 'utf8')
    );
    assert.ok(messages.item.open_e2ee);
  }
});
