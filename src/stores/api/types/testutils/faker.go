package testutils

import (
	"encoding/base64"
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"
	testutils2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities/testutils"
)

func FakeSetSecretRequest() *types2.SetSecretRequest {
	return &types2.SetSecretRequest{
		ID:    "my-secret",
		Value: "my-secret-value",
		Tags:  testutils2.FakeTags(),
	}
}

func FakeCreateKeyRequest() *types2.CreateKeyRequest {
	return &types2.CreateKeyRequest{
		ID:               "my-key",
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags:             testutils2.FakeTags(),
	}
}

func FakeImportKeyRequest() *types2.ImportKeyRequest {
	return &types2.ImportKeyRequest{
		ID:               "my-key",
		Curve:            "secp256k1",
		SigningAlgorithm: "ecdsa",
		Tags:             testutils2.FakeTags(),
		PrivateKey:       "2zN8oyleQFBYZ5PyUuZB87OoNzkBj6TM4BqBypIOfhw=",
	}
}

func FakeSignPayloadRequest() *types2.SignBase64PayloadRequest {
	return &types2.SignBase64PayloadRequest{
		Data: base64.URLEncoding.EncodeToString([]byte("my data to sign")),
	}
}
