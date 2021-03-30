package testutils

import (
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
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
	return &entities.Key{
		ID:        "my-key-id",
		PublicKey: "0x0433d7f005495fb6c0a34e22336dc3adcf4064553d5e194f77126bcac6da19491e0bab2772115cd284605d3bba94b69dc8c7a215021b58bcc87a70c9a440a3ff83",
		Algo:      FakeAlgorithm(),
		Metadata:  FakeMetadata(),
		Tags:      FakeTags(),
	}
}

func FakeMetadata() *entities.Metadata {
	return &entities.Metadata{
		Version:   1,
		Disabled:  false,
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
