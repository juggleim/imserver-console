let users = [
  { 
    category: '用户相关', 
    isFold: false, 
    isActive: false,
    children: [
      { 
        url: '/users/register',
        doc: 'user/register/',
        method: 'post', 
        name: '用户注册', 
        isActive: false,
        body: {
          "user_id": "userid1",
          "nickname": "nickname",
          "user_portrait": "https://portrait.example.com/avatar.png",
          "ext_fields":{
            "k1":"v1",
            "k2":"v2"
          }
        }
      },
      { 
        url: '/users/banusers/ban',
        doc: 'user/addbanuser/',
        method: 'post', 
        name: '封禁用户', 
        isActive: false,
        body: {
          "items":[
              {
               "user_id":"user2",
                "ban_type":1,
                "end_time_offset":300000,
                "scope_key":"platform",
                "scope_value":"iOS,Android",
                "ext":"aabbcc"
              }
            ]
        }
      },
      { 
        url: '/users/banusers/unban',
        doc: 'user/unbanuser/',
        method: 'post', 
        name: '解禁用户', 
        isActive: false,
        body: {
          "items":[
            {
              "user_id":"user1",
              "scope_key":"platform"
            },
            {
              "user_id":"user2"
            }
          ]
        }
      },
      { 
        url: '/users/banusers/query?limit=20&offset=""',
        doc: 'user/qrybanusers/',
        method: 'get', 
        name: '查询封禁列表', 
        isActive: false,
        body: {
          "//": "请修改上方 URL 参数进行调试"
        }
      },
      { 
        url: '/users/onlinestatus/query',
        doc: 'user/useronlinestatus/',
        method: 'post', 
        name: '查询在线状态', 
        isActive: false,
        body: {
          "user_ids":["userid1","userid2"]
        }
      },
      { 
        url: '/users/blockusers/block',
        doc: 'user/addblockuser/',
        method: 'post', 
        name: '设置单聊禁言', 
        isActive: false,
        body: {
          "user_id":"user1",
          "block_user_ids":["user2","user3"]
        }
      },
      { 
        url: '/users/blockusers/unblock',
        doc: 'user/unblockuser/',
        method: 'post', 
        name: '解除单聊禁言', 
        isActive: false,
        body: {
          "user_id":"user1",
          "block_user_ids":["user2","user3"]
        }
      },
      { 
        url: '/users/blockusers/query?user_id=user1&limit=20&offset=""',
        doc: 'user/qryblockusers/',
        method: 'get', 
        name: '单聊禁言列表', 
        isActive: false,
        body: {
          "//": "请修改上方 URL 参数进行调试"
        }
      },
      { 
        url: '/users/update',
        doc: 'user/updateuser/',
        method: 'post', 
        name: '更新用户信息', 
        isActive: false,
        body: {
          "user_id":"userid1",
          "nickname":"user1",
          "user_portrait":"xxxxx",
          "ext_fields":{
            "field1":"aaa",
            "field2":"bbb"
          }
        }
      },
      { 
        url: '/users/info?user_id=userid1',
        doc: 'user/qryuserinfo/',
        method: 'get', 
        name: '查询用户信息', 
        isActive: false,
        body: {
          "//": "请修改上方 URL 参数进行调试"
        }
      },
      { 
        url: '/users/kick',
        doc: 'user/kickuser/',
        method: 'post', 
        name: '踢用户下线', 
        isActive: false,
        body: {
          "user_id":"userid1",
          "platforms":["iOS","Android"],
          "device_ids":["xxxxx","yyyyy"]
        }
      },
    ] 
  }
]

export { users };