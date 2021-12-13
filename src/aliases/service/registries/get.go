package registries

import (
	"context"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Registries) Get(ctx context.Context, name string, userInfo *auth.UserInfo) (*entities.AliasRegistry, error) {
	logger := s.logger.With("name", name)

	registry, err := s.db.FindOne(ctx, name, userInfo.Tenant)
	if err != nil {
		return nil, err
	}

	logger.Debug("alias registry retrieved successfully")
	return registry, nil
}
