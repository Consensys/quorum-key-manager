package aliasmanager

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	aliaspg "github.com/consensys/quorum-key-manager/src/aliases/store/postgres"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

// RegisterService creates and register the alias service in the app.
func RegisterService(a *app.App, pgClient postgres.Client) error {
	db := aliaspg.NewDatabase(pgClient)
	m := New(db)
	err := a.RegisterService(m)
	if err != nil {
		return err
	}

	return nil
}
