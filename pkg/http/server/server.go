package server

import (
	"fmt"
	"net/http"
)

func New(cfg *Config) *http.Server {
	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}

	if cfg.TLSConfig != nil {
		server.TLSConfig = cfg.TLSConfig
	}

	return server
}

func NewHealthz(cfg *Config) *http.Server {
	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.HealthzPort),
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}

	return server
}
