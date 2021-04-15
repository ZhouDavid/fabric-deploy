package main

import "Data_Bank/fabric-deploy-tools/example/fabric"

func main() {
	/* Create
	ChainName   string   `json:"chainName"`
		Account     string   `json:"account"`     //用户帐号
		Consensus   string   `json:"consensus"`   //共识
		PeersOrgs   []string `json:"peersOrgs"`   //参与组织 除了orderer
		OrderCount  int      `json:"orderCount"`  //orderer节点个数
		PeerCount   int      `json:"peerCount"`   //每个组织节点个数
		ChannelName string   `json:"channelName"` //channel 名
		TlsEnabled  string   `json:"tlsEnabled"`  //是否开启tls  true or false
	*/
	c := fabric.FabricChain{
		Account:     "leixw",
		Consensus:   "etcdraft",
		PeersOrgs:   []string{"org1", "org2"},
		OrderCount:  1,
		PeerCount:   2,
		ChannelName: "test_channel",
		TlsEnabled:  "true",
	}

	build := fabric.NewConfigBuilder(c, "./")
	build.BuildTxFile()
	build.BuildCryptoFile()
}
