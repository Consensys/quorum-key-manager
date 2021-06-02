package crypto

import (
	"github.com/consensys/gnark-crypto/crypto/hash"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

func VerifyEDDSASignature(publicKey, message, signature []byte) (bool, error) {
	pubKey := eddsa.PublicKey{}
	_, err := pubKey.SetBytes(publicKey)
	if err != nil {
		return false, err
	}

	verified, err := pubKey.Verify(signature, message, hash.MIMC_BN254.New("seed"))
	if err != nil {
		return false, err
	}

	return verified, nil
}
