package models

import (
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/ethereum/go-ethereum/common"
	"time"
)

type ETH1Account struct {
	tableName struct{} `pg:"eth1-accounts"` // nolint:unused,structcheck // reason

	Address   string `pg:",pk"`
	KeyID     string `pg:"alias:key_id,notnull"`
	Key       *Key
	Tags      map[string]string
	Disabled  bool
	CreatedAt time.Time `pg:"default:now()"`
	UpdatedAt time.Time `pg:"default:now()"`
	DeletedAt time.Time `pg:",soft_delete"`
}

func NewETH1Account(eth1 *entities.ETH1Account) *ETH1Account {
	return &ETH1Account{
		KeyID:     eth1.KeyID,
		Address:   eth1.Address.Hex(),
		Tags:      eth1.Tags,
		Key:       NewKey(eth1.Key),
		Disabled:  eth1.Metadata.Disabled,
		CreatedAt: eth1.Metadata.CreatedAt,
		UpdatedAt: eth1.Metadata.UpdatedAt,
		DeletedAt: eth1.Metadata.DeletedAt,
	}
}

func (eth1 *ETH1Account) ToEntity() *entities.ETH1Account {
	return &entities.ETH1Account{
		Address: common.HexToAddress(eth1.Address),
		KeyID:   eth1.KeyID,
		Key:     eth1.Key.ToEntity(),
		Metadata: &entities.Metadata{
			Disabled:  eth1.Disabled,
			CreatedAt: eth1.CreatedAt,
			UpdatedAt: eth1.UpdatedAt,
			DeletedAt: eth1.DeletedAt,
		},
		Tags: eth1.Tags,
	}
}
