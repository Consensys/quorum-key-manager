package local

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ethereum/go-ethereum/crypto"
)

func parseKey(key *entities.Key) (*entities.ETH1Account, error) {
	pubKey, err := crypto.UnmarshalPubkey(key.PublicKey)
	if err != nil {
		return nil, errors.EncodingError("failed to unmarshal public key")
	}

	return &entities.ETH1Account{
		ID:                  key.ID,
		Address:             crypto.PubkeyToAddress(*pubKey).Hex(),
		Metadata:            key.Metadata,
		Tags:                key.Tags,
		PublicKey:           crypto.FromECDSAPub(pubKey),
		CompressedPublicKey: crypto.CompressPubkey(pubKey),
	}, nil
}
