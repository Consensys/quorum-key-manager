package stores

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/app"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager"
)

func RegisterService(a *app.App) error {
	// Load manifests service
	m := new(manager.Manager)
	err := a.Service(m)
	if err != nil {
		return err
	}

	// Create and register the stores service
	stores := storemanager.New(*m)
	err = a.RegisterService(stores)
	if err != nil {
		return err
	}

	// Create and register stores API
	api.New(stores).Register(a.Router())

	return nil
}
