package aliasmanager

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
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

	db := aliaspg.NewDatabase(pgClient)
	m := New(db)
	err = a.RegisterService(m)
	if err != nil {
		return err
	}

	return nil
}
