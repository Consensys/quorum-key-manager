package testutils

import (
	"encoding/base64"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
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

func FakeSignPayloadRequest() *types.SignPayloadRequest {
	return &types.SignPayloadRequest{
		Data: base64.URLEncoding.EncodeToString([]byte("my data to sign")),
	}
}
