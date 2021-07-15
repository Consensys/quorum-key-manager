package tls

import (
	"crypto/x509"
)

type Config struct {
	Certificates []*x509.Certificate
}

func NewConfig(certs ...*x509.Certificate) *Config {
	return &Config{
		Certificates: certs,
	}
}
