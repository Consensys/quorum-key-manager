package app

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	storesapi "github.com/consensys/quorum-key-manager/src/stores/api"
	"github.com/consensys/quorum-key-manager/src/stores/database/postgres"
	storesmanager "github.com/consensys/quorum-key-manager/src/stores/manager"
)

func RegisterService(a *app.App, logger log.Logger) error {
	cfg := new(Config)
	err := a.ServiceConfig(cfg)
	if err != nil {
		return err
	}

	// Create Postgres DB
	postgresClient, err := client.NewClient(cfg.Postgres)
	if err != nil {
		return err
	}
	db := postgres.New(logger, postgresClient)

	// Load manifests service
	m := new(manifestsmanager.Manager)
	err = a.Service(m)
	if err != nil {
		return err
	}

	// Load policy manager service
	policyMngr := new(auth.Manager)
	err = a.Service(policyMngr)
	if err != nil {
		return err
	}

	// Create and register the stores service
	stores := storesmanager.New(*m, *policyMngr, db, logger)
	err = a.RegisterService(stores)
	if err != nil {
		return err
	}

	// Create and register stores API
	storesapi.New(stores).Register(a.Router())

	return nil
}
