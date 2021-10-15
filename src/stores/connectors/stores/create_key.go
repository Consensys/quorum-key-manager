package stores

import (
	"context"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/json"
	akvclient "github.com/consensys/quorum-key-manager/src/infra/akv/client"
	"github.com/consensys/quorum-key-manager/src/infra/aws/client"
	hashicorpclient "github.com/consensys/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/token"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/akv"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/aws"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func (c *Connector) CreateKey(_ context.Context, storeName string, vaultType manifest.VaultType, specs interface{}, allowedTenants []string) error {
	logger := c.logger.With("vault_type", vaultType, "store_name", storeName)
	logger.Debug("creating key store")

	store, err := c.createKeyStore(storeName, vaultType, specs)
	if err != nil {
		return err
	}

	c.mux.Lock()
	c.stores[storeName] = &entities.StoreInfo{AllowedTenants: allowedTenants, Store: store, StoreType: manifest.Keys}
	c.mux.Unlock()

	logger.Info("key store created successfully")
	return nil
}

func (c *Connector) createKeyStore(storeName string, vaultType manifest.VaultType, specs interface{}) (interface{}, error) {
	logger := c.logger.With("vault_type", vaultType, "store_name", storeName)

	switch vaultType {
	case manifest.HashicorpKeys:
		spec := &entities.HashicorpSpecs{}
		if err := json.UnmarshalJSON(specs, spec); err != nil {
			errMessage := "failed to unmarshal Hashicorp key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		return newHashicorpKeyStore(spec, logger)
	case manifest.AKVKeys:
		spec := &entities.AkvSpecs{}
		if err := json.UnmarshalJSON(specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		return newAkvKeyStore(spec, logger)
	case manifest.AWSKeys:
		spec := &entities.AwsSpecs{}
		if err := json.UnmarshalJSON(specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		return newAwsKeyStore(spec, logger)
	case manifest.LocalKeys:
		localKeySpecs := &entities.LocalKeySpecs{}
		if err := json.UnmarshalJSON(specs, localKeySpecs); err != nil {
			errMessage := "failed to unmarshal local key store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		return c.createSecretStore(storeName, localKeySpecs.SecretStore, localKeySpecs.Specs)
	default:
		errMessage := "invalid store type for key store"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
}

func newHashicorpKeyStore(specs *entities.HashicorpSpecs, logger log.Logger) (*hashicorp.Store, error) {
	cli, err := hashicorpclient.NewClient(hashicorpclient.NewConfig(specs))
	if err != nil {
		errMessage := "failed to instantiate Hashicorp client (keys)"
		logger.WithError(err).Error(errMessage)
		return nil, errors.ConfigError(errMessage)
	}

	if specs.SkipVerify {
		logger.Warn("skipping certs verification will make your connection insecure and is not recommended in production")
	}

	if specs.Token != "" {
		cli.SetToken(specs.Token)
	} else if specs.TokenPath != "" {
		tokenWatcher, err := token.NewRenewTokenWatcher(cli, specs.TokenPath, logger)
		if err != nil {
			return nil, err
		}

		go func() {
			err = tokenWatcher.Start(context.Background())
			if err != nil {
				logger.WithError(err).Error("token watcher has exited with errors")
			} else {
				logger.Warn("token watcher has exited gracefully")
			}
		}()

		// If the client token is read from filesystem, wait for it to be loaded before we continue
		maxRetries := 3
		retries := 0
		for retries < maxRetries {
			err = cli.HealthCheck()
			if err == nil {
				break
			}

			logger.WithError(err).Debug("waiting for hashicorp client to be ready...", "retries", retries)
			time.Sleep(100 * time.Millisecond)
			retries++

			if retries == maxRetries {
				errMessage := "failed to reach hashicorp vault (keys). Please verify that the server is reachable"
				logger.WithError(err).Error(errMessage)
				return nil, errors.ConfigError(errMessage)
			}
		}

	}

	store := hashicorp.New(cli, specs.MountPoint, logger)
	return store, nil
}

func newAkvKeyStore(spec *entities.AkvSpecs, logger log.Logger) (*akv.Store, error) {
	cli, err := akvclient.NewClient(akvclient.NewConfig(spec.VaultName, spec.TenantID, spec.ClientID, spec.ClientSecret))
	if err != nil {
		errMessage := "failed to instantiate AKV client (keys)"
		logger.WithError(err).Error(errMessage, "specs", spec)
		return nil, errors.ConfigError(errMessage)
	}

	store := akv.New(cli, logger)
	return store, nil
}

func newAwsKeyStore(specs *entities.AwsSpecs, logger log.Logger) (*aws.Store, error) {
	cfg := client.NewConfig(specs.Region, specs.AccessID, specs.SecretKey, specs.Debug)
	cli, err := client.NewKmsClient(cfg)
	if err != nil {
		errMessage := "failed to instantiate AWS client (keys)"
		logger.WithError(err).Error(errMessage, "specs", specs)
		return nil, errors.ConfigError(errMessage)
	}

	store := aws.New(cli, logger)
	return store, nil
}
