# FabricNetwork部署文档
fabricNetwork是一个方便用户进行fabric联盟链管理的命令行工具，能够实现联盟链的创建、启动与关闭，通道创建，链码安装与初始化等核心功能。

本文档将介绍如何用fabricNetwork部署一个简单的联盟链系统。

## 1.fabric 相关概念

## 2.部署流程

### 2.1 环境准备
**机器配置**
- 若干台 CentOS7/Ubuntu 18.04
- 每台机器均需安装[git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git), [wget](https://www.tecmint.com/install-wget-in-linux/), [curl](https://www.tecmint.com/install-curl-in-linux/), [vim](https://www.tecmint.com/install-vim-in-linux/)
- 良好的网络连接

**获取命令行工具**

 执行 ./build.sh 后将fabricNetwork.tar.gz 放置在机器的固定位置，
 ```
 example:
   path: /opt 目录下
 ```

- 解压

    ```
    mkdir -p /opt/fabric-workspace
    tar -zxvf fabricNetwork.tar.gz -C /opt/fabric-workspace
    cd /opt/fabric-workspace
    ```

- 安装包目录说明

    ```
        ls /opt/fabric-workspace

        --------------------
        - env    // 环境依赖安装包以及脚本
        - fabricNetwork  // 可执行脚本，主要用来执行安装命令
        - test-network   // fabric 工作目录以及安装脚本
            - chaincode  // 智能合约位置
            - channel-artifacts  // 通道位置
            - configtx  // fabric 配置文件生成位置
            - docker-compose  // fabric dc.yaml 生成
            - explorer //智能合约浏览器
            - shell 可执行脚本
        - networkconfig.json // 应用配置文件，主要包括 机器IP,USER,PWD,以及对fabric 节点的规划

        总结：用户只需要对networkconfig.json 进行配置，具体参考《修改配置文件》部分
    ```


- 预览命令

    ```
    ./fabricNetwork -h
    ```
    - 说明
        - env
            - 远程执行fabric环境依赖
        - scp 
            - 软件包分发

### 修改配置文件
    > 当前位置 /opt/fabric-workspace
    //TODO 
    ```
        vim networkconfig.json

        {
    "ordererType": "etcdedraft",   # etcdedraft /  solo
    "orgs": [
        {
            "name": "Org1",    # 名称
            "domain": "org1.example.com",  # 域名
            "userNum": 1, 
            "peers": [
                {
                    "isOrderer": false,
                    "ip": "172.38.50.211",
                    "peerPort": 7051,
                    "sshPort": 22,
                    "username": "root",
                    "password": "ehl1234",
                    "isInstallExplorer": true
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


### 环境安装

- 软件分发
    - 查看

    ```
        ./fabricNetwork scp -h
    ```

    - 分发

    ```
        ./fabricNetwork scp --addr ip:22 --user root --pwd ehl --scpSource /opt/fabric-workspace -dpath /opt/fabric-workspace

        说明: 
        --addr 远程IP:22 
        --user 用户名
        --pwd 密码
        --scpSource  需要分发的软件目录
        --dpath  目标机器位置，建议与scpSource 保持一致

    ```


## 验证
- 区块链浏览器访问

```

```