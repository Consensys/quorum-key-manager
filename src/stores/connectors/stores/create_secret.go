package stores

import (
	"context"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	akvinfra "github.com/consensys/quorum-key-manager/src/infra/akv"
	awsinfra "github.com/consensys/quorum-key-manager/src/infra/aws"
	hashicorpinfra "github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/akv"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/aws"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"

	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func (c *Connector) CreateSecret(ctx context.Context, name, vault string, allowedTenants []string, _ *auth.UserInfo) error {
	logger := c.logger.With("name", name, "vault", vault)
	logger.Debug("creating secret store")

	store, err := c.createSecretStore(ctx, vault, name)
	if err != nil {
		return err
	}

	c.createStore(name, entities.SecretStoreType, store, allowedTenants)

	logger.Info("secret store created successfully")
	return nil
}

func (c *Connector) createSecretStore(ctx context.Context, vaultName, storeName string) (stores.SecretStore, error) {
	logger := c.logger.With("vault", vaultName)

	vault, err := c.vaults.Get(ctx, vaultName)
	if err != nil {
		return nil, err
	}

	switch vault.VaultType {
	case entities2.HashicorpVaultType:
		return hashicorp.New(vault.Client.(hashicorpinfra.VaultClient), c.db.Secrets(storeName), logger), nil
	case entities2.AzureVaultType:
		return akv.New(vault.Client.(akvinfra.SecretClient), logger), nil
	case entities2.AWSVaultType:
		return aws.New(vault.Client.(awsinfra.SecretsManagerClient), logger), nil
	default:
		errMessage := "invalid vault for secret store"
		logger.Error(errMessage)
		return nil, errors.InvalidParameterError(errMessage)
	}
}
