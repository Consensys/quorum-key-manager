package aliases

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Aliases) Update(ctx context.Context, registry, key, kind string, value interface{}, userInfo *auth.UserInfo) (*entities.Alias, error) {
	logger := s.logger.With("registry", registry, "key", key, "type", kind)

	alias, err := entities.NewAlias(registry, key, kind, value)
	if err != nil {
		return nil, err
	}

	_, err = s.Get(ctx, registry, key, userInfo)
	if err != nil {
		return nil, err
	}

	a, err := s.aliasDB.Update(ctx, alias)
	if err != nil {
		errMessage := "failed to update alias"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	logger.Info("alias updated successfully")
	return a, nil
}
