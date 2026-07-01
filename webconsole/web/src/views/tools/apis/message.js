let messages = [
  { 
    category: '消息相关', 
    isFold: false, 
    isActive: false,
    children: [
      { 
        url: '/messages/private/send',
        doc: 'message/privatemsg/',
        method: 'post', 
        name: '发送单聊消息', 
        isActive: false,
        body: {
          "sender_id":"userid1",
          "target_ids":["userid2","userid3"],
          "msg_type":"cim:text",
          "msg_content":"{\"content\":\"hello IM\"}",
          "is_storage":true,
          "is_notify_sender":true,
          "is_state":false
        }
      },
      { 
        url: '/messages/group/send',
        doc: 'message/groupmsg/',
        method: 'post', 
        name: '发送群聊消息', 
        isActive: false,
        body: {
          "sender_id":"userid1",
          "target_ids":["groupid1","groupid2"],
          "msg_type":"text",
          "msg_content":"{\"content\":\"aabbcc\"}",
          "is_storage":true,
          "is_notify_sender":true,
          "is_state":false,
          "mention_info":{
            "mention_type":"mention_all",
            "target_user_ids":["userid1","userid2"]
          },
          "refer_msg":{
            "msg_id":"xxx",
            "sender_id":"xxx",
            "target_id":"xxx",
            "channel_type":1,
            "msg_type":"xxx",
            "msg_content":"xxxxx"
          }
        }
      },
      { 
        url: '/messages/system/send',
        doc: 'message/sysmsg/',
        method: 'post', 
        name: '发送系统消息', 
        isActive: false,
        body: {
          "sender_id":"sys1",
          "target_ids":["userid1","userid2"],
          "msg_type":"text",
          "msg_content":"{\"content\":\"aabbcc\"}",
          "is_storage":true,
          "is_state":false
        }
      },
      { 
        url: '/messages/groupcast/send',
        doc: 'message/groupcastmsg/',
        method: 'post', 
        name: '发送群发消息', 
        isActive: false,
        body: {
          "sender_id":"userid1",
          "target_id":"groupcastid1",
          "msg_type":"cim:text",
          "msg_content":"{\"content\":\"aabbcc\"}",
          "target_convers":[
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
        url: '/messages/broadcast/send',
        doc: 'message/broadcastmsg/',
        method: 'post', 
        name: '发送广播消息', 
        isActive: false,
        body: {
          "sender_id":"userid1",
          "msg_type":"text",
          "msg_content":"{\"content\":\"hello im\"}",
        }
      },
      { 
        url: '/messages/markread',
        doc: 'message/markread/',
        method: 'post', 
        name: '标记消息已读', 
        isActive: false,
        body: {
          "user_id":"userid1",
          "target_id":"userid2",
          "channel_type":1,
          "msg_ids":["xxxxxxx"]
        }
      },
    ] 
  }
]

export { messages };