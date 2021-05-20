package api

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/handlers"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/middleware"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

const (
	storesPrefix   = "/stores"
	secretsPrefix  = "/stores/" + middleware.StoreURLPlaceholder + "/secrets"
	keysPrefix     = "/stores/" + middleware.StoreURLPlaceholder + "/keys"
	accountsPrefix = "/stores/" + middleware.StoreURLPlaceholder + "/accounts"
	jsonRPCPrefix  = "/nodes"
)

func New(backend core.Backend) http.Handler {
	r := mux.NewRouter()

	r.PathPrefix(secretsPrefix).Handler(
		middleware.StoreSelector(storesPrefix,
			middleware.StripPrefix(secretsPrefix, handlers.NewSecretsHandler(backend)),
		))
	r.PathPrefix(keysPrefix).Handler(
		middleware.StoreSelector(storesPrefix,
			middleware.StripPrefix(keysPrefix, handlers.NewKeysHandler(backend)),
		))
	r.PathPrefix(accountsPrefix).Handler(
		middleware.StoreSelector(storesPrefix,
			middleware.StripPrefix(accountsPrefix, handlers.NewAccountsHandler(backend)),
		))
	r.PathPrefix(storesPrefix).Handler(middleware.StripPrefix(storesPrefix, handlers.NewStoresHandler(backend)))
	r.PathPrefix(jsonRPCPrefix).Methods(http.MethodPost).Handler(middleware.StripPrefix(jsonRPCPrefix, handlers.NewJSONRPCHandler(backend)))

	return middleware.New(backend)(r)
}
