package oidc

import (
	"crypto/x509"
)

type Config struct {
	Certificates []*x509.Certificate
	Claims       *ClaimsConfig
}

type ClaimsConfig struct {
	Username string
	Group    string
}

func NewConfig(usernameClaim, groupClaims string, certs ...*x509.Certificate) *Config {
	return &Config{
		Certificates: certs,
		Claims: &ClaimsConfig{
			Username: usernameClaim,
			Group:    groupClaims,
		},
	}
}
