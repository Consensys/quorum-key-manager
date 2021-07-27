package testutils

import (
	"encoding/base64"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
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
		ID:        "my-key-id",
		PublicKey: pubKey,
		Algo:      FakeAlgorithm(),
		Metadata:  FakeMetadata(),
		Tags:      FakeTags(),
	}
}

func FakeETH1Account() *entities.ETH1Account {
	key := FakeKey()
	return &entities.ETH1Account{
		Address:  common.HexToAddress("0x83a0254be47813BBff771F4562744676C4e793F0"),
		KeyID:    key.ID,
		Key:      key,
		Metadata: FakeMetadata(),
		Tags:     FakeTags(),
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
