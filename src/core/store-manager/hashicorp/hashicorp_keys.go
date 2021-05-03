package hashicorp

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/hashicorp"
)

// KeySpecs is the specs format for an Hashicorp Vault key store
type KeySpecs struct {
	MountPoint string `json:"mountPoint"`
	Address    string `json:"address"`
	Token      string `json:"token"`
	Namespace  string `json:"namespace"`
}

func NewKeyStore(specs *KeySpecs) (*hashicorp.Store, error) {
	cfg := client.NewBaseConfig(specs.Address, specs.Token, specs.Namespace)
	cli, err := client.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	store := hashicorp.New(cli, specs.MountPoint)
	return store, nil
}
