package api

import (
	"github.com/consensys/quorum-key-manager/src/stores/api/handlers"
	storesmanager "github.com/consensys/quorum-key-manager/src/stores/manager"
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
