package api

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/api/handlers"
	storesmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/manager"
	"github.com/gorilla/mux"
)

type StoresAPI struct {
	stores storesmanager.Manager
}

func New(m storesmanager.Manager) *StoresAPI {
	return &StoresAPI{
		stores: m,
	}
}

func (api *StoresAPI) Register(r *mux.Router) {
	handlers.NewStoresHandler(api.stores).Register(r)
}
