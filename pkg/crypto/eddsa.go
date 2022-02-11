package crypto

import (
	"crypto/ed25519"

	babyjubjub "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func VerifyEDDSABabyJubJubSignature(publicKey, message, signature []byte) (bool, error) {
	pubKey := babyjubjub.PublicKey{}
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

func VerifyEDDSA25519Signature(publicKey, message, signature []byte) (bool, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return false, errors.InvalidParameterError("invalid ED25519 public key size")
	} 
	pubKey := ed25519.PublicKey(publicKey)
	return ed25519.Verify(pubKey, message, signature), nil
}
