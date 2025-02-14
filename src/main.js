import { createApp } from 'vue';
import App from './App.vue';
import { setupRouter } from './router';
import VCalendar from 'v-calendar';
// import { create, NButton, NIcon, NUpload, NUploadDragger } from 'naive-ui';

import Toast from './components/toast';
import Chat from './components/chat';

import 'v-calendar/style.css';
import './assets/scss/iconfont.css';
import './assets/scss/style.scss';
import './assets/scss/custom/_header.css';
import './assets/scss/custom/_main.css';
import './assets/scss/custom/_login.css';
import './assets/scss/custom/_order.css';
import './assets/scss/custom/_dashboard.css';
import './assets/scss/custom/_invoice.css';
import './assets/scss/_custom.css';
import naive from 'naive-ui'

// if(location.search == '?debug'){
//   var vConsole = new window.VConsole();
// }

async function init() {
  // const naive = create({
  //   components: [
  //     NUpload,
  //     NUploadDragger,
  //     NButton,
  //     NIcon,
  //   ],
  // });

  const app = createApp(App);
  Toast.install(app);
  Chat.install(app);
  await setupRouter(app);
  app.use(VCalendar, {});
  app.use(naive)
  app.mount('#app');
}

init();
