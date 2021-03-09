package api

import (
	"net/http"

	accountsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/accounts"
	jsonrpcapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/jsonrpc"
	keysapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/keys"
	secretsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/secrets"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

// New creates the http.Handler processing all http requests
func New(bcknd core.Backend) http.Handler {
	// Create HTTP Middleware
	mid := NewHTTPMiddleware(bcknd)

	// Create router
	r := mux.NewRouter()
	r.PathPrefix("/secrets").Handler(secretsapi.New(bcknd))
	r.PathPrefix("/keys").Handler(keysapi.New(bcknd))
	r.PathPrefix("/accounts").Handler(accountsapi.New(bcknd))
	r.Path("/jsonrpc").Handler(jsonrpcapi.New(bcknd))

	// Return wrapped router
	return mid(r)
}

func NewHTTPMiddleware(bcknd core.Backend) func(http.Handler) http.Handler {
	// TODO: implement the sequence of middlewares to apply before routing
	return func(h http.Handler) http.Handler {
		return h
	}
}
