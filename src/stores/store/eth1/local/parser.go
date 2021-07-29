package local

import (
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/models"
	"github.com/ethereum/go-ethereum/crypto"
)

func parseKey(key *models.Key, attr *entities.Attributes) *models.ETH1Account {
	pubKey, _ := crypto.UnmarshalPubkey(key.PublicKey)
	return &models.ETH1Account{
		KeyID:               key.ID,
		Address:             crypto.PubkeyToAddress(*pubKey).Hex(),
		Tags:                attr.Tags,
		PublicKey:           key.PublicKey,
		CompressedPublicKey: crypto.CompressPubkey(pubKey),
	}
}
