package cmd

import (
	"Data_Bank/fabric-deploy-tools/fabricNetwork/config"
	"Data_Bank/fabric-deploy-tools/fabricNetwork/utils"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var channelName, channelOrgs string
var channelUseScp bool

var createChannelCmd = &cobra.Command{
	Use:   "createChannel",
	Short: "create fabric network application channel",
	Long:  `create fabric network application channel`,
	Run: func(cmd *cobra.Command, args []string) {
		runCreateChannel()
	},
}

func init() {
	rootCmd.AddCommand(createChannelCmd)
	createChannelCmd.Flags().StringVar(&channelName, "channel-name", "", "")
	createChannelCmd.Flags().StringVar(&channelOrgs, "channel-orgs", "Org1", "")
	createChannelCmd.Flags().BoolVar(&channelUseScp, "useScp", true, "")
	installChaincodeCmd.MarkFlagRequired("channel-name")
}

func runCreateChannel() {
	if !filepath.IsAbs(targetPath) {
		panic("Fatal error, dPath is not absolute path!!! Please use absolute path.")
	}
	initConfig()

	conf := &config.NetworkConfig{}
	if err := viper.Unmarshal(conf); err != nil {
		fmt.Printf("error when unmarshal config file%s,error message:%v", viper.ConfigFileUsed(), err)
		panic(err)
	}
	fmt.Printf("Create channel config using %s\n", viper.ConfigFileUsed())
	configBuilder := config.NewFabricConfigBuilder(conf)

	channelTxProfileFilename := path.Join(outputPath, "configtx", channelName+".yaml")
	fabricTxProfileFilename := path.Join(outputPath, "configtx", "configtx.yaml")
	channelTxOutputFilename := path.Join(outputPath, "channel-artifacts", channelName+".tx")
	// channelTxBlockFilename := path.Join(outputPath, "channel-artifacts", channelName+".block")
	configPath := path.Join(outputPath, "configtx")

	orgNames := strings.Split(channelOrgs, ",")
	configBuilder.BuildTxFile(configPath, channelName, orgNames...)
	os.Rename(channelTxProfileFilename, fabricTxProfileFilename) // Rename channel.yaml to configtx.yaml
	utils.ExecuteCommand("env/fabric_bin/configtxgen", "-profile", "ApplicationChannel", "-channelID", channelName, "-outputCreateChannelTx", channelTxOutputFilename, "-configPath", configPath)
	os.Rename(fabricTxProfileFilename, channelTxProfileFilename) // Rename back

	// Scp generated tx file to remote target machine
	ipMap, err := config.NewIPMap(path.Join(outputPath, "ipmap.txt"))
	if channelUseScp {
		if err != nil {
			panic(err)
		}
		if err != nil {
			os.Exit(1)
		}
		clientIP, _ := utils.GetClientIp()
		for _, value := range ipMap {
			if clientIP != value.IP {
				fmt.Printf("Scp %s to %s...\n", channelTxOutputFilename, value.Domain)
				scp(value.UserName, value.Password, value.GetAddress(), channelTxOutputFilename, filepath.Join(targetPath, "channel-artifacts"), true)
				//scp(value.UserName, value.Password, value.GetAddress(), filepath.Join(targetPath, "shell"), targetPath, false)
			}
		}
	}

	// Run peer channel create command on target machine to generate .block file
	orderer := ipMap.GetOrderer()
	ordererAddress := orderer.GetHostDomainIP()

	// Execute createChannel.sh at peer0 of first org
	targetPeer,_ := ipMap.GetPeerFromOrg(orgNames[0])
	shellPath, _ := filepath.Abs(path.Join(outputPath, "shell", "peer", "createChannel.sh"))
	channelCreateCmd := fmt.Sprintf("%s %s %s %s %s", shellPath, ordererAddress, orderer.Domain, orderer.OrgDomain, channelName)

	fmt.Printf("Executing create channel command at %s \n", targetPeer.Domain)
	client, err := utils.Dial(targetPeer.UserName, targetPeer.Password, fmt.Sprintf("%s:%d", targetPeer.IP, targetPeer.SSHPort))
	if err != nil {
		fmt.Printf("Fail to ssh to %s \n", targetPeer.Domain)
		panic(err)
	}
	fmt.Println(channelCreateCmd)
	if stdout, err := utils.RunCommand(client, channelCreateCmd, true); err != nil {
		fmt.Printf("Fail to execute create channel command, stdout:%s, error:%v\n", stdout.String(), err)
	} else {
		fmt.Printf("Run create channel command successfully \n.stdout:%s \n", stdout.String())
	}

	// Execute joinChannel.sh at peer0 of each channel org
	shellPath, _ = filepath.Abs(path.Join(outputPath, "shell", "peer", "joinChannel.sh"))

	for _, orgName := range orgNames {
		targetPeer,err:= ipMap.GetPeerFromOrg(orgName)
		if err!=nil{
			fmt.Printf("Fail to add %s to channel %s \n",orgName,channelName)
		}
		client, err := utils.Dial(targetPeer.UserName, targetPeer.Password, fmt.Sprintf("%s:%d", targetPeer.IP, targetPeer.SSHPort))
		defer client.Close()
		if err != nil {
			fmt.Printf("Fail to ssh to %s \n", targetPeer.Domain)
			panic(err)
		}
		// envs := targetPeer.GetPeerEnvVariables(outputPath)
		channelJoinCmd := fmt.Sprintf("%s %s", shellPath, channelName)
		fmt.Printf("Executing join channel command:%s at %s \n", channelJoinCmd, targetPeer.Domain)
		fmt.Println(channelJoinCmd)
		if stdout, err := utils.RunCommand(client, channelJoinCmd, true); err != nil {
			fmt.Printf("Fail to execute join channel command, stdout:%s, error:%v\n", stdout.String(), err)
		} else {
			fmt.Printf("Run join channel command successfully \n.stdout:%s \n", stdout.String())
		}
	}
}
