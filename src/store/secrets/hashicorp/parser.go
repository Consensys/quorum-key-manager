package hashicorp

import (
	"encoding/json"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

func formatHashicorpSecret(id, value string, tags map[string]string, metadata *entities.Metadata) *entities.Secret {
	return &entities.Secret{
		ID:       id,
		Value:    value,
		Tags:     tags,
		Metadata: metadata,
	}
}

func extractMetadata(data map[string]interface{}) (*entities.Metadata, error) {
	version := data[versionLabel].(json.Number)
	metadata := &entities.Metadata{
		Version: version.String(),
	}

	var err error

	metadata.CreatedAt, err = time.Parse(time.RFC3339, data["created_time"].(string))
	if err != nil {
		return nil, err
	}

	metadata.UpdatedAt = metadata.CreatedAt

	if data["deletion_time"].(string) != "" {
		deletionTime, err := time.Parse(time.RFC3339, data["deletion_time"].(string))
		if err != nil {
			return nil, err
		}

		// If deletion time is in the future, we populate the expireAt property, otherwise it has been deleted
		if deletionTime.After(time.Now()) {
			metadata.ExpireAt = deletionTime
		} else {
			metadata.DeletedAt = deletionTime
			metadata.Disabled = true
		}

		// If secret has been destroyed, deletion time is the destroyed time
		if data["destroyed"].(bool) {
			metadata.DestroyedAt = deletionTime
		}
	}

	return metadata, nil
}
