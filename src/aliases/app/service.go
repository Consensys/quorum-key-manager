package aliasapp

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	aliasapi "github.com/consensys/quorum-key-manager/src/aliases/api"
	aliaspg "github.com/consensys/quorum-key-manager/src/aliases/database/postgres"
	aliasconn "github.com/consensys/quorum-key-manager/src/aliases/interactors/aliases"
	aliasmgr "github.com/consensys/quorum-key-manager/src/aliases/manager"
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

	db := aliaspg.NewDatabase(pgClient, logger)

	aliasSrv, err := aliasconn.NewInteractor(db.Alias(), logger)
	if err != nil {
		return err
	}

	m := aliasmgr.New(aliasSrv)
	err = a.RegisterService(m)
	if err != nil {
		return err
	}

	api := aliasapi.New(aliasSrv)
	api.Register(a.Router())

	return nil
}
