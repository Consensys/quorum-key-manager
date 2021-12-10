package registries

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (s *Registries) Create(ctx context.Context, name string) (*entities.AliasRegistry, error) {
	logger := s.logger.With("name", name)

	registry, err := s.db.Insert(ctx, name)
	if err != nil {
		return nil, err
	}

	logger.Info("alias registry created successfully")
	return registry, nil
}
