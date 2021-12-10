package app

import (
	"github.com/consensys/quorum-key-manager/src/aliases/api/http"
	db "github.com/consensys/quorum-key-manager/src/aliases/database/postgres"
	"github.com/consensys/quorum-key-manager/src/aliases/service/aliases"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/gorilla/mux"

	"github.com/consensys/quorum-key-manager/src/infra/log"
)

// RegisterService creates and register the alias service in the app.
func RegisterService(router *mux.Router, logger log.Logger, postgresClient postgres.Client) *aliases.Aliases {
	// Data layer
	aliasDB := db.NewDatabase(postgresClient, logger)

	// Business layer
	aliasService := aliases.New(aliasDB.Alias(), logger)

	// Service layer
	http.NewRegistryHandler(aliasService).Register(router)
	http.NewAliasHandler(aliasService).Register(router)

	return aliasService
}
