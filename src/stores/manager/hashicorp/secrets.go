package hashicorp

import (
	"context"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/pkg/log-old"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/hashicorp/client"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/hashicorp/token"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/secrets/hashicorp"
)

// SecretSpecs is the specs format for an Hashicorp Vault secret store
type SecretSpecs struct {
	MountPoint string `json:"mountPoint"`
	Address    string `json:"address"`
	Token      string `json:"token"`
	TokenPath  string `json:"tokenPath"`
	Namespace  string `json:"namespace"`
}

func NewSecretStore(specs *SecretSpecs, logger *log_old.Logger) (*hashicorp.Store, error) {
	cfg := client.NewConfig(specs.Address, specs.Namespace)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, errors.HashicorpVaultError(err.Error())
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
