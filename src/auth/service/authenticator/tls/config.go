package tls

import (
	"crypto/x509"
)

type Config struct {
	CAs *x509.CertPool
}

func NewConfig(cas *x509.CertPool) *Config {
	return &Config{CAs: cas}
}
