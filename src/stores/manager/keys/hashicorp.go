package keys

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp/token"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/consensys/quorum-key-manager/pkg/errors"

	"github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"
)

func NewHashicorpKeyStore(specs *entities.HashicorpSpecs, logger log.Logger) (*hashicorp.Store, error) {
	cfg := client.NewConfig(specs)
	cli, err := client.NewClient(cfg)
	if err != nil {
		errMessage := "failed to instantiate Hashicorp client (keys)"
		logger.WithError(err).Error(errMessage, "specs", specs)
		return nil, errors.ConfigError(errMessage)
	}

	if cfg.SkipVerify {
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
	}

	store := hashicorp.New(cli, specs.MountPoint, logger)
	return store, nil
}
