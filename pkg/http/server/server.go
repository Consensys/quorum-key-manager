package server

import (
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

func NewHealthz(cfg *Config) *http.Server {
	return &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.HealthzPort),
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}
}
