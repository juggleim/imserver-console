import utils from "./utils";

export let INVOICE_TYPE = {
  ONLINE: 1,
  PAPER: 2
};

let ErrorMessages = [
  { code: 200, key: 'errors.api.SUCCESS', name: 'SUCCESS' },
  { code: 0, key: 'errors.api.SUCCESS_0', name: 'SUCCESS_0' },
  { code: 1003, key: 'errors.api.USER_LOGIN_FAILED', name: 'USER_LOGIN_FAILED' },
  { code: 1016, key: 'errors.api.USER_BLOCKED', name: 'USER_BLOCKED' },
  { code: 1012, key: 'errors.api.USER_EXISTS', name: 'USER_EXISTS' },
  { code: 1001, key: 'errors.api.USER_OLDPWD_WRONG', name: 'USER_OLDPWD_WRONG' },
  { code: 1006, key: 'errors.api.APP_EXISTS', name: 'APP_EXISTS' },
  { code: 1001, key: 'errors.api.USER_TOKEN_EXPIRE', name: 'USER_TOKEN_EXPIRE' },
];

export let Errors = ErrorMessages;

function getErrorType() {
  let errors = {};
  utils.forEach(ErrorMessages, (error) => {
    let { name, code, key } = error;
    errors[name] = { code, key };
  });
  return errors;
}
export let ErrorType = getErrorType();

export let STORAGE = {
  PREFIX: 'jgadmin',
  USER_TOKEN: 'user_auth_token',
  APP_KEY: 'app_key',
  LOCALE: 'locale',
}
export let USER_STATE = {
  ENABLE: 0,
  DISABLE: 1
};
export let DEPLOY_TYPE = {
  PUBLIC: 1,
  PRIVATE: 2
};
export let APP_TYPE = {
  PRIVATE: 0,
  SPECIAL: 1,
  PUBLIC: 2,
};
export let EVENT_NAME = {
  CREATE_APP: 'app_create',
  CREATE_APP_FINISHED: 'app_create_finished',

  ON_LOGOUT: 'app_logout',
};
export let FUNC_TYPE = {
  INPUT: 'input',
  INPUT_TEXT: 'input_text',
  INPUT_MODAL: 'input_modal',
  SELECT: 'select',
  SWITCH: 'switch',
  UPLOAD_FILE: 'upload_file',
};
export let APP_STATUS = {
  NORMAL: 0,
  ONLINE: 1,
  BLOCKED: 2,
};
export let RESPONSE = {
  SUCCESS: 0,
  UNATHORIZED: 401,
  PUSH_CONF_EXISTED: 1019,
  SERVER_LOGS_FAILED: 1021,
  REQUEST_LIMIT: 1022,
  USER_SENDCODE_ERROR: 10514,
  USER_SEARCH_ERROR: 10509,
};

export let DIALOG_TYPE = {
  UPGRADE: 1,
  DELAY: 2,
  CREATE: 3
};
export let LOG_PULL_STATUS = {
  NONE: 0,
  SUCCESS: 1,
  FAIL: 2,
  COMPLETE: 3,
  UPLOAD_FAIL: 4,
  WITHOUT: 5,
};

export let SIGNAL_TYPE = {
  CONNECTED: 1,
  DISCONNECTED: 2,
  USER: 3,
  SERVER: 4,
  REPLY: 5,
}

export let STAT_TYPE = {
  UP: 1,
  DISPATCH: 2,
  DOWN: 3,
};

export let METHOD_MAP = {
  connect: { titleKey: 'tools.connection.method.connect' },
  disconnect: { titleKey: 'tools.connection.method.disconnect' },
  qry_ack: { title: 'Q_ACK' },
  u_pub_ack: { title: 'U_PUB_ACK' },
  s_pub_ack: { title: 'S_PUB_ACK' },
  qry_hismsgs: { titleKey: 'tools.connection.method.qry_hismsgs' },
  qry_convers: { titleKey: 'tools.connection.method.qry_convers' },
  qry_top_convers: { titleKey: 'tools.connection.method.qry_top_convers' },
  sync_convers: { titleKey: 'tools.connection.method.sync_convers' },
  sync_msgs: { titleKey: 'tools.connection.method.sync_msgs' },
  recall_msg: { titleKey: 'tools.connection.method.recall_msg' },
  qry_mention_msgs: { titleKey: 'tools.connection.method.qry_mention_msgs' },
  ntf: { titleKey: 'tools.connection.method.ntf' },
  msg: { titleKey: 'tools.connection.method.msg' },
  c_user_ntf: { titleKey: 'tools.connection.method.c_user_ntf' },
  g_msg: { titleKey: 'tools.connection.method.g_msg' },
  p_msg: { titleKey: 'tools.connection.method.p_msg' },
  c_msg: { titleKey: 'tools.connection.method.c_msg' },
  qry_merged_msgs: { titleKey: 'tools.connection.method.qry_merged_msgs' },
  qry_first_unread_msg: { titleKey: 'tools.connection.method.qry_first_unread_msg' },
  clear_unread: { titleKey: 'tools.connection.method.clear_unread' },
  del_convers: { titleKey: 'tools.connection.method.del_convers' },
  add_conver: { titleKey: 'tools.connection.method.add_conver' },
  qry_conver: { titleKey: 'tools.connection.method.qry_conver' },
  undisturb_convers: { titleKey: 'tools.connection.method.undisturb_convers' },
  top_convers: { titleKey: 'tools.connection.method.top_convers' },
  qry_total_unread_count: { titleKey: 'tools.connection.method.qry_total_unread_count' },
  clear_total_unread: { titleKey: 'tools.connection.method.clear_total_unread' },
  mark_unread: { titleKey: 'tools.connection.method.mark_unread' },
  set_user_undisturb: { titleKey: 'tools.connection.method.set_user_undisturb' },
  get_user_undisturb: { titleKey: 'tools.connection.method.get_user_undisturb' },
  mark_read: { titleKey: 'tools.connection.method.mark_read' },
  qry_read_detail: { titleKey: 'tools.connection.method.qry_read_detail' },
  modify_msg: { titleKey: 'tools.connection.method.modify_msg' },
  clean_hismsg: { titleKey: 'tools.connection.method.clean_hismsg' },
  del_hismsg: { titleKey: 'tools.connection.method.del_hismsg' },
  qry_hismsg_by_ids: { titleKey: 'tools.connection.method.qry_hismsg_by_ids' },
  file_cred: { titleKey: 'tools.connection.method.file_cred' },
  qry_user_info: { titleKey: 'tools.connection.method.qry_user_info' },
  c_join: { titleKey: 'tools.connection.method.c_join' },
  c_quit: { titleKey: 'tools.connection.method.c_quit' },
  c_sync_msgs: { titleKey: 'tools.connection.method.c_sync_msgs' },
  c_sync_atts: { titleKey: 'tools.connection.method.c_sync_atts' },
  c_batch_add_att: { titleKey: 'tools.connection.method.c_batch_add_att' },
  c_batch_del_att: { titleKey: 'tools.connection.method.c_batch_del_att' },
};
export let IM_ERRORS = [
{ code: 0, key: 'errors.im.0' },
{ code: 11000, key: 'errors.im.11000' },
{ code: 11001, key: 'errors.im.11001' },
{ code: 11002, key: 'errors.im.11002' },
{ code: 11003, key: 'errors.im.11003' },
{ code: 11004, key: 'errors.im.11004' },
{ code: 11005, key: 'errors.im.11005' },
{ code: 11006, key: 'errors.im.11006' },
{ code: 11007, key: 'errors.im.11007' },
{ code: 11008, key: 'errors.im.11008' },
{ code: 11009, key: 'errors.im.11009' },
{ code: 11010, key: 'errors.im.11010' },
{ code: 11011, key: 'errors.im.11011' },
{ code: 11012, key: 'errors.im.11012' },
{ code: 11013, key: 'errors.im.11013' },
{ code: 11014, key: 'errors.im.11014' },
{ code: 11015, key: 'errors.im.11015' },
{ code: 10100, key: 'errors.im.10100' },
{ code: 10101, key: 'errors.im.10101' },
{ code: 10102, key: 'errors.im.10102' },
{ code: 10103, key: 'errors.im.10103' },
{ code: 12000, key: 'errors.im.12000' },
{ code: 12001, key: 'errors.im.12001' },
{ code: 12002, key: 'errors.im.12002' },
{ code: 12003, key: 'errors.im.12003' },
{ code: 12004, key: 'errors.im.12004' },
{ code: 12005, key: 'errors.im.12005' },
{ code: 12006, key: 'errors.im.12006' },
{ code: 12007, key: 'errors.im.12007' },
{ code: 13000, key: 'errors.im.13000' },
{ code: 13001, key: 'errors.im.13001' },
{ code: 13002, key: 'errors.im.13002' },
{ code: 13003, key: 'errors.im.13003' },
{ code: 13004, key: 'errors.im.13004' },
{ code: 13005, key: 'errors.im.13005' },
{ code: 13006, key: 'errors.im.13006' },
{ code: 14000, key: 'errors.im.14000' },
{ code: 14001, key: 'errors.im.14001' },
{ code: 14002, key: 'errors.im.14002' },
{ code: 14003, key: 'errors.im.14003' },
{ code: 14004, key: 'errors.im.14004' },
{ code: 14005, key: 'errors.im.14005' },
{ code: 14006, key: 'errors.im.14006' },
{ code: 14007, key: 'errors.im.14007' },
{ code: 14008, key: 'errors.im.14008' },
{ code: 15000, key: 'errors.im.15000' },
{ code: 15001, key: 'errors.im.15001' },
]

export let CONVERSATION_TYPE = { 
  PRIVATE: 1,
  GROUP: 2,
  CHATROOM: 3
}

export let ANA_DATE_RANGES = [
  { titleKey: 'analysis.range.7days', name: 7, isActive: true },
  { titleKey: 'analysis.range.14days', name: 14, isActive: false },
  { titleKey: 'analysis.range.30days', name: 30, isActive: false },
]

export let PLATFORMAS = [
  { name: 'Android', labelKey: '', value: 'Android' },
  { name: 'iOS', labelKey: '', value: 'iOS' },
  { name: 'Web', labelKey: '', value: 'Web' },
  { name: 'PC', labelKey: '', value: 'PC' }, 
];

export let MENU_UID = {
  APP: 1,
  SENTSIVE: 2,
  ANALYSE: 3,
  LOG: 4,
  DEV_TOOL: 5,
  USER_MANGER: 6,
  
};


export let R3D_USE_TYPE = {
  ENABLE: 1,
  DISABLE: 2
};
export let TRANSLATE_CHANNELS = [
  { uid: 'baidu', icon: 'baidu', isUsed: false,
    children: [
      { type: 'text', key: 'api_key', value: '' },
      { type: 'text', key: 'secret_key', value: '', secretValue: '**************' },
      { type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
        { value: R3D_USE_TYPE.ENABLE },
        { value: R3D_USE_TYPE.DISABLE },
      ] },
    ] 
  },
  { uid: 'deepl', icon: 'deepl', isUsed: false,
    children: [
      { type: 'text', key: 'auth_key', value: '' },
      { type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
        { value: R3D_USE_TYPE.ENABLE },
        { value: R3D_USE_TYPE.DISABLE },
      ] },
    ] 
  },
];

export let RTC_CHANNELS = [
  { uid: 'zego_conf', icon: 'zego', isUsed: false,
    children: [
      { type: 'number', key: 'app_id', value: '' },
      { type: 'text', key: 'secret', value: '', secretValue: '**************' },
      // { name: '是否启用', type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
      //   { label: '启用', value: R3D_USE_TYPE.ENABLE },
      //   { label: '禁用', value: R3D_USE_TYPE.DISABLE },
      // ] },
    ] 
  },
  { uid: 'agora_conf', icon: 'agora', isUsed: false,
    children: [
      { type: 'text', key: 'app_id', value: '' },
      { type: 'text', key: 'app_certificate', value: '', secretValue: '**************' },
      // { name: '是否启用', type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
      //   { label: '启用', value: R3D_USE_TYPE.ENABLE },
      //   { label: '禁用', value: R3D_USE_TYPE.DISABLE },
      // ] },
    ] 
  },
  { uid: 'livekit_conf', icon: 'livekit', isUsed: false,
    children: [
      { type: 'text', key: 'app_key', value: '' },
      { type: 'text', key: 'app_secret', value: '', secretValue: '**************' },
      { type: 'text', key: 'service_url', value: '' },
      // { name: '是否启用', type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
      //   { label: '启用', value: R3D_USE_TYPE.ENABLE },
      //   { label: '禁用', value: R3D_USE_TYPE.DISABLE },
      // ] },
    ] 
  }
];
export let EMAIL_CHANNELS = [
  { uid: 'ali', icon: 'ali', isUsed: false,
    children: [
      { type: 'text', key: 'access_key', value: '' },
      { type: 'text', key: 'access_secret', value: '', secretValue: '**************' },
      { type: 'text', key: 'from_email', value: '' },
      { type: 'text', key: 'from_alias', value: '' },
      { type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
        { value: R3D_USE_TYPE.ENABLE },
        { value: R3D_USE_TYPE.DISABLE },
      ] },
    ] 
  },
  { uid: 'engagelab', icon: 'engagelab', isUsed: false,
    children: [
      { type: 'text', key: 'url', value: '' },
      { type: 'text', key: 'api_user', value: '' },
      { type: 'text', key: 'api_key', value: '', secretValue: '**************' },
      { type: 'text', key: 'from_email', value: '' },
      { type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
        { value: R3D_USE_TYPE.ENABLE },
        { value: R3D_USE_TYPE.DISABLE },
      ] },
    ] 
  },
];
export let SMS_CHANNELS = [
  { uid: 'baidu', icon: 'baidu', isUsed: false,
    children: [
      { type: 'text', key: 'api_key', value: '' },
      { type: 'text', key: 'secret_key', value: '', secretValue: '**************' },
      { type: 'text', key: 'endpoint', value: '' },
      { type: 'text', key: 'template', value: '' },
      { type: 'text', key: 'signature_id', value: '' },
      { type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
        { value: R3D_USE_TYPE.ENABLE },
        { value: R3D_USE_TYPE.DISABLE },
      ] },
    ] 
  },
  { uid: 'smsbao', icon: 'smsbao', isUsed: false,
    children: [
      { type: 'text', key: 'username', value: '' },
      { type: 'text', key: 'password', value: '', secretValue: '**************' },
      { type: 'text', key: 'template', value: '' },
      { type: 'select', key: 'is_used', value: R3D_USE_TYPE.DISABLE, children: [
        { value: R3D_USE_TYPE.ENABLE },
        { value: R3D_USE_TYPE.DISABLE },
      ] },
    ] 
  },
];
export let USER_ROLE_TYPE = {
  ADMIN: 0,
  USER: 1,
}
export let ROLES = [
  { labelKey: 'userManager.role.admin', value: USER_ROLE_TYPE.ADMIN, menuIds: [ MENU_UID.APP, MENU_UID.SENTSIVE, MENU_UID.ANALYSE, MENU_UID.LOG, MENU_UID.DEV_TOOL, MENU_UID.USER_MANGER] },
  { labelKey: 'userManager.role.user', value: USER_ROLE_TYPE.USER, menuIds: [MENU_UID.ANALYSE] },
];
