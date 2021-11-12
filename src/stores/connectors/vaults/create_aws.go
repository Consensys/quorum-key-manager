package vaults

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/aws/client"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

func (c *Connector) CreateAWS(_ context.Context, name string, config *entities.AWSConfig) error {
	logger := c.logger.With("name", name)
	logger.Debug("creating aws vault client")

	cli, err := client.NewSecretsClient(client.NewConfig(config))
	if err != nil {
		errMessage := "failed to instantiate AWS client"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	c.createVault(name, entities.AWSVaultType, cli)

	logger.Info("aws vault created successfully")
	return nil
}
