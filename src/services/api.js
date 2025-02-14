import { CONFIG } from './config'
import utils from '../common/utils'
let SERVER_PATH = {
  USER_LOGIN: 'login',
  USER_ADD: 'accounts/add',
  USER_DELETE: 'accounts/delete',
  USER_DISABLE: 'accounts/disable',
  USER_UPDATE_PWD: 'accounts/updpass',
  USER_LIST: 'accounts/list',

  APP_CREATE: 'apps/create',
  APP_IMPORT: 'apps/active',
  APP_UPDATE_LICENSE: 'apps/updlicense',
  APP_SETTING_UPDATE: 'apps/configs/set',
  APP_SETTING_GET: 'apps/configs/get',
  APP_GET_LIST: 'apps/list',
  APP_GET_ONE: 'apps/info',
  
  APP_EVENT_SET: 'apps/eventsubconfig/set',
  APP_EVENT_GET: 'apps/eventsubconfig/get',
  APP_ANDROID_UPDATE: 'apps/androidpushconf/set',
  APP_ANDROID_GET: 'apps/androidpushconf/get',

  APP_IOS_UPDATE: 'apps/iospushcer/set',
  APP_IOS_UPLOAD: 'apps/iospushcer/upload',
  APP_IOS_GET: 'apps/iospushcer/get',

  APP_FCM_GET: 'apps/fcmpushconf/get',
  APP_FCM_UPLOAD: 'apps/fcmpushconf/upload',

  UPLOAD_FILE: 'common/upload',

  STORAGE_GET: 'apps/fileconf/get',
  STORAGE_SET: 'apps/fileconf/set',
  ENABLE_STORAGE_GET: 'apps/fileconf/switch/get',
  ENABLE_STORAGE_SET: 'apps/fileconf/switch/set',

  APP_SENSITIVE_SET: 'apps/sensitiveconf/set',
  APP_SENSITIVE_GET: 'apps/sensitiveconf/get',
  APP_SENSITIVE_WORDS_GET: 'apps/sensitivewords/list',
  APP_SENSITIVE_WORDS_ADD: 'apps/sensitivewords/add',
  APP_SENSITIVE_WORDS_UPDATE: 'apps/sensitivewords/update',
  APP_SENSITIVE_WORDS_DELETE: 'apps/sensitivewords/delete',
  APP_SENSITIVE_WORDS_IMPORT: 'apps/sensitivewords/import',

  MISC_REQUEST: 'imapiagent',

  LOG_GET_LIST: 'apps/clientlogs/list',
  LOG_CREATE_PULL: 'apps/clientlogs/notify',
  LOG_DOWNLOAD: 'apps/clientlogs/download',
  
  CONN_GET_LIST: 'apps/serverlogs/userconnect',
  CONN_GET_ONE: 'apps/serverlogs/connect',
  
  R3D_TRANSLTE_GET: 'apps/translate/get',
  R3D_TRANSLTE_SET: 'apps/translate/set',

  R3D_RTC_GET: 'apps/rtcconf/get',
  R3D_RTC_SET: 'apps/rtcconf/set',
  
  ANALIYSIS_MESSAGE: 'apps/statistic/msg',
  ANALIYSIS_DAU: 'apps/statistic/useractivity',
};
utils.forEach(SERVER_PATH, (url, name) => {
  SERVER_PATH[name] = `${CONFIG.API}${url}`;
});

export default SERVER_PATH;