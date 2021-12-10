package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Aliases) Update(ctx context.Context, registry, key, kind string, value interface{}) (*entities.Alias, error) {
	logger := s.logger.With("registry", registry, "key", key, "type", kind)

	alias, err := entities.NewAlias(registry, key, kind, value)
	if err != nil {
		return nil, err
	}

	a, err := s.db.Update(ctx, registry, alias)
	if err != nil {
		return nil, err
	}

	logger.Info("alias updated successfully")
	return a, nil
}
