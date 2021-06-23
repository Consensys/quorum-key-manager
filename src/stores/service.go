package stores

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/app"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	manifestsmanager "github.com/consensysquorum/quorum-key-manager/src/manifests/manager"
	storesapi "github.com/consensysquorum/quorum-key-manager/src/stores/api"
	storesmanager "github.com/consensysquorum/quorum-key-manager/src/stores/manager"
)

func RegisterService(a *app.App, logger log.Logger) error {
	// Load manifests service
	m := new(manifestsmanager.Manager)
	err := a.Service(m)
	if err != nil {
		return err
	}

	// Create and register the stores service
	stores := storesmanager.New(*m, logger)
	err = a.RegisterService(stores)
	if err != nil {
		return err
	}

	// Create and register stores API
	storesapi.New(stores).Register(a.Router())

	return nil
}
