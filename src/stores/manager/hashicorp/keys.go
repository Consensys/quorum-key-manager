package hashicorp

import (
	"context"
	client2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp/client"
	token2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/hashicorp/token"
	hashicorp2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/hashicorp"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

// KeySpecs is the specs format for an Hashicorp Vault key store
type KeySpecs struct {
	MountPoint string `json:"mountPoint"`
	Address    string `json:"address"`
	Token      string `json:"token"`
	TokenPath  string `json:"tokenPath"`
	Namespace  string `json:"namespace"`
}

func NewKeyStore(specs *KeySpecs, logger *log.Logger) (*hashicorp2.Store, error) {
	cfg := client2.NewConfig(specs.Address, specs.Namespace)
	cli, err := client2.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	if specs.Token != "" {
		cli.Client().SetToken(specs.Token)
	} else if specs.TokenPath != "" {
		tokenWatcher, err := token2.NewRenewTokenWatcher(cli, specs.TokenPath, logger)
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

	store := hashicorp2.New(cli, specs.MountPoint, logger)
	return store, nil
}
