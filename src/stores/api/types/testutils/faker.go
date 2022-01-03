package testutils

import (
	"encoding/base64"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
)

func FakeSetSecretRequest() *types.SetSecretRequest {
	return &types.SetSecretRequest{
		Value: "my-secret-value",
		Tags:  testutils.FakeTags(),
	}
}

func FakeCreateKeyRequest() *types.CreateKeyRequest {
	return &types.CreateKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags:             testutils.FakeTags(),
	}
}

func FakeImportKeyRequest() *types.ImportKeyRequest {
	privKey, _ := base64.StdEncoding.DecodeString("2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw=")
	return &types.ImportKeyRequest{
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags:             testutils.FakeTags(),
		PrivateKey:       privKey,
	}
}

func FakeSignBase64PayloadRequest() *types.SignBase64PayloadRequest {
	return &types.SignBase64PayloadRequest{
		Data: []byte("my data to sign"),
	}
}

func FakeCreateEthAccountRequest() *types.CreateEthAccountRequest {
	return &types.CreateEthAccountRequest{
		KeyID: "my-key-account",
		Tags:  testutils.FakeTags(),
	}
}

func FakeImportEthAccountRequest() *types.ImportEthAccountRequest {
	return &types.ImportEthAccountRequest{
		KeyID:      "my-import-key-account",
		PrivateKey: hexutil.MustDecode("0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"),
		Tags:       testutils.FakeTags(),
	}
}

func FakeUpdateEthAccountRequest() *types.UpdateEthAccountRequest {
	return &types.UpdateEthAccountRequest{
		Tags: testutils.FakeTags(),
	}
}

func FakeSignMessageRequest() *types.SignMessageRequest {
	return &types.SignMessageRequest{
		Message: []byte("any message goes here"),
	}
}

func FakeSignTypedDataRequest() *types.SignTypedDataRequest {
	return &types.SignTypedDataRequest{
		DomainSeparator: types.DomainSeparator{
			Name:              "quorum",
			Version:           "v2.6.0",
			ChainID:           1,
			VerifyingContract: "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			Salt:              "mySalt",
		},
		Types: map[string][]types.Type{
			"Mail": {
				{Name: "sender", Type: "address"},
				{Name: "recipient", Type: "address"},
				{Name: "content", Type: "string"},
			},
		},
		Message: map[string]interface{}{
			"sender":    "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
			"recipient": "0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73",
			"content":   "my content",
		},
		MessageType: "Mail",
	}
}

func FakeSignETHTransactionRequest(txType string) *types.SignETHTransactionRequest {
	toAddress := common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")
	accessList := []ethtypes.AccessTuple{
		{
			Address:     toAddress,
			StorageKeys: []common.Hash{common.HexToHash("0xfeee")},
		},
	}

	req := &types.SignETHTransactionRequest{
		TransactionType: txType,
		Nonce:           0,
		To:              &toAddress,
		Value:           hexutil.Big(*hexutil.MustDecodeBig("0xfeee")),
		GasLimit:        21000,
		ChainID:         hexutil.Big(*hexutil.MustDecodeBig("0x1")),
		Data:            hexutil.MustDecode("0xfeee"),
	}

	switch txType {
	case types.LegacyTxType:
		req.GasPrice = hexutil.Big(*hexutil.MustDecodeBig("0xfeee"))
	case types.AccessListTxType:
		req.AccessList = accessList
	case "", types.DynamicFeeTxType:
		baseFee := hexutil.Big(*hexutil.MustDecodeBig("0xfeee"))
		minerTip := hexutil.Big(*hexutil.MustDecodeBig("0xfeee"))
		req.GasFeeCap = &baseFee
		req.GasTipCap = &minerTip
		req.AccessList = accessList
	default:
		return nil
	}

	return req
}

func FakeSignQuorumPrivateTransactionRequest() *types.SignQuorumPrivateTransactionRequest {
	toAddress := common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")

	return &types.SignQuorumPrivateTransactionRequest{
		Nonce:    0,
		To:       &toAddress,
		Value:    hexutil.Big(*hexutil.MustDecodeBig("0xfeee")),
		GasPrice: hexutil.Big(*hexutil.MustDecodeBig("0xfeee")),
		GasLimit: 21000,
		Data:     hexutil.MustDecode("0xfeee"),
	}
}

func FakeSignEEATransactionRequest() *types.SignEEATransactionRequest {
	toAddress := common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")

	return &types.SignEEATransactionRequest{
		Nonce:       0,
		To:          &toAddress,
		ChainID:     hexutil.Big(*hexutil.MustDecodeBig("0x1")),
		PrivateFrom: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
		PrivateFor:  []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
		Data:        hexutil.MustDecode("0xfeee"),
	}
}
