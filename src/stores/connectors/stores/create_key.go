package stores

import (
	"context"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	akvinfra "github.com/consensys/quorum-key-manager/src/infra/akv"
	awsinfra "github.com/consensys/quorum-key-manager/src/infra/aws"
	hashicorpinfra "github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores"

	"github.com/consensys/quorum-key-manager/src/stores/store/keys/akv"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/aws"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
	localkeys "github.com/consensys/quorum-key-manager/src/stores/store/keys/local"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func (c *Connector) CreateKey(ctx context.Context, name, vault, secretStore string, allowedTenants []string, userInfo *auth.UserInfo) error {
	logger := c.logger.With("name", name, "vault", vault, "secret_store", secretStore)
	logger.Debug("creating key store")

	if name != "" && secretStore != "" {
		errMessage := "cannot specify vault and secret store simultaneously. Please choose one option"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	// If vault is specified, it is a remote key store, otherwise it's a local key store
	var store stores.KeyStore
	var err error
	switch {
	case vault != "":
		store, err = c.createKeyStore(ctx, vault)
		if err != nil {
			return err
		}
	case secretStore != "":
		// TODO: Uncomment when authManager no longer a runnable
		// permissions := c.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(userInfo.Permissions, userInfo.Tenant, c.logger)

		secretstore, err := c.getSecretStore(ctx, secretStore, resolver)
		if err != nil {
			return err
		}

		store = localkeys.New(secretstore, c.db.Secrets(name), c.logger)
	default:
		errMessage := "either vault or secret store must be specified. Please choose one option"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	c.createStore(name, entities.KeyStoreType, store, allowedTenants)

	logger.Info("key store created successfully")
	return nil
}

func (c *Connector) createKeyStore(ctx context.Context, vaultName string) (stores.KeyStore, error) {
	logger := c.logger.With("vault", vaultName)

	vault, err := c.vaults.Get(ctx, vaultName)
	if err != nil {
		return nil, err
	}

	switch vault.VaultType {
	case entities.HashicorpVaultType:
		return hashicorp.New(vault.Client.(hashicorpinfra.VaultClient), specs.MountPoint, logger), nil
	case entities.AzureVaultType:
		return akv.New(vault.Client.(akvinfra.KeysClient), logger), nil
	case entities.AWSVaultType:
		return aws.New(vault.Client.(awsinfra.KmsClient), logger), nil
	default:
		errMessage := "invalid vault for key store"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
}
