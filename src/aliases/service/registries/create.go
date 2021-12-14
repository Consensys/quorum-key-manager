package registries

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Registries) Create(ctx context.Context, name string, allowedTenants []string, _ *auth.UserInfo) (*entities.AliasRegistry, error) {
	logger := s.logger.With("name", name)

	registry, err := s.db.Insert(ctx, &entities.AliasRegistry{Name: name, AllowedTenants: allowedTenants})
	if err != nil {
		errMessage := "failed to create registry"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	logger.Info("alias registry created successfully")
	return registry, nil
}
