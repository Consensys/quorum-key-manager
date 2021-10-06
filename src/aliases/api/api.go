package aliasapi

import (
	"github.com/gorilla/mux"

	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/api/handlers"
)

// AliasAPI expose the alias service as an HTTP REST API
type AliasAPI struct {
	alias aliases.Repository
}

func New(repo aliases.Repository) *AliasAPI {
	return &AliasAPI{
		alias: repo,
	}
}

// Register registers HTTP endpoints to the HTTP router.
func (api *AliasAPI) Register(r *mux.Router) {
	handlers.NewAliasHandler(api.alias).Register(r)
}
