package stores

import (
	"context"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/pkg/json"

	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func (c *Connector) CreateEthereum(_ context.Context, storeName string, vaultType manifest.VaultType, specs interface{}, allowedTenants []string) error {
	logger := c.logger.With("vault_type", vaultType, "store_name", storeName)
	logger.Debug("creating ethereum store")

	store, err := c.createEthStore(storeName, vaultType, specs)
	if err != nil {
		return err
	}

	c.mux.Lock()
	c.stores[storeName] = &entities.StoreInfo{AllowedTenants: allowedTenants, Store: store, StoreType: manifest.Ethereum}
	c.mux.Unlock()

	logger.Info("Ethereum store created successfully")
	return nil
}

func (c *Connector) createEthStore(storeName string, vaultType manifest.VaultType, specs interface{}) (interface{}, error) {
	logger := c.logger.With("vault_type", vaultType, "store_name", storeName)

	switch vaultType {
	case manifest.LocalEthereum:
		localETHSpecs := &entities.LocalEthSpecs{}
		if err := json.UnmarshalJSON(specs, localETHSpecs); err != nil {
			errMessage := "failed to unmarshal ethereum store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidFormatError(errMessage)
		}

		return c.createKeyStore(storeName, localETHSpecs.Keystore, localETHSpecs.Specs)
	default:
		errMessage := "invalid store type for ethereum store"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
}
