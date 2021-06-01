package ethereum

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

const EthSignatureLength = 65

func VerifySignature(signature, msg, privKeyB []byte) (bool, error) {
	privKey, err := crypto.ToECDSA(privKeyB)
	if err != nil {
		return false, err
	}

	if len(signature) == EthSignatureLength {
		retrievedPubkey, err := crypto.SigToPub(msg, signature)
		if err != nil {
			return false, err
		}

		return privKey.PublicKey.Equal(retrievedPubkey), nil
	}

	r := new(big.Int).SetBytes(signature[0:32])
	s := new(big.Int).SetBytes(signature[32:64])
	return ecdsa.Verify(&privKey.PublicKey, msg, r, s), nil
}
