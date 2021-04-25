package ethereum

import (
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	From     *ethcommon.Address // the sender of the 'transaction'
	To       *ethcommon.Address // the destination contract (nil for contract creation)
	Gas      *uint64            // if 0, the call executes with near-infinite gas
	GasPrice *big.Int           // wei <-> gas exchange ratio
	Value    *big.Int           // amount of wei sent along with the call
	Data     *[]byte            // input data, usually an ABI-encoded contract method invocation
}

func (msg *CallMsg) WithFrom(addr ethcommon.Address) *CallMsg {
	msg.From = &addr
	return msg
}

func (msg *CallMsg) WithTo(addr ethcommon.Address) *CallMsg {
	msg.To = &addr
	return msg
}

func (msg *CallMsg) WithGas(gas uint64) *CallMsg {
	msg.Gas = &gas
	return msg
}

func (msg *CallMsg) WithGasPrice(p *big.Int) *CallMsg {
	msg.GasPrice = p
	return msg
}

func (msg *CallMsg) WithValue(v *big.Int) *CallMsg {
	msg.Value = v
	return msg
}

func (msg *CallMsg) WithData(d []byte) *CallMsg {
	if len(d) != 0 {
		b := make([]byte, len(d))
		copy(b, d)
		msg.Data = &b
	}
	return msg
}

type jsonCallMsg struct {
	From     *ethcommon.Address `json:"from,omitempty"`
	To       *ethcommon.Address `json:"to,omitempty"`
	Gas      *hexutil.Uint64    `json:"gas,omitempty"`
	GasPrice *hexutil.Big       `json:"gasPrice,omitempty"`
	Value    *hexutil.Big       `json:"value,omitempty"`
	Data     *hexutil.Bytes     `json:"data,omitempty"`
}

func (msg *CallMsg) UnmarshalJSON(b []byte) error {
	raw := new(jsonCallMsg)
	err := json.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	*msg = CallMsg{
		From:     raw.From,
		To:       raw.To,
		Gas:      (*uint64)(raw.Gas),
		GasPrice: (*big.Int)(raw.GasPrice),
		Value:    (*big.Int)(raw.Value),
		Data:     (*[]byte)(raw.Data),
	}

	if raw.Data != nil {
		msg.Data = (*[]byte)(raw.Data)
	}

	return nil
}

func (msg *CallMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonCallMsg{
		From:     msg.From,
		To:       msg.To,
		Gas:      (*hexutil.Uint64)(msg.Gas),
		GasPrice: (*hexutil.Big)(msg.GasPrice),
		Value:    (*hexutil.Big)(msg.Value),
		Data:     (*hexutil.Bytes)(msg.Data),
	})
}

type PrivacyFlag uint64

const (
	StandardPrivatePrivacyFlag PrivacyFlag = iota                              // 0
	PartyProtectionPrivacyFlag PrivacyFlag = 1 << PrivacyFlag(iota-1)          // 1
	StateValidationPrivacyFlag             = iota | PartyProtectionPrivacyFlag // 3 which includes PrivacyFlagPartyProtection
)

// PrivateArgs arguments for private transactions
type PrivateArgs struct {
	PrivateFrom    *string      `json:"privateFrom,omitempty"`
	PrivateFor     *[]string    `json:"privateFor,omitempty"`
	PrivateType    *string      `json:"restriction,omitempty"`
	PrivacyFlag    *PrivacyFlag `json:"privacyFlag,omitempty"`
	PrivacyGroupID *string      `json:"privacyGroupId,omitempty"`
}

func (args *PrivateArgs) WithPrivateFrom(pubKey string) *PrivateArgs {
	args.PrivateFrom = &pubKey
	return args
}

func (args *PrivateArgs) WithPrivateFor(pubKeys []string) *PrivateArgs {
	args.PrivateFor = &pubKeys
	return args
}

func (args *PrivateArgs) WithPrivateType(pubKey string) *PrivateArgs {
	args.PrivateFrom = &pubKey
	return args
}

func (args *PrivateArgs) WithPrivacyFlag(flag PrivacyFlag) *PrivateArgs {
	args.PrivacyFlag = &flag
	return args
}

func (args *PrivateArgs) WithPrivacyGroupID(id string) *PrivateArgs {
	args.PrivacyGroupID = &id
	return args
}

type SendTxMsg struct {
	From     ethcommon.Address
	To       *ethcommon.Address
	Gas      *uint64
	GasPrice *big.Int
	Value    *big.Int
	Nonce    *uint64
	Data     *[]byte

	PrivateArgs
}

func (msg *SendTxMsg) IsPrivate() bool {
	return msg.PrivateArgs != PrivateArgs{}
}

func (msg *SendTxMsg) TxData() (*TxData, error) {
	if msg.Gas == nil {
		return nil, fmt.Errorf("gas not specified")
	}

	if msg.GasPrice == nil {
		return nil, fmt.Errorf("gasPrice not specified")
	}

	if msg.Nonce == nil {
		return nil, fmt.Errorf("nonce not specified")
	}

	tx := &TxData{
		To:       msg.To,
		GasLimit: *msg.Gas,
		GasPrice: msg.GasPrice,
		Value:    msg.Value,
		Nonce:    *msg.Nonce,
	}

	if msg.Data != nil {
		tx.Data = *msg.Data
	}

	tx.SetDefault()

	return tx, nil
}

type jsonCallSendTxMsg struct {
	From     ethcommon.Address  `json:"from,omitempty"`
	To       *ethcommon.Address `json:"to,omitempty"`
	Gas      *hexutil.Uint64    `json:"gas,omitempty"`
	GasPrice *hexutil.Big       `json:"gasPrice,omitempty"`
	Value    *hexutil.Big       `json:"value,omitempty"`
	Nonce    *hexutil.Uint64    `json:"nonce,omitempty"`
	Data     *hexutil.Bytes     `json:"data,omitempty"`
	Input    *hexutil.Bytes     `json:"input,omitempty"`

	PrivateArgs
}

func (msg *SendTxMsg) UnmarshalJSON(b []byte) error {
	raw := new(jsonCallSendTxMsg)
	err := json.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	*msg = SendTxMsg{
		From:        raw.From,
		To:          raw.To,
		Gas:         (*uint64)(raw.Gas),
		GasPrice:    (*big.Int)(raw.GasPrice),
		Value:       (*big.Int)(raw.Value),
		Nonce:       (*uint64)(raw.Nonce),
		PrivateArgs: raw.PrivateArgs,
		Data:        (*[]byte)(raw.Input),
	}

	if raw.Data != nil {
		msg.Data = (*[]byte)(raw.Data)
	}

	return nil
}

func (msg *SendTxMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonCallSendTxMsg{
		From:        msg.From,
		To:          msg.To,
		Gas:         (*hexutil.Uint64)(msg.Gas),
		GasPrice:    (*hexutil.Big)(msg.GasPrice),
		Value:       (*hexutil.Big)(msg.Value),
		Nonce:       (*hexutil.Uint64)(msg.Nonce),
		Data:        (*hexutil.Bytes)(msg.Data),
		PrivateArgs: msg.PrivateArgs,
	})
}
