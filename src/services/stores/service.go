package stores

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/app"
	manifestsmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/manager"
	storesapi "github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/api"
	storesmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/manager"
)

func RegisterService(a *app.App) error {
	// Load manifests service
	m := new(manifestsmanager.Manager)
	err := a.Service(m)
	if err != nil {
		return err
	}

	// Create and register the stores service
	stores := storesmanager.New(*m)
	err = a.RegisterService(stores)
	if err != nil {
		return err
	}

	// Create and register stores API
	storesapi.New(stores).Register(a.Router())

	return nil
}
