import { getCurrentInstance } from 'vue';

export function showToast({ title = '', text = '', icon = 'success' }) {
  const context = getCurrentInstance();
  window.$toast({ icon, text, title });
}
