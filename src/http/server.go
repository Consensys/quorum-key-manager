package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/http/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/http/middleware"
)

type Server struct {
	cfg    *Config
	server *http.Server
}

func NewServer(cfg *Config, bckend core.Backend) *Server {
	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:     middleware.New(bckend)(api.New(bckend)),
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}

	return &Server{
		cfg:    cfg,
		server: server,
	}
}

func (h *Server) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.WithField("addr", h.server.Addr).Info("starting server")
	return h.server.ListenAndServe()
}

func (h *Server) Stop(_ context.Context) error {
	return h.server.Close()
}

func (h *Server) Close() error {
	return h.server.Close()
}

func (h *Server) Error() error {
	return nil
}
