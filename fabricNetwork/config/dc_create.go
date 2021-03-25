package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// NewOrdererComposeService create
func NewOrdererComposeService(orgDomain string, peerDomain string, hosts []string) []byte {
	var orderer = map[string]interface{}{
		"version": "2",
		"services": map[string]interface{}{
			peerDomain: map[string]interface{}{
				"container_name": peerDomain,
				"image":          "hyperledger/fabric-orderer:2.2.0",
				"working_dir":    "/opt/gopath/src/github.com/hyperledger/fabric",
				"command":        "orderer",
				"volumes": []string{
					fmt.Sprintf("../../system-genesis-block/genesis.block:/var/hyperledger/orderer/orderer.genesis.block"),
					fmt.Sprintf("../../organizations/ordererOrganizations/%s/orderers/%s/msp:/var/hyperledger/orderer/msp", orgDomain, peerDomain),
					fmt.Sprintf("../../organizations/ordererOrganizations/%s/orderers/%s/tls/:/var/hyperledger/orderer/tls", orgDomain, peerDomain),
				},
				"ports":       []string{"7050:7050"},
				"extra_hosts": hosts,
				"environment": []string{
					"FABRIC_LOGGING_SPEC=INFO",
					"ORDERER_GENERAL_LISTENADDRESS=0.0.0.0",
					"ORDERER_GENERAL_BOOTSTRAPMETHOD=file",
					"ORDERER_GENERAL_BOOTSTRAPFILE=/var/hyperledger/orderer/orderer.genesis.block",
					"ORDERER_GENERAL_LOCALMSPID=OrdererMSP",
					"ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp",
					"ORDERER_GENERAL_TLS_ENABLED=true",
					"ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
					"ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
					"ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
					"ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
					"ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
					"ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
				},
			},
		},
	}
	return toYaml(&orderer)
}

func toYaml(input interface{}) []byte {
	yaml, err := yaml.Marshal(input)
	if err != nil {
		panic(err)
	}
	return yaml

}

func NewPeerComposeService(orgDomain, peerDomain, mspID string, hosts []string, ports []string) []byte {
	var peer = map[string]interface{}{
		"version": "2",
		"services": map[string]interface{}{
			peerDomain: map[string]interface{}{
				"container_name": peerDomain,
				"image":          "hyperledger/fabric-peer:2.2.0",
				"working_dir":    "/opt/gopath/src/github.com/hyperledger/fabric/peer",
				"command":        "peer node start",
				"volumes": []string{
					"/var/run/:/host/var/run/",
					fmt.Sprintf("../../organizations/peerOrganizations/%s/peers/%s/msp:/etc/hyperledger/fabric/msp", orgDomain, peerDomain),
					fmt.Sprintf("../../organizations/peerOrganizations/%s/peers/%s/tls:/etc/hyperledger/fabric/tls", orgDomain, peerDomain),
				},
				"ports":       ports,
				"extra_hosts": hosts,
				"environment": []string{
					"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
					fmt.Sprintf("CORE_PEER_ID=%s", peerDomain),
					fmt.Sprintf("CORE_PEER_ADDRESS=%s:7051", peerDomain),
					"CORE_PEER_LISTENADDRESS=0.0.0.0:7051",
					"CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:7052",
					"CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=peers_default",
					fmt.Sprintf("CORE_PEER_CHAINCODEADDRESS=%s:7052", peerDomain),
					fmt.Sprintf("CORE_PEER_GOSSIP_EXTERNALENDPOINT=%s:7051", peerDomain),
					fmt.Sprintf("CORE_PEER_GOSSIP_BOOTSTRAP=%s:7051", peerDomain),
					fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", mspID),
					"FABRIC_LOGGING_SPEC=INFO",
					"CORE_PEER_TLS_ENABLED=true",
					"CORE_PEER_GOSSIP_USELEADERELECTION=true",
					"CORE_PEER_GOSSIP_ORGLEADER=false",
					"CORE_PEER_PROFILE_ENABLED=true",
					"CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt",
					"CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key",
					"CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt",
					"CORE_CHAINCODE_EXECUTETIMEOUT=300s",
				},
			},
			"cli": map[string]interface{}{
				"container_name": "cli",
				"image":          "hyperledger/fabric-tools:2.2.0",
				"tty":            true,
				"stdin_open":     true,
				"environment": []string{
					"GOPROXY=https://goproxy.io,direct",
					"GO111MODULE=on",
					"GOPATH=/opt/gopath",
					"CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
					"FABRIC_LOGGING_SPEC=INFO",
					"CORE_PEER_ID=cli",
					fmt.Sprintf("CORE_PEER_ADDRESS=%s:7051", peerDomain),
					fmt.Sprintf("CORE_PEER_LOCALMSPID=%s", mspID),
					"CORE_PEER_TLS_ENABLED=true",
					fmt.Sprintf("CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/server.crt", orgDomain, peerDomain),
					fmt.Sprintf("CORE_PEER_TLS_KEY_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/server.key", orgDomain, peerDomain),
					fmt.Sprintf("CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/peers/%s/tls/ca.crt", orgDomain, peerDomain),
					fmt.Sprintf("CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/%s/users/Admin@%s/msp", orgDomain, orgDomain),
				},
				"working_dir": "/opt/gopath/src/github.com/hyperledger/fabric/peer",
				"command":     "/bin/bash",
				"volumes": []string{
					"/var/run/:/host/var/run/",
					"../../chaincode/go/:/opt/gopath/src/github.com/hyperledger/chaincode/go",
					"../../organizations:/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/",
					"../../channel-artifacts:/opt/gopath/src/github.com/hyperledger/fabric/peer/channel-artifacts",
					"../../chaincode:/opt/gopath/src/github.com/hyperledger/fabric/peer/chaincode",
				},
				"depends_on": []string{
					peerDomain,
				},
				"extra_hosts": hosts,
			},
		},
	}
	return toYaml(peer)
}

type ExplorerVolume struct {
	PgData      interface{} `yaml:"pgdata"`
	WalletStore interface{} `yaml:"walletstore"`
}

func NewExplorerDockerComposeService(hosts []string, crypto string) []byte {
	var explorerTemp = map[string]interface{}{
		"version": "2.1",
		"volumes": &ExplorerVolume{},
		"services": map[string]interface{}{
			"explorerdb.mynetwork.com": map[string]interface{}{
				"image":          "hyperledger/explorer-db:1.1.4",
				"container_name": "explorerdb.mynetwork.com",
				"hostname":       "explorerdb.mynetwork.com",
				"environment": []string{
					"DATABASE_DATABASE=fabricexplorer",
					"DATABASE_USERNAME=hppoc",
					"DATABASE_PASSWORD=password",
				},
				"healthcheck": map[string]interface{}{
					"test":     "pg_isready -h localhost -p 5432 -q -U postgres",
					"interval": "30s",
					"timeout":  "10s",
					"retries":  5,
				},
				"extra_hosts": hosts,
				"volumes": []string{
					"pgdata:/var/lib/postgresql/data",
				},
			},
			"explorer.mynetwork.com": map[string]interface{}{
				"image":          "hyperledger/explorer:1.1.4",
				"container_name": "explorer.mynetwork.com",
				"hostname":       "explorer.mynetwork.com",
				"depends_on": map[string]interface{}{
					"explorerdb.mynetwork.com": map[string]interface{}{
						"condition": "service_healthy",
					},
				},
				"environment": []string{
					"DATABASE_HOST=explorerdb.mynetwork.com",
					"DATABASE_DATABASE=fabricexplorer",
					"DATABASE_USERNAME=hppoc",
					"DATABASE_PASSWD=password",
					"LOG_LEVEL_APP=debug",
					"LOG_LEVEL_DB=debug",
					"LOG_LEVEL_CONSOLE=info",
					"LOG_CONSOLE_STDOUT=true",
					"DISCOVERY_AS_LOCALHOST=false",
				},
				"volumes": []string{
					"../config.json:/opt/explorer/app/platform/fabric/config.json",
					"../connection-profile:/opt/explorer/app/platform/fabric/connection-profile",
					"walletstore:/opt/explorer/wallet",
					fmt.Sprintf("%s:/tmp/crypto", crypto),
				},
				"extra_hosts": hosts,
				"ports": []string{
					"8080:8080",
				},
			},
		},
	}
	return toYaml(explorerTemp)
}
