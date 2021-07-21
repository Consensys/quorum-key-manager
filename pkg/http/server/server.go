package server

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func New(cfg *Config) *http.Server {
	return &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}
}

func NewTLS(cfg *Config) *http.Server {
	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.TLSHost, cfg.TLSPort),
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}
	// This will require peer cert during TLS handshake
	server.TLSConfig = &tls.Config{ClientAuth: tls.RequireAnyClientCert}
	return server
}

func NewHealthz(cfg *Config) *http.Server {
	return &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.HealthzPort),
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}
}
