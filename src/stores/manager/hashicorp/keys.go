package hashicorp

import (
	"context"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp/token"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/hashicorp"
)

const maxRetries = 3

// KeySpecs is the specs format for an Hashicorp Vault key store
type KeySpecs struct {
	MountPoint string `json:"mountPoint"`
	Address    string `json:"address"`
	Token      string `json:"token"`
	TokenPath  string `json:"tokenPath"`
	Namespace  string `json:"namespace"`
}

func NewKeyStore(specs *KeySpecs, logger *log.Logger) (*hashicorp.Store, error) {
	cfg := client.NewConfig(specs.Address, specs.Namespace)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if specs.Token != "" {
		cli.Client().SetToken(specs.Token)
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

		// We wait for the token to be set before we continue
		for currRetries := 1; currRetries <= maxRetries; currRetries++ {
			if tokenWatcher.IsTokenLoaded() {
				break
			}

			if currRetries == maxRetries {
				errMessage := "failed to load token from file"
				logger.WithError(err).Error(errMessage)
				return nil, errors.ConfigError(errMessage)
			}

			time.Sleep(time.Second)
		}
	}

	store := hashicorp.New(cli, specs.MountPoint, logger)
	return store, nil
}
