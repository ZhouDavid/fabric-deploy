# Fabric 部署管理程序设计文档

## 功能设计
0. 读取网络配置 (Go)
从networkConfig.json(yaml)读取各节点相关信息，包括：
- Orderer section
        - orderer type
        - orderer addresses

- Org section
        - 组织名称
        - 域名名称
        - peer节点数量
        - admin之外的user节点数量
        - 组织内各节点ip+port

读取完毕后将以上信息**写入configtx.yaml，crypto.yaml，docker-compose.yaml** 三个template
```bash
fabricNetwork loadConfig --config-path <config_path>
```

1. 查询、下载、删除所需要的依赖，包括：(sh)
- 下载docker, docker-compose 
- 下载go1.14.12
- 下载fabric 2.2.2的各类docker image (peer, orderer, ccenv, tools)
- 下载fabric 2.2.2的各类binary (https://github.com/hyperledger/fabric/releases/download/v2.2.2/hyperledger-fabric-linux-amd64-2.2.2.tar.gz)
以上所有需在所有节点上安装，需要master节点与网络内其他各节点通信，**以交互形式输入密码**。


#### 命令定义
- env  



```bash

# 从configtx.yaml读取IP信息，用户交互式输入密码，远程ssh到相关节点并安装相关依赖
fabricNetwork installDependencies # default is 2.2.2
fabricNetwork installDependency --binary <binary-name> --version <version> # default is peer+orderer+ccenv+tools(2.2.2)
fabricNetwork installDependency --docker-image <image-name> # default is 2.2.2
fabricNetwork installDependency --docker-compose <docker-compose version> # default is latest
fabricNetwork installDependency --docker <docker version> # default is latest
```

```bash
fabricNetwork CheckDependencies
```

```bash
fabricNetwork RemoveDependencies
```

2. 创建组织 (go)
- 从configtx.yaml文件读取各节点的ip，所属的组织，在组织中的角色，证书(MSP)位置，若无相关信息提示`run loadConfig first`
- 将信息写入crypto.yaml，作为cryptogen命令的参数文件
- 按照cryptogen格式生成各组织所需要的证书，并根据角色分发到各个ip (sh)

```bash
# 从configtx.yaml 读取ip, 组织名称,用cryptogen/fabric-ca生成组织证书
fabricNetwork createOrgs # by default is cryptogen
fabricNetwork createOrgs --crypto=fabric-ca --ca-address <ca-address> --ca-name <ca-name> # otherwise use fabric-ca
```

3. 启动网络 (sh)
- 创建创世区块，`configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block`
- 若orderer type为raft, scp genesis.block到各个orderer节点
- 从docker-compose.yaml文件读取各节点配置信息
- 环境变量写入.env文件
- ssh到各个节点执行`docker-compose up -d`启动相应docker container

```bash
fabricNetwork startNetwork
```

4. 创建通道 (sh)
- Channel Policy全走默认configtx.yaml template, 即只能创建平权channel
- 检查网络是否启动，等待查看各节点是否正常
- 读取通道名称，读取参与通道的组织名称，结果写入configtx.yaml
- 从configtx.yaml读取通道配置信息并生成block并分发到各相关节点
- ssh到Orderer执行`peer channel create`
- ssh到各组织节点执行`peer channel join`
```bash
fabricNetwork createChannel --channelName <channel-name> --orgNames <space-seperated-string>
# fabricNetwork createChannel --channelName mychannel --orgNames "org1 org2 org3"
```

5. 部署智能合约
- 读取chaincode位置，语言，所关联的channel名称
- 在chaincode根目录下执行`go mod`下载依赖
- 执行`peer lifecycle chaincode package` 打包链码
- 执行`peer lifecycle chaincode install` 安装链码
- ssh到相应节点执行`peer lifecycle chaincode approveformyorg`
- ssh到任一channel内节点执行`peer lifecycle chaincode commit` 提交链码

6. 执行智能合约
- 读取智能合约名称，channel名称
- ssh到任一channel内节点执行`peer chaincode invoke`

7. 查询智能合约
- 读取智能合约名称，channel名称
- ssh到人一channel内节点执行`peer chaincode query`
