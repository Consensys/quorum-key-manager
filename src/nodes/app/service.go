package app

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/auth"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/nodes/api"
	nodesmanager "github.com/consensys/quorum-key-manager/src/nodes/manager"
	"github.com/consensys/quorum-key-manager/src/stores"
)

func RegisterService(a *app.App, logger log.Logger, manifests []entities2.Manifest, storesService stores.Stores, aliasService aliases.Service) error {
	// Business layer
	authManager := new(auth.Manager)
	err := a.Service(authManager)
	if err != nil {
		return err
	}

	nodes := nodesmanager.New(storesService, manifests, *authManager, aliasService, logger)

	// Service layer
	api.New(nodes).Register(a.Router())

	return nil
}
