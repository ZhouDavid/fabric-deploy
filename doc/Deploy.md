# FabricNetwork部署文档
fabricNetwork是一个方便用户进行fabric联盟链管理的命令行工具，能够实现联盟链的创建、启动与关闭，通道创建，链码安装与初始化等核心功能。

本文档将介绍如何用fabricNetwork部署一个简单的联盟链系统。

## 1.fabric 相关概念

## 2.部署流程
我们将部署一个仅有唯一组织的简单联盟链系统，为简化部署，规定该组织只有一个peer, 同时整个联盟链有一个orderer。

### 2.1 环境准备
**2.1.1 集群配置**
- 一台CentOS7/Ubuntu (以下简称master)：用于执行fabricNetwork相关命令
- 一台CentOS7/Ubuntu (以下简称peer)：代表组织的唯一peer, 执行peer相关命令。
- 一台CentOS7/Ubuntu (以下简称orderer)：联盟内唯一排序节点，与peer通信执行排序相关命令。
- 每台机器均需安装[git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git), [wget](https://www.tecmint.com/install-wget-in-linux/), [curl](https://www.tecmint.com/install-curl-in-linux/), [vim](https://www.tecmint.com/install-vim-in-linux/)
- 良好的网络连接

**注**： 本文档中master与peer为同一台物理机，且下述所有操作均只在**master**上进行

**2.1.2 获取命令行工具及配置文件**
- 从gitlab获取 [fabricNetwork v1.0](http://10.20.5.5:10080/Data_Bank/fabric-deploy-tools/repository/archive.tar.gz?ref=v1.0) (请确保以连接数据湖IDC内网), 以tar.gz格式将项目下载至master的/opt/fabric-workspace目录下
    ```bash
    curl http://10.20.5.5:10080/Data_Bank/fabric-deploy-tools/repository/archive.tar.gz?ref=v1.0
    ```
    **注**：若无法下载tar包，请联系zhoujianyu@ehualu.com获取

- 解压至/opt/fabric-workspace目录
    ```bash
    mkdir -p /opt/fabric-workspace
    tar -zxvf fabric-deploy-tools-v1.0.tar.gz.tar.gz -C /opt/fabric-workspace
    cd /opt/fabric-workspace
    ```

- 解压后的包目录说明
    ```
    fabric-workspace/
    ├── test.sh     // 端到端测试脚本
    ├── env         // docker-compose,docker-image,fabric-binary依赖及安装脚本
    ├── fabricNetwork // 命令行工具，所有命令均由此执行
    └── test-network
        ├── chaincode
        ├── configtx // 通道配置文件生成位置
        ├── docker-compose // docker-compose.yaml生成位置
        ├── explorer //智能合约浏览器
        ├── networkconfig.json  // 默认联盟链配置文件，可根据实际情况修改
        └── shell  // 辅助脚本，将由fabricNetwork间接调用
    ```
**注**： 用户只需要对networkconfig.json进行相应修改，无需修改其他文件。

**下述所有操作均默认在master的/opt/fabric-workspace下进行**

### 2.1.2 下载环境依赖
为所有脚本和二进制文件设置可执行权限
```bash
> find . -type f -name "*.sh" -exec chmod +x {} \+
> chmod +x ./env/fabric_bin/*
> chmod +x fabricNetwork
```

```bash
```

### 2.2 修改配置文件
项目自带配置文件如下所示，该配置文件定义了两个组织Org1和Orderer,其中Orderer组织仅参与排序服务，不实际存储区块链账本, **测试时请将ip设置为测试机器的ip地址。**
```json
{
    "ordererType": "etcdedraft",
    "orgs": [
        {
            "name": "Org1",     // 组织名称
            "domain": "org1.example.com", //组织域名
            "userNum": 1,   // 组织中user的数量
            "peers": [      // 组织中每个peer的详细信息
                {
                    "isOrderer": false, // 该peer是否实际是orderer
                    "ip": "172.38.50.211", // host ip地址
                    "peerPort": 7051,   // peer端口
                    "sshPort": 22,  // 远程ssh到该peer的端口
                    "username": "root", // 远程ssh到该peer的username
                    "password": "ehl1234", // 远程ssh到该peer的password
                    "isInstallExplorer": true // 是否在该peer上安装区块链浏览器
                }
            ]
        },
        {
            "name": "Orderer",
            "domain": "orderer.example.com",
            "userNum": 1,
            "peers": [
                {
                    "isOrderer": true,
                    "ip": "172.38.50.210",
                    "peerPort": 7050,
                    "sshPort": 22,
                    "username": "root",
                    "password": "ehl1234"
                }
            ]
        }
    ]
}
```

### 2.3 测试网络

测试网络有两种方法，运行脚本测试和手动运行命令测试

**2.3.1 运行脚本测试**
直接运行test.sh
```bash
> ./test.sh
```
若无问题应该在最后的输出看到
```bash
2021-03-25 09:27:27.687 UTC [chaincodeCmd] chaincodeInvokeOrQuery -> INFO 001 Chaincode invoke successful. result: status:200
```

**2.3.2 手动运行命令测试**

生成fabric初始化所需的所有配置文件,并分发至各个机器
```bash
> ./fabricNetwork loadConfig --config $PWD/test-network/networkconfig.json --dPath /opt/fabric_install/test-network
# --config: 用户自定义fabric配置文件路径
# --dPath: 生成的所有配置文件目标路径根目录，该路径在联盟链所有机器中统一
```

启动网络
```bash
> ./fabricNetwork startNetwork --dPath /opt/fabric_install/test-network --sPath $PWD/test-network
# --dPath：同上
# --sPath: 配置文件源路径根目录，若未指定则为$PWD/test-network
```

创建应用通道
```bash
> ./fabricNetwork createChannel --dPath /opt/fabric_install/test-network --channel-name mychannel
# --dPath：同上
# --channel-name: 要创建通道的名称
```


安装智能合约
```bash
> ./fabricNetwork installChaincode --ccPath /opt/fabric_install/test-network/chaincode/go/basic --ccName basic --ccVersion 1.0 --channelName mychannel --hosts peer0.org1.example.com
# --ccPath: 智能合约代码所在路径根目录
# --ccName: 智能合约名称
# --ccVersion: 智能合约版本号
# --channelName: 智能合约安装所在的通道名称
# --hosts: 将安装该智能合约的peer名称
```

## 验证
- 区块链浏览器访问

```bash

```