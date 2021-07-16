package aliasservice

import (
	"github.com/go-pg/pg/v10"

	"github.com/consensys/quorum-key-manager/pkg/app"
	aliasstore "github.com/consensys/quorum-key-manager/src/aliases/store"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
)

func RegisterService(a *app.App, logger log.Logger, db database.Database) error {
	// Create and register the stores service

	//TODO replace by the database.Database abstraction
	pgdb := &pg.DB{}
	//TODO replace by the database.Database abstraction
	store := aliasstore.New(pgdb)
	err := a.RegisterService(store)
	if err != nil {
		return err
	}

	return nil
}
