package aliasapi

import (
	"github.com/gorilla/mux"

	"github.com/consensys/quorum-key-manager/src/aliases/api/handlers"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

// AliasAPI expose the alias service as an HTTP REST API
type AliasAPI struct {
	alias aliasent.AliasBackend
}

func New(alias aliasent.AliasBackend) *AliasAPI {
	return &AliasAPI{
		alias: alias,
	}
}

// Register registers HTTP endpoints to the HTTP router.
func (api *AliasAPI) Register(r *mux.Router) {
	handlers.NewAliasHandler(api.alias).Register(r)
}
