package hashicorp

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/hashicorp"
)

// SecretSpecs is the specs format for an Hashicorp Vault secret store
type SecretSpecs struct {
	MountPoint string `json:"mountPoint"`
	Address    string `json:"address"`
	Token      string `json:"token"`
	TokenPath  string `json:"tokenPath"`
	Namespace  string `json:"namespace"`
}

func NewSecretStore(specs *SecretSpecs, logger *log.Logger) (*hashicorp.Store, error) {
	cfg := client.NewConfig(specs.Address, specs.Namespace)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, errors.HashicorpVaultConnectionError(err.Error())
	}

	if specs.Token != "" {
		cli.Client().SetToken(specs.Token)
	} else if specs.TokenPath != "" {
		tokenWatcher, err := NewRenewTokenWatcher(cli.Client(), specs.TokenPath, logger)
		if err != nil {
			return nil, err
		}

		go func() {
			err = tokenWatcher.Run(context.Background())
			if err != nil {
				logger.WithError(err).Error("token watcher has exited with errors")
			}
			logger.Warn("token watcher has exited gracefully")
		}()
	}

	store := hashicorp.New(cli, specs.MountPoint)
	return store, nil
}
