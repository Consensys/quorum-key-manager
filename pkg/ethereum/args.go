package ethereum

import (
	"encoding/json"
	"math/big"

	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

type PrivacyFlag uint64
type PrivateType string

const (
	StandardPrivatePrivacyFlag PrivacyFlag = iota                              // 0
	PartyProtectionPrivacyFlag PrivacyFlag = 1 << PrivacyFlag(iota-1)          // 1
	StateValidationPrivacyFlag             = iota | PartyProtectionPrivacyFlag // 3 which includes PrivacyFlagPartyProtection
)

const (
	PrivateTypeRestricted   PrivateType = "restricted"
	PrivateTypeUnrestricted PrivateType = "unrestricted"
)

// TODO: Delete usage of unnecessary pointers: https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/consensys/quorum-key-manager/96
// PrivateArgs arguments for private transactions
type PrivateArgs struct {
	PrivateFrom    *string      `json:"privateFrom,omitempty"`
	PrivateFor     *[]string    `json:"privateFor,omitempty"`
	PrivateType    *PrivateType `json:"restriction,omitempty"`
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

// TODO: Delete usage of unnecessary pointers: https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/consensys/quorum-key-manager/96
type SendTxMsg struct {
	From       ethcommon.Address
	To         *ethcommon.Address
	Gas        *uint64
	GasPrice   *big.Int
	Value      *big.Int
	Nonce      *uint64
	Data       *[]byte
	GasFeeCap  *big.Int
	GasTipCap  *big.Int
	AccessList types.AccessList

	PrivateArgs
}

func (msg *SendTxMsg) IsPrivate() bool {
	return msg.PrivateArgs != PrivateArgs{}
}

func (msg *SendTxMsg) IsLegacy() bool {
	return msg.GasPrice != nil
}

func (msg *SendTxMsg) TxData(txType int, chainID *big.Int) *types.Transaction {
	var txData types.TxData

	switch txType {
	case types.LegacyTxType:
		txData = &types.LegacyTx{
			Nonce:    *msg.Nonce,
			GasPrice: msg.GasPrice,
			Gas:      *msg.Gas,
			To:       msg.To,
			Value:    msg.Value,
			Data:     *msg.Data,
		}
	case types.DynamicFeeTxType:
		txData = &types.DynamicFeeTx{
			ChainID:    chainID,
			Nonce:      *msg.Nonce,
			GasTipCap:  msg.GasTipCap,
			GasFeeCap:  msg.GasFeeCap,
			Gas:        *msg.Gas,
			To:         msg.To,
			Value:      msg.Value,
			Data:       *msg.Data,
			AccessList: msg.AccessList,
		}
	}

	return types.NewTx(txData)
}

// TODO: Delete this function and use only go-quorum types when
func (msg *SendTxMsg) TxDataQuorum() *quorumtypes.Transaction {
	if msg.To == nil {
		return quorumtypes.NewContractCreation(*msg.Nonce, msg.Value, *msg.Gas, msg.GasPrice, *msg.Data)
	}

	return quorumtypes.NewTransaction(*msg.Nonce, *msg.To, msg.Value, *msg.Gas, msg.GasPrice, *msg.Data)
}

// TODO: Delete usage of unnecessary pointers: https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/consensys/quorum-key-manager/96
type jsonSendTxMsg struct {
	From       ethcommon.Address  `json:"from,omitempty"`
	To         *ethcommon.Address `json:"to,omitempty"`
	Gas        *hexutil.Uint64    `json:"gas,omitempty"`
	GasPrice   *hexutil.Big       `json:"gasPrice,omitempty"`
	Value      *hexutil.Big       `json:"value,omitempty"`
	Nonce      *hexutil.Uint64    `json:"nonce,omitempty"`
	Data       *hexutil.Bytes     `json:"data,omitempty"`
	Input      *hexutil.Bytes     `json:"input,omitempty"`
	GasFeeCap  *hexutil.Big       `json:"maxFeePerGas,omitempty"`
	GasTipCap  *hexutil.Big       `json:"maxPriorityFeePerGas,omitempty"`
	AccessList types.AccessList   `json:"accessList,omitempty"`

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
		GasTipCap:   (*big.Int)(raw.GasTipCap),
		GasFeeCap:   (*big.Int)(raw.GasFeeCap),
		AccessList:  raw.AccessList,
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
		GasTipCap:   (*hexutil.Big)(msg.GasTipCap),
		GasFeeCap:   (*hexutil.Big)(msg.GasFeeCap),
		AccessList:  msg.AccessList,
		PrivateArgs: msg.PrivateArgs,
	})
}

// TODO: Delete usage of unnecessary pointers: https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/consensys/quorum-key-manager/96
type SendEEATxMsg struct {
	From     ethcommon.Address
	To       *ethcommon.Address
	Nonce    *uint64
	Data     *[]byte
	Gas      *uint64
	GasPrice *big.Int

	PrivateArgs
}

func (msg *SendEEATxMsg) TxData() *types.Transaction {
	var data []byte
	if msg.Data != nil {
		data = *msg.Data
	}

	return types.NewTx(&types.LegacyTx{
		Nonce:    *msg.Nonce,
		GasPrice: msg.GasPrice,
		Gas:      *msg.Gas,
		To:       msg.To,
		Value:    nil,
		Data:     data,
	})
}

// TODO: Delete usage of unnecessary pointers: https://app.zenhub.com/workspaces/orchestrate-5ea70772b186e10067f57842/issues/consensys/quorum-key-manager/96
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
