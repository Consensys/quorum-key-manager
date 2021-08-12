package aliasapi

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/api/handlers"
	"github.com/gorilla/mux"
)

type AliasAPI struct {
	alias aliases.Alias
}

func New(alias aliases.Alias) *AliasAPI {
	return &AliasAPI{
		alias: alias,
	}
}

func (api *AliasAPI) Register(r *mux.Router) {
	aliasSub := r.PathPrefix("/aliases").Subrouter()
	handlers.NewAliasHandler(api.alias).Register(aliasSub)
}
