package app

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/api"
	"github.com/consensys/quorum-key-manager/src/aliases/database/postgres"
	interactor "github.com/consensys/quorum-key-manager/src/aliases/interactors/aliases"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	postgresinfra "github.com/consensys/quorum-key-manager/src/infra/postgres"
)

// RegisterService creates and register the alias service in the app.
func RegisterService(a *app.App, logger log.Logger, postgresClient postgresinfra.Client) (aliases.Service, error) {
	// Data layer
	db := postgres.NewDatabase(postgresClient, logger)

	// Business layer
	aliasService, err := interactor.NewInteractor(db.Alias(), logger)
	if err != nil {
		return nil, err
	}

	// Service layer
	api.New(aliasService).Register(a.Router())

	return aliasService, nil
}
