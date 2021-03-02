package ethereum

import (
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	Nonce    uint64
	To       *ethcommon.Address
	Value    *big.Int
	GasPrice *big.Int
	GasLimit uint64
	Data     []byte
}

type EEAPrivateArgs struct {
	PrivateFrom    string
	PrivateFor     []string
	PrivacyGroupID string
}
