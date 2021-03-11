package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

type Server struct {
	cfg    *Config
	server *http.Server
	logger *log.Logger
}

const Component = "http"

func NewServer(cfg *Config, handler http.Handler, logger *log.Logger) *Server {
	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:     handler,
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}

	return &Server{
		cfg:    cfg,
		server: server,
		logger: logger.SetComponent(Component),
	}
}

func (h *Server) Start(ctx context.Context) error {
	var cerr = make(chan error, 1)
	defer close(cerr)

	go func() {
		h.logger.WithField("addr", h.server.Addr).Info("started server")
		cerr <- h.server.ListenAndServe()
	}()

	select {
	case err := <-cerr:
		if err == http.ErrServerClosed {
			return nil
		}

		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *Server) Stop(ctx context.Context) error {
	h.logger.Info("shutting down server")
	return h.server.Shutdown(ctx)
}

func (h *Server) Close() error {
	return h.server.Close()
}

func (h *Server) Error() error {
	return nil
}
