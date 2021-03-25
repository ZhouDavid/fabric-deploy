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

// 变量
var envType, shellUser, shellPwd, shellDpath string

var isOnline, isLocal, isShell bool
var shellMap = map[string]string{
	DOCKER:    "docker.sh",
	GO:        "go.sh",
	DC:        "docker_compose.sh",
	LOADIMAGE: "images.sh",
}

var tarMap = map[string]string{
	DOCKER:    "docker-19.03.9.tgz",
	GO:        "go.tar.gz",
	DC:        "docker_compose",
	LOADIMAGE: "null",
}

const (
	DOCKER    string = "docker"
	GO        string = "go"
	DC        string = "dc"
	LOADIMAGE        = "loadImage"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		runRemoteShellCommand()
	},
}

func init() {

	rootCmd.AddCommand(envCmd)
	//default envType

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	envCmd.Flags().StringVar(&envType, "type", "", "类型：docker,go,dc(docker_compose),loadImage")
	envCmd.Flags().BoolVar(&isOnline, "isOnline", true, "true,false")
	envCmd.Flags().StringVar(&shellUser, "user", "", "ssh user")
	envCmd.Flags().StringVar(&shellPwd, "pwd", "", "ssh pwd")
	envCmd.Flags().StringVar(&shellDpath, "dpath", "d", "workspace")
	envCmd.Flags().BoolVar(&isLocal, "isLocal", false, "local script")
	envCmd.Flags().BoolVar(&isShell, "isShell", true, "use /bin/bash cmd")
	// envCmd.Flags().StringVar(&dFile, "dfile", "tar.gz", "")
	// envCmd.MarkFlagRequired("user")
	// envCmd.MarkFlagRequired("pwd")
	envCmd.MarkFlagRequired("type")
	// envCmd.MarkFlagRequired("dpath")
}

//runRemoteShellCommand exec TODO log
func runRemoteShellCommand() {
	var cmd string
	cmd = filepath.Join(shellDpath, "env", "shell", shellMap[envType])
	installPackage := filepath.Join(shellDpath, "env", "lib", tarMap[envType])
	if isLocal {

		var opts []string
		if isOnline {
			opts = []string{"online"}
		} else {
			opts = []string{"offline", installPackage}
		}
		fmt.Println(cmd, opts, isLocal, isShell, isOnline)
		runLocalShell(cmd, opts...)
	} else {
		if err := CheckArgs(); err != nil {
			fmt.Println("args check failed:", err)
			os.Exit(1)
		}

		if isOnline {
			cmd += " online"
		} else {
			cmd += " offline " + installPackage
		}
		fmt.Println(cmd)
		runShell(shellUser, shellPwd, addr, cmd, true)
	}
}

func runLocalShell(cmd string, opts ...string) {
	outByteBuffer, errByteBuffer, err := utils.ExecuteCommand(cmd, opts...)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("cmd out", outByteBuffer.String())
	fmt.Println("cmd err", errByteBuffer.String())
}

func runShell(user, pwd, addr, cmd string, isOnline bool) {
	client, err := utils.Dial(user, pwd, addr)
	fmt.Println(user, pwd, addr, cmd, isOnline)
	defer client.Close()
	if err != nil {
		fmt.Println(err, ";dial error")
		os.Exit(1)

	}

	bf, err := utils.RunCommand(client, cmd, true)
	if err != nil {
		fmt.Println("exec cmd error", err)
		os.Exit(1)
	}
	fmt.Printf("exec result: %s \r\n", bf.String())
}
