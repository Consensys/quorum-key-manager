package stores

import (
	"context"

	arrays "github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/auth/authorizator"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/stores/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c *Connector) ImportEthereum(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) error {
	logger := c.logger.With("store_name", storeName)
	logger.Info("importing ethereum accounts...")

	// TODO: Uncomment when authManager no longer a runnable
	// permissions := c.authManager.UserPermissions(userInfo)
	resolver := authorizator.New(userInfo.Permissions, userInfo.Tenant, c.logger)

	store, err := c.getKeyStore(ctx, storeName, resolver)
	if err != nil {
		return err
	}

	storeIDs, err := store.List(ctx, 0, 0)
	if err != nil {
		return err
	}

	db := c.db.ETHAccounts(storeName)
	dbAddresses, err := db.SearchAddresses(ctx, false, 0, 0)
	if err != nil {
		return err
	}
	addressMap := arrays.ToMap(dbAddresses)

	var nSuccesses uint
	var nFailures uint
	for _, id := range storeIDs {
		key, err := store.Get(ctx, id)
		if err != nil {
			nFailures++
			continue
		}

		if !key.IsETHAccount() {
			continue
		}

		acc := models.NewETHAccountFromKey(key, &entities.Attributes{})
		if _, found := addressMap[acc.Address.Hex()]; !found {
			_, err = db.Add(ctx, acc)
			if err != nil {
				nFailures++
				continue
			}

			nSuccesses++
		}
	}

	logger.Info("ethereum accounts import completed", "n_successes", nSuccesses, "n_failures", nFailures)
	return nil
}
