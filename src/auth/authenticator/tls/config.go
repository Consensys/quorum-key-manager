package tls

import (
	"crypto/x509"
)

type Config struct {
	CertPool *x509.CertPool
}

func NewConfig(certPool *x509.CertPool) *Config {
	return &Config{CertPool: certPool}
}
