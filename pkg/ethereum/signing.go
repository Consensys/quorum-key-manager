package ethereum

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

func EIP155Signature(sig []byte, chainID *big.Int) ([]byte, error) {
	r, s, v, err := homesteadSignatureValues(sig)
	if err != nil {
		return nil, err
	}

	if chainID.Sign() != 0 {
		v = big.NewInt(int64(sig[64] + 35))
		v.Add(v, new(big.Int).Mul(chainID, big.NewInt(2)))
	}

	return append(append(r.Bytes(), s.Bytes()...), v.Bytes()...), nil
}

func homesteadSignatureValues(sig []byte) (r, s, v *big.Int, err error) {
	if len(sig) != crypto.SignatureLength {
		return nil, nil, nil, fmt.Errorf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength)
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return
}
