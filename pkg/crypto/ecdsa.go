package crypto

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

func VerifyECDSASignature(publicKey, message, signature []byte) (bool, error) {
	pubKey, err := crypto.UnmarshalPubkey(publicKey)
	if err != nil {
		return false, err
	}

	r := new(big.Int).SetBytes(signature[0:32])
	s := new(big.Int).SetBytes(signature[32:64])

	return ecdsa.Verify(pubKey, message, r, s), nil
}
