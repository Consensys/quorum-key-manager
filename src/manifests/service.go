package manifests

import (
	"github.com/consensysquorum/quorum-key-manager/pkg/app"
	"github.com/consensysquorum/quorum-key-manager/pkg/log"
	manifestsmanager "github.com/consensysquorum/quorum-key-manager/src/manifests/manager"
)

func RegisterService(a *app.App, logger log.Logger) error {
	// Load configuration
	cfg := new(manifestsmanager.Config)
	err := a.ServiceConfig(cfg)
	if err != nil {
		return err
	}

	// Create and register the stores service
	manifests, err := manifestsmanager.NewLocalManager(cfg, logger)
	if err != nil {
		return err
	}

	err = a.RegisterService(manifests)
	if err != nil {
		return err
	}

	return nil
}
