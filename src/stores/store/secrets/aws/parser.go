package aws

import (
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
)

func formatAwsSecret(id, value string, tags map[string]string, metadata *entities2.Metadata) *entities2.Secret {
	return &entities2.Secret{
		ID:       id,
		Value:    value,
		Tags:     tags,
		Metadata: metadata,
	}
}
