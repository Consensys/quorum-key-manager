package stores

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"

	arrays "github.com/consensys/quorum-key-manager/pkg/common"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
)

func (c *Connector) ImportKeys(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) error {
	logger := c.logger.With("store_name", storeName)
	logger.Info("importing keys...")

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

	db := c.db.Keys(storeName)
	dbIDs, err := db.SearchIDs(ctx, false, 0, 0)
	if err != nil {
		return err
	}

	var nSuccesses uint
	var nFailures uint
	for _, id := range arrays.Diff(storeIDs, dbIDs) {
		secret, err := store.Get(ctx, id)
		if err != nil {
			nFailures++
			continue
		}

		_, err = db.Add(ctx, secret)
		if err != nil {
			nFailures++
			continue
		}

		nSuccesses++
	}

	logger.Info("keys import completed", "n_successes", nSuccesses, "n_failures", nFailures)
	return nil
}
