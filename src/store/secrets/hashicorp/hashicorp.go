package hashicorp

import (
	"context"
	"path"
	"time"

	"github.com/hashicorp/vault/api"

	hashicorpclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

const (
	dataLabel        = "data"
	metadataLabel    = "metadata"
	valueLabel       = "value"
	deleteAfterLabel = "delete_version_after"
	tagsLabel        = "tags"
	versionLabel     = "version"
)

// Store is an implementation of secret store relying on Hashicorp Vault kv-v2 secret engine
type SecretStore struct {
	client     hashicorp.VaultClient
	mountPoint string
}

// New creates an HashiCorp secret store
func New(client hashicorp.VaultClient, mountPoint string) *SecretStore {
	return &SecretStore{
		client:     client,
		mountPoint: mountPoint,
	}
}

func (s *SecretStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

// Set a secret
func (s *SecretStore) Set(_ context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	data := map[string]interface{}{
		dataLabel: map[string]interface{}{
			valueLabel: value,
			tagsLabel:  attr.Tags,
		},
	}

	hashicorpSecret, err := s.client.Write(s.pathData(id), data)
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	// Hashicorp only returns metadata as the "data" field when creating a new secret
	metadata, err := extractMetadata(hashicorpSecret.Data)
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(id, value, attr.Tags, metadata), nil
}

// Get a secret
func (s *SecretStore) Get(_ context.Context, id, version string) (*entities.Secret, error) {
	var hashicorpSecret *api.Secret
	var err error
	if version != "" {
		hashicorpSecret, err = s.client.ReadWithData(s.pathData(id), map[string][]string{
			versionLabel: {version},
		})
	} else {
		hashicorpSecret, err = s.client.Read(s.pathData(id))
	}
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	if hashicorpSecret == nil {
		return nil, errors.NotFoundError("secret not found")
	}

	// Hashicorp returns metadata in "data.metadata" and data in "data.data" fields
	data := hashicorpSecret.Data[dataLabel].(map[string]interface{})
	metadata, err := extractMetadata(hashicorpSecret.Data[metadataLabel].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	tags := make(map[string]string)
	if data[tagsLabel] != nil {
		tags = data[tagsLabel].(map[string]string)
	}

	return formatHashicorpSecret(id, data[valueLabel].(string), tags, metadata), nil
}

// Get all secret ids
func (s *SecretStore) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathMetadata(""))
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	if res == nil {
		return []string{}, nil
	}

	keysInterface := res.Data["keys"].([]interface{})
	keysStr := make([]string, len(keysInterface))
	for i, key := range keysInterface {
		keysStr[i] = key.(string)
	}

	return keysStr, nil
}

// Refresh an existing secret by extending its TTL
func (s *SecretStore) Refresh(_ context.Context, id, _ string, expirationDate time.Time) error {
	data := make(map[string]interface{})
	if !expirationDate.IsZero() {
		data[deleteAfterLabel] = time.Until(expirationDate).String()
	}

	_, err := s.client.Write(s.pathMetadata(id), data)
	if err != nil {
		return hashicorpclient.ParseErrorResponse(err)
	}

	return nil
}

// Delete a secret
func (s *SecretStore) Delete(_ context.Context, id string, versions ...string) (*entities.Secret, error) {
	return nil, errors.NotImplementedError
}

// Gets a deleted secret
func (s *SecretStore) GetDeleted(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.NotImplementedError
}

// Lists all deleted secrets
func (s *SecretStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.NotImplementedError
}

// Undelete a previously deleted secret
func (s *SecretStore) Undelete(ctx context.Context, id string) error {
	return errors.NotImplementedError
}

// Destroy a secret permanently
func (s *SecretStore) Destroy(ctx context.Context, id string, versions ...string) error {
	return errors.NotImplementedError
}

// path compute path from hashicorp mount
func (s *SecretStore) pathData(id string) string {
	return path.Join(s.mountPoint, dataLabel, id)
}

func (s *SecretStore) pathMetadata(id string) string {
	return path.Join(s.mountPoint, metadataLabel, id)
}
