# Fabric前后端接口设计

## Server & User 
### 用户注册
```json
POST:/signup
request:
{
    "username": "ehualu",
    "role": "admin",
    "password": "8888",
    "reenter-password":"8888"
}
response:
{
    "code":200,
    "msg":"succeed"
}
content-type:application/json
```

### 用户登录
```json
POST:/login
request:
{
    "username":"ehualu",
    "password":"8888",
    ""
}
response:
{
    "code":200,
    "msg":"succeed"
}
```

### 组织创建
- 组织创建: 仅管理员可以创建组织
    - 组织名称
    - 组织域名
需要在需要加入组织的机器上预先agent，由agent发现这些机器的ip并上报server,
```json
POST:/createorg
content-type:application/json
request:
{
    "userid":"1234abcd",
    "org":"Org1",
    "domain":"org1.com",
    "ca":"ca.org1.com"  // 可选，若不指定则可让agent直接用cryptogen生成。
}
```

### 联盟链配置
- 配置模式
    - 随机模式：会根据后续创建的组织，随机选择节点加入每个组织
    - 自定义模式：需手动指定每个组织内的peer节点

- 账本存储方式：leveldb/couchdb

```json
POST:/mode
content-type:application/json
request:
{
    "mode":"random",
    "userid":"1234abcd",
    "consensus":"raft",
    "db":"leveldb",
    "orderer_num":1,
    "peer_num":2
},
{
    "mode":"manual",
    "userid":"1234abcd",
    "consensus":"raft",
    "db":"couchdb",
    "roles":[
        {
            "ip":"127.0.0.1",
            "role":"peer"
        },
        {
            "ip":"127.0.0.2",
            "role":"orderer"
        }
    ]
}
```

### 通道创建
- 通道名称
- 通道内的组织
    - 组织名1
    - 组织名2
```json
POST:/createchannel
request:
{
    "name":"mychannel",
    "orgs":["org1","org2"]
}
```

### 邀请组织加入通道
```json
POST:/inviteorg
content-type:application/json
request:
{
    "org":"Org3",
    "channel":"mychannel"
}
```

### 链码实例化
- 链码名称
- 链码本地路径
- 将要部署的通道名称
- 版本号
- 将要安装的peer节点
- 初始化函数名称，默认为不需要初始化
- 背书策略，默认为通道内所有组织均需同意

```json
POST:/chaincode
request:
{
    "name":"basic",
    "path":"/user/local/chaincode/basic.tar.gz",
    "channel":"mychannel",
    "version":"1.0",
    "peers":["peer0.org1","peer0.org2"],
    "init":"init", // 默认为空
    "endorsement": "org1.admin,org2.admin", // 默认为channel所有组织的admin均需同意
}
```

### 链码调用
- 通道名称
- 被调用链码名称
- 被调用函数名称
```json
POST:/invoke
request:
{
    "peers":["peer0.org1","peer0.org1"],
    "channel":"mychannel",
    "name":"basic",
    "func":"InitLedger",
    "args":{
        "arg1":"a",
        "arg2":"b"
    }
}
```

### 区块链浏览器