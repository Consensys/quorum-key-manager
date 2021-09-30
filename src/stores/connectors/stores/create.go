package stores

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
	eth "github.com/consensys/quorum-key-manager/src/stores/manager/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/manager/keys"
	"github.com/consensys/quorum-key-manager/src/stores/manager/secrets"
)

func (c *Connector) Create(_ context.Context, mnf *manifest.Manifest) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	logger := c.logger.With("kind", mnf.Kind).With("name", mnf.Name)
	logger.Debug("loading store manifest")

	switch mnf.Kind {
	case manifest.HashicorpSecrets:
		spec := &entities.HashicorpSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := secrets.NewHashicorpSecretStore(spec, c.db.Secrets(mnf.Name), logger)
		if err != nil {
			return err
		}

		c.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case manifest.HashicorpKeys:
		spec := &entities.HashicorpSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := keys.NewHashicorpKeyStore(spec, logger)
		if err != nil {
			return err
		}

		c.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case manifest.AKVSecrets:
		spec := &secrets.AkvSecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AKV secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := secrets.NewAkvSecretStore(spec, logger)
		if err != nil {
			return err
		}

		c.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case manifest.AKVKeys:
		spec := &keys.AkvKeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AKV key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := keys.NewAkvKeyStore(spec, logger)
		if err != nil {
			return err
		}

		c.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case manifest.AWSSecrets:
		spec := &secrets.AwsSecretSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AWS secret store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := secrets.NewAwsSecretStore(spec, logger)
		if err != nil {
			return err
		}

		c.secrets[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case manifest.AWSKeys:
		spec := &keys.AwsKeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal AWS key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := keys.NewAwsKeyStore(spec, logger)
		if err != nil {
			return err
		}

		c.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case manifest.LocalKeys:
		spec := &keys.LocalKeySpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal local key store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := keys.NewLocalKeyStore(spec, c.db.Secrets(mnf.Name), logger)
		if err != nil {
			return err
		}

		c.keys[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	case manifest.Ethereum:
		spec := &eth.LocalEthSpecs{}
		if err := mnf.UnmarshalSpecs(spec); err != nil {
			errMessage := "failed to unmarshal Eth store specs"
			logger.WithError(err).Error(errMessage)
			return errors.InvalidFormatError(errMessage)
		}

		store, err := eth.NewLocalEth(spec, c.db.Secrets(mnf.Name), logger)
		if err != nil {
			return err
		}

		c.ethAccounts[mnf.Name] = &storeBundle{manifest: mnf, store: store, logger: logger}
	default:
		errMessage := "invalid manifest kind"
		logger.Error(errMessage, "kind", mnf.Kind)
		return errors.InvalidFormatError(errMessage)
	}

	logger.Info("store manifest loaded successfully")
	return nil
}
