package manifests

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/app"
	manifestsmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/manifests/manager"
)

func RegisterService(a *app.App) error {
	// Load configuration
	cfg := new(manifestsmanager.Config)
	err := a.ServiceConfig(cfg)
	if err != nil {
		return err
	}

	// Create and register the stores service
	manifests, err := manifestsmanager.NewLocalManager(cfg)
	if err != nil {
		return err
	}

	err = a.RegisterService(manifests)
	if err != nil {
		return err
	}

	return nil
}
