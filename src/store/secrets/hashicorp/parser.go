package hashicorp

import (
	"encoding/json"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"

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
		return nil, errors.HashicorpVaultConnectionError("failed to parse hashicorp created time from data")
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

	secretVersion := jsonMetadata["versions"].(map[string]interface{})[version].(map[string]interface{})
	if secretVersion["deletion_time"].(string) != "" {
		deletionTime, err := time.Parse(time.RFC3339, secretVersion["deletion_time"].(string))
		if err != nil {
			return nil, errors.HashicorpVaultConnectionError("failed to parse deletion time from metadata")
		}

		metadata.DeletedAt = deletionTime
		metadata.Disabled = true

		// If secret has been destroyed, deletion time is the destroyed time
		if secretVersion["destroyed"].(bool) {
			metadata.DestroyedAt = deletionTime
		}
	}

	var err error
	metadata.CreatedAt, err = time.Parse(time.RFC3339, secretVersion["created_time"].(string))
	if err != nil {
		return nil, errors.HashicorpVaultConnectionError("failed to parse created time from metadata")
	}
	metadata.UpdatedAt = metadata.CreatedAt

	expirationDurationStr := jsonMetadata["delete_version_after"].(string)
	if expirationDurationStr != "0s" {
		expirationDuration, der := time.ParseDuration(expirationDurationStr)
		if der != nil {
			return nil, errors.HashicorpVaultConnectionError("failed to parse expiration time from metadata")
		}

		metadata.ExpireAt = metadata.CreatedAt.Add(expirationDuration)
	}

	return metadata, nil
}
