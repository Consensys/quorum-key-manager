package stores

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	authmanager "github.com/consensys/quorum-key-manager/src/auth/policy"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	storesapi "github.com/consensys/quorum-key-manager/src/stores/api"
	storesmanager "github.com/consensys/quorum-key-manager/src/stores/manager"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
)

func RegisterService(a *app.App, logger log.Logger, db database.Database) error {
	// Load manifests service
	m := new(manifestsmanager.Manager)
	err := a.Service(m)
	if err != nil {
		return err
	}
	
	// Load policy manager service
	policyMngr := new(authmanager.Manager)
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
