package eddsa

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"

	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func CreateX25519(importedPrivKey []byte) (privKey, pubKey []byte, err error) {
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

func SignX25519(privKeyB, data []byte) ([]byte, error) {
	if len(privKeyB) != ed25519.PrivateKeySize {
		return nil, errors.InvalidParameterError("invalid ED25519 private key size")
	}
	privKey := ed25519.PrivateKey(privKeyB)
	signature := ed25519.Sign(privKey, data)
	return signature, nil
}

func VerifyX25519Signature(publicKey, message, signature []byte) (bool, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return false, errors.InvalidParameterError("invalid ED25519 public key size")
	}
	pubKey := ed25519.PublicKey(publicKey)
	return ed25519.Verify(pubKey, message, signature), nil
}
