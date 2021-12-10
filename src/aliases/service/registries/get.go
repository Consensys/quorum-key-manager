package registries

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Registries) Get(ctx context.Context, name string) (*entities.AliasRegistry, error) {
	logger := s.logger.With("name", name)

	registry, err := s.db.FindOne(ctx, name)
	if err != nil {
		return nil, err
	}

	logger.Debug("alias registry retrieved successfully")
	return registry, nil
}
