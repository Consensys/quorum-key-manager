package aliases

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	aliasmgr "github.com/consensys/quorum-key-manager/src/aliases/manager"
	aliaspg "github.com/consensys/quorum-key-manager/src/aliases/store/postgres"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

// RegisterService creates and register the alias service in the app.
func RegisterService(a *app.App, pgClient postgres.Client) error {
	store := aliaspg.New(pgClient)
	m := aliasmgr.NewManager(store)
	err := a.RegisterService(m)
	if err != nil {
		return err
	}

	return nil
}
