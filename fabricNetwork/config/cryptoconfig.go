package config

type CryptoConfig struct {
	OrdererOrgs []OrdererOrg `yaml:"OrdererOrgs"`
	PeerOrgs    []PeerOrg    `yaml:"PeerOrgs"`
}

type OrdererOrg struct {
	Name          string `yaml:"Name"`
	Domain        string `yaml:"Domain"`
	EnableNodeOUs bool   `yaml:"EnableNodeOUs"`
	Specs         []Spec `yaml:"Specs"`
}

type Spec struct {
	Hostname string   `yaml:"Hostname"`
	SANS     []string `yaml:"SANS"`
}

type PeerOrg struct {
	Name          string   `yaml:"Name"`
	Domain        string   `yaml:"Domain"`
	EnableNodeOUs bool     `yaml:"EnableNodeOUs"`
	Template      Template `yaml:"Template"`
	Users         Users    `yaml:"Users"`
}

type Template struct {
	Count int      `yaml:"Count"`
	SANS  []string `yaml:"SANS"`
}

type Users struct {
	Count int `yaml:"Count"`
}
