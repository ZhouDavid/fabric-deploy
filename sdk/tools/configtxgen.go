/*
Copyright IBM Corp. 2017 All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"Data_Bank/fabric-deploy-tools/sdk/fabric/common/encoder"
	"Data_Bank/fabric-deploy-tools/sdk/fabric/common/genesisconfig"
	"Data_Bank/fabric-deploy-tools/sdk/fabric/common/metadata"
	"Data_Bank/fabric-deploy-tools/sdk/fabric/common/update"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-config/protolator"
	"github.com/hyperledger/fabric-config/protolator/protoext/ordererext"
	cb "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/hyperledger/fabric/common/channelconfig"
	"github.com/hyperledger/fabric/common/flogging"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("common.tools.configtxgen")

func doOutputBlock(config *genesisconfig.Profile, channelID string, outputBlock string) error {
	pgen, err := encoder.NewBootstrapper(config)
	if err != nil {
		return errors.WithMessage(err, "could not create bootstrapper")
	}
	logger.Info("Generating genesis block")
	if config.Orderer == nil {
		return errors.New("refusing to generate block which is missing orderer section")
	}
	if config.Consortiums != nil {
		logger.Info("Creating system channel genesis block")
	} else {
		if config.Application == nil {
			return errors.New("refusing to generate application channel block which is missing application section")
		}
		logger.Info("Creating application channel genesis block")
	}
	genesisBlock := pgen.GenesisBlockForChannel(channelID)
	logger.Info("Writing genesis block")
	err = writeFile(outputBlock, protoutil.MarshalOrPanic(genesisBlock), 0640)
	if err != nil {
		return fmt.Errorf("error writing genesis block: %s", err)
	}
	return nil
}

func doOutputChannelCreateTx(conf, baseProfile *genesisconfig.Profile, channelID string, outputChannelCreateTx string) error {
	logger.Info("Generating new channel configtx")

	var configtx *cb.Envelope
	var err error
	if baseProfile == nil {
		configtx, err = encoder.MakeChannelCreationTransaction(channelID, nil, conf)
	} else {
		configtx, err = encoder.MakeChannelCreationTransactionWithSystemChannelContext(channelID, nil, conf, baseProfile)
	}
	if err != nil {
		return err
	}

	logger.Info("Writing new channel tx")
	err = writeFile(outputChannelCreateTx, protoutil.MarshalOrPanic(configtx), 0640)
	if err != nil {
		return fmt.Errorf("error writing channel create tx: %s", err)
	}
	return nil
}

func doOutputAnchorPeersUpdate(conf *genesisconfig.Profile, channelID string, outputAnchorPeersUpdate string, asOrg string) error {
	logger.Info("Generating anchor peer update")
	if asOrg == "" {
		return fmt.Errorf("must specify an organization to update the anchor peer for")
	}

	if conf.Application == nil {
		return fmt.Errorf("cannot update anchor peers without an application section")
	}

	original, err := encoder.NewChannelGroup(conf)
	if err != nil {
		return errors.WithMessage(err, "error parsing profile as channel group")
	}
	original.Groups[channelconfig.ApplicationGroupKey].Version = 1

	updated := proto.Clone(original).(*cb.ConfigGroup)

	originalOrg, ok := original.Groups[channelconfig.ApplicationGroupKey].Groups[asOrg]
	if !ok {
		return errors.Errorf("org with name '%s' does not exist in config", asOrg)
	}

	if _, ok = originalOrg.Values[channelconfig.AnchorPeersKey]; !ok {
		return errors.Errorf("org '%s' does not have any anchor peers defined", asOrg)
	}

	delete(originalOrg.Values, channelconfig.AnchorPeersKey)

	updt, err := update.Compute(&cb.Config{ChannelGroup: original}, &cb.Config{ChannelGroup: updated})
	if err != nil {
		return errors.WithMessage(err, "could not compute update")
	}
	updt.ChannelId = channelID

	newConfigUpdateEnv := &cb.ConfigUpdateEnvelope{
		ConfigUpdate: protoutil.MarshalOrPanic(updt),
	}

	updateTx, err := protoutil.CreateSignedEnvelope(cb.HeaderType_CONFIG_UPDATE, channelID, nil, newConfigUpdateEnv, 0, 0)
	if err != nil {
		return errors.WithMessage(err, "could not create signed envelope")
	}

	logger.Info("Writing anchor peer update")
	err = writeFile(outputAnchorPeersUpdate, protoutil.MarshalOrPanic(updateTx), 0640)
	if err != nil {
		return fmt.Errorf("Error writing channel anchor peer update: %s", err)
	}
	return nil
}

func doInspectBlock(inspectBlock string) error {
	logger.Info("Inspecting block")
	data, err := ioutil.ReadFile(inspectBlock)
	if err != nil {
		return fmt.Errorf("could not read block %s", inspectBlock)
	}

	logger.Info("Parsing genesis block")
	block, err := protoutil.UnmarshalBlock(data)
	if err != nil {
		return fmt.Errorf("error unmarshaling to block: %s", err)
	}
	err = protolator.DeepMarshalJSON(os.Stdout, block)
	if err != nil {
		return fmt.Errorf("malformed block contents: %s", err)
	}
	return nil
}

func doInspectChannelCreateTx(inspectChannelCreateTx string) error {
	logger.Info("Inspecting transaction")
	data, err := ioutil.ReadFile(inspectChannelCreateTx)
	if err != nil {
		return fmt.Errorf("could not read channel create tx: %s", err)
	}

	logger.Info("Parsing transaction")
	env, err := protoutil.UnmarshalEnvelope(data)
	if err != nil {
		return fmt.Errorf("Error unmarshaling envelope: %s", err)
	}

	err = protolator.DeepMarshalJSON(os.Stdout, env)
	if err != nil {
		return fmt.Errorf("malformed transaction contents: %s", err)
	}

	return nil
}

func doPrintOrg(t *genesisconfig.TopLevel, printOrg string) error {
	for _, org := range t.Organizations {
		if org.Name == printOrg {
			og, err := encoder.NewOrdererOrgGroup(org)
			if err != nil {
				return errors.Wrapf(err, "bad org definition for org %s", org.Name)
			}

			if err := protolator.DeepMarshalJSON(os.Stdout, &ordererext.DynamicOrdererOrgGroup{ConfigGroup: og}); err != nil {
				return errors.Wrapf(err, "malformed org definition for org: %s", org.Name)
			}
			return nil
		}
	}
	return errors.Errorf("organization %s not found", printOrg)
}

func writeFile(filename string, data []byte, perm os.FileMode) error {
	dirPath := filepath.Dir(filename)
	exists, err := dirExists(dirPath)
	if err != nil {
		return err
	}
	if !exists {
		err = os.MkdirAll(dirPath, 0750)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(filename, data, perm)
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type Configtxgen struct {
	outputBlock                string //The path to write the genesis block to (if set)
	outputChannelCreateTx      string //The path to write a channel creation configtx to (if set)
	channelCreateTxBaseProfile string //"Specifies a profile to consider as the orderer system channel current state to allow modification of non-application parameters during channel create tx generation. Only valid in conjunction with 'outputCreateChannelTx'."
	profile                    string //The profile from configtx.yaml to use for generation.
	configPath                 string //The path containing the configuration to use (if set)
	channelID                  string //The channel ID to use in the configtx
	inspectChannelCreateTx     string //Prints the configuration contained in the transaction at the specified path
	outputAnchorPeersUpdate    string //[DEPRECATED] Creates a config update to update an anchor peer (works only with the default channel creation, and only for the first update)
	asOrg                      string //Performs the config generation as a particular organization (by name), only including values in the write set that org (likely) has privilege to set
	// printOrg                   string //Prints the definition of an organization as JSON. (useful for adding an org to a channel manually)
	inspectBlock string //Prints the configuration contained in the block at the specified path
}

func NewConfigtxgen() *Configtxgen {
	return &Configtxgen{}
}

// func (c *Configtxgen)

func (c *Configtxgen) Exec() {

	if c.channelID == "" && (c.outputBlock != "" || c.outputChannelCreateTx != "" || c.outputAnchorPeersUpdate != "") {
		logger.Fatalf("Missing channelID, please specify it with '-channelID'")
	}

	// // show version
	// if c.*version {
	// 	printVersion()
	// 	os.Exit(0)
	// }

	// don't need to panic when running via command line
	defer func() {
		if err := recover(); err != nil {
			if strings.Contains(fmt.Sprint(err), "Error reading configuration: Unsupported Config Type") {
				logger.Error("Could not find configtx.yaml. " +
					"Please make sure that FABRIC_CFG_PATH or -configPath is set to a path " +
					"which contains configtx.yaml")
				os.Exit(1)
			}
			if strings.Contains(fmt.Sprint(err), "Could not find profile") {
				logger.Error(fmt.Sprint(err) + ". " +
					"Please make sure that FABRIC_CFG_PATH or -configPath is set to a path " +
					"which contains configtx.yaml with the specified profile")
				os.Exit(1)
			}
			logger.Panic(err)
		}
	}()

	logger.Info("Loading configuration")
	err := factory.InitFactories(nil)
	if err != nil {
		logger.Fatalf("Error on initFactories: %s", err)
	}
	var profileConfig *genesisconfig.Profile
	if c.outputBlock != "" || c.outputChannelCreateTx != "" || c.outputAnchorPeersUpdate != "" {
		if c.profile == "" {
			logger.Fatalf("The '-profile' is required when '-outputBlock', '-outputChannelCreateTx', or '-outputAnchorPeersUpdate' is specified")
		}

		if c.configPath != "" {
			profileConfig = genesisconfig.Load(c.profile, c.configPath)
		} else {
			profileConfig = genesisconfig.Load(c.profile)
		}
	}

	var baseProfile *genesisconfig.Profile
	if c.channelCreateTxBaseProfile != "" {
		if c.outputChannelCreateTx == "" {
			logger.Warning("Specified 'channelCreateTxBaseProfile', but did not specify 'outputChannelCreateTx', 'channelCreateTxBaseProfile' will not affect output.")
		}
		if c.configPath != "" {
			baseProfile = genesisconfig.Load(c.channelCreateTxBaseProfile, c.configPath)
		} else {
			baseProfile = genesisconfig.Load(c.channelCreateTxBaseProfile)
		}
	}

	if c.outputBlock != "" {
		if err := doOutputBlock(profileConfig, c.channelID, c.outputBlock); err != nil {
			logger.Fatalf("Error on outputBlock: %s", err)
		}
	}

	if c.outputChannelCreateTx != "" {
		if err := doOutputChannelCreateTx(profileConfig, baseProfile, c.channelID, c.outputChannelCreateTx); err != nil {
			logger.Fatalf("Error on outputChannelCreateTx: %s", err)
		}
	}

	if c.inspectBlock != "" {
		if err := doInspectBlock(c.inspectBlock); err != nil {
			logger.Fatalf("Error on inspectBlock: %s", err)
		}
	}

	if c.inspectChannelCreateTx != "" {
		if err := doInspectChannelCreateTx(c.inspectChannelCreateTx); err != nil {
			logger.Fatalf("Error on inspectChannelCreateTx: %s", err)
		}
	}

	if c.outputAnchorPeersUpdate != "" {
		if err := doOutputAnchorPeersUpdate(profileConfig, c.channelID, c.outputAnchorPeersUpdate, c.asOrg); err != nil {
			logger.Fatalf("Error on inspectChannelCreateTx: %s", err)
		}
	}

	// if c.printOrg != "" {
	// 	var topLevelConfig *genesisconfig.TopLevel
	// 	if c.configPath != "" {
	// 		topLevelConfig = genesisconfig.LoadTopLevel(configPath)
	// 	} else {
	// 		topLevelConfig = genesisconfig.LoadTopLevel()
	// 	}

	// 	if err := doPrintOrg(topLevelConfig, printOrg); err != nil {
	// 		logger.Fatalf("Error on printOrg: %s", err)
	// 	}
	// }
}

func printVersion() {
	fmt.Println(metadata.GetVersionInfo())
}
