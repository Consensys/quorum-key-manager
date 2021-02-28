package api

import (
	"net/http"

	accountsapi "github.com/ConsenSysQuorum/quorum-key-manager/api/accounts"
	jsonrpcapi "github.com/ConsenSysQuorum/quorum-key-manager/api/jsonrpc"
	keysapi "github.com/ConsenSysQuorum/quorum-key-manager/api/keys"
	secretsapi "github.com/ConsenSysQuorum/quorum-key-manager/api/secrets"
	"github.com/ConsenSysQuorum/quorum-key-manager/backend"
	"github.com/gorilla/mux"
)

// New creates the http.Handler processing all http requests
func New(bcknd backend.Backend) http.Handler {
	r := mux.NewRouter()
	r.PathPrefix("/secrets").Handler(secretsapi.New(bcknd))
	r.PathPrefix("/keys").Handler(keysapi.New(bcknd))
	r.PathPrefix("/accounts").Handler(accountsapi.New(bcknd))
	r.Path("/jsonrpc").Handler(jsonrpcapi.New(bcknd))
	return r
}
