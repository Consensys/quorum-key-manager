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

func (c *Connector) CreateSecret(ctx context.Context, name, vaultName string, allowedTenants []string, userInfo *auth.UserInfo) error {
	logger := c.logger.With("name", name, "vault", vaultName)
	logger.Debug("creating secret store")

	vault, err := c.vaults.Get(ctx, vaultName, userInfo)
	if err != nil {
		return err
	}

	var store stores.SecretStore
	switch vault.VaultType {
	case entities2.HashicorpVaultType:
		store, err = hashicorp.New(vault.Client.(hashicorpinfra.Kvv2Client), c.db.Secrets(name), logger), nil
	case entities2.AzureVaultType:
		store, err = akv.New(vault.Client.(akvinfra.SecretClient), logger), nil
	case entities2.AWSVaultType:
		store, err = aws.New(vault.Client.(awsinfra.SecretsManagerClient), logger), nil
	default:
		errMessage := "invalid vault for secret store"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}
	if err != nil {
		return err
	}

	c.createStore(name, entities.SecretStoreType, store, allowedTenants)

	logger.Info("secret store created successfully")
	return nil
}
