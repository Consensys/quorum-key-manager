package vaults

import (
	"context"

	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (c *Vaults) Get(ctx context.Context, name string, userInfo *auth.UserInfo) (*entities.Vault, error) {
	logger := c.logger.With("name", name)

	permissions := c.roles.UserPermissions(ctx, userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

	vault, err := c.getVault(ctx, name, resolver)
	if err != nil {
		return nil, err
	}

	logger.Debug("vault found successfully")
	return vault, nil
}
