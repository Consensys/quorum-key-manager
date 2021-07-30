package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"

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
		key, err = c.db.Keys().Update(ctx, key)
		if err != nil {
			return err
		}

		_, err = c.store.Update(ctx, id, attr)
		if err != nil && !errors.IsNotSupportedError(err) { // If the underlying store does not support updating, we only update in DB
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	logger.Info("key updated successfully")
	return key, nil
}
