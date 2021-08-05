package tls

import (
	"crypto/x509"
)

type Config struct {
}

func NewConfig(certs ...*x509.Certificate) *Config {
	return &Config{}
}
