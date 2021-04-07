package hashicorp

import (
	"encoding/json"
	"time"

	"github.com/hashicorp/vault/api"

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

func formatHashicorpSecretData(jsonData map[string]interface{}) (*entities.Metadata, error) {
	metadata := &entities.Metadata{
		Version:  jsonData[versionLabel].(json.Number).String(),
		Disabled: false,
	}

	var err error

	metadata.CreatedAt, err = time.Parse(time.RFC3339, jsonData["created_time"].(string))
	if err != nil {
		return nil, err
	}

	metadata.UpdatedAt = metadata.CreatedAt

	return metadata, nil
}

func formatHashicorpSecretMetadata(secret *api.Secret, version string) (*entities.Metadata, error) {
	jsonMetadata := secret.Data

	if version == "" {
		version = jsonMetadata["current_version"].(json.Number).String()
	}

	metadata := &entities.Metadata{
		Version: version,
	}

	var err error
	metadata.CreatedAt, err = time.Parse(time.RFC3339, jsonMetadata["created_time"].(string))
	if err != nil {
		return nil, err
	}

	metadata.UpdatedAt, err = time.Parse(time.RFC3339, jsonMetadata["updated_time"].(string))
	if err != nil {
		return nil, err
	}

	expirationDurationStr := jsonMetadata["delete_version_after"].(string)
	if expirationDurationStr != "0s" {
		expirationDuration, err := time.ParseDuration(expirationDurationStr)
		if err != nil {
			return nil, err
		}

		metadata.ExpireAt = metadata.CreatedAt.Add(expirationDuration)
	}

	secretVersion := jsonMetadata["versions"].(map[string]interface{})[version].(map[string]interface{})
	if secretVersion["deletion_time"].(string) != "" {
		deletionTime, err := time.Parse(time.RFC3339, jsonMetadata["deletion_time"].(string))
		if err != nil {
			return nil, err
		}

		metadata.DeletedAt = deletionTime
		metadata.Disabled = true

		// If secret has been destroyed, deletion time is the destroyed time
		if jsonMetadata["destroyed"].(bool) {
			metadata.DestroyedAt = deletionTime
		}
	}

	return metadata, nil
}
