package ethereum

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

type SigningFunc func(data []byte) ([]byte, error)

func ECDSASigner(privKey *ecdsa.PrivateKey) SigningFunc {
	return func(data []byte) ([]byte, error) {
		return crypto.Sign(data, privKey)
	}
}

func HomesteadSign(tx *Transaction, sign SigningFunc) ([]byte, error) {
	return sign(FrontierHash(tx).Bytes())
}

func EIP155Sign(tx *Transaction, chainID *big.Int, sign SigningFunc) ([]byte, error) {
	return sign(EIP155Hash(tx, chainID).Bytes())
}

func EEASign(tx *Transaction, chainID *big.Int, args *EEAPrivateArgs, sign SigningFunc) ([]byte, error) {
	return sign(EEAHash(tx, chainID, args).Bytes())
}

func HomesteadSignatureValues(sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) != crypto.SignatureLength {
		return nil, nil, nil, fmt.Errorf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength)
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return
}

func EIP155SignatureValues(sig []byte, chainID *big.Int) (r, s, v *big.Int, err error) {
	r, s, v, err = HomesteadSignatureValues(sig)
	if err != nil {
		return nil, nil, nil, err
	}

	if chainID.Sign() != 0 {
		v = big.NewInt(int64(sig[64] + 35))
		v.Add(v, new(big.Int).Mul(chainID, big.NewInt(2)))
	}

	return
}

func EEASignatureValues(sig []byte, chainID *big.Int) (r, s, v *big.Int, err error) {
	// TODO: to be implemented
	return
}
