package cmd

import (
	"Data_Bank/fabric-deploy-tools/fabricNetwork/config"
	"Data_Bank/fabric-deploy-tools/fabricNetwork/utils"
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loadConfigCmd = &cobra.Command{
	Use:   "loadConfig",
	Short: "load fabric network configuration",
	Long:  `load fabric network configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		runLoadConfig()
	},
}

func init(){
	rootCmd.AddCommand(loadConfigCmd)
}

func runLoadConfig() {
	conf := &config.NetworkConfig{}
	if err := viper.Unmarshal(conf); err != nil {
		fmt.Printf("error when unmarshal config file %s,error message:%v", viper.ConfigFileUsed(), err)
		panic(err)
	}
	fmt.Println(conf)
	// Generate config files
	configBuilder := config.NewFabricConfigBuilder(conf)
	cryptoPath, _ := configBuilder.BuildCryptoFile(outputPath)   // crypto.yaml
	configBuilder.BuildDockerComposeFiles(outputPath)            // docker-compose.yaml
	configBuilder.BuildTxFile(outputPath+"/configtx", "genesis") // conigtx.yaml

	// Generate ipmap
	configBuilder.BuildIPMap(outputPath) // ipmap.txt

	// Generate organization crypto
	fmt.Println("Generating organizations' crypto materials...")
	utils.ExecuteCommand("env/fabric_bin/cryptogen", "generate", fmt.Sprintf("--config=%s", cryptoPath), fmt.Sprintf("--output=%s/organizations", outputPath))

	// Generate genesis block
	fmt.Println("Generating genesis block...")
	genesisBlockPath := fmt.Sprintf("%s/system-genesis-block/genesis.block", outputPath)
	os.Rename(outputPath+"/configtx/genesis.yaml", outputPath+"/configtx/configtx.yaml") // Rename genesis to configtx.yaml
	utils.ExecuteCommand("env/fabric_bin/configtxgen", "-profile", "GenesisChannel", "-channelID", "system-channel", "-outputBlock", genesisBlockPath, "-configPath", outputPath+"/configtx")
	os.Rename(outputPath+"/configtx/configtx.yaml", outputPath+"/configtx/genesis.yaml") // Rename back
	
	// Scp all config files to target machines
	// for now, it broadcast config files to every host in the network
	ipMapPath := path.Join(outputPath, "ipmap.txt")
	ipMap, err := config.NewIPMap(ipMapPath)
	if err!=nil {
		logrus.Fatalf("Fail to read from %s ", ipMapPath)
	}
	for _, host:=range ipMap{
		logrus.Infof("Scp config folder %s to %s:%s ... \n", outputPath, host.Domain, targetPath)
		scp(host.UserName, host.Password, host.GetAddress(), outputPath, targetPath, false)
	}
}
