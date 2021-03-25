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
	"Data_Bank/fabric-deploy-tools/fabricNetwork/explorer"
	"Data_Bank/fabric-deploy-tools/fabricNetwork/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var explorerWorkDir, explorerDPath string
var isPreInstall bool = false

// explorerCmd represents the explorer command
var explorerCmd = &cobra.Command{
	Use:   "explorer",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("explorer called")
		// exec dep
		runRunExplorer()
	},
}

func init() {
	rootCmd.AddCommand(explorerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// explorerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// explorerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	explorerCmd.Flags().StringVar(&explorerWorkDir, "dir", "test-network", "")
	explorerCmd.Flags().StringVar(&explorerDPath, "dpath", "", "目标位置")
	explorerCmd.Flags().BoolVar(&isPreInstall, "preInstall", false, "")
	explorerCmd.MarkFlagRequired("dir")
}
func check() {
	//docker network
}
func CheckScp() error {
	if explorerDPath == "" {
		return errors.New("null dpath")
	}
	if !filepath.IsAbs(explorerDPath) {
		return errors.New("dpath 为绝对路径")
	}
	return nil
}
func runRunExplorer() {
	exportWorkDirAbs, _ := filepath.Abs(explorerWorkDir)
	ips, err := config.NewIPMap(filepath.Join(exportWorkDirAbs, "ipmap.txt"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	entry, err := ips.GetInstallExplorer()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	domain := explorer.NewInput(entry)
	jsonPath := filepath.Join(exportWorkDirAbs, "explorer", "connection-profile")
	os.MkdirAll(jsonPath, os.ModePerm)
	if err := domain.CreateExplorerConfig(filepath.Join(jsonPath, "test-network.json")); err != nil {
		panic(err)
	}

	//docker-compose n
	// var hosts = make([]string, 0)
	// hosts = append(hosts, entry.Domain+":"+entry.IP)
	explorerContent := config.NewExplorerDockerComposeService(ips.GetAllDomainIP(), filepath.Join(exportWorkDirAbs, "organizations"))

	dockerComposePath := filepath.Join(exportWorkDirAbs, "explorer", "docker")
	os.MkdirAll(dockerComposePath, os.ModePerm)
	ioutil.WriteFile(filepath.Join(dockerComposePath, "docker-compose-ehl.yaml"), explorerContent, os.ModePerm)
	//
	cmdPath := filepath.Join(exportWorkDirAbs, "explorer", "shell")
	os.MkdirAll(cmdPath, os.ModePerm)
	clientIP, err := utils.GetClientIp()
	if isPreInstall {
		if clientIP != entry.IP {
			//sccp

			if err := CheckScp(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			scp(entry.UserName, entry.Password, entry.GetAddress(), filepath.Join(explorerDPath, "explorer"), filepath.Join(exportWorkDirAbs, "explorer"), false)
		}
		runShell(entry.UserName, entry.Password, entry.GetAddress(), filepath.Join(cmdPath, "pre_install.sh"), true)
	}
	if clientIP != entry.IP {
		//scp
		if err := CheckScp(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		scp(entry.UserName, entry.Password, entry.GetAddress(), filepath.Join(explorerDPath, "explorer"), filepath.Join(exportWorkDirAbs, "explorer"), false)
	}
	runShell(entry.UserName, entry.Password, entry.GetAddress(), filepath.Join(cmdPath, "start_explorer.sh"), true)
}
