package stores

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/json"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	eth "github.com/consensys/quorum-key-manager/src/stores/manager/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/manager/keys"
	"github.com/consensys/quorum-key-manager/src/stores/manager/secrets"
)

func (c *Connector) Create(_ context.Context, mnf *manifest.Manifest) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	logger := c.logger.With("name", mnf.Name)

	switch mnf.Kind {
	case manifest.HashicorpSecrets:
		spec := &entities.HashicorpSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := secrets.NewHashicorpSecretStore(spec, c.db.Secrets(mnf.Name), logger)
		if err != nil {
			return err
		}

		c.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}

		logger.Info("Hashicorp secret store created successfully")
	case manifest.HashicorpKeys:
		spec := &entities.HashicorpSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := keys.NewHashicorpKeyStore(spec, logger)
		if err != nil {
			return err
		}

		c.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}

		logger.Info("Hashicorp key store created successfully")
	case manifest.AKVSecrets:
		spec := &entities.AkvSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := secrets.NewAkvSecretStore(spec, logger)
		if err != nil {
			return err
		}

		c.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}

		logger.Info("Azure secret store created successfully")
	case manifest.AKVKeys:
		spec := &entities.AkvSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := keys.NewAkvKeyStore(spec, logger)
		if err != nil {
			return err
		}

		c.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}

		logger.Info("Azure key store created successfully")
	case manifest.AWSSecrets:
		spec := &entities.AwsSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := secrets.NewAwsSecretStore(spec, logger)
		if err != nil {
			return err
		}

		c.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}

		logger.Info("AWS secret store created successfully")
	case manifest.AWSKeys:
		spec := &entities.AwsSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := keys.NewAwsKeyStore(spec, logger)
		if err != nil {
			return err
		}

		c.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}

		logger.Info("AWS key store created successfully")
	case manifest.LocalKeys:
		spec := &entities.LocalKeySpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "failed to unmarshal local key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := keys.NewLocalKeyStore(spec, c.db.Secrets(mnf.Name), logger)
		if err != nil {
			return err
		}

		c.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}

		logger.Info("Local key store created successfully")
	case manifest.Ethereum:
		spec := &entities.LocalEthSpecs{}
		if err := json.UnmarshalJSON(mnf.Specs, spec); err != nil {
			errMessage := "failed to unmarshal Eth store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := eth.NewLocalEth(spec, c.db.Secrets(mnf.Name), logger)
		if err != nil {
			return err
		}

		c.ethAccounts[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
		logger.Info("Ethereum store created successfully")
	}

	return nil
}
