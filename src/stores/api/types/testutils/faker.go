package testutils

import (
	"encoding/base64"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities/testutils"
)

func FakeSetSecretRequest() *types.SetSecretRequest {
	return &types.SetSecretRequest{
		ID:    "my-secret",
		Value: "my-secret-value",
		Tags:  testutils.FakeTags(),
	}
}

func FakeCreateKeyRequest() *types.CreateKeyRequest {
	return &types.CreateKeyRequest{
		ID:               "my-key",
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags:             testutils.FakeTags(),
	}
}

func FakeImportKeyRequest() *types.ImportKeyRequest {
	return &types.ImportKeyRequest{
		ID:               "my-key",
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags:             testutils.FakeTags(),
		PrivateKey:       "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw=",
	}
}

func FakeSignBase64PayloadRequest() *types.SignBase64PayloadRequest {
	return &types.SignBase64PayloadRequest{
		Data: base64.URLEncoding.EncodeToString([]byte("my data to sign")),
	}
}

func FakeCreateEth1AccountRequest() *types.CreateEth1AccountRequest {
	return &types.CreateEth1AccountRequest{
		ID:   "my-account",
		Tags: testutils.FakeTags(),
	}
}

func FakeImportEth1AccountRequest() *types.ImportEth1AccountRequest {
	return &types.ImportEth1AccountRequest{
		ID:         "my-account",
		PrivateKey: hexutil.MustDecode("0xdb337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"),
		Tags:       testutils.FakeTags(),
	}
}

func FakeUpdateEth1AccountRequest() *types.UpdateEth1AccountRequest {
	return &types.UpdateEth1AccountRequest{
		Tags: testutils.FakeTags(),
	}
}

func FakeSignHexPayloadRequest() *types.SignHexPayloadRequest {
	return &types.SignHexPayloadRequest{
		Data: hexutil.MustDecode("0xfeee"),
	}
}

func FakeSignTypedDataRequest() *types.SignTypedDataRequest {
	return &types.SignTypedDataRequest{
		DomainSeparator: types.DomainSeparator{
			Name:              "orchestrate",
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

func FakeSignETHTransactionRequest() *types.SignETHTransactionRequest {
	return &types.SignETHTransactionRequest{
		Nonce:    0,
		To:       "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		Value:    hexutil.Big(*hexutil.MustDecodeBig("0xfeee")),
		GasPrice: hexutil.Big(*hexutil.MustDecodeBig("0xfeee")),
		GasLimit: 21000,
		ChainID:  hexutil.Big(*hexutil.MustDecodeBig("0x1")),
		Data:     hexutil.MustDecode("0xfeee"),
	}
}

func FakeSignQuorumPrivateTransactionRequest() *types.SignQuorumPrivateTransactionRequest {
	return &types.SignQuorumPrivateTransactionRequest{
		Nonce:    0,
		To:       "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		Value:    hexutil.Big(*hexutil.MustDecodeBig("0xfeee")),
		GasPrice: hexutil.Big(*hexutil.MustDecodeBig("0xfeee")),
		GasLimit: 21000,
		Data:     hexutil.MustDecode("0xfeee"),
	}
}

func FakeSignEEATransactionRequest() *types.SignEEATransactionRequest {
	return &types.SignEEATransactionRequest{
		Nonce:       0,
		To:          "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
		ChainID:     hexutil.Big(*hexutil.MustDecodeBig("0x1")),
		PrivateFrom: "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=",
		PrivateFor:  []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="},
		Data:        hexutil.MustDecode("0xfeee"),
	}
}

func FakeECRecoverRequest() *types.ECRecoverRequest {
	return &types.ECRecoverRequest{
		Data:      hexutil.MustDecode("0xfeee"),
		Signature: hexutil.MustDecode("0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201"),
	}
}

func FakeVerifyTypedDataPayloadRequest() *types.VerifyTypedDataRequest {
	return &types.VerifyTypedDataRequest{
		TypedData: *FakeSignTypedDataRequest(),
		Signature: hexutil.MustDecode("0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201"),
		Address:   "0x5Cc634233E4a454d47aACd9fC68801482Fb02610",
	}
}

func FakeVerifyEth1SignatureRequest() *types.VerifyEth1SignatureRequest {
	return &types.VerifyEth1SignatureRequest{
		Data:      hexutil.MustDecode("0xfeee"),
		Signature: hexutil.MustDecode("0x3399aeb23d6564b3a0b220447e9f1bb2057ffb82cfb766147620aa6bc84938e26941e7583d6460fea405d99da897e88cab07a7fd0991c6c2163645c45d25e4b201"),
		Address:   "0x5Cc634233E4a454d47aACd9fC68801482Fb02610",
	}
}
