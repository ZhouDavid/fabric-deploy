package config

type NetworkConfig struct {
	Orgs        []Org  `yaml:"orgs"`
	OrdererType string `yaml:"ordererType"`
}
type Peers struct {
	IsOrderer         bool   `yaml:"isOrderer"`
	IP                string `yaml:"ip"`
	PeerPort          int    `yaml:"peerPort"`
	SSHPort           int    `yaml:"sshPort"`
	UserName          string `yaml:"username"`
	Password          string `yaml:"password"`
	IsInstallExplorer bool   `yaml:"isInstallExplorer"`
}
type Org struct {
	Name    string  `yaml:"name"`
	Domain  string  `yaml:"domain"`
	UserNum int     `yaml:"userNum"`
	Peers   []Peers `yaml:"peers"`
}

func (o Org) HasOrderer() bool {
	for _, peer := range o.Peers {
		if peer.IsOrderer {
			return true
		}
	}
	return false
}
func (o Org) HasPeer() bool {
	for _, peer := range o.Peers {
		if !peer.IsOrderer {
			return true
		}
	}
	return false
}
