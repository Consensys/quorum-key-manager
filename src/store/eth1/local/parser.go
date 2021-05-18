package local

import (
	"crypto/ecdsa"
	"encoding/base64"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func parseKey(key *entities.Key) (*entities.ETH1Account, error) {
	pubKey, err := parsePubKey(key.PublicKey)
	if err != nil {
		return nil, err
	}

	return &entities.ETH1Account{
		ID:                  key.ID,
		Address:             crypto.PubkeyToAddress(*pubKey).Hex(),
		Metadata:            key.Metadata,
		Tags:                key.Tags,
		PublicKey:           key.PublicKey,
		CompressedPublicKey: hexutil.Encode(crypto.CompressPubkey(pubKey)),
	}, nil
}

func parseRecID(pubKeyS string) (string, error) {
	pubKey, err := parsePubKey(pubKeyS)
	if err != nil {
		return "", err
	}

	if pubKey.Y.Mod(pubKey.Y, big.NewInt(2)) == big.NewInt(0) {
		return "00", nil
	}

	return "01", nil
}

func parsePubKey(pubKeyS string) (*ecdsa.PublicKey, error) {
	pubKeyB, err := base64.URLEncoding.DecodeString(pubKeyS)
	if err != nil {
		return nil, errors.EncodingError("failed to decode public key")
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyB)
	if err != nil {
		return nil, errors.EncodingError("failed to unmarshal public key")
	}

	return pubKey, nil
}
