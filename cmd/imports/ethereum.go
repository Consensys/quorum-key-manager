package imports

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

func ImportEthereum(ctx context.Context, db database.ETHAccounts, mnf *manifest.Manifest, logger log.Logger) error {
	store, err := getKeyStore(mnf, logger)
	if err != nil {
		return err
	}

	ids, err := store.List(ctx, 0, 0)
	if err != nil {
		return err
	}

	var total uint
	for _, id := range ids {
		key, err := store.Get(ctx, id)
		if err != nil {
			return err
		}

		// If key already exists in DB, we skip. This allows idempotency of the import script (run it multiple times)
		// No need to treat the error here
		dbKey, _ := db.Get(ctx, id)
		if dbKey != nil {
			logger.Debug("ethereum account already exists, skipping", "id", id)
			continue
		}

		_, err = db.Add(ctx, key)
		if err != nil {
			return err
		}
	}

	logger.Info("ethereum accounts successfully imported", "n", total)
	return nil
}
