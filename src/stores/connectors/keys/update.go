package keys

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

func (c Connector) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	logger := c.logger.With("id", id)
	logger.Debug("updating key")

	key, err := c.db.Keys().Get(ctx, id)
	if err != nil {
		return nil, err
	}
	key.Tags = attr.Tags

	err = c.db.RunInTransaction(ctx, func(dbtx database.Database) error {
		derr := c.db.Keys().Update(ctx, key)
		if derr != nil {
			return derr
		}

		_, derr = c.store.Update(ctx, id, attr)
		if derr != nil {
			return derr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Info("key updated successfully")
	return key.ToEntity(), nil
}
