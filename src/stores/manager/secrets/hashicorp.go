package secrets

import (
	"context"
	"time"

	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/token"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets/hashicorp"
)

const MaxRetries = 3

func NewHashicorpSecretStore(specs *entities.HashicorpSpecs, db database.Secrets, logger log.Logger) (*hashicorp.Store, error) {
	cli, err := client.NewClient(client.NewConfig(specs))
	if err != nil {
		errMessage := "failed to instantiate Hashicorp client (secrets)"
		logger.WithError(err).Error(errMessage, "specs", specs)
		return nil, errors.ConfigError(errMessage)
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
		retries := 0
		for retries < MaxRetries {
			err = cli.HealthCheck()
			if err == nil {
				break
			}

			logger.WithError(err).Debug("waiting for hashicorp client to be ready...", "retries", retries)
			time.Sleep(100 * time.Millisecond)
			retries++

			if retries == MaxRetries {
				errMessage := "failed to reach hashicorp vault (secrets). Please verify that the server is reachable"
				logger.WithError(err).Error(errMessage)
				return nil, errors.ConfigError(errMessage)
			}
		}
	}

	store := hashicorp.New(cli, db, specs.MountPoint, logger)
	return store, nil
}
