package local

import (
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

func parseKey(key *entities.Key, attr *entities.Attributes) *entities.ETH1Account {
	pubKey, _ := crypto.UnmarshalPubkey(key.PublicKey)
	return &entities.ETH1Account{
		KeyID:               key.ID,
		Address:             crypto.PubkeyToAddress(*pubKey),
		Tags:                attr.Tags,
		PublicKey:           key.PublicKey,
		CompressedPublicKey: crypto.CompressPubkey(pubKey),
	}
}
