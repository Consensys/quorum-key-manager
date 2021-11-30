package testutils

import (
	"encoding/base64"
	"fmt"
	entities2 "github.com/consensys/quorum-key-manager/src/entities"
	"time"

	common2 "github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/ethereum/go-ethereum/common"

	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func FakeSecret() *entities.Secret {
	return &entities.Secret{
		ID:       fmt.Sprintf("my-secret-%d", common2.RandInt(100)),
		Value:    fmt.Sprintf("my-secret-value-%s", common2.RandString(4)),
		Tags:     FakeTags(),
		Metadata: FakeMetadata(),
	}
}

func FakeKey() *entities.Key {
	pubKey, _ := base64.URLEncoding.DecodeString("BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=")
	return &entities.Key{
		ID:          fmt.Sprintf("my-key-%d", common2.RandInt(100)),
		PublicKey:   pubKey,
		Algo:        FakeAlgorithm(),
		Metadata:    FakeMetadata(),
		Tags:        FakeTags(),
		Annotations: FakeAnnotations(),
	}
}

func FakeETHAccount() *entities.ETHAccount {
	return &entities.ETHAccount{
		KeyID:               fmt.Sprintf("my-account-%d", common2.RandInt(100)),
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

func FakeAlgorithm() *entities2.Algorithm {
	return &entities2.Algorithm{
		Type:          entities2.Ecdsa,
		EllipticCurve: entities2.Secp256k1,
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
