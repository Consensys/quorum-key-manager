package stores

import (
	"context"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/pkg/json"

	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func (c *Connector) CreateEthereum(_ context.Context, storeName string, storeType manifest.StoreType, specs interface{}, allowedTenants []string) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	logger := c.logger.With("store_name", storeName)
	logger.Debug("creating ethereum store")

	localETHSpecs := &entities.LocalEthSpecs{}
	if err := json.UnmarshalJSON(specs, localETHSpecs); err != nil {
		errMessage := "failed to unmarshal ethereum store specs"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidFormatError(errMessage)
	}

	store, err := c.createKeyStore(storeName, localETHSpecs.Keystore, localETHSpecs.Specs)
	if err != nil {
		return err
	}

	c.keys[storeName] = storeBundle{allowedTenants: allowedTenants, store: store, storeType: storeType}

	logger.Info("Ethereum store created successfully")
	return nil
}
