let usertags = [
  { 
    category: '标签推送相关', 
    isFold: false, 
    isActive: false,
    children: [
      { 
        url: '/usertags/add',
        doc: 'push/addusertags/',
        method: 'post', 
        name: '添加用户标签', 
        isActive: false,
        body: {
          "user_tags":[
            {
              "user_id":"userid1",
              "tags":["aa","bb"]
            }
          ]
        }
      },
      { 
        url: '/usertags/del',
        doc: 'push/delusertags/',
        method: 'post', 
        name: '删除用户标签', 
        isActive: false,
        body: {
          "user_tags":[
            {
              "user_id":"userid1",
              "tags":["aa","bb"]
            }
          ]
        }
      },
      { 
        url: '/usertags/clear',
        doc: 'push/clearusertags/',
        method: 'post', 
        name: '清除用户标签', 
        isActive: false,
        body: {
         "user_ids":["userid1","userid2"]
        }
      },
      { 
        url: '/usertags/quer?user_ids=userid1,userid2',
        doc: 'push/qryusertags/',
        method: 'get', 
        name: '查询用户标签', 
        isActive: false,
        body: {
         "//": "请修改上方 URL 参数进行调试"
        }
      },
      { 
        url: '/push',
        doc: 'push/pushwithtags/',
        method: 'post', 
        name: '全员/标签推送', 
        isActive: false,
        body: {
          "from_user_id":"userid1",
          "condition":{
            "tags_and":["tag1","tag2"],
            "tags_or":["tag1","tag2"]
          },
          "msg_body":{
            "msg_type":"cim:text",
            "msg_content":"{\"content\":\"Hello World!\"}"
          },
          "notification":{
            "title":"title",
            "push_text":"推送详情内容"
          }
        }
      },
    ] 
  }
]

export { usertags };