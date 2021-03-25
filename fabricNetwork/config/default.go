package config

import "fmt"

// crypto.yaml
const (
	EnableNodeOUs = true
	SANS          = "localhost"
)

// tx policy
func GetDefaultOrgPolicies(orgMSP string) map[string]*Policy {
	return map[string]*Policy{
		"Readers":     &Policy{Type: "Signature", Rule: fmt.Sprintf("OR('%s.admin','%s.peer','%s.client')", orgMSP, orgMSP, orgMSP)},
		"Writers":     &Policy{Type: "Signature", Rule: fmt.Sprintf("OR('%s.admin','%s.client')", orgMSP, orgMSP)},
		"Admins":      &Policy{Type: "Signature", Rule: fmt.Sprintf("OR('%s.admin')", orgMSP)},
		"Endorsement": &Policy{Type: "Signature", Rule: fmt.Sprintf("OR('%s.peer')", orgMSP)},
	}
}


func GetDefaultOrdererOrgPolicies(orgMSP string) map[string]*Policy{
	return map[string]*Policy{
		"Readers":     &Policy{Type: "Signature", Rule: fmt.Sprintf("OR('%s.member')", orgMSP)},
		"Writers":     &Policy{Type: "Signature", Rule: fmt.Sprintf("OR('%s.member')", orgMSP)},
		"Admins":      &Policy{Type: "Signature", Rule: fmt.Sprintf("OR('%s.admin')", orgMSP)},
	}
}

func GetDefaultOrdererPolicies() map[string]*Policy {
	return map[string]*Policy{
		"Readers":         &Policy{Type: "ImplicitMeta", Rule: "ANY Readers"},
		"Writers":         &Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
		"Admins":          &Policy{Type: "ImplicitMeta", Rule: "MAJORITY Admins"},
		"BlockValidation": &Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
	}
}

func GetDefaultApplicationPolicies() map[string]*Policy {
	return map[string]*Policy{
		"Readers":              &Policy{Type: "ImplicitMeta", Rule: "ANY Readers"},
		"Writers":              &Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
		"Admins":               &Policy{Type: "ImplicitMeta", Rule: "MAJORITY Admins"},
		"LifecycleEndorsement": &Policy{Type: "ImplicitMeta", Rule: "MAJORITY Endorsement"},
		"Endorsement":          &Policy{Type: "ImplicitMeta", Rule: "MAJORITY Endorsement"},
	}
}

func GetDefaultChannelPolicies() map[string]*Policy {
	return map[string]*Policy{
		"Readers": &Policy{Type: "ImplicitMeta", Rule: "ANY Readers"},
		"Writers": &Policy{Type: "ImplicitMeta", Rule: "ANY Writers"},
		"Admins":  &Policy{Type: "ImplicitMeta", Rule: "MAJORITY Admins"},
	}
}

// tx capability
func GetDefaultCapabilities() map[string]map[string]bool {
	return map[string]map[string]bool{
		"Channel": {
			"V2_0": true,
		},
		"Orderer": {
			"V2_0": true,
		},
		"Application": {
			"V2_0": true,
		},
	}
}

// tx orderer
const (
	MaxMessageCount   = 10
	BatchTimeout      = "2s"
	AbsoluteMaxBytes  = "99 MB"
	PreferredMaxBytes = "512 KB"
	MaxChannels       = 10
)
