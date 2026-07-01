let chatrooms = [
  {
    category: '聊天室相关',
    isFold: false,
    isActive: false,
    children: [{
      url: '/chatrooms/create',
      doc: 'chatroom/createchatroom/',
      method: 'post',
      name: '创建聊天室',
      isActive: false,
      body: {
        "chat_id": "chatroom1000",
        "chat_name": "chatroom1000"
      }
    },
    {
      url: '/chatrooms/destroy',
      doc: 'chatroom/destroychatroom/',
      method: 'post',
      name: '销毁聊天室',
      isActive: false,
      body: {
        "chat_id": "chatroom1000"
      }
    },
    {
      url: '/chatrooms/info?chat_id=chatroom1000&order=0&count=100',
      doc: 'chatroom/qrychrminfo/',
      method: 'get',
      name: '查询聊天室信息',
      isActive: false,
      body: {
        "//": "请修改上方 URL 参数进行调试"
      }
    },
    ]
  },
  {
    category: '聊天室属性',
    isFold: false,
    isActive: false,
    children: [{
      url: '/chatrooms/atts/add',
      doc: 'chatroom/chrmatt/addchrmatt/',
      method: 'post',
      name: '设置聊天室属性',
      isActive: false,
      body: {
        "from_id": "userid1",
        "chat_id": "chatroom1000",
        "atts": [{
          "key": "k1",
          "value": "v1",
          "is_force": false
        }]
      }
    },
    {
      url: '/chatrooms/atts/del',
      doc: 'chatroom/chrmatt/delchrmatt/',
      method: 'post',
      name: '删除聊天室属性',
      isActive: false,
      body: {
        "from_id": "userid1",
        "chat_id": "chatroom1000",
        "atts": [{
          "key": "k1",
          "is_force": false
        }]
      }
    },
    {
      url: '/chatrooms/atts/qry',
      doc: 'chatroom/chrmatt/qrychrmatt/',
      method: 'post',
      name: '查询聊天室属性',
      isActive: false,
      body: {
        "chat_id": "chatroom1000",
        "att_keys": ["k1", "k2"]
      }
    },
    {
      url: '/chatrooms/atts/list?chat_id=chatroom1000',
      doc: 'chatroom/chrmatt/listchrmatt/',
      method: 'get',
      name: '查询聊天室全量属性',
      isActive: false,
      body: {
        "//": "请修改上方 URL 参数进行调试"
      }
    },
    ]
  },
  {
    category: '聊天室禁言',
    isFold: false,
    isActive: false,
    children: [{
      url: '/chatrooms/mutemembers/add',
      doc: 'chatroom/chrmmute/addchrmmute/',
      method: 'post',
      name: '添加聊天室成员禁言',
      isActive: false,
      body: {
        "chat_id":"chatroom1000",
        "member_ids":["member1","member2"]
      }
    },
    {
      url: '/chatrooms/mutemembers/del',
      doc: 'chatroom/chrmmute/delchrmmute/',
      method: 'post',
      name: '移除聊天室成员禁言',
      isActive: false,
      body: {
        "chat_id":"chatroom1000",
        "member_ids":["member1","member2"]
      }
    },
    {
      url: '/chatrooms/mutemembers/query?chat_id=chatroom1000&offset=1&limit=100',
      doc: 'chatroom/chrmmute/qrychrmmute/',
      method: 'get',
      name: '查询禁言聊天室成员',
      isActive: false,
      body: {
        "//": "请修改上方 URL 参数进行调试"
      }
    },
    ]
  },
  {
    category: '聊天室封禁',
    isFold: false,
    isActive: false,
    children: [{
      url: '/chatrooms/banmembers/add',
      doc: 'chatroom/chrmban/addchrmban/',
      method: 'post',
      name: '添加聊天室成员封禁',
      isActive: false,
      body: {
        "chat_id":"chatroom1000",
        "member_ids":["member1","member2"]
      }
    },
    {
      url: '/chatrooms/banmembers/del',
      doc: 'chatroom/chrmban/delchrmban/',
      method: 'post',
      name: '移除聊天室成员封禁',
      isActive: false,
      body: {
        "chat_id":"chatroom1000",
        "member_ids":["member1","member2"]
      }
    },
    {
      url: '/chatrooms/banmembers/query?chat_id=chatroom1000&offset=1&limit=100',
      doc: 'chatroom/chrmban/qrychrmban/',
      method: 'get',
      name: '查询封禁的聊天室成员',
      isActive: false,
      body: {
        "//": "请修改上方 URL 参数进行调试"
      }
    },
    ]
  },
  {
    category: '聊天室白名单',
    isFold: false,
    isActive: false,
    children: [{
      url: '/chatrooms/allowmembers/add',
      doc: 'chatroom/chrmallow/addchrmallow/',
      method: 'post',
      name: '添加聊天室成员白名单',
      isActive: false,
      body: {
        "chat_id":"chatroom1000",
        "member_ids":["member1","member2"]
      }
    },
    {
      url: '/chatrooms/allowmembers/del',
      doc: 'chatroom/chrmallow/delchrmallow/',
      method: 'post',
      name: '移除聊天室成员白名单',
      isActive: false,
      body: {
        "chat_id":"chatroom1000",
        "member_ids":["member1","member2"]
      }
    },
    {
      url: '/chatrooms/allowmembers/query?chat_id=chatroom1000&offset=1&limit=100',
      doc: 'chatroom/chrmallow/qrychrmallow/',
      method: 'get',
      name: '查询禁言聊天室成员',
      isActive: false,
      body: {
        "//": "请修改上方 URL 参数进行调试"
      }
    },
    ]
  }
]

export { chatrooms };