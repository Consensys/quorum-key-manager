package api

import (
	"net/http"

	accountsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/accounts"
	jsonrpcapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/jsonrpc"
	keysapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/middleware"
	secretsapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/secrets"
	storesapi "github.com/ConsenSysQuorum/quorum-key-manager/src/api/stores"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core"
	"github.com/gorilla/mux"
)

const (
	secretsPrefix  = "/secrets"
	keysPrefix     = "/keys"
	accountsPrefix = "/accounts"
	storesPrefix   = "/stores"
	jsonRPCPrefix  = "/jsonrpc"
)

func New(bcknd core.Backend) http.Handler {
	r := mux.NewRouter()
	r.PathPrefix(secretsPrefix).Handler(middleware.StripPrefix(secretsPrefix, secretsapi.New(bcknd)))
	r.PathPrefix(keysPrefix).Handler(middleware.StripPrefix(keysPrefix, keysapi.New(bcknd)))
	r.PathPrefix(accountsPrefix).Handler(middleware.StripPrefix(accountsPrefix, accountsapi.New(bcknd)))
	r.PathPrefix(storesPrefix).Handler(middleware.StripPrefix(storesPrefix, storesapi.New(bcknd)))
	r.PathPrefix(jsonRPCPrefix).Handler(middleware.StripPrefix(jsonRPCPrefix, jsonrpcapi.New(bcknd)))
	return middleware.New(bcknd)(r)
}
