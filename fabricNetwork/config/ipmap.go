package config

import (
	"Data_Bank/fabric-deploy-tools/fabricNetwork/utils"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type IPMap map[string]*RoleEntry

// RoleEntry TODO 写入boltdb or leveldb
type RoleEntry struct {
	Domain            string `json:"domain"`
	IP                string `json:"ip"`
	ListenPorts       []int  `json:"ports"`
	SSHPort           int    `json:"sshPort"`
	UserName          string `json:"userName"`
	Password          string `json:"password"`
	Role              string `json:"role"`
	IsInstallExplorer bool   `json:"isInstallExplorer"`
	Org               string `json:"org"`
	OrgDomain         string `json:"orgDomain"`
}

func (r RoleEntry) GetListenPort() int {
	return r.ListenPorts[0]
}

func (r RoleEntry) GetAddress() string {
	return fmt.Sprintf("%s:%d", r.IP, r.SSHPort)
}

func (r RoleEntry) GetHostDomainIP() string {
	return r.Domain + ":" + strconv.Itoa(r.ListenPorts[0])
}

// // outputPath is like ${projectRoot}/test-network, it should contain core.yaml
// func (r RoleEntry) GetPeerEnvVariables(outputPath string) map[string]string {
// 	if r.Role != "peer" {
// 		panic(fmt.Errorf("%s is not a peer!!!", r.Domain))
// 	}
// 	return map[string]string{
// 		"FABRIC_CFG_PATH":             outputPath,
// 		"CORE_PEER_LOCALMSPID":        r.Org, // such as Org1MSP
// 		"CORE_PEER_TLS_ROOTCERT_FILE": fmt.Sprintf("%s/organizations/peerOrganizations/%s/peers/%s/tls/ca.crt", outputPath, r.OrgDomain, r.Domain),
// 		"CORE_PEER_MSPCONFIGPATH":     fmt.Sprintf("%s/organizations/peerOrganizations/%s/users/Admin@%s/msp", outputPath, r.OrgDomain, r.OrgDomain),
// 		"CORE_PEER_ADDRESS":           fmt.Sprintf("%s:%d", r.Domain, r.GetListenPort()),
// 	}
// }

func NewRoleEntry(input string) (*RoleEntry, error) {
	items := strings.Split(input, " ")
	if len(items) != 9 {
		fmt.Errorf("Fail to parse input record, input record cannot be split into 6 parts, %d parts instead", len(items))
		return nil, errors.New("fail to parse input record")
	}
	listenPorts := make([]int, 0)
	ports := strings.Split(items[2], ",")
	for _, port := range ports[:len(ports)-1] {
		if p, err := strconv.Atoi(port); err != nil {
			fmt.Errorf("port parser error %w", err)
			continue
		} else {
			listenPorts = append(listenPorts, p)
		}
	}
	SSHPort, _ := strconv.Atoi(ports[len(ports)-1])
	isInstallExplorer, _ := strconv.ParseBool(items[6])
	return &RoleEntry{
		Domain:            items[0],
		IP:                items[1],
		ListenPorts:       listenPorts,
		SSHPort:           SSHPort,
		UserName:          items[3],
		Password:          items[4],
		Role:              items[5],
		IsInstallExplorer: isInstallExplorer,
		Org:               items[7],
		OrgDomain:         items[8],
	}, nil
}

func NewIPMap(ipmapFileName string) (IPMap, error) {
	if lines, err := utils.ReadLines(ipmapFileName); err != nil {
		fmt.Printf("Fail to read ipmap file from %s, error message: %v\n", ipmapFileName, err)
		return nil, err
	} else {
		m := IPMap{}
		for _, line := range lines {
			elem, err := NewRoleEntry(line)
			if err != nil {
				fmt.Errorf("new Role entry err %s", err)
				continue
			}
			m[elem.Domain] = elem
		}
		return m, nil
	}
}

func (m IPMap) GetPeerFromOrg(orgName string) (*RoleEntry,error) {
	for _, v := range m {
		if strings.HasPrefix(v.Org, orgName) {
			return v,nil
		}
	}
	return nil,fmt.Errorf("Fail to get a peer from org:%s \n", orgName)
}

// GetOrderer return the first role entry that is an orderer
func (m IPMap) GetOrderer() *RoleEntry {
	for _, e := range m {
		if e.Role == "orderer" {
			return e
		}
	}
	return nil
}

func (m IPMap) GetInstallExplorer() (*RoleEntry, error) {
	if m == nil {
		return nil, errors.New("map nil")
	}
	for _, v := range m {
		if v.IsInstallExplorer {
			return v, nil
		}
	}
	return nil, errors.New("no found ")
}

//GetAllDomainIP return []string(doamin:ip)
func (m IPMap) GetAllDomainIP() []string {
	if m == nil {
		return nil
	}
	domainIPs := make([]string, 0)
	for _, v := range m {
		domainIPs = append(domainIPs, v.Domain+":"+v.IP)
	}
	return domainIPs
}
