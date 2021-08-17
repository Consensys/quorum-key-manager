package tls

import (
	"crypto/x509"
)

type Config struct {
	CAs []*x509.Certificate
}

func NewConfig(cas []*x509.Certificate) *Config {
	return &Config{CAs: cas}
}
