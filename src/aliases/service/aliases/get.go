package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Aliases) Get(ctx context.Context, registry, key string) (*entities.Alias, error) {
	logger := s.logger.With("registry", registry, "key", key)

	alias, err := s.db.FindOne(ctx, registry, key)
	if err != nil {
		return nil, err
	}

	logger.Debug("alias registry retrieved successfully")
	return alias, nil
}
