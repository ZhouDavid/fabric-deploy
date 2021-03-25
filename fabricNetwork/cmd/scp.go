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
	"Data_Bank/fabric-deploy-tools/fabricNetwork/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var user, pwd, dPath string

// scpCmd represents the scp command
var scpCmd = &cobra.Command{
	Use:   "scp",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		runScp()
	},
}

func init() {
	rootCmd.AddCommand(scpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	scpCmd.Flags().StringVar(&user, "user", "", "ssh user")
	scpCmd.Flags().StringVar(&pwd, "pwd", "", "ssh pwd")
	scpCmd.Flags().StringVar(&dPath, "dpath", "", "")
	scpCmd.Flags().StringVar(&scpSource, "scpSource", "", "scp ")
	scpCmd.MarkFlagRequired("user")
	scpCmd.MarkFlagRequired("pwd")
	scpCmd.MarkFlagRequired("dpath")
}

func runScp() {
	if err := CheckArgs(); err != nil {
		fmt.Println("args check failed:", err)
		os.Exit(1)
	}
	_, err := filepath.Abs(scpSource)
	if err != nil {
		os.Exit(1)
	}
	scp(user, pwd, addr, scpSource, dPath, false)
	// files, err := utils.GetFiles(sourceAbs)
	// if err != nil {
	// 	os.Exit(1)
	// }

	// ssh and scp
	// client, err := utils.Dial(user, pwd, addr)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// for _, file := range files {
	// 	fmt.Printf("scp file:%s \r\n", file.Info.Name())
	// 	utils.Scp(client, file.F, file.Info.Size(), file.Info.Mode(), dPath, file.Info.Name())
	// 	fmt.Printf("scp file %s end \r\n", file.Info.Name())
	// }
}

func scp(user, pwd, addr, sPath, dPath string, isFile bool) {
	fmt.Println("username=", user, " password=", pwd, " addr=", addr, " sourcePath=", sPath, " dPath=", dPath)
	client, err := utils.Connect(user, pwd, addr)
	defer client.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if isFile {
		utils.UploadFile(client, sPath, dPath)
	} else {
		utils.UploadDirectory(client, sPath, dPath)
	}
}
