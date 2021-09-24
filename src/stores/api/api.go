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

func New(stores stores.Stores, utils stores.Utilities) *StoresAPI {
	return &StoresAPI{
		stores: stores,
		utils:  utils,
	}
}

func (api *StoresAPI) Register(r *mux.Router) {
	handlers.NewStoresHandler(api.stores).Register(r)
	handlers.NewUtilsHandler(api.utils).Register(r)
}
