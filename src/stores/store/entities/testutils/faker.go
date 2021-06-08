package testutils

import (
	"encoding/base64"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func FakeSecret() *entities2.Secret {
	return &entities2.Secret{
		ID:       "my-secret-id",
		Value:    "my-secret-value",
		Tags:     FakeTags(),
		Metadata: FakeMetadata(),
	}
}

func FakeKey() *entities2.Key {
	pubKey, _ := base64.URLEncoding.DecodeString("BFVSFJhqUh9DQJwcayNtsWdDMvqq8R_EKnBHqwd4Hr5vCXTyJlqKfYIgj4jCGixVZjsz5a-S2RklJRFjjoLf-LI=")
	return &entities2.Key{
		ID:        "my-key-id",
		PublicKey: pubKey,
		Algo:      FakeAlgorithm(),
		Metadata:  FakeMetadata(),
		Tags:      FakeTags(),
	}
}

func FakeETH1Account() *entities2.ETH1Account {
	return &entities2.ETH1Account{
		ID:                  "my-account",
		Address:             "0x83a0254be47813BBff771F4562744676C4e793F0",
		Metadata:            FakeMetadata(),
		PublicKey:           hexutil.MustDecode("0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2"),
		CompressedPublicKey: hexutil.MustDecode("0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f"),
		Tags:                FakeTags(),
	}
}

func FakeMetadata() *entities2.Metadata {
	return &entities2.Metadata{
		Version:   "1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func FakeAttributes() *entities2.Attributes {
	return &entities2.Attributes{
		Operations: []entities2.CryptoOperation{
			entities2.Signing, entities2.Encryption,
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
