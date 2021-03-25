/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"Data_Bank/fabric-deploy-tools/fabricNetwork/config"
	"Data_Bank/fabric-deploy-tools/fabricNetwork/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var netCfg, spath, dpath string
var useScp bool

//
// startNetworkCmd represents the startNetwork command
var startNetworkCmd = &cobra.Command{
	Use:   "startNetwork",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStartNetwork()
	},
}

func init() {
	rootCmd.AddCommand(startNetworkCmd)
	//peer0.org1.example.com 127.0.0.1 7050,22 ehualu 8888 peer
	//orderer0.orderer.example.com 127.0.0.1 7050,22 ehualu 8888 orderer
	// startNetworkCmd.Flags().StringVar(&netCfg, "startCfg", "ipmap.txt", "ipmap.txt")
	startNetworkCmd.Flags().StringVar(&spath, "sPath", "", "source: test-network dir ")
	startNetworkCmd.Flags().StringVar(&dpath, "dPath", "test-network", "target :test-network dir")
	startNetworkCmd.Flags().BoolVar(&useScp, "useScp", true, "")
	// startNetworkCmd.MarkFlagRequired("sPath")
	startNetworkCmd.MarkFlagRequired("dPath")
}

const (
	defaultCfg    string = ""
	startShellPwd        = "start_network"
)

func runStartNetwork() {

	ips, err := config.NewIPMap(filepath.Join(dpath, "ipmap.txt"))
	if err != nil {
		panic(err)
	}
	clientIP, _ := utils.GetClientIp()
	for _, value := range ips {
		if useScp {
			// scp 分发
			if spath == "" {
				fmt.Println("use scp cmd,but source directory =nil")
				os.Exit(1)
			}
			if value.IP != clientIP {
				//user,pwd,addr,sPath,dPath  存疑 或者test-network 放到其他位置
				//scp(inputs[3], inputs[4], inputs[1]+":"+ports[len(ports)-1], spath, dpath)
				scp(value.UserName, value.Password, value.GetAddress(), spath, dpath,false)
			}

		}
		//run shell cmd
		cmd := getStartCmd(value.Role, value.IP)
		fmt.Println("start network shell ", cmd)
		runShell(value.UserName, value.Password, value.GetAddress(), cmd, true)
	}

}

// cmd $1=ip $2 =dpath
func getStartCmd(nodeType, ip string) string {
	switch nodeType {
	case "orderer":
		return filepath.Join(dpath, "shell", startShellPwd, "orderer_start.sh") + " " + ip + " " + dpath //常量后面统一提取
	case "peer":
		return filepath.Join(dpath, "shell", startShellPwd, "peer_start.sh") + " " + ip + " " + dpath
	default:
		return filepath.Join(dpath, "shell")
	}
}
