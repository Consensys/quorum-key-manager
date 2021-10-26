package app

import (
	manifestreader "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/pkg/app"
	authmanager "github.com/consensys/quorum-key-manager/src/auth/manager"
)

func RegisterService(a *app.App, logger log.Logger) error {
	// Load configuration
	cfg := new(Config)
	err := a.ServiceConfig(cfg)
	if err != nil {
		return err
	}

	manifestReader, err := manifestreader.New(cfg.Manifest)
	if err != nil {
		return err
	}

	// Create and register the stores service
	policyMngr := authmanager.New(manifestReader, logger)
	err = a.RegisterService(policyMngr)
	if err != nil {
		return err
	}

	return nil
}
