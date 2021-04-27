package api

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/handlers"

	jsonrpcapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/middleware"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

const (
	secretsPrefix  = "/secrets"
	keysPrefix     = "/keys"
	accountsPrefix = "/accounts"
	storesPrefix   = "/stores"
	jsonRPCPrefix  = "/"
)

func New(backend core.Backend) http.Handler {
	r := mux.NewRouter()

	r.PathPrefix(secretsPrefix).Handler(middleware.StripPrefix(secretsPrefix, handlers.NewSecretsHandler(backend)))
	r.PathPrefix(keysPrefix).Handler(middleware.StripPrefix(keysPrefix, handlers.NewKeysHandler(backend)))
	r.PathPrefix(accountsPrefix).Handler(middleware.StripPrefix(accountsPrefix, handlers.NewAccountsHandler(backend)))
	r.PathPrefix(storesPrefix).Handler(middleware.StripPrefix(storesPrefix, handlers.NewStoresHandler(backend)))
	r.PathPrefix(jsonRPCPrefix).Methods(http.MethodPost).Handler(middleware.StripPrefix(jsonRPCPrefix, jsonrpcapi.New(backend)))

	return middleware.New(backend)(r)
}
