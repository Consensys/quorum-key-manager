package stores

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"

	"github.com/consensys/quorum-key-manager/src/auth"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	eth "github.com/consensys/quorum-key-manager/src/stores/connectors/ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/stores"
)

func (c *Connector) Ethereum(ctx context.Context, storeName string, userInfo *authtypes.UserInfo) (stores.EthStore, error) {
	permissions := c.authManager.UserPermissions(userInfo)
	resolver := authorizator.New(permissions, userInfo.Tenant, c.logger)

	store, err := c.getEthStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	c.logger.Debug("ethereum store found successfully", "store_name", storeName)
	return eth.NewConnector(store, c.db.ETHAccounts(storeName), resolver, c.logger), nil
}

func (c *Connector) EthereumByAddr(ctx context.Context, addr common.Address, userInfo *authtypes.UserInfo) (stores.EthStore, error) {
	logger := c.logger.With("address", addr.Hex())

	ethStores, err := c.List(ctx, manifest.Ethereum, userInfo)
	if err != nil {
		return nil, err
	}

	for _, storeName := range ethStores {
		ethStore, err := c.Ethereum(ctx, storeName, userInfo)
		if err != nil {
			return nil, err
		}

		// If the account is not found in this store, continue to next one
		if _, err = ethStore.Get(ctx, addr); err != nil && errors.IsNotFoundError(err) {
			continue
		}
		if err != nil {
			return nil, err
		}

		logger.Debug("ethereum store found successfully", "store_name", storeName)
		return ethStore, nil
	}

	errMessage := "ethereum store was not found for the given address"
	logger.Error(errMessage)
	return nil, errors.NotFoundError(errMessage)
}

func (c *Connector) getEthStore(ctx context.Context, storeName string, resolver auth.Authorizator) (stores.KeyStore, error) {
	storeInfo, err := c.getStore(ctx, storeName, resolver)
	if err != nil {
		return nil, err
	}

	if storeInfo.StoreType != manifest.Ethereum {
		errMessage := "not an ethereum store"
		c.logger.Error(errMessage, "store_name", storeName)
		return nil, errors.NotFoundError(errMessage)
	}

	return storeInfo.Store.(stores.KeyStore), nil
}
