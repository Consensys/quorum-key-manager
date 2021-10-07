package nodes

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	aliasmgr "github.com/consensys/quorum-key-manager/src/aliases/manager"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/infra/manifests/filesystem"
	nodesapi "github.com/consensys/quorum-key-manager/src/nodes/api"
	nodesmanager "github.com/consensys/quorum-key-manager/src/nodes/manager"
	"github.com/consensys/quorum-key-manager/src/stores"
)

func RegisterService(a *app.App, logger log.Logger) error {
	cfg := new(Config)
	err := a.ServiceConfig(cfg)
	if err != nil {
		return err
	}

	// Create manifest reader
	manifestReader, err := manifestsmanager.New(cfg.Manifest)
	if err != nil {
		return err
	}

	// Load stores service
	storeManager := new(stores.Manager)
	err = a.Service(storeManager)
	if err != nil {
		return err
	}

	// Load auth manager service
	authManager := new(auth.Manager)
	err = a.Service(authManager)
	if err != nil {
		return err
	}

	aliasManager := new(aliasmgr.BaseManager)
	err = a.Service(aliasManager)
	if err != nil {
		return err
	}

	// Create and register nodes service
	nodes := nodesmanager.New(*storeManager, manifestReader, *authManager, aliasManager.Aliases, logger)
	err = a.RegisterService(nodes)
	if err != nil {
		return err
	}

	// Create and register nodes API
	nodesapi.New(nodes).Register(a.Router())

	return nil
}
