// nolint
package hashicorp

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/hashicorp"
)

// HashicorpSecretSpecs is the specs format for an Hashicorp Vault secret store
type SecretSpecs struct {
	MountPoint string `json:"mount_point"`
	Address    string `json:"address"`
	Token      string `json:"token"`
}

func NewSecretStore(specs *SecretSpecs) (*hashicorp.SecretStore, error) {
	cfg := client.NewBaseConfig(specs.Address, specs.Token)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	store := hashicorp.New(cli, specs.MountPoint)
	return store, nil
}
