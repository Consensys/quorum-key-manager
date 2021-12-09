package aliases

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
)

func (i *Aliases) Update(ctx context.Context, registry string, alias *entities.Alias) (*entities.Alias, error) {
	logger := i.logger.With("registry", registry, "key", alias.Key)

	err := alias.Validate()
	if err != nil {
		return nil, err
	}

	a, err := i.db.Update(ctx, registry, alias)
	if err != nil {
		return nil, err
	}

	logger.Info("alias updated successfully")
	return a, nil
}
