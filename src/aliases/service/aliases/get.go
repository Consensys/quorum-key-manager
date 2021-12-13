package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Aliases) Get(ctx context.Context, registry, key string) (*entities.Alias, error) {
	logger := s.logger.With("registry", registry, "key", key)

	alias, err := s.db.FindOne(ctx, registry, key)
	if err != nil {
		errMessage := "failed to get alias"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	logger.Debug("alias registry retrieved successfully")
	return alias, nil
}
