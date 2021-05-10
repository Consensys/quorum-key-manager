package ethereum

import (
	"encoding/json"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

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

func (msg *SendTxMsg) TxData() *TxData {
	data := &TxData{
		To:       msg.To,
		GasPrice: msg.GasPrice,
		Value:    msg.Value,
	}

	if msg.Gas != nil {
		data.GasLimit = *msg.Gas
	}

	if msg.Nonce != nil {
		data.Nonce = *msg.Nonce
	}

	if msg.Data != nil {
		data.Data = *msg.Data
	}

	return data.SetDefault()
}

type jsonSendTxMsg struct {
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
	raw := new(jsonSendTxMsg)
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
	return json.Marshal(&jsonSendTxMsg{
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

type EEATxData struct {
	Nonce uint64
	To    *ethcommon.Address
	Data  []byte
}

type SendEEATxMsg struct {
	From  ethcommon.Address
	To    *ethcommon.Address
	Nonce *uint64
	Data  *[]byte

	PrivateArgs
}

func (msg *SendEEATxMsg) TxData() *EEATxData {
	data := &EEATxData{
		To: msg.To,
	}

	if msg.Nonce != nil {
		data.Nonce = *msg.Nonce
	}

	if msg.Data != nil {
		data.Data = *msg.Data
	}

	return data
}

type jsonSendEEATxMsg struct {
	From  ethcommon.Address  `json:"from,omitempty"`
	To    *ethcommon.Address `json:"to,omitempty"`
	Nonce *hexutil.Uint64    `json:"nonce,omitempty"`
	Data  *hexutil.Bytes     `json:"data,omitempty"`

	PrivateArgs
}

func (msg *SendEEATxMsg) UnmarshalJSON(b []byte) error {
	raw := new(jsonSendEEATxMsg)
	err := json.Unmarshal(b, raw)
	if err != nil {
		return err
	}

	*msg = SendEEATxMsg{
		From:        raw.From,
		To:          raw.To,
		Nonce:       (*uint64)(raw.Nonce),
		PrivateArgs: raw.PrivateArgs,
		Data:        (*[]byte)(raw.Data),
	}

	return nil
}

func (msg *SendEEATxMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonSendEEATxMsg{
		From:        msg.From,
		To:          msg.To,
		Nonce:       (*hexutil.Uint64)(msg.Nonce),
		Data:        (*hexutil.Bytes)(msg.Data),
		PrivateArgs: msg.PrivateArgs,
	})
}