package stores

import (
	"context"
	"time"

	akvclient "github.com/consensys/quorum-key-manager/src/infra/akv/client"
	awsclient "github.com/consensys/quorum-key-manager/src/infra/aws/client"
	hashicorpclient "github.com/consensys/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/token"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/akv"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/aws"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"

	"github.com/consensys/quorum-key-manager/pkg/json"

	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"

	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func (c *Connector) CreateSecret(_ context.Context, storeName string, storeType manifest.StoreType, specs interface{}, allowedTenants []string) error {
	logger := c.logger.With("store_type", storeType, "store_name", storeName)
	logger.Debug("creating secret store")

	store, err := c.createSecretStore(storeName, storeType, specs)
	if err != nil {
		return nil
	}

	c.secrets[storeName] = storeBundle{allowedTenants: allowedTenants, store: store, storeType: storeType}

	logger.Info("secret store created successfully")
	return nil
}

func (c *Connector) createSecretStore(storeName string, storeType manifest.StoreType, specs interface{}) (stores.SecretStore, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	logger := c.logger.With("store_type", storeType, "store_name", storeName)

	switch storeType {
	case manifest.HashicorpSecrets:
		hashicorpSpecs := &entities.HashicorpSpecs{}
		if err := json.UnmarshalJSON(specs, hashicorpSpecs); err != nil {
			errMessage := "failed to unmarshal Hashicorp secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		return newHashicorpSecretStore(hashicorpSpecs, c.db.Secrets(storeName), logger)
	case manifest.AKVSecrets:
		spec := &entities.AkvSpecs{}
		if err := json.UnmarshalJSON(specs, spec); err != nil {
			errMessage := "failed to unmarshal AKV secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		return newAkvSecretStore(spec, logger)
	case manifest.AWSSecrets:
		spec := &entities.AwsSpecs{}
		if err := json.UnmarshalJSON(specs, spec); err != nil {
			errMessage := "failed to unmarshal AWS secret store specs"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		return newAwsSecretStore(spec, logger)
	default:
		errMessage := "invalid store type for secret store"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
}

func newHashicorpSecretStore(specs *entities.HashicorpSpecs, db database.Secrets, logger log.Logger) (*hashicorp.Store, error) {
	cli, err := hashicorpclient.NewClient(hashicorpclient.NewConfig(specs))
	if err != nil {
		errMessage := "failed to instantiate Hashicorp client (secrets)"
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
				errMessage := "failed to reach hashicorp vault (secrets). Please verify that the server is reachable"
				logger.WithError(err).Error(errMessage)
				return nil, errors.ConfigError(errMessage)
			}
		}
	}

	return hashicorp.New(cli, db, specs.MountPoint, logger), nil
}

func newAkvSecretStore(spec *entities.AkvSpecs, logger log.Logger) (*akv.Store, error) {
	cli, err := akvclient.NewClient(akvclient.NewConfig(spec.VaultName, spec.TenantID, spec.ClientID, spec.ClientSecret))
	if err != nil {
		errMessage := "failed to instantiate AKV client (secrets)"
		logger.WithError(err).Error(errMessage, "specs", spec)
		return nil, errors.ConfigError(errMessage)
	}

	return akv.New(cli, logger), nil
}

func newAwsSecretStore(specs *entities.AwsSpecs, logger log.Logger) (*aws.Store, error) {
	cli, err := awsclient.NewSecretsClient(awsclient.NewConfig(specs.Region, specs.AccessID, specs.SecretKey, specs.Debug))
	if err != nil {
		errMessage := "failed to instantiate AWS client (secrets)"
		logger.WithError(err).Error(errMessage, "specs", specs)
		return nil, errors.ConfigError(errMessage)
	}

	store := aws.New(cli, logger)
	return store, nil
}
