package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"

	babyjubjub "github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

func EdDSABabyjubjub(importedPrivKey []byte) (privKey []byte, pubKey []byte, err error) {
	babyJubJubPrivKey := babyjubjub.PrivateKey{}
	if importedPrivKey != nil {
		_, err = babyJubJubPrivKey.SetBytes(importedPrivKey)
		if err != nil {
			return nil, nil, err
		}
	}

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

	privKey = babyJubJubPrivKey.Bytes()
	pubKey = babyJubJubPrivKey.Public().Bytes()
	return pubKey, privKey, nil
}

func SignECDSA256k1(privKey, data []byte) ([]byte, error) {
	if len(data) != crypto.DigestLength {
		return nil, errors.InvalidParameterError(fmt.Sprintf("data is required to be exactly %d bytes (%d)", crypto.DigestLength, len(data)))
	}

	ecdsaPrivKey, err := crypto.ToECDSA(privKey)
	if err != nil {
		return nil, errors.InvalidParameterError(fmt.Sprintf("failed to parse private key. %s", err.Error()))
	}

	signature, err := crypto.Sign(data, ecdsaPrivKey)
	if err != nil {
		return nil, errors.CryptoOperationError(fmt.Sprintf("failed to sign. %s", err.Error()))
	}

	// We remove the recID from the signature (last byte).
	return signature[:len(signature)-1], nil
}

func VerifyECDSASignature(publicKey, message, signature []byte) (bool, error) {
	pubKey, err := crypto.UnmarshalPubkey(publicKey)
	if err != nil {
		return false, err
	}

	r := new(big.Int).SetBytes(signature[0:32])
	s := new(big.Int).SetBytes(signature[32:64])

	return ecdsa.Verify(pubKey, message, r, s), nil
}
