package eddsa

import (
	"bytes"
	"crypto/rand"
	"fmt"

	babyjubjub "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/quorum-key-manager/pkg/errors"
)

func CreateBabyjubjub(importedPrivKey []byte) (privKey, pubKey []byte, err error) {
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

func SignBabyjubjub(privKeyB, data []byte) ([]byte, error) {
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

func VerifyBabyJubJubSignature(publicKey, message, signature []byte) (bool, error) {
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
