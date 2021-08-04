package api

import (
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/stores/api/handlers"
	"github.com/gorilla/mux"
)

type StoresAPI struct {
	alias aliases.Alias
}

func New(m aliases.Alias) *StoresAPI {
	return &StoresAPI{
		alias: m,
	}
}

func (api *StoresAPI) Register(r *mux.Router) {
	handlers.NewAliasHandler(api.alias).Register(r)
}
