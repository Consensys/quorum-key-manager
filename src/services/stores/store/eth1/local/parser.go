package local

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

func parseKey(key *entities.Key) *entities.ETH1Account {
	pubKey, _ := crypto.UnmarshalPubkey(key.PublicKey)
	return &entities.ETH1Account{
		ID:                  key.ID,
		Address:             crypto.PubkeyToAddress(*pubKey).Hex(),
		Metadata:            key.Metadata,
		Tags:                key.Tags,
		PublicKey:           crypto.FromECDSAPub(pubKey),
		CompressedPublicKey: crypto.CompressPubkey(pubKey),
	}
}
