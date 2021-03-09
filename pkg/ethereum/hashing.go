package ethereum

import (
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

func RLPHash(x interface{}) (h ethcommon.Hash) {
	hw := sha3.NewLegacyKeccak256()
	err := rlp.Encode(hw, x)
	if err != nil {
		panic(fmt.Sprintf("can not rlp encode: %v", x))
	}
	hw.Sum(h[:0])
	return h
}

func FrontierHash(tx *Transaction) ethcommon.Hash {
	return RLPHash([]interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
	})
}

func EIP155Hash(tx *Transaction, chainID *big.Int) ethcommon.Hash {
	return RLPHash([]interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
		chainID,
		uint(0),
		uint(0),
	})
}

func EEAHash(tx *Transaction, chainID *big.Int, args *EEAPrivateArgs) ethcommon.Hash {
	// TODO: to be implemented
	return ethcommon.Hash{}
}
