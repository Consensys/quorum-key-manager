package nodes

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	"sort"
)

func (i *Nodes) List(ctx context.Context, userInfo *entities.UserInfo) ([]string, error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	var nodeNames []string
	for name, nodeInfo := range i.nodes {
		permissions := i.roles.UserPermissions(ctx, userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, i.logger)

		if err := resolver.CheckAccess(nodeInfo.AllowedTenants); err != nil {
			continue
		}
		nodeNames = append(nodeNames, name)
	}

	sort.Strings(nodeNames)

	return nodeNames, nil
}
