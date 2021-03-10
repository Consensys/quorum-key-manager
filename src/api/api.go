package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	accountsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/accounts"
	jsonrpcapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/jsonrpc"
	keysapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/keys"
	secretsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/secrets"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

type apiServer struct {
	cfg    *Config
	server *http.Server
}

func New(cfg *Config, bcknd core.Backend) *apiServer {
	// Create HTTP Middleware
	mid := newHTTPMiddleware(bcknd)

	// Create router
	r := mux.NewRouter()
	r.PathPrefix(secretsapi.RoutePrefix).Handler(secretsapi.New(bcknd))
	r.PathPrefix(keysapi.RoutePrefix).Handler(keysapi.New(bcknd))
	r.PathPrefix(accountsapi.RoutePrefix).Handler(accountsapi.New(bcknd))
	r.Path("/jsonrpc").Handler(jsonrpcapi.New(bcknd))

	server := &http.Server{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:     mid(r),
		IdleTimeout: cfg.IdleConnTimeout,
		ReadTimeout: cfg.Timeout,
	}

	return &apiServer{
		cfg:    cfg,
		server: server,
	}
}

func (h *apiServer) Start(ctx context.Context) error {
	logger := log.FromContext(ctx)
	logger.WithField("addr", h.server.Addr).Info("starting server")
	return h.server.ListenAndServe()
}

func (h *apiServer) Stop(ctx context.Context) error {
	return h.server.Close()
}

func (h *apiServer) Close() error {
	return h.server.Close()
}

func (h *apiServer) Error() error {
	panic("implement me")
}
