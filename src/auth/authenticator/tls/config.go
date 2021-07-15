package tls

import (
	"crypto/x509"
)

type Config struct {
	Certificates []*x509.Certificate
}

type ClaimsConfig struct {
	Username string
	Group    string
}

func NewConfig(certs ...*x509.Certificate) *Config {
	return &Config{
		Certificates: certs,
	}
}
