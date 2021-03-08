package testutils

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/models"
	"time"
)

func FakeSecret() *models.Secret {
	return &models.Secret{
		Value:    "my-secret",
		Disabled: false,
		Recovery: nil,
		Tags: map[string]string{
			"tag1": "tagValue1",
			"tag2": "tagValue2",
		},
		Version:   1,
		CreatedAt: time.Now(),
	}
}

func FakeAttributes() *models.Attributes {
	return &models.Attributes{
		Operations: []models.CryptoOperation{
			models.Signing, models.Encryption,
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
