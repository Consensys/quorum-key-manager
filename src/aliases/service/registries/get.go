package registries

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Registries) Get(ctx context.Context, name string, userInfo *auth.UserInfo) (*entities.AliasRegistry, error) {
	logger := s.logger.With("name", name)

	err := s.authorizator.CheckPermission(&auth.Operation{Action: auth.ActionRead, Resource: auth.ResourceAlias})
	if err != nil {
		return nil, err
	}

	registry, err := s.db.FindOne(ctx, name, userInfo.Tenant)
	if err != nil {
		errMessage := "failed to get registry"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	logger.Debug("alias registry retrieved successfully")
	return registry, nil
}
