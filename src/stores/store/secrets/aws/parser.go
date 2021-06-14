package aws

import (
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
)

func formatAwsSecret(id, value string, tags map[string]string, metadata *entities.Metadata) *entities.Secret {
	return &entities.Secret{
		ID:       id,
		Value:    value,
		Tags:     tags,
		Metadata: metadata,
	}
}
