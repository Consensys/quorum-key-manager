package stores

import (
	"context"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/auth/service/authorizator"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
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

func (c *Connector) CreateKey(ctx context.Context, name, vaultName, secretStore string, allowedTenants []string, userInfo *authtypes.UserInfo) error {
	logger := c.logger.With("name", name, "vault", vaultName, "secret_store", secretStore)
	logger.Debug("creating key store")

	if vaultName != "" && secretStore != "" {
		errMessage := "cannot specify vault and secret store simultaneously. Please choose one option"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	// TODO: Uncomment when authManager no longer a runnable
	// permissions := c.authManager.UserPermissions(userInfo)
	resolver := authorizator.New(userInfo.Permissions, userInfo.Tenant, c.logger)

	// If vault is specified, it is a remote key store, otherwise it's a local key store
	var store stores.KeyStore
	switch {
	case vaultName != "":
		vault, err := c.vaults.Get(ctx, vaultName, userInfo)
		if err != nil {
			return err
		}

		switch vault.VaultType {
		case entities2.HashicorpVaultType:
			store, err = hashicorp.New(vault.Client.(hashicorpinfra.PluginClient), logger), nil
		case entities2.AzureVaultType:
			store, err = akv.New(vault.Client.(akvinfra.KeysClient), logger), nil
		case entities2.AWSVaultType:
			store, err = aws.New(vault.Client.(awsinfra.KmsClient), logger), nil
		default:
			errMessage := "invalid vault for key store"
			logger.Error(errMessage)
			return errors.InvalidParameterError(errMessage)
		}
		if err != nil {
			return err
		}
	case secretStore != "":
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
