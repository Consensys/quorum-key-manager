package server

import (
	"net/http"
)

func New(addr string, cfg *Config) *http.Server {
	return &http.Server{
		Addr:        addr,
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}
}
