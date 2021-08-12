package api

import (
	stores "github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/handlers"
	"github.com/gorilla/mux"
)

type StoresAPI struct {
	stores stores.Manager
}

func New(m stores.Manager) *StoresAPI {
	return &StoresAPI{
		stores: m,
	}
}

func (api *StoresAPI) Register(r *mux.Router) {
	handlers.NewStoresHandler(api.stores).Register(r)
}
