package oicd

import (
	testutils2 "github.com/consensys/quorum-key-manager/pkg/tls/testutils"
)

const (
	defaultUsernameClaim = "auth.username"
	defaultGroupClaim    = "auth.group"
)

type Config struct {
	Certificate       string
	CertificateServer string
	Claims            *ClaimsConfig
}

type ClaimsConfig struct {
	Username string
	Group    string
}

func NewDefaultConfig() *Config {
	return &Config{
		Certificate:  testutils2.OneLineRSACertPEMA,
		Claims: &ClaimsConfig{
			Username: defaultUsernameClaim,
			Group:    defaultGroupClaim,
		},
	}
}
