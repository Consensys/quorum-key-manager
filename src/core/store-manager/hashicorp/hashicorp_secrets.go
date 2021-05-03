package hashicorp

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/hashicorp"
)

// SecretSpecs is the specs format for an Hashicorp Vault secret store
type SecretSpecs struct {
	MountPoint string `json:"mountPoint"`
	Address    string `json:"address"`
	Token      string `json:"token"`
	Namespace  string `json:"namespace"`
}

func NewSecretStore(specs *SecretSpecs) (*hashicorp.SecretStore, error) {
	cfg := client.NewBaseConfig(specs.Address, specs.Token, specs.Namespace)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, errors.HashicorpVaultConnectionError(err.Error())
	}

	store := hashicorp.New(cli, specs.MountPoint)
	return store, nil
}
