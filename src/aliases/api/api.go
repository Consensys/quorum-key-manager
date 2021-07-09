package aliasapi

import (
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/gorilla/mux"
)

type AliasAPI struct {
	store *aliasstore.Store
}

func New(store *aliasstore.Store) *AliasAPI {
	return &AliasAPI{
		store: store,
	}
}

func (a *AliasAPI) Register(r *mux.Router) {
	//TODO: the: implement aliashandlers
	//	aliashandlers.NewStoresHandler(a.store).Register(r)
}
