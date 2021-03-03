package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/models"
	"github.com/ConsenSysQuorum/quorum-key-manager/libs/vault"
	"path"
	"time"

	hashicorp "github.com/hashicorp/vault/api"
)

const (
	valueLabel          = "value"
	expirationDateLabel = "expirationDate"
	tagsLabel           = "tags"
	enabledLabel        = "enabled"
)

// Store is an implementation of secret store relying on HashiCorp Vault kv-v2 secret engine
type hashicorpSecretStore struct {
	client     vault.HashicorpVaultClient
	mountPoint string
}

// New creates an HashiCorp secret store
func New(client vault.HashicorpVaultClient, mountPoint string) (*hashicorpSecretStore, error) {
	return &hashicorpSecretStore{
		client:     client,
		mountPoint: mountPoint,
	}, nil
}

func (s *hashicorpSecretStore) Info(context.Context) (*models.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

// Set a secret
func (s *hashicorpSecretStore) Set(ctx context.Context, id, value string, attr *models.Attributes) (*models.Secret, error) {
	data := map[string]interface{}{
		valueLabel:          value,
		expirationDateLabel: attr.ExpireAt.UTC().Format(time.UnixDate),
		tagsLabel:           attr.Tags,
		enabledLabel:        attr.Enabled,
	}

	secret, err := s.client.Write(s.pathData(id), data)
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(secret), err
}

// Get a secret
func (s *hashicorpSecretStore) Get(ctx context.Context, id, version string) (*models.Secret, error) {
	// Get latest version by default if no version is specified, otherwise get specific version
	pathData := s.pathData(id)
	if version != "" {
		pathData = fmt.Sprintf("%v?version=%v", pathData, version)
	}

	secret, err := s.client.Read(pathData)
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(secret), err
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

	secrets := res.Data["keys"].([]interface{})
	ids := make([]string, len(secrets))
	for i, elem := range secrets {
		ids[i] = fmt.Sprintf("%v", elem)
	}

	return ids, nil
}

// Update a secret
func (s *hashicorpSecretStore) Update(ctx context.Context, id string, newValue string, attr *models.Attributes) (*models.Secret, error) {
	// Update simply overrides a secret
	return s.Set(ctx, id, newValue, attr)
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
	return path.Join(s.mountPoint, "data", id)
}

func (s *hashicorpSecretStore) pathMetadata(id string) string {
	return path.Join(s.mountPoint, "metadata", id)
}

func formatHashicorpSecret(secret *hashicorp.Secret) *models.Secret {
	for k, _ := range secret.Data {
		fmt.Println(json.MarshalIndent(secret.Data[k], "", "  "))
	}

	return &models.Secret{
		Value:    secret.Data[valueLabel].(string),
		Attr:     nil,
		Metadata: nil,
	}
}
