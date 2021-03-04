package secrets

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/models"
	"github.com/ConsenSysQuorum/quorum-key-manager/infra/vault"
	"path"
	"time"
)

const (
	dataLabel        = "data"
	metadataLabel    = "metadata"
	valueLabel       = "value"
	deleteAfterLabel = "delete_version_after"
	tagsLabel        = "tags"
)

// Store is an implementation of secret store relying on HashiCorp Vault kv-v2 secret engine
type hashicorpSecretStore struct {
	client     vault.HashicorpVaultClient
	mountPoint string
}

// New creates an HashiCorp secret store
func New(client vault.HashicorpVaultClient, mountPoint string) *hashicorpSecretStore {
	return &hashicorpSecretStore{
		client:     client,
		mountPoint: mountPoint,
	}
}

func (s *hashicorpSecretStore) Info(context.Context) (*models.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

// Set a secret
func (s *hashicorpSecretStore) Set(ctx context.Context, id, value string, attr *models.Attributes) (*models.Secret, error) {
	data := map[string]interface{}{
		valueLabel: value,
		tagsLabel:  attr.Tags,
	}

	hashicorpSecret, err := s.client.Write(s.pathData(id), data)
	if err != nil {
		return nil, err
	}

	// Hashicorp only returns metadata as the "data" field when creating a new secret
	metadata, err := extractMetadata(hashicorpSecret.Data)
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(value, attr.Tags, metadata), nil
}

// Get a secret
func (s *hashicorpSecretStore) Get(ctx context.Context, id string, version int) (*models.Secret, error) {
	// Get latest version by default if no version is specified, otherwise get specific version
	pathData := s.pathData(id)
	if version != 0 {
		pathData = fmt.Sprintf("%v?version=%v", pathData, version)
	}

	hashicorpSecret, err := s.client.Read(pathData)
	if err != nil {
		return nil, err
	}

	// Hashicorp returns metadata in "data.metadata" and data in "data.data" fields
	data := hashicorpSecret.Data[dataLabel].(map[string]interface{})
	metadata, err := extractMetadata(hashicorpSecret.Data[metadataLabel].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(data[valueLabel].(string), data[tagsLabel].(map[string]string), metadata), nil
}

// Get all secret ids
func (s *hashicorpSecretStore) List(ctx context.Context) ([]string, error) {
	res, err := s.client.List(path.Join(s.mountPoint, "metadata"))
	if err != nil {
		return nil, err
	}

	if res == nil {
		return []string{}, nil
	}

	return res.Data["keys"].([]string), nil
}

// Refresh an existing secret by extending its TTL
func (s *hashicorpSecretStore) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	data := make(map[string]interface{})
	if !expirationDate.IsZero() {
		data[deleteAfterLabel] = time.Since(expirationDate).String()
	}

	_, err := s.client.Write(s.pathMetadata(id), data)
	if err != nil {
		return err
	}

	return nil
}

// Delete a secret
func (s *hashicorpSecretStore) Delete(ctx context.Context, id string, versions ...int) (*models.Secret, error) {
	return nil, errors.NotImplementedError
}

// Gets a deleted secret
func (s *hashicorpSecretStore) GetDeleted(ctx context.Context, id string) (*models.Secret, error) {
	return nil, errors.NotImplementedError
}

// Lists all deleted secrets
func (s *hashicorpSecretStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.NotImplementedError
}

// Undelete a previously deleted secret
func (s *hashicorpSecretStore) Undelete(ctx context.Context, id string) error {
	return errors.NotImplementedError
}

// Destroy a secret permanenty
func (s *hashicorpSecretStore) Destroy(ctx context.Context, id string, versions ...int) error {
	return errors.NotImplementedError
}

// path compute path from hashicorp mount
func (s *hashicorpSecretStore) pathData(id string) string {
	return path.Join(s.mountPoint, dataLabel, id)
}

func (s *hashicorpSecretStore) pathMetadata(id string) string {
	return path.Join(s.mountPoint, metadataLabel, id)
}

func formatHashicorpSecret(value string, tags map[string]string, metadata *models.Metadata) *models.Secret {
	return &models.Secret{
		Value:       value,
		Tags:        tags,
		Disabled:    metadata.Disabled,
		Version:     metadata.Version,
		ExpireAt:    metadata.ExpireAt,
		CreatedAt:   metadata.CreatedAt,
		DeletedAt:   metadata.DeletedAt,
		DestroyedAt: metadata.DestroyedAt,
	}
}

func extractMetadata(data map[string]interface{}) (*models.Metadata, error) {
	metadata := &models.Metadata{
		Version: data["version"].(int),
	}

	var err error

	metadata.CreatedAt, err = time.Parse(time.RFC3339, data["created_time"].(string))
	if err != nil {
		return nil, err
	}

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
