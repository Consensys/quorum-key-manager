package testutils

import (
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

func FakeSecret() *entities.Secret {
	return &entities.Secret{
		Value:    "my-secret",
		Disabled: false,
		Recovery: nil,
		Tags: map[string]string{
			"tag1": "tagValue1",
			"tag2": "tagValue2",
		},
		Metadata: FakeMetadata(),
	}
}

func FakeMetadata() *entities.Metadata {
	return &entities.Metadata{
		Version:   1,
		CreatedAt: time.Now(),
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
		Tags: map[string]string{
			"tag1": "tagValue1",
			"tag2": "tagValue2",
		},
	}
}
