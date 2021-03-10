package testutils

import "github.com/ConsenSysQuorum/quorum-key-manager/pkg/sdk/types"

func FakeCreateSecretRequest() *types.CreateSecretRequest {
	return &types.CreateSecretRequest{
		ID:    "my-privateKey",
		Value: "fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249",
		Tags: map[string]string{
			"label1": "labelValue1",
			"label2": "labelValue2",
		},
	}
}
