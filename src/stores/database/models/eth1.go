package models

import (
	"time"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/ethereum/go-ethereum/common"
)

type ETH1Account struct {
	tableName struct{} `pg:"eth_accounts"` // nolint:unused,structcheck // reason

	Address             string `pg:",pk"`
	StoreID             string `pg:",pk"`
	KeyID               string
	PublicKey           []byte
	CompressedPublicKey []byte
	Tags                map[string]string
	Disabled            bool
	CreatedAt           time.Time `pg:"default:now()"`
	UpdatedAt           time.Time `pg:"default:now()"`
	DeletedAt           time.Time `pg:",soft_delete"`
}

func NewETH1Account(account *entities.ETH1Account) *ETH1Account {
	return &ETH1Account{
		Address:             account.Address.Hex(),
		KeyID:               account.KeyID,
		PublicKey:           account.PublicKey,
		CompressedPublicKey: account.CompressedPublicKey,
		Tags:                account.Tags,
		Disabled:            account.Metadata.Disabled,
		CreatedAt:           account.Metadata.CreatedAt,
		UpdatedAt:           account.Metadata.UpdatedAt,
		DeletedAt:           account.Metadata.DeletedAt,
	}
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
