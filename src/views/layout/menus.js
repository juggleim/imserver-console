import utils from '../../common/utils';
import Storage from '../../common/storage';
import { STORAGE, USER_ROLE, ROLES, MENU_UID } from '../../common/enum';
import appTools from '../../common/app-tools';

function MenuFactory() {
  let homePage = {
    title: '首页导航',
    name: 'Dashboard',
    icon: 'cicon-home',
  };
  let isHidden = false;

  let privateMenus = {
    getApps: ({ router }) => {
      return [
        homePage,
        {
          title: '应用配置',
          icon: 'cicon-set',
          isHidden: isHidenMenu(MENU_UID.APP),
          isUnfold: isFold('Argu', router),
          children: [
            { title: '基本信息', name: 'ArguBase' },
            { title: '功能配置', name: 'ArguSwitch' },
            { title: '事件设置', name: 'ArguCallback' },
            { title: '推送设置', name: 'ArguPush' },
            { title: '存储配置', name: 'ArguStorage' },
            { title: '翻译配置', name: 'ArguTranslte' },
            { title: '音视频配置', name: 'ArguRTC' },
          ],
        },
        {
          title: '敏感词管理',
          icon: 'cicon-set',
          isHidden: isHidenMenu(MENU_UID.SENTSIVE),
          isUnfold: isFold('sensitive', router),
          children: [
            { title: '敏感词配置', name: 'sensitiveConfig' },
          ],
        },
        {
          title: '数据统计',
          icon: 'cicon-analysis',
          isHidden: isHidenMenu(MENU_UID.ANALYSE),
          isUnfold: isFold('Analysis', router),
          children: [
            { title: '用户统计【日活】', name: 'AnalysisUser' },
            { title: '消息统计【单聊】', name: 'AnalysisMessage' },
            { title: '消息统计【群聊】', name: 'AnalysisGroup' },
            { title: '消息统计【聊天室】', name: 'AnalysisChatroom' },
          ],
        },
        {
          title: '日志管理',
          icon: 'cicon-logs',
          isHidden: isHidenMenu(MENU_UID.LOG),
          isUnfold: isFold('Logs', router),
          children: [
            { title: '日志列表', name: 'Logs' },
          ]
        },
        {
          title: '开发工具',
          icon: 'cicon-dev',
          isHidden: isHidenMenu(MENU_UID.DEV_TOOL),
          isUnfold: isFold('Tools', router),
          children: [
            { title: 'API 调试', name: 'ToolsAPI' },
            { title: '连接排查', name: 'ToolsConnection' },
          ]
        },
      ];
    },
    getUsers: ({ router }) => {
      return [
        {
          title: '账户信息',
          icon: 'cicon-user',
          isUnfold: isFold('User', router),
          children: [
            { title: '账户设置', name: 'UserSetting' },
            { title: '用户管理', name: 'UserManader', isHidden: isHidenMenu(MENU_UID.USER_MANGER), },
          ],
        },
      ];
    },
  };

  function isFold(type, router) {
    if (utils.isUndefined(router)) {
      return false;
    }
    let {
      currentRoute: {
        _rawValue: { name },
      },
    } = router;
    return utils.isInclude(name, type);
  }

  function isHidenMenu(menuId){
    let user = Storage.get(STORAGE.USER_TOKEN);
    let roleId = user.role_id || 0;
    let role = ROLES.find((_role) => { return utils.isEqual(roleId, _role.value);}) || { menuIds: [] }
    let index = utils.find(role.menuIds, (_menuId) => {
      return utils.isEqual(menuId, _menuId);
    })
    return index == -1;
  }

  return {
    private: privateMenus,
  };
}

export default MenuFactory;
