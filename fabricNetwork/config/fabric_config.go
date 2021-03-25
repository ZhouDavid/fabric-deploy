package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

func createEmptyDockerComposeFile(outputRoot string, peer *PeerInfo) string {
	dir := ""
	if peer.IsOrderer {
		dir = path.Join(outputRoot, "docker-compose", "orderers")
	} else {
		dir = path.Join(outputRoot, "docker-compose", "peers")
	}
	os.MkdirAll(dir, os.ModePerm)
	return path.Join(dir, fmt.Sprintf("%s.yaml", peer.IP))
}

type PeerInfo struct {
	OrgInfo           *OrgInfo
	PeerName          string // peer0, orderer0
	IP                string
	IsOrderer         bool
	ListenPort        int
	SSHPort           int
	UserName          string
	Password          string
	IsInstallExplorer bool
}

func (p PeerInfo) GetPorts() string {
	return fmt.Sprintf("%d,%d", p.ListenPort, p.SSHPort)
}

func (p *PeerInfo) createEmptyDockerComposeFile(outputRoot string) string {
	dir := ""
	if p.IsOrderer {
		dir = path.Join(outputRoot, "docker-compose", "orderers")
	} else {
		dir = path.Join(outputRoot, "docker-compose", "peers")
	}
	os.MkdirAll(dir, os.ModePerm)
	return path.Join(dir, fmt.Sprintf("%s.yaml", p.IP))

}

func (p PeerInfo) GetAddress() string {
	return fmt.Sprintf("%s:%d", p.GetPeerDomain(), p.ListenPort)
}

func (p PeerInfo) GetPeerDomain() string {
	return p.PeerName + "." + p.GetOrgDomain()
}

func (p PeerInfo) GetOrgDomain() string {
	return p.OrgInfo.Domain
}

// GetTLSCertPath is used to form consenters' tls cert path in configtx.yaml
// So use relative path
func (p PeerInfo) GetTLSCertPath() string {
	if p.IsOrderer {
		return fmt.Sprintf("../organizations/ordererOrganizations/%s/orderers/%s/tls/server.crt", p.GetOrgDomain(), p.GetPeerDomain())
	}
	return fmt.Sprintf("../organizations/peerOrganizations/%s/peers/%s/tls/server.crt", p.GetOrgDomain(), p.GetPeerDomain())
}

// // GetTLSCACertPath is used to form orderers' tls ca cert path 
// // when calling "peer creat command" at remote server
// func (p PeerInfo) GetTLSCACertPath() string {
// 	if p.IsOrderer{
// 		return fmt.Sprintf("")
// 	}
// }

type OrgInfo struct {
	Name       string
	Domain     string
	Orderers   []*PeerInfo
	PeerNum    int
	UserNum    int
	HasPeer    bool
	HasOrderer bool
}

func (o OrgInfo) GetMSPID() string {
	return o.Name + "MSP"
}
func (o OrgInfo) GetOrderers() []*PeerInfo {
	return o.Orderers
}

// GetOrdererAddresses get all orderer endpoints in an organization
func (o OrgInfo) GetOrdererAddresses() []string {
	addresses := make([]string, 0)
	for _, orderer := range o.Orderers {
		addresses = append(addresses, orderer.GetAddress())
	}
	return addresses
}

func (o OrgInfo) GetOrdererSpecs() []Spec {
	specs := make([]Spec, 0)
	for _, orderer := range o.Orderers {
		specs = append(specs, Spec{
			Hostname: orderer.PeerName,
			SANS:     []string{SANS},
		})
	}
	return specs
}

func (o OrgInfo) GetMSPDir() string {
	if o.HasPeer {
		return fmt.Sprintf("../organizations/peerOrganizations/%s/msp", o.Domain)
	}
	return fmt.Sprintf("../organizations/ordererOrganizations/%s/msp", o.Domain)
}

func (o OrgInfo) GetOrdererConsenters() []Consenter {
	consenters := make([]Consenter, 0)
	for _, orderer := range o.Orderers {
		consenters = append(consenters, Consenter{
			Host:          orderer.GetPeerDomain(),
			Port:          orderer.ListenPort,
			ClientTLSCert: orderer.GetTLSCertPath(),
			ServerTLSCert: orderer.GetTLSCertPath(),
		})
	}
	return consenters
}

// FabricConfigBuilder builds crypto.yaml, docker-compose.yaml and configtx.yaml
type FabricConfigBuilder interface {
	// build crypto.yaml
	SetCryptoOrdererOrgs() FabricConfig
	SetCryptoPeerOrgs() FabricConfig
	BuildCryptoFile(string) (string, error)

	// build docker-copmose.yaml
	BuildDockerComposeFiles(string)

	// build configtx.yaml
	SetTxOrganizations() FabricConfig
	SetTxCapabilities() FabricConfig
	SetTxApplication(...string) FabricConfig
	SetTxOrderer() FabricConfig
	SetTxChannel() FabricConfig
	SetTxProfiles() FabricConfig
	BuildTxFile(string, string, ...string) (string, error)

	// build ipmap.txt
	BuildIPMap(string) error
}

func NewFabricConfigBuilder(config *NetworkConfig) FabricConfig {
	fmt.Println("Parsing network config...")
	f := FabricConfig{}
	f.OrdererType = config.OrdererType
	if f.OrgInfos == nil {
		f.OrgInfos = make([]*OrgInfo, 0)
	}
	if f.PeerInfos == nil {
		f.PeerInfos = make([]*PeerInfo, 0)
	}
	for _, org := range config.Orgs {
		fmt.Printf("Parsing org:%s ...\n", org.Name)
		orgInfo := &OrgInfo{
			Name:       org.Name,
			Domain:     org.Domain,
			PeerNum:    len(org.Peers),
			UserNum:    org.UserNum,
			HasPeer:    org.HasPeer(),
			HasOrderer: org.HasOrderer(),
			Orderers:   make([]*PeerInfo, 0),
		}
		f.OrgInfos = append(f.OrgInfos, orgInfo)
		peerID := 0
		ordererID := 0
		for _, peer := range org.Peers {
			var name string
			if peer.IsOrderer {
				name = fmt.Sprintf("orderer%d", ordererID)
				ordererID++
			} else {
				name = fmt.Sprintf("peer%d", peerID)
				peerID++
			}

			peerInfo := &PeerInfo{
				OrgInfo:           orgInfo,
 				PeerName:          name,
 				IP:                peer.IP,
 				ListenPort:        peer.PeerPort,
 				SSHPort:           peer.SSHPort,
 				IsOrderer:         peer.IsOrderer,
 				UserName:          peer.UserName,
 				Password:          peer.Password,
 				IsInstallExplorer: peer.IsInstallExplorer,
			}
			if peer.IsOrderer {
				orgInfo.Orderers = append(orgInfo.Orderers, peerInfo)
			}
			f.PeerInfos = append(f.PeerInfos, peerInfo)
		}
	}
	return f
}

type FabricConfig struct {
	OrdererType  string
	PeerInfos    []*PeerInfo
	OrgInfos     []*OrgInfo
	peerOrgs     []*Organization
	ordererOrgs  []*Organization
	cryptoConfig CryptoConfig
	txConfig     TxConfig
}


func (f FabricConfig) BuildIPMap(outputPath string) error {
	fmt.Println("Building ipmap.txt...")
	lines := make([]string, 0)
	for _, peer := range f.PeerInfos {
		role := "peer"
		if peer.IsOrderer {
			role = "orderer"
		}
		lines = append(lines, fmt.Sprintf("%s %s %s %s %s %s %v %s %s", peer.GetPeerDomain(), peer.IP, peer.GetPorts(), peer.UserName, peer.Password, role, peer.IsInstallExplorer, peer.OrgInfo.GetMSPID(), peer.GetOrgDomain()))	}
	return ioutil.WriteFile(path.Join(outputPath, "ipmap.txt"), []byte(strings.Join(lines, "\n")), os.ModePerm) //orderString+"\n"+peerString
}

func (f FabricConfig) GetOrderers() []*PeerInfo {
	orderers := make([]*PeerInfo, 0)
	for _, org := range f.OrgInfos {
		for _, orderer := range org.GetOrderers() {
			orderers = append(orderers, orderer)
		}
	}
	return orderers
}

// GetOrdererAddresses gets all orderer addresses defined in network configuration
func (f FabricConfig) GetOrdererAddresses() []string {
	addresses := make([]string, 0)
	for _, org := range f.OrgInfos {
		addresses = append(addresses, org.GetOrdererAddresses()...)
	}
	return addresses
}

// GetPeerOrgs gets all organizations that has at least 1 peer defined in network configuration
func (f FabricConfig) GetPeerOrgs() []*OrgInfo {
	peerOrgs := make([]*OrgInfo, 0)
	for _, orgInfo := range f.OrgInfos {
		if orgInfo.HasPeer {
			peerOrgs = append(peerOrgs, orgInfo)
		}
	}
	return peerOrgs
}

// GetOrdererOrgs gets all organizations that has at least 1 orderer defined in network configuration
func (f FabricConfig) GetOrdererOrgs() []*OrgInfo {
	ordererOrgs := make([]*OrgInfo, 0)
	for _, orgInfo := range f.OrgInfos {
		if orgInfo.HasOrderer {
			ordererOrgs = append(ordererOrgs, orgInfo)
		}
	}
	return ordererOrgs
}

func (f FabricConfig) GetOrdererConsenters() []Consenter {
	consenters := make([]Consenter, 0)
	for _, org := range f.OrgInfos {
		consenters = append(consenters, org.GetOrdererConsenters()...)
	}
	return consenters
}

func (f FabricConfig) SetCryptoOrdererOrgs() FabricConfig {
	for _, org := range f.GetOrdererOrgs() {
		f.cryptoConfig.OrdererOrgs = append(f.cryptoConfig.OrdererOrgs, OrdererOrg{
			Name:          org.Name,
			Domain:        org.Domain,
			EnableNodeOUs: EnableNodeOUs,
			Specs:         org.GetOrdererSpecs(),
		})
	}
	return f
}
func (f FabricConfig) SetCryptoPeerOrgs() FabricConfig {
	for _, org := range f.GetPeerOrgs() {
		f.cryptoConfig.PeerOrgs = append(f.cryptoConfig.PeerOrgs, PeerOrg{
			Name:          org.Name,
			Domain:        org.Domain,
			EnableNodeOUs: EnableNodeOUs,
			Template:      Template{Count: org.PeerNum, SANS: []string{"localhost"}},
			Users:         Users{Count: org.UserNum},
		})
	}
	return f
}

func (f FabricConfig) BuildCryptoFile(outputRoot string) (string, error) {
	fmt.Println("Building crypto.yaml...")
	os.MkdirAll(outputRoot, os.ModePerm)
	f = f.SetCryptoOrdererOrgs().SetCryptoPeerOrgs()
	crypto, err := yaml.Marshal(&f.cryptoConfig)
	if err != nil {
		return "", err
	}
	cryptoPath := path.Join(outputRoot, "crypto.yaml")
	if err := ioutil.WriteFile(cryptoPath, crypto, os.ModePerm); err != nil {
		fmt.Printf("Fail to write crypto.yaml, error message:%v\n", err)
		return "", err
	}
	return cryptoPath, nil
}

func (f FabricConfig) buildDockerComposeFile(peer *PeerInfo, outputPath string) error {
	var yaml []byte
	hosts := make([]string, 0)
	filename := peer.createEmptyDockerComposeFile(outputPath)

	if peer.IsOrderer {
		// Get all orderer hosts
		for _, orderer := range f.GetOrderers() {
			hosts = append(hosts, fmt.Sprintf("%s:%s", orderer.GetPeerDomain(), orderer.IP))
		}
		yaml = NewOrdererComposeService(peer.GetOrgDomain(), peer.GetPeerDomain(), hosts)
	} else {
		// Get all peer hosts
		for _, peer := range f.PeerInfos {
			hosts = append(hosts, fmt.Sprintf("%s:%s", peer.GetPeerDomain(), peer.IP))
		}
		ports := []string{"7051:7051", "7052:7052"}
		yaml = NewPeerComposeService(peer.GetOrgDomain(), peer.GetPeerDomain(), peer.OrgInfo.GetMSPID(), hosts, ports)
	}
	return ioutil.WriteFile(filename, yaml, os.ModePerm)
}

func (f FabricConfig) BuildDockerComposeFiles(outputRoot string) {
	fmt.Println("Building docker-compose files...")
	for _, peer := range f.PeerInfos {
		if err := f.buildDockerComposeFile(peer, outputRoot); err != nil {
			fmt.Printf("Fail to build docker-compose file for %s, error message: %v", peer.PeerName, err)
		} else {
			fmt.Printf("Build docker-compose file for %s successfully\n", peer.PeerName)
		}
	}
}

func (f FabricConfig) SetTxOrganizations() FabricConfig {
	orgs := make([]*Organization, 0)
	peerOrgs := make([]*Organization, 0)
	ordererOrgs := make([]*Organization, 0)
	for _, orgInfo := range f.OrgInfos {
		org := &Organization{
			Name:     orgInfo.Name,
			ID:       orgInfo.GetMSPID(),
			MSPDir:   orgInfo.GetMSPDir(),
			Policies: GetDefaultOrgPolicies(orgInfo.GetMSPID()),
		}
		orgs = append(orgs, org)
		if orgInfo.HasOrderer {
			org.Policies = GetDefaultOrdererOrgPolicies(orgInfo.GetMSPID())
			ordererOrgs = append(ordererOrgs, org)
		}
		if orgInfo.HasPeer {
			peerOrgs = append(peerOrgs, org)
		}
	}
	f.peerOrgs = peerOrgs
	f.ordererOrgs = ordererOrgs
	f.txConfig.Organizations = orgs

	return f
}

func (f FabricConfig) SetTxCapabilities() FabricConfig {
	f.txConfig.Capabilities = GetDefaultCapabilities()
	return f
}

func (f FabricConfig) SetTxApplication(orgNames ...string) FabricConfig {
	// if orgNames is not given then means it's a genesis channel, no need to setup application section
	if len(orgNames) == 0 {
		return f
	}
	application := &Application{}
	orgSet := make(map[string]interface{})
	application.Organizations = make([]*Organization, 0)
	for _, name := range orgNames {
		orgSet[name] = nil
	}
	for _, org := range f.txConfig.Organizations {
		if _, ok := orgSet[org.Name]; ok {
			application.Organizations = append(application.Organizations, org)
		}
	}
	application.Capabilities = f.txConfig.Capabilities["Application"]
	application.Policies = GetDefaultApplicationPolicies()
	f.txConfig.Application = application
	return f
}

func (f FabricConfig) SetTxOrderer() FabricConfig {
	orderer := &Orderer{}
	orderer.OrdererType = f.OrdererType
	orderer.Addresses = f.GetOrdererAddresses()
	orderer.BatchTimeout = BatchTimeout
	orderer.BatchSize = BatchSize{
		MaxMessageCount:   MaxMessageCount,
		AbsoluteMaxBytes:  AbsoluteMaxBytes,
		PreferredMaxBytes: PreferredMaxBytes,
	}
	orderer.MaxChannels = MaxChannels
	orderer.Capabilities = f.txConfig.Capabilities["Orderer"]
	orderer.Policies = GetDefaultOrdererPolicies()
	switch f.OrdererType {
	case "etcdraft":
		orderer.EtcdRaft.Consenters = f.GetOrdererConsenters()
	}
	orderer.Organizations = f.ordererOrgs
	f.txConfig.Orderer = orderer
	return f
}

func (f FabricConfig) SetTxChannel() FabricConfig {
	channel := Channel{
		Policies:     GetDefaultChannelPolicies(),
		Capabilities: f.txConfig.Capabilities["Channel"],
	}
	f.txConfig.Channel = channel
	return f
}

func (f FabricConfig) SetTxProfiles() FabricConfig {
	profiles := make(map[string]*Profile)
	// If it's setting up a genesis channel
	if f.txConfig.Application == nil {
		profiles["GenesisChannel"] = &Profile{
			Orderer: f.txConfig.Orderer,
			Consortiums: map[string]*Consortium{
				"SampleConsortium": {
					Organizations: f.peerOrgs,
				},
			},
			Policies:     f.txConfig.Channel.Policies,
			Capabilities: f.txConfig.Channel.Capabilities,
		}
	} else {
		profiles["ApplicationChannel"] = &Profile{
			Consortium:   "SampleConsortium",
			Application:  f.txConfig.Application,
			Policies:     f.txConfig.Channel.Policies,
			Capabilities: f.txConfig.Channel.Capabilities,
		}
	}
	f.txConfig.Profiles = profiles
	return f
}

func (f FabricConfig) BuildTxFile(outputRoot string, channelName string, orgNames ...string) (string, error) {
	fmt.Printf("Building configtx.yaml for channel:%s...\n", channelName)
	if err := os.MkdirAll(outputRoot, os.ModePerm); err != nil {
		fmt.Printf("Fail to make dir:%s, error message:%v\n", outputRoot, err)
	}
	f = f.SetTxOrganizations().SetTxCapabilities().SetTxChannel().SetTxOrderer().SetTxApplication(orgNames...).SetTxProfiles()
	tx, err := yaml.Marshal(&f.txConfig)
	if err != nil {
		fmt.Printf("Fail to marshal configtx.yaml, error message:%v", err)
		return "", nil
	}
	txPath := path.Join(outputRoot, fmt.Sprintf("%s.yaml", channelName))
	ioutil.WriteFile(txPath, tx, os.ModePerm)
	return txPath, nil
}
