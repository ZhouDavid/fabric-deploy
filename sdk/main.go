package main

import (
	"Data_Bank/fabric-deploy-tools/sdk/tools"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"

	// "github.com/hyperledger/fabric/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

type InitInfo struct {
	ChannelID      string
	ChannelConfig  string
	OrgAdmin       string
	OrgName        string
	OrdererOrgName string
	OrgResMgmt     *resmgmt.Client
}

const (
	channelID      = "testchannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
	peer1          = "peer0.org1.example.com"
)

func main() {
	//tx
	gen := tools.NewConfigtxgen()
	gen.Exec()
	// privider := config.FromFile("config_test.yaml")
	// sdk, err := fabsdk.New(privider)

	// if err != nil {
	// 	println("...")
	// 	// println(err)
	// }
	// defer sdk.Close()
	// println("hello")
	//管理员账号才能进行Hyperledger fabric网络的管理操作，所以创建资源管理客户端一定要使用管理员账号。
	// rcp := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
	// //然后通过resmgmt.New创建fabric go sdk资源管理客户端。
	// rc, err := resmgmt.New(rcp)
	// if err != nil {
	// 	log.Panicf("failed to create resource client: %s", err)
	// }
	// pb, err := rc.QueryConfigFromOrderer("mychannel")
	// if err != nil {
	// 	println(err)
	// }
	// println(pb.Orderers())
	// config := os.Args[1]
	// initInfo := &InitInfo{

	// 	ChannelID:     "test",
	// 	ChannelConfig: config,

	// 	OrgAdmin:       "Admin",
	// 	OrgName:        "Org1",
	// 	OrdererOrgName: "orderer.example.com",
	// }

	// sdk, err := SetupSDK("config_test.yaml", false)
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// 	return
	// }

	// defer sdk.Close()

	// err = CreateChannel(sdk, initInfo)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// Invoke(sdk, os.Args[1], os.Args[2], os.Args[3])

}

func Invoke(sdk *fabsdk.FabricSDK, chaincodeId, funcName, id string) {
	clientChannelContext := sdk.ChannelContext("mychannel", fabsdk.WithUser("User1"), fabsdk.WithOrg("Org1"))
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		// return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}

	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")
	req := channel.Request{ChaincodeID: chaincodeId, Fcn: funcName, Args: [][]byte{[]byte(id)}}
	respone, err := channelClient.Query(req)
	if err != nil {
		println(err.Error(), "....")
		os.Exit(1)

	}
	println(string(respone.Payload))

}

const ChaincodeVersion = "1.0"

func SetupSDK(ConfigFile string, initialized bool) (*fabsdk.FabricSDK, error) {

	if initialized {
		return nil, fmt.Errorf("Fabric SDK已被实例化")
	}

	sdk, err := fabsdk.New(config.FromFile(ConfigFile))
	if err != nil {
		return nil, fmt.Errorf("实例化Fabric SDK失败: %v", err)
	}

	fmt.Println("Fabric SDK初始化成功")
	return sdk, nil
}
func CreateChannel(sdk *fabsdk.FabricSDK, info *InitInfo) error {

	clientContext := sdk.Context(fabsdk.WithUser(info.OrgAdmin), fabsdk.WithOrg(info.OrgName))
	if clientContext == nil {
		return fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}

	// New returns a resource management client instance.
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		return fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
	}

	// New creates a new Client instance
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(info.OrgName))
	if err != nil {
		return fmt.Errorf("根据指定的 OrgName 创建 Org MSP 客户端实例失败: %v", err)
	}

	//  Returns: signing identity
	adminIdentity, err := mspClient.GetSigningIdentity(info.OrgAdmin)
	if err != nil {
		return fmt.Errorf("获取指定id的签名标识失败: %v", err)
	}
	fmt.Println(adminIdentity)
	// SaveChannelRequest holds parameters for save channel request
	channelReq := resmgmt.SaveChannelRequest{ChannelID: info.ChannelID, ChannelConfigPath: info.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	// save channel response with transaction ID
	_, err = resMgmtClient.SaveChannel(channelReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("创建应用通道失败: %v", err)
	}

	fmt.Println("通道已成功创建，")

	info.OrgResMgmt = resMgmtClient

	// allows for peers to join existing channel with optional custom options (specific peers, filtered peers). If peer(s) are not specified in options it will default to all peers that belong to client's MSP.
	err = info.OrgResMgmt.JoinChannel(info.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(info.OrdererOrgName))
	if err != nil {
		return fmt.Errorf("Peers加入通道失败: %v", err)
	}

	fmt.Println("peers 已成功加入通道.")
	return nil
}
