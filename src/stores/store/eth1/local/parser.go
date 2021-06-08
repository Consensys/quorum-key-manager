package local

import (
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

func parseKey(key *entities2.Key) *entities2.ETH1Account {
	pubKey, _ := crypto.UnmarshalPubkey(key.PublicKey)
	return &entities2.ETH1Account{
		ID:                  key.ID,
		Address:             crypto.PubkeyToAddress(*pubKey).Hex(),
		Metadata:            key.Metadata,
		Tags:                key.Tags,
		PublicKey:           crypto.FromECDSAPub(pubKey),
		CompressedPublicKey: crypto.CompressPubkey(pubKey),
	}
}
