package nodes

import (
	"context"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
)

func (i *Interactor) Get(ctx context.Context, name string, userInfo *authtypes.UserInfo) (*proxynode.Node, error) {
	permissions := i.roles.UserPermissions(ctx, userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, i.logger)

	err := resolver.CheckPermission(&authtypes.Operation{Action: authtypes.ActionProxy, Resource: authtypes.ResourceNode})
	if err != nil {
		return nil, err
	}

	node, err := i.getNode(ctx, name)
	if err != nil {
		return nil, err
	}

	err = resolver.CheckAccess(node.AllowedTenants)
	if err != nil {
		return nil, err
	}

	return node.Node, nil
}
