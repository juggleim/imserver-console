let conversations = [
  { 
    category: '会话相关', 
    isFold: false, 
    isActive: false,
    children: [
      { 
        url: '/convers/undisturb',
        doc: 'convers/undisturbconvers/',
        method: 'post', 
        name: '设置免打扰', 
        isActive: false,
        body: {
          "user_id":"userid1",
          "items":[
            {
              "target_id":"userid2",
              "channel_type":1,
              "undisturb_type":1
            },
            {
              "target_id":"groupid1",
              "channel_type":2,
              "undisturb_type":1
            }
          ]
        }
      },
      { 
        url: '/convers/del',
        doc: 'convers/delconvers/',
        method: 'post', 
        name: '删除会话', 
        isActive: false,
        body: {
          "user_id":"userid1",
          "items":[
            {
              "target_id":"userid2",
              "channel_type":1
            },
            {
              "target_id":"groupid1",
              "channel_type":2
            }
          ]
        }
      },
      { 
        url: '/convers/add',
        doc: 'convers/addconver/',
        method: 'post', 
        name: '添加会话', 
        isActive: false,
        body: {
          "user_id":"userid1",
          "target_id":"userid2",
          "channel_type":1
        }
      },
      { 
        url: '/globalconvers/query?start=1&count=20&user_id=userid1',
        doc: 'convers/qryconvers/',
        method: 'get', 
        name: '查询全局会话列表', 
        isActive: false,
        body: {
          "//": "请修改上方 URL 参数进行调试"
        }
      },
    ] 
  }
]

export { conversations };