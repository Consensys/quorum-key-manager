package aliasapi

import (
	"github.com/gorilla/mux"

	"github.com/consensys/quorum-key-manager/src/aliases/api/handlers"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

type AliasAPI struct {
	alias aliasent.AliasBackend
}

func New(alias aliasent.AliasBackend) *AliasAPI {
	return &AliasAPI{
		alias: alias,
	}
}

func (api *AliasAPI) Register(r *mux.Router) {
	aliasSub := r.PathPrefix("/aliases").Subrouter()
	handlers.NewAliasHandler(api.alias).Register(aliasSub)
}
