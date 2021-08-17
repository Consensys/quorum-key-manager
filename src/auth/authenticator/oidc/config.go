package oidc

import (
	"crypto/x509"
)

type Config struct {
	Certificates []*x509.Certificate
	Claims       *ClaimsConfig
}

type ClaimsConfig struct {
	Subject string
	Scope   string
}

func NewConfig(subject, scope string, certs ...*x509.Certificate) *Config {
	return &Config{
		Certificates: certs,
		Claims: &ClaimsConfig{
			Subject: subject,
			Scope:   scope,
		},
	}
}
