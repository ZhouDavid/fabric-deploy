package cmd

import (
	"Data_Bank/fabric-deploy-tools/fabricNetwork/config"
	"Data_Bank/fabric-deploy-tools/fabricNetwork/utils"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)


var installChaincodeCmd = &cobra.Command{
	Use:   "installChaincode",
	Short: "install chaincode",
	Long:  `install chaincode`,
	Run: func(cmd *cobra.Command, args []string) {
		runInstallChaincode()
	},
}

var chaincodeLocalPath string
var chaincodeName string
var chaincodeVersion string
var hosts []string

func init() {
	rootCmd.AddCommand(installChaincodeCmd)
	installChaincodeCmd.Flags().StringVar(&chaincodeLocalPath, "ccPath", "", "chaincode local path")
	installChaincodeCmd.Flags().StringVar(&chaincodeName, "ccName", "", "chaincode name")
	installChaincodeCmd.Flags().StringVar(&chaincodeVersion, "ccVersion", "1.0", "chaincode version")
	installChaincodeCmd.Flags().StringVar(&channelName, "channelName","","channel to install chaincode on")
	installChaincodeCmd.Flags().StringSliceVar(&hosts, "hosts", nil, "domain name of hosts to install chaincode on")
	installChaincodeCmd.MarkFlagRequired("ccPath")
	installChaincodeCmd.MarkFlagRequired("ccName")
	installChaincodeCmd.MarkFlagRequired("channelName")
	installChaincodeCmd.MarkFlagRequired("hosts")
	chaincodeLocalPath,_=filepath.Abs(chaincodeLocalPath)
}

func runInstallChaincode() {
	// Load ipmap
	ipMap, err := config.NewIPMap(path.Join(targetPath, "ipmap.txt"))
	if err != nil{
		logrus.Errorf("Fail to load ipmap")
		os.Exit(1)
	}

	chaincodeRemotePath := filepath.Join(targetPath, "chaincode","go", filepath.Base(chaincodeLocalPath))	
	fmt.Println("Chaincode remote path: ", chaincodeRemotePath)
	ordererAddress := ipMap.GetOrderer().GetHostDomainIP()
	for _,host:=range hosts{
		h,ok:=ipMap[host]
		if !ok {
			logrus.Errorf("Fail to find host domain: %s in ipmap, skip chaincode installation \n", host)
			continue
		}
		client, err := utils.Dial(h.UserName, h.Password, fmt.Sprintf("%s:%d", h.IP, h.SSHPort))
		defer client.Close()
		if err != nil {
			logrus.Errorf("Fail to ssh to %s \n", h.Domain)
			os.Exit(1)
		}
		// Scp chaincode folder
		logrus.Infof("Scp chaincode folder:%s to remote path:%s \n", chaincodeLocalPath, chaincodeRemotePath)
		scp(h.UserName, h.Password, h.GetAddress(), chaincodeLocalPath, chaincodeRemotePath, false)
		// Install chaincode remotely
		peerAddress := h.GetHostDomainIP()
		chaincodePathInContainer := filepath.Join("chaincode","go", filepath.Base(chaincodeLocalPath))	
		installCmd := getChaincodeInstallCmd("cli", chaincodePathInContainer, peerAddress,ordererAddress)
		fmt.Println(installCmd)
		if stdout, err:=utils.RunCommand(client, installCmd, true);err!=nil{
			logrus.Errorf("Fail to execute chaincode install command, stdout:%s, error:%v\n", stdout.String(), err)
			os.Exit(1)
		} else {
			logrus.Errorf("Run chaincode install command successfully \n.stdout:%s \n", stdout.String())
		}
	}
}

func getChaincodeInstallCmd(dockerName, chaincodePath, peerAddress, ordererAddress string) string {
	shellPath := filepath.Join(targetPath, "shell", "chaincode", "chaincode_install.sh")
	return fmt.Sprintf("%s %s %s %s %s %s %s %s", shellPath, dockerName, chaincodePath, chaincodeName, chaincodeVersion, peerAddress, ordererAddress, channelName)
}