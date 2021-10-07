package app

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/aliases/api"
	"github.com/consensys/quorum-key-manager/src/aliases/database/postgres"
	"github.com/consensys/quorum-key-manager/src/aliases/interactors/aliases"
	"github.com/consensys/quorum-key-manager/src/aliases/manager"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
)

// RegisterService creates and register the alias service in the app.
func RegisterService(a *app.App, logger log.Logger) error {
	var cfg Config
	err := a.ServiceConfig(&cfg)
	if err != nil {
		return err
	}

	pgClient, err := client.NewClient(cfg.Postgres)
	if err != nil {
		return err
	}

	db := postgres.NewDatabase(pgClient, logger)

	aliasSrv, err := aliases.NewInteractor(db.Alias(), logger)
	if err != nil {
		return err
	}

	m := manager.New(aliasSrv)
	err = a.RegisterService(m)
	if err != nil {
		return err
	}

	api := api.New(aliasSrv)
	api.Register(a.Router())

	return nil
}
