package testutils

import (
	"encoding/base64"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func FakeSecret() *entities.Secret {
	return &entities.Secret{
		ID:       "my-secret-id",
		Value:    "my-secret-value",
		Tags:     FakeTags(),
		Metadata: FakeMetadata(),
	}
}

func FakeKey() *entities.Key {
	pubKey, _ := base64.URLEncoding.DecodeString("BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=")
	return &entities.Key{
		ID:          "my-key-id",
		PublicKey:   pubKey,
		Algo:        FakeAlgorithm(),
		Metadata:    FakeMetadata(),
		Tags:        FakeTags(),
		Annotations: FakeAnnotations(),
	}
}

func FakeETH1Account() *entities.ETH1Account {
	return &entities.ETH1Account{
		KeyID:               "my-account",
		Address:             common.HexToAddress("0x83a0254be47813BBff771F4562744676C4e793F0"),
		Metadata:            FakeMetadata(),
		PublicKey:           hexutil.MustDecode("0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2"),
		CompressedPublicKey: hexutil.MustDecode("0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f"),
		Tags:                FakeTags(),
	}
}

func FakeMetadata() *entities.Metadata {
	return &entities.Metadata{
		Version:   "1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func FakeAttributes() *entities.Attributes {
	return &entities.Attributes{
		Operations: []entities.CryptoOperation{
			entities.Signing, entities.Encryption,
		},
		Disabled: false,
		TTL:      24 * time.Hour,
		Recovery: nil,
		Tags:     FakeTags(),
	}
}

func FakeAlgorithm() *entities.Algorithm {
	return &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}
}

func FakeTags() map[string]string {
	return map[string]string{
		"tag1": "tagValue1",
		"tag2": "tagValue2",
	}
}

func FakeAnnotations() *entities.Annotation {
	return &entities.Annotation{
		AWSKeyID:             "awsKeyID",
		AWSCustomKeyStoreID:  "awsCustomKeyStoreID",
		AWSCloudHsmClusterID: "awsCloudHsmClusterID",
		AWSAccountID:         "awsAccountID",
		AWSArn:               "awsARN",
	}
}
