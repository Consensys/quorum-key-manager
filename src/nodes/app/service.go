package app

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/nodes/api"
	"github.com/consensys/quorum-key-manager/src/nodes/api/manifest"
	"github.com/consensys/quorum-key-manager/src/nodes/service/nodes"
	"github.com/consensys/quorum-key-manager/src/stores"
)

func RegisterService(
	ctx context.Context,
	a *app.App,
	logger log.Logger,
	manifests map[string][]entities.Manifest,
	authService auth.Roles,
	storesService stores.Stores,
	aliasService aliases.Service,
) (*nodes.Interactor, error) {
	// Business layer
	nodesService := nodes.New(storesService, authService, aliasService, logger)

	// Service layer
	api.New(nodesService).Register(a.Router())

	manifestNodesHandler := manifest.NewNodesHandler(nodesService) // Manifest reading is synchronous, similar to a config file
	for _, mnf := range manifests[entities.NodeKind] {
		err := manifestNodesHandler.Create(ctx, mnf.Specs)
		if err != nil {
			return nil, err
		}
	}

	return nodesService, nil
}
