let groups = [
  { 
    category: '群组相关', 
    isFold: false, 
    isActive: false,
    children: [
      { 
        url: '/groups/add',
        doc: 'group/groupcreate/',
        method: 'post', 
        name: '创建群组', 
        isActive: false,
        body: {
          "group_id":"groupid1",
          "group_name":"group1",
          "member_ids":["userid1","userid2"],
          "ext_fields":{
            "k1":"v1",
            "k2":"v2"
          }
        }
      },
      { 
        url: '/groups/del',
        doc: 'group/groupdissolve/',
        method: 'post', 
        name: '解散群组', 
        isActive: false,
        body: {
          "group_id": "groupid1"
        }
      },
      { 
        url: '/groups/update',
        doc: 'group/updategroup/',
        method: 'post', 
        name: '更新群信息', 
        isActive: false,
        body: {
          "group_id":"group1",
          "group_name":"group1",
          "group_portrait":"xxx",
          "ext_fields":{
            "field1":"aaa",
            "field2":"bbb"
          }
        }
      },
      { 
        url: '/groups/info?group_id=group1',
        doc: 'group/qrygroupinfo/',
        method: 'get', 
        name: '查询群信息', 
        isActive: false,
        body: {
          "//": "请修改上方 URL 参数进行调试"
        }
      },
      { 
        url: '/groups/groupmute/set',
        doc: 'group/groupmute/',
        method: 'post', 
        name: '设置群禁言', 
        isActive: false,
        body: {
          "group_id":"groupid1",
          "is_mute":1
        }
      },
      { 
        url: '/groups/groups/settings/set',
        doc: 'group/groupsetting/',
        method: 'post', 
        name: '更新群设置', 
        isActive: false,
        body: {
          "group_id":"groupid1",
          "settings":{
            "hide_grp_msg":"1"
          }
        }
      },
      { 
        url: '/groups/members/add',
        doc: 'group/groupaddmember/',
        method: 'post', 
        name: '添加群成员', 
        isActive: false,
        body: {
          "group_id":"groupid1",
          "group_name":"group1",
          "member_ids":["userid1","userid2"]
        }
      },
      { 
        url: '/groups/members/del',
        doc: 'group/groupdelmember/',
        method: 'post', 
        name: '移除群成员', 
        isActive: false,
        body: {
          "group_id":"groupid1",
          "member_ids":["userid1","userid2"]
        }
      },
      { 
        url: '/groups/members/query?group_id=group1&limit=50&offset=aabb',
        doc: 'group/qrygroupmember/',
        method: 'get', 
        name: '查询群成员', 
        isActive: false,
        body: {
          "//": "请修改上方 URL 参数进行调试"
        }
      },
      { 
        url: '/groups/groupmembermute/set',
        doc: 'group/groupmembermute/',
        method: 'post', 
        name: '设置群成员禁言', 
        isActive: false,
        body: {
          "group_id":"groupid1",
          "member_ids":["member1","member2"],
          "is_mute":1
        }
      },
      { 
        url: '/groups/groupmemberallow/set',
        doc: 'group/groupmemberallow/',
        method: 'post', 
        name: '设置群成员白名单', 
        isActive: false,
        body: {
          "group_id":"groupid1",
          "member_ids":["member1","member2"],
          "is_allow":1
        }
      },
      { 
        url: '/groups/members/querybyids',
        doc: 'group/qrygroupmemberbyids/',
        method: 'post', 
        name: '按成员 ID 查询信息', 
        isActive: false,
        body: {
          "group_id":"groupid1",
          "member_ids":["member1","member2"]
        }
      },
    ] 
  }
]

export { groups };