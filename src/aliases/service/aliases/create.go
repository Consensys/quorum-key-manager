package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (i *Aliases) Create(ctx context.Context, registry string, alias *entities.Alias) (*entities.Alias, error) {
	logger := i.logger.With("registry", registry, "key", alias.Key)

	err := alias.Validate()
	if err != nil {
		return nil, err
	}

	alias, err = i.db.Create(ctx, registry, alias)
	if err != nil {
		return nil, err
	}

	logger.Info("alias created successfully")
	return alias, nil
}
