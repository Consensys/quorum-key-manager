package api

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/secrets"
	"net/http"

	accountsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/accounts"
	jsonrpcapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/jsonrpc"
	keysapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

// New creates the http.Handler processing all http requests
func New(backend core.Backend) http.Handler {
	// Create HTTP Middleware
	mid := NewHTTPMiddleware(backend)

	// Create router
	r := mux.NewRouter()
	secrets.New(backend.StoreManager()).Append(r)

	r.PathPrefix("/keys").Handler(keysapi.New(backend))
	r.PathPrefix("/accounts").Handler(accountsapi.New(backend))
	r.Path("/jsonrpc").Handler(jsonrpcapi.New(backend))

	// Return wrapped router
	return mid(r)
}

func NewHTTPMiddleware(bcknd core.Backend) func(http.Handler) http.Handler {
	// TODO: implement the sequence of middlewares to apply before routing
	return func(h http.Handler) http.Handler {
		return h
	}
}
