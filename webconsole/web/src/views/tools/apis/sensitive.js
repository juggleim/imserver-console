let sensitivewords = [
  { 
    category: '敏感词相关', 
    isFold: false, 
    isActive: false,
    children: [
      { 
        url: '/sensitivewords/add',
        doc: 'sensitivewords/addsensitivewords/',
        method: 'post', 
        name: '添加敏感词', 
        isActive: false,
        body: {
          "items":[
            {
              "word":"xxxxx",
              "word_type":1
            }
        ]
        }
      },
      { 
        url: '/sensitivewords/del',
        doc: 'sensitivewords/delsensitivewords/',
        method: 'post', 
        name: '删除敏感词', 
        isActive: false,
        body: {
          "words":["xxxxx","yyyyy"]
        }
      },
      { 
        url: '/sensitivewords/list?limit=20&offset=1',
        doc: 'sensitivewords/qrysensitivewords/',
        method: 'get', 
        name: '查询敏感词', 
        isActive: false,
        body: {
          "//": "请修改上方 URL 参数进行调试"
        }
      },
    ] 
  }
]

export { sensitivewords };