package testutils

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
		PrivateKey:       "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c",
	}
}

func FakeSignPayloadRequest() *types.SignPayloadRequest {
	return &types.SignPayloadRequest{
		Data: hexutil.Encode([]byte("my data to sign")),
	}
}
