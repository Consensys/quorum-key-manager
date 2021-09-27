package nodes

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	aliasmgr "github.com/consensys/quorum-key-manager/src/aliases/manager"
	"github.com/consensys/quorum-key-manager/src/aliases/placeholder"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	nodesapi "github.com/consensys/quorum-key-manager/src/nodes/api"
	nodesmanager "github.com/consensys/quorum-key-manager/src/nodes/manager"
	"github.com/consensys/quorum-key-manager/src/stores"
)

func RegisterService(a *app.App, logger log.Logger) error {
	// Load manifests service
	manifestManager := new(manifestsmanager.Manager)
	err := a.Service(manifestManager)
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

	aliasParser, err := placeholder.New()
	if err != nil {
		return err
	}

	// Create and register nodes service
	nodes := nodesmanager.New(*storeManager, *manifestManager, *authManager, aliasManager.Aliases, aliasParser, logger)
	err = a.RegisterService(nodes)
	if err != nil {
		return err
	}

	// Create and register nodes API
	nodesapi.New(nodes).Register(a.Router())

	return nil
}
