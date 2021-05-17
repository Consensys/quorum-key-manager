package ethereum

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type BlockNumber int64

const (
	PendingBlockNumber  = BlockNumber(-2)
	LatestBlockNumber   = BlockNumber(-1)
	EarliestBlockNumber = BlockNumber(0)
)

func (n BlockNumber) MarshalText() ([]byte, error) {
	switch n {
	case -2:
		return []byte(`pending`), nil
	case -1:
		return []byte(`latest`), nil
	case 0:
		return []byte(`earliest`), nil
	default:
		buf := make([]byte, 2, 10)
		copy(buf, `0x`)
		buf = strconv.AppendUint(buf, uint64(n), 16)
		return buf, nil
	}
}

// UnmarshalJSON parses the given JSON fragment into a BlockNumber. It supports:
// - "latest", "earliest" or "pending" as string arguments
// - the block number
// Returned errors:
// - an invalid block number error when the given argument isn't a known strings
// - an out of range error when the given block number is either too little or too large
func (n *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case "earliest":
		*n = EarliestBlockNumber
		return nil
	case "latest":
		*n = LatestBlockNumber
		return nil
	case "pending":
		*n = PendingBlockNumber
		return nil
	}

	blckNum, err := hexutil.DecodeUint64(input)
	if err != nil {
		return err
	}
	if blckNum > math.MaxInt64 {
		return fmt.Errorf("block number larger than int64")
	}
	*n = BlockNumber(blckNum)
	return nil
}

func (n BlockNumber) Int64() int64 {
	return int64(n)
}
