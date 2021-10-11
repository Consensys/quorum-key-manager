package imports

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func ImportEthereum(ctx context.Context, db database.ETHAccounts, mnf *manifest.Manifest, logger log.Logger) error {
	logger.Info("importing ethereum accounts...", "store", mnf.Kind, "store_name", mnf.Name)

	store, err := getKeyStore(mnf, logger)
	if err != nil {
		return err
	}

	storeIDs, err := store.List(ctx, 0, 0)
	if err != nil {
		return err
	}

	dbAddresses, err := db.SearchAddresses(ctx, false, 0, 0)
	if err != nil {
		return err
	}

	var n uint
	for _, id := range storeIDs {
		key, err := store.Get(ctx, id)
		if err != nil {
			return err
		}

		if key.IsETHAccount() {
			acc := models.NewETHAccountFromKey(key, &entities.Attributes{})

			if !contains(acc.Address.Hex(), dbAddresses) {
				_, err = db.Add(ctx, acc)
				if err != nil {
					return err
				}

				n++
			}
		}
	}

	logger.Info("ethereum accounts imported successfully", "n", n)
	return nil
}
