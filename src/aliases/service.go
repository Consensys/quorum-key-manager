package aliases

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	aliasapi "github.com/consensys/quorum-key-manager/src/aliases/api"
	aliasconn "github.com/consensys/quorum-key-manager/src/aliases/connectors/aliases"
	aliasmgr "github.com/consensys/quorum-key-manager/src/aliases/manager"
	aliaspg "github.com/consensys/quorum-key-manager/src/aliases/store/postgres"
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

	aliasSrv := aliasconn.NewConnector(db, logger)
	m := aliasmgr.New(aliasSrv)
	err = a.RegisterService(m)
	if err != nil {
		return err
	}

	api := aliasapi.New(aliasSrv)
	api.Register(a.Router())

	return nil
}
