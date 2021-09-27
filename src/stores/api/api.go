package api

import (
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/handlers"
	"github.com/gorilla/mux"
)

type StoresAPI struct {
	stores stores.Stores
	utils  stores.Utilities
}

func New(storesConnector stores.Stores, utilsConnector stores.Utilities) *StoresAPI {
	return &StoresAPI{
		stores: storesConnector,
		utils:  utilsConnector,
	}
}

func (api *StoresAPI) Register(r *mux.Router) {
	handlers.NewStoresHandler(api.stores).Register(r)
	handlers.NewUtilsHandler(api.utils).Register(r)
}
