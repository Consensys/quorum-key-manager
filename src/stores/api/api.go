package api

import (
	handlers2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/handlers"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager"
	"github.com/gorilla/mux"
)

type StoresAPI struct {
	stores storemanager.Manager
}

func New(m storemanager.Manager) *StoresAPI {
	return &StoresAPI{
		stores: m,
	}
}

func (api *StoresAPI) Register(r *mux.Router) {
	handlers2.NewStoresHandler(api.stores).Register(r)
}
