package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (s *Aliases) Delete(ctx context.Context, registry, key string, userInfo *auth.UserInfo) error {
	logger := s.logger.With("registry", registry, "key", key)

	resolver := authorizator.New(s.roles.UserPermissions(ctx, userInfo), userInfo.Tenant, logger)
	err := resolver.CheckPermission(&auth.Operation{Action: auth.ActionDelete, Resource: auth.ResourceAlias})
	if err != nil {
		return err
	}

	_, err = s.Get(ctx, registry, key, userInfo)
	if err != nil {
		return err
	}

	err = s.aliasDB.Delete(ctx, registry, key)
	if err != nil {
		errMessage := "failed to delete alias"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	logger.Info("alias deleted successfully")
	return nil
}
