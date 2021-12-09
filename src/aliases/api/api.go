package api

import (
	"github.com/gorilla/mux"

	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/api/http"
)

// AliasAPI expose the alias service as an HTTP REST API
type AliasAPI struct {
	alias aliases.Interactor
}

func New(alias aliases.Interactor) *AliasAPI {
	return &AliasAPI{
		alias: alias,
	}
}

// Register registers HTTP endpoints to the HTTP router.
func (api *AliasAPI) Register(r *mux.Router) {
	http.NewAliasHandler(api.alias).Register(r)
}
