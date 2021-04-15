package main

import "gopkg.in/yaml.v2"

type Config struct {
}
type ExplorerVolume struct {
	PgData      interface{} `yaml:"pgdata"`
	WalletStore string      `yaml:"walletstore"`
}

func main() {

	out, _ := yaml.Marshal(map[string]interface{}{
		"a": &ExplorerVolume{},
	})
	println(string(out))
}
