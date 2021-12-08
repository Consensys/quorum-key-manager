package nodes

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
)

func (i *Nodes) Get(ctx context.Context, name string, userInfo *authtypes.UserInfo) (*proxynode.Node, error) {
	permissions := i.roles.UserPermissions(ctx, userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, i.logger)

	err := resolver.CheckPermission(&authtypes.Operation{Action: authtypes.ActionProxy, Resource: authtypes.ResourceNode})
	if err != nil {
		return nil, err
	}

	node := i.getNode(ctx, name)
	if node == nil {
		errMessage := "node was not found"
		i.logger.Error(errMessage, "name", name)
		return nil, errors.NotFoundError(errMessage)
	}

	err = resolver.CheckAccess(node.AllowedTenants)
	if err != nil {
		return nil, err
	}

	return node.Node, nil
}
