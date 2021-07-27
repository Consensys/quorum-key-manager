package local

import (
	"github.com/consensys/quorum-key-manager/src/stores/store/database/models"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

func parseKey(key *entities.Key, attr *entities.Attributes) *models.ETH1Account {
	pubKey, _ := crypto.UnmarshalPubkey(key.PublicKey)
	return &models.ETH1Account{
		Address: crypto.PubkeyToAddress(*pubKey).Hex(),
		KeyID:   key.ID,
		Key:     models.NewKey(key),
		Tags:    attr.Tags,
	}
}
