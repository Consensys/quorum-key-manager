package local

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

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

func parseRecID(pubKeyB []byte) (*byte, error) {
	pubKey, err := crypto.UnmarshalPubkey(pubKeyB)
	if err != nil {
		return nil, errors.EncodingError("failed to unmarshal public key")
	}

	if pubKey.Y.Mod(pubKey.Y, big.NewInt(2)) == big.NewInt(0) {
		b := byte(0)
		return &b, nil
	}

	b := byte(1)
	return &b, nil
}

func parseSignatureValues(tx *types.Transaction, sig []byte, signer types.Signer) ([]byte, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}

	return append(append(r.Bytes(), s.Bytes()...), v.Bytes()...), nil
}
