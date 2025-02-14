let histories = [
  { 
    category: '历史消息', 
    isFold: false, 
    isActive: false,
    children: [
      { 
        url: '/hismsgs/del',
        doc: 'hismsg/delhismsgs/',
        method: 'post', 
        name: '删除消息', 
        isActive: false,
        body: {
          "from_id":"xxx",
          "target_id":"xxx",
          "channel_type":1,
          "del_scope":0,
          "msgs":[
            {
              "msg_id":"xxxxx"
            }
          ]
        }
      },
      { 
        url: '/hismsgs/recall',
        doc: 'hismsg/recallmsg/',
        method: 'post', 
        name: '撤回消息', 
        isActive: false,
        body: {
          "from_id":"xxx",
          "target_id":"xxx",
          "channel_type":1,
          "msg_id":"xxxx",
          "msg_time":1569345643212,
          "exts":{
            "k1":"v1"
          }
        }
      },
      { 
        url: '/apigateway/hismsgs/clean',
        doc: 'hismsg/cleanhismsgs/',
        method: 'post', 
        name: '清空历史消息', 
        isActive: false,
        body: {
          "from_id":"xxx",
          "target_id":"xxx",
          "channel_type":1,
          "clean_time":1569345643212,
          "clean_time_offset":0,
          "clean_scope":0,
          "sender_id":"user1"
        }
      },
      { 
        url: '/hismsgs/query?from_id=xxx&target_id=xxx&channel_type=1&start=1&count=20',
        doc: 'hismsg/qryhismsgs/',
        method: 'get', 
        name: '查询历史消息', 
        isActive: false,
        body: {
          "//": "请修改上方 URL 参数进行调试"
        }
      },
    ] 
  }
]

export { histories };