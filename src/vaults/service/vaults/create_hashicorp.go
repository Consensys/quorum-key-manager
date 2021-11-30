package vaults

import (
	"context"
	"github.com/consensys/quorum-key-manager/src/entities"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/token"
)

func (c *Vaults) CreateHashicorp(_ context.Context, name string, config *entities.HashicorpConfig) error {
	logger := c.logger.With("name", name)
	logger.Debug("creating hashicorp vault client")

	cli, err := client.NewClient(client.NewConfig(config))
	if err != nil {
		errMessage := "failed to instantiate Hashicorp client"
		logger.WithError(err).Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	if config.SkipVerify {
		logger.Warn("skipping certs verification will make your connection insecure and is not recommended in production")
	}

	if config.Token != "" {
		cli.SetToken(config.Token)
	} else if config.TokenPath != "" {
		tokenWatcher, err := token.NewRenewTokenWatcher(cli, config.TokenPath, logger)
		if err != nil {
			return err
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
				errMessage := "failed to reach hashicorp vault. Please verify that the server is reachable"
				logger.WithError(err).Error(errMessage)
				return errors.InvalidFormatError(errMessage)
			}
		}
	}

	c.createVault(name, entities.HashicorpVaultType, cli)

	logger.Info("hashicorp vault created successfully")
	return nil
}
