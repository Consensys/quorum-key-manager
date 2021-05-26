package aws

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

func formatAwsSecret(id, value string, tags map[string]string, metadata *entities.Metadata) *entities.Secret {
	return &entities.Secret{
		ID:       id,
		Value:    value,
		Tags:     tags,
		Metadata: metadata,
	}
}
