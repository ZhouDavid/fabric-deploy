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
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	stopShellPwd = "stop_network"
)

var stopDpath string
var isDeleteFiles bool = false

// stopNetworkCmd represents the stopNetwork command
var stopNetworkCmd = &cobra.Command{
	Use:   "stopNetwork",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStopNetwork()
	},
}

func init() {
	rootCmd.AddCommand(stopNetworkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stopNetworkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stopNetworkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	stopNetworkCmd.Flags().StringVar(&stopDpath, "dpath", "", "stop cmd dpath")
	stopNetworkCmd.Flags().BoolVar(&isDeleteFiles, "deleteFiles", false, "清理生产成的内容")
	stopNetworkCmd.MarkFlagRequired("dpath")
}

func runStopNetwork() {
	workdirAbs, _ := filepath.Abs(stopDpath)
	ips, err := config.NewIPMap(filepath.Join(workdirAbs, "ipmap.txt"))
	if err != nil {
		println(err)
		os.Exit(1)
	}

	for _, value := range ips {
		cmd := getStopCmd(value.Role, value.IP, workdirAbs)
		fmt.Println("start network shell ", cmd)
		runShell(value.UserName, value.Password, value.GetAddress(), cmd, true)
		if isDeleteFiles {
			deleteCmd := filepath.Join(workdirAbs, "shell", "cleanup.sh")
			runShell(value.UserName, value.Password, value.GetAddress(), deleteCmd, true)
		}
	}

}

func getStopCmd(nodeType, ip, workdirAbs string) string {
	switch nodeType {
	case "orderer":
		return filepath.Join(workdirAbs, "shell", stopShellPwd, "orderer_stop.sh") + " " + ip + " " + stopDpath //常量后面统一提取
	case "peer":
		return filepath.Join(workdirAbs, "shell", stopShellPwd, "peer_start.sh") + " " + ip + " " + stopDpath
	default:
		return filepath.Join(workdirAbs, "shell")
	}
}
