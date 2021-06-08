package manifests

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/app"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
)

func RegisterService(a *app.App) error {
	// Load configuration
	cfg := new(manager.Config)
	err := a.ServiceConfig(cfg)
	if err != nil {
		return err
	}

	// Create and register the stores service
	manifests, err := manager.NewLocalManager(cfg)
	if err != nil {
		return err
	}

	err = a.RegisterService(manifests)
	if err != nil {
		return err
	}

	return nil
}
