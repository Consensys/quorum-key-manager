package nodes

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/nodes/interceptor"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
)

func (i *Interactor) Create(ctx context.Context, name string, config *proxynode.Config, allowedTenants []string, _ *entities.UserInfo) error {
	i.mux.Lock()
	defer i.mux.Unlock()

	logger := i.logger.With("name", name, "allowed_tenants", allowedTenants)

	// TODO: Add authorization checks

	if _, ok := i.nodes[name]; ok {
		errMessage := "node already exists"
		logger.Error(errMessage)
		return errors.AlreadyExistsError(errMessage)
	}

	prxNode, err := proxynode.New(config, i.logger)
	if err != nil {
		logger.WithError(err).Error("failed to create node")
		return err
	}

	// Set interceptor on proxy node
	prxNode.Handler = interceptor.New(i.storesService, i.aliasParser, i.logger)

	// Start node
	err = prxNode.Start(ctx)
	if err != nil {
		logger.WithError(err).Error("error starting node")
		return err
	}

	i.createNode(ctx, name, prxNode, allowedTenants)

	logger.Info("node created successfully")
	return nil
}
