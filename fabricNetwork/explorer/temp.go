package explorer

import (
	"Data_Bank/fabric-deploy-tools/fabricNetwork/config"
	"fmt"
	"os"
	"text/template"
)

// parameters: username, password,organizationName
var ExplorerTemplate = `{
	"name": "test-network",
	"version": "1.0.0",
	"client": {
		"tlsEnable": true,
		"adminCredential": {
			"id": "{{ .UserName}}",
			"password": "{{ .Password}}"
		},
		"enableAuthentication": true,
		"organization": "{{ .OrganizationName}}",
		"connection": {
			"timeout": {
				"peer": {
					"endorser": "300"
				},
				"orderer": "300"
			}
		}
	},
	"channels": {
		"mychannel": {
			"peers": {
				"{{ .PeerName}}": {}
			}
		}
	},
	"organizations": {
		"{{ .OrganizationName}}": {
			"mspid": "{{ .OrganizationName}}",
			"adminPrivateKey": {
				"path": "{{ .AdminPrivateKeyPath}}"
			},
			"peers": [
				"{{ .PeerName}}"
			],
			"signedCert": {
				"path": "{{ .SigningCertPath}}"
			}
		}
	},
	"peers": {
		"{{ .PeerName}}": {
			"tlsCACerts": {
				"path": "{{ .PeerTLSCACertsPath}}"
			},
			"url": "{{ .PeerURL}}"
		}
	}
}`

const (
	DefaultUsername            = "admin"
	DefaultPwd                 = "exploreradminpw"
	DefaultOrgName             = "Org1MSP"
	DefaultAdminPrivateKeyPath = "/tmp/crypto/peerOrganizations/%s/users/Admin@%s/msp/keystore/priv_sk" //domain domain
	DefaultPeerName            = "peer0.org1.example.com"
	DefaultSingCertPath        = "/tmp/crypto/peerOrganizations/%s/users/Admin@%s/msp/signcerts/Admin@%s-cert.pem" //domain domain domain
	DefaultPeerTLSCaCertsPath  = "/tmp/crypto/peerOrganizations/%s/peers/%s/tls/ca.crt"                            //domain peerDomain
	DefaultPeerURL             = "grpcs://%s:7051"                                                                 // peer
)

//ExplorerInput input
type ExplorerInput struct {
	UserName            string
	Password            string
	OrganizationName    string
	AdminPrivateKeyPath string
	PeerName            string
	SigningCertPath     string
	PeerTLSCACertsPath  string
	PeerURL             string
}

//New default explorer
// func New() *ExplorerInput {

// 	return &ExplorerInput{
// 		UserName:            DefaultUsername,
// 		Password:            "ehl1234",
// 		OrganizationName:    DefaultOrgName,
// 		AdminPrivateKeyPath: DefaultAdminPrivateKeyPath,
// 		PeerName:            DefaultPeerName,
// 		SigningCertPath:     DefaultSingCertPath,
// 		PeerTLSCACertsPath:  DefaultPeerTLSCaCertsPath,
// 		PeerURL:             DefaultPeerURL,
// 	}
// }

//CreateExplorerConfig test
func (i *ExplorerInput) CreateExplorerConfig(file string) error {
	temp, err := template.New("explorer").Parse(ExplorerTemplate)
	if err != nil {
		fmt.Println("template err")
	}

	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return temp.Execute(f, i)
}

func NewInput(entry *config.RoleEntry) *ExplorerInput {

	return &ExplorerInput{
		UserName:            "admin",
		Password:            "ehl1234",
		OrganizationName:    entry.Org,
		AdminPrivateKeyPath: fmt.Sprintf(DefaultAdminPrivateKeyPath, entry.OrgDomain, entry.OrgDomain),
		PeerName:            entry.Domain,
		SigningCertPath:     fmt.Sprintf(DefaultSingCertPath, entry.OrgDomain, entry.OrgDomain, entry.OrgDomain),
		PeerTLSCACertsPath:  fmt.Sprintf(DefaultPeerTLSCaCertsPath, entry.OrgDomain, entry.Domain),
		PeerURL:             fmt.Sprintf(DefaultPeerURL, entry.Domain),
	}
}
