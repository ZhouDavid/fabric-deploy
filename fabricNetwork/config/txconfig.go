package config

// TxConfig consists of the structs used by the configtxgen tool.
type TxConfig struct {
	Profiles      map[string]*Profile        `yaml:"Profiles"`
	Organizations []*Organization            `yaml:"-"`
	Channel       Channel                    `yaml:"-"`
	Application   *Application               `yaml:"-"`
	Orderer       *Orderer                   `yaml:"-"`
	Capabilities  map[string]map[string]bool `yaml:"-"`
}

// Profile encodes orderer/application configuration combinations for the
// configtxgen tool.
type Profile struct {
	Consortium   string                 `yaml:"Consortium,omitempty"`
	Application  *Application           `yaml:"Application,omitempty"`
	Orderer      *Orderer               `yaml:"Orderer,omitempty"`
	Consortiums  map[string]*Consortium `yaml:"Consortiums,omitempty"`
	Capabilities map[string]bool        `yaml:"Capabilities,omitempty"`
	Policies     map[string]*Policy     `yaml:"Policies,omitempty"`
}

// Policy encodes a channel config policy
type Policy struct {
	Type string `yaml:"Type"`
	Rule string `yaml:"Rule"`
}

type Channel struct {
	Policies     map[string]*Policy `yaml:"Policies"`
	Capabilities map[string]bool    `yaml:"Capabilities"`
}

// Consortium represents a group of organizations which may create channels
// with each other
type Consortium struct {
	Organizations []*Organization `yaml:"Organizations"`
}

// Application encodes the application-level configuration needed in config
// transactions.
type Application struct {
	Organizations []*Organization    `yaml:"Organizations"`
	Capabilities  map[string]bool    `yaml:"Capabilities"`
	Policies      map[string]*Policy `yaml:"Policies"`
	ACLs          map[string]string  `yaml:"ACLs"`
}

// Organization encodes the organization-level configuration needed in
// config transactions.
type Organization struct {
	Name        string             `yaml:"Name"`
	ID          string             `yaml:"ID"`
	MSPDir      string             `yaml:"MSPDir"`
	Policies    map[string]*Policy `yaml:"Policies"`
	AnchorPeers []AnchorPeer       `yaml:"AnchorPeers"`
}

// AnchorPeer encodes the necessary fields to identify an anchor peer.
type AnchorPeer struct {
	Host string `yaml:"Host"`
	Port int    `yaml:"Port"`
}

// Orderer contains configuration associated to a channel.
type Orderer struct {
	OrdererType   string             `yaml:"OrdererType"`
	Addresses     []string           `yaml:"Addresses"`
	BatchTimeout  string             `yaml:"BatchTimeout"`
	BatchSize     BatchSize          `yaml:"BatchSize"`
	EtcdRaft      EtcdRaft           `yaml:"EtcdRaft"`
	Organizations []*Organization    `yaml:"Organizations"`
	MaxChannels   uint64             `yaml:"MaxChannels"`
	Capabilities  map[string]bool    `yaml:"Capabilities"`
	Policies      map[string]*Policy `yaml:"Policies"`
}

// BatchSize contains configuration affecting the size of batches.
type BatchSize struct {
	MaxMessageCount   uint32 `yaml:"MaxMessageCount"`
	AbsoluteMaxBytes  string `yaml:"AbsoluteMaxBytes"`
	PreferredMaxBytes string `yaml:"PreferredMaxBytes"`
}

type Consenter struct {
	Host          string `yaml:"Host"`
	Port          int    `yaml:"Port"`
	ClientTLSCert string `yaml:"ClientTLSCert"`
	ServerTLSCert string `yaml:"ServerTLSCert"`
}
type EtcdRaft struct {
	Consenters []Consenter `yaml:"Consenters"`
}
