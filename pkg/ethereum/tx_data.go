package ethereum

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type TxData struct {
	Nonce    uint64
	To       *ethcommon.Address
	Value    *big.Int
	GasPrice *big.Int
	GasLimit uint64
	Data     []byte
}

func (tx *TxData) SetDefault() {
	if tx.Value == nil {
		tx.Value = big.NewInt(0)
	}

	if tx.GasPrice == nil {
		tx.GasPrice = big.NewInt(0)
	}
}
