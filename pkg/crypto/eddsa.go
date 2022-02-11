package crypto

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"

	babyjubjub "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func EdDSA25519(importedPrivKey []byte) (privKey, pubKey []byte, err error) {
	// https://pkg.go.dev/crypto/ed25519#section-documentation
	if importedPrivKey != nil {
		if len(importedPrivKey) != ed25519.PrivateKeySize {
			return nil, nil, errors.InvalidParameterError("invalid private key value")
		}
		ed25519PrivKey := ed25519.PrivateKey(importedPrivKey)
		pubKey = ed25519PrivKey.Public().(ed25519.PublicKey)
		privKey = ed25519PrivKey
	} else {
		seed := make([]byte, 32)
		if _, err = rand.Read(seed); err != nil {
			return nil, nil, err
		}

		pubKey, privKey, err = ed25519.GenerateKey(bytes.NewReader(seed))
		if err != nil {
			return nil, nil, err
		}
	}

	return privKey, pubKey, nil
}

func EdDSABabyjubjub(importedPrivKey []byte) (privKey, pubKey []byte, err error) {
	babyJubJubPrivKey := babyjubjub.PrivateKey{}
	if importedPrivKey != nil {
		_, err = babyJubJubPrivKey.SetBytes(importedPrivKey)
		if err != nil {
			return nil, nil, err
		}
	} else {
		seed := make([]byte, 32)
		_, err = rand.Read(seed)
		if err != nil {
			return nil, nil, err
		}

		// Usually standards implementations of eddsa do not require the choice of a specific hash function (usually it's SHA256).
		// Here we needed to allow the choice of the hash, so we can choose a hash function that is easily programmable in a snark circuit.
		// Same hFunc should be used for sign and verify
		babyJubJubPrivKey, err = babyjubjub.GenerateKey(bytes.NewReader(seed))
		if err != nil {
			return nil, nil, err
		}
	}

	privKey = babyJubJubPrivKey.Bytes()
	pubKey = babyJubJubPrivKey.Public().Bytes()
	return privKey, pubKey, nil
}

func SignEDDSABabyjubjub(privKeyB, data []byte) ([]byte, error) {
	privKey := babyjubjub.PrivateKey{}
	_, err := privKey.SetBytes(privKeyB)
	if err != nil {
		return nil, errors.InvalidParameterError(fmt.Sprintf("failed to parse private key. %s", err.Error()))
	}

	signature, err := privKey.Sign(data, hash.MIMC_BN254.New("seed"))
	if err != nil {
		return nil, errors.CryptoOperationError(fmt.Sprintf("failed to sign. %s", err.Error()))
	}

	return signature, nil
}

func SignEDDSA25519(privKeyB, data []byte) ([]byte, error) {
	if len(privKeyB) != ed25519.PrivateKeySize {
		return nil, errors.InvalidParameterError("invalid ED25519 private key size")
	}
	privKey := ed25519.PrivateKey(privKeyB)
	signature := ed25519.Sign(privKey, data)
	return signature, nil
}

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
