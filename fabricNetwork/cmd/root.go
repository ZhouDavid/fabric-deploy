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
	"Data_Bank/fabric-deploy-tools/fabricNetwork/log"
	"errors"
	"fmt"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	ADDR        = "EHL_ADDR"
	SCP_SOURCES = "EHL_SCP_SOURCE"
)

var addr, scpSource string
var targetPath string
var networkCfg string
var outputPath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fabricNetwork",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&addr, "addr", "", "remote ssh :ip and port")
	log.InitLogger("network.log", "./", "info")
	rootCmd.PersistentFlags().StringVar(&scpSource, "scpSource", "", "scp ")
	// Path of CLI master node to read network config json
	rootCmd.PersistentFlags().StringVar(&networkCfg, "config", "./test-network/networkconfig.json", "")
	// Path of CLI master node where loadConfig command output config files
	rootCmd.PersistentFlags().StringVar(&outputPath, "output", "./test-network", "")
	// Path of peer/orderer node where createChannel/loadConfig command Scp config files to
	rootCmd.PersistentFlags().StringVar(&targetPath, "dPath", "/opt/fabric_install/test-network", "")

	// Convert paths to absolute paths
	outputPath, _ = filepath.Abs(outputPath)
	networkCfg, _ = filepath.Abs(networkCfg)
	targetPath, _ = filepath.Abs(targetPath)
}


// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if networkCfg != "" {
		// Use config file fom the flag.
		viper.SetConfigFile(networkCfg)
	} else {
		// Findhome directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)
		// Search config inhome directory with name ".fabricNetwork" (without extension).
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Printf("Fail to read config file: %s, error message:%v\n", viper.ConfigFileUsed(), err)
		panic(err)
	}
}
//CheckArgs 检查
func CheckArgs() error {
	if addr == "" {
		return errors.New("remote addr  is empty")
	}
	return nil
}
