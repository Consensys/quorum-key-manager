package ethereum

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	From       *ethcommon.Address // the sender of the 'transaction'
	To         *ethcommon.Address // the destination contract (nil for contract creation)
	Gas        *uint64            // if 0, the call executes with near-infinite gas
	GasPrice   *big.Int           // wei <-> gas exchange ratio
	Value      *big.Int           // amount of wei sent along with the call
	Data       *[]byte            // input data, usually an ABI-encoded contract method invocation
	GasFeeCap  *big.Int
	GasTipCap  *big.Int
	AccessList types.AccessList
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
	From       *ethcommon.Address `json:"from,omitempty"`
	To         *ethcommon.Address `json:"to,omitempty"`
	Gas        *hexutil.Uint64    `json:"gas,omitempty"`
	GasPrice   *hexutil.Big       `json:"gasPrice,omitempty"`
	Value      *hexutil.Big       `json:"value,omitempty"`
	Data       *hexutil.Bytes     `json:"data,omitempty"`
	GasFeeCap  *hexutil.Big       `json:"maxFeePerGas,omitempty"`
	GasTipCap  *hexutil.Big       `json:"maxPriorityFeePerGas,omitempty"`
	AccessList types.AccessList   `json:"accessList,omitempty"`
}

func (msg *CallMsg) UnmarshalJSON(b []byte) error {
	raw := new(jsonCallMsg)
	err := json.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	*msg = CallMsg{
		From:       raw.From,
		To:         raw.To,
		Gas:        (*uint64)(raw.Gas),
		GasPrice:   (*big.Int)(raw.GasPrice),
		Value:      (*big.Int)(raw.Value),
		Data:       (*[]byte)(raw.Data),
		GasTipCap:  (*big.Int)(raw.GasTipCap),
		GasFeeCap:  (*big.Int)(raw.GasFeeCap),
		AccessList: raw.AccessList,
	}

	if raw.Data != nil {
		msg.Data = (*[]byte)(raw.Data)
	}

	return nil
}

func (msg *CallMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonCallMsg{
		From:       msg.From,
		To:         msg.To,
		Gas:        (*hexutil.Uint64)(msg.Gas),
		GasPrice:   (*hexutil.Big)(msg.GasPrice),
		Value:      (*hexutil.Big)(msg.Value),
		Data:       (*hexutil.Bytes)(msg.Data),
		GasTipCap:  (*hexutil.Big)(msg.GasTipCap),
		GasFeeCap:  (*hexutil.Big)(msg.GasFeeCap),
		AccessList: msg.AccessList,
	})
}
