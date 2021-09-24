package entities

import "time"

// HashicorpSpecs is the specs format for a Hashicorp Vault client
type HashicorpSpecs struct {
	MountPoint    string
	Address       string
	Token         string
	TokenPath     string
	Namespace     string
	CACert        string
	CAPath        string
	ClientCert    string
	ClientKey     string
	TLSServerName string
	ClientTimeout time.Duration
	RateLimit     float64
	BurstLimit    int
	MaxRetries    int
	SkipVerify    bool
}
