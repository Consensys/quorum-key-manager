package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"

	babyjubjub "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

func EdDSA25519(importedPrivKey []byte) (privKey []byte, pubKey []byte, err error) {
	// https://pkg.go.dev/crypto/ed25519#section-documentation
	if importedPrivKey != nil {
		if len(importedPrivKey) != ed25519.PrivateKeySize {
			return nil, nil, errors.InvalidParameterError("invalid private key value")
		}
		ed25519PrivKey := ed25519.PrivateKey(importedPrivKey)
		pubKey = ed25519PrivKey.Public().(ed25519.PublicKey)
		return ed25519PrivKey, pubKey, nil
	}

	seed := make([]byte, 32)
	if _, err = rand.Read(seed); err != nil {
		return nil, nil, err
	}

	return ed25519.GenerateKey(bytes.NewReader(seed))
}

func ECDSASecp256k1(importedPrivKey []byte) (privKey []byte, pubKey []byte, err error) {
	ecdsaKey := &ecdsa.PrivateKey{}
	if importedPrivKey != nil {
		ecdsaKey, err = crypto.ToECDSA(importedPrivKey)
		if err != nil {
			return nil, nil, err
		}
	}

	ecdsaKey, err = crypto.GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	privKey = crypto.FromECDSA(ecdsaKey)
	pubKey = crypto.FromECDSAPub(&ecdsaKey.PublicKey)
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
