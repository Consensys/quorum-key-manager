package hashicorp

import (
	"context"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp/token"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/hashicorp"
)

const maxRetries = 5

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

		for currRetries := 0; currRetries < maxRetries; currRetries++ {
			if tokenWatcher.IsTokenLoaded() {
				break
			}

			time.Sleep(time.Second)
		}
	}

	store := hashicorp.New(cli, specs.MountPoint, logger)
	return store, nil
}
