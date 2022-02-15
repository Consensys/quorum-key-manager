package ecdsa

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

func CreateSecp256k1(importedPrivKey []byte) (privKey, pubKey []byte, err error) {
	var ecdsaKey *ecdsa.PrivateKey
	if importedPrivKey != nil {
		ecdsaKey, err = crypto.ToECDSA(importedPrivKey)
		if err != nil {
			return nil, nil, err
		}
	} else {
		ecdsaKey, err = crypto.GenerateKey()
		if err != nil {
			return nil, nil, err
		}
	}

	privKey = crypto.FromECDSA(ecdsaKey)
	pubKey = crypto.FromECDSAPub(&ecdsaKey.PublicKey)
	return privKey, pubKey, nil
}

func SignSecp256k1(privKey, data []byte) ([]byte, error) {
	if len(data) != crypto.DigestLength {
		return nil, fmt.Errorf("data is required to be exactly %d bytes (%d)", crypto.DigestLength, len(data))
	}

	ecdsaPrivKey, err := crypto.ToECDSA(privKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key. %s", err.Error())
	}

	signature, err := crypto.Sign(data, ecdsaPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign. %s", err.Error())
	}

	// We remove the recID from the signature (last byte).
	return signature[:len(signature)-1], nil
}

func VerifySecp256k1Signature(publicKey, message, signature []byte) (bool, error) {
	pubKey, err := crypto.UnmarshalPubkey(publicKey)
	if err != nil {
		return false, err
	}
	if len(signature) != 64 {
		return false, fmt.Errorf("invalid secp256k1 signature length")
	}

	r := new(big.Int).SetBytes(signature[0:32])
	s := new(big.Int).SetBytes(signature[32:64])

	return ecdsa.Verify(pubKey, message, r, s), nil
}
