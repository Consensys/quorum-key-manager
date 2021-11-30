package app

import (
	"github.com/consensys/quorum-key-manager/pkg/app"
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/nodes/api"
	"github.com/consensys/quorum-key-manager/src/nodes/service/nodes"
	"github.com/consensys/quorum-key-manager/src/stores"
)

func RegisterService(
	a *app.App,
	logger log.Logger,
	authService auth.Roles,
	storesService stores.Stores,
	aliasService aliases.Service,
) *nodes.Nodes {
	// Business layer
	nodesService := nodes.New(storesService, authService, aliasService, logger)

	// Service layer
	api.New(nodesService).Register(a.Router())

	return nodesService
}
