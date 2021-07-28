package models

import (
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/ethereum/go-ethereum/common"
)

type ETH1Account struct {
	Address             string `pg:",pk"`
	KeyID               string
	PublicKey           []byte
	CompressedPublicKey []byte
	Tags                map[string]string
	Disabled            bool
	CreatedAt           time.Time `pg:"default:now()"`
	UpdatedAt           time.Time `pg:"default:now()"`
	DeletedAt           time.Time `pg:",soft_delete"`
}

func (eth1 *ETH1Account) ToEntity() *entities.ETH1Account {
	return &entities.ETH1Account{
		Address:             common.HexToAddress(eth1.Address),
		KeyID:               eth1.KeyID,
		PublicKey:           eth1.PublicKey,
		CompressedPublicKey: eth1.CompressedPublicKey,
		Metadata: &entities.Metadata{
			Disabled:  eth1.Disabled,
			CreatedAt: eth1.CreatedAt,
			UpdatedAt: eth1.UpdatedAt,
			DeletedAt: eth1.DeletedAt,
		},
		Tags: eth1.Tags,
	}
}
