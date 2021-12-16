package aliases

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Aliases) Get(ctx context.Context, registry, key string, userInfo *auth.UserInfo) (*entities.Alias, error) {
	logger := s.logger.With("registry", registry, "key", key)

	resolver := authorizator.New(s.roles.UserPermissions(ctx, userInfo), userInfo.Tenant, logger)
	err := resolver.CheckPermission(&auth.Operation{Action: auth.ActionRead, Resource: auth.ResourceAlias})
	if err != nil {
		return nil, err
	}

	alias, err := s.aliasDB.FindOne(ctx, registry, key, userInfo.Tenant)
	if err != nil {
		errMessage := "failed to get alias"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	logger.Debug("alias registry retrieved successfully")
	return alias, nil
}
