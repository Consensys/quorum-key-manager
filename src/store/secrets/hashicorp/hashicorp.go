package hashicorp

import (
	"context"
	"path"
	"strconv"
	"time"

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

	metadata, err := formatHashicorpSecretData(hashicorpSecret.Data)
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(id, value, attr.Tags, metadata), nil
}

// Get a secret
func (s *SecretStore) Get(_ context.Context, id, version string) (*entities.Secret, error) {
	var callData map[string][]string
	if version != "" {
		_, err := strconv.Atoi(version)
		if err != nil {
			return nil, errors.InvalidParameterError("version must be a number")
		}

		callData = map[string][]string{
			versionLabel: {version},
		}
	}

	hashicorpSecretData, err := s.client.Read(s.pathData(id), callData)
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	} else if hashicorpSecretData == nil {
		return nil, errors.NotFoundError("secret not found")
	}

	data := hashicorpSecretData.Data[dataLabel].(map[string]interface{})
	value := data[valueLabel].(string)
	tags := make(map[string]string)
	if data[tagsLabel] != nil {
		tags = data[tagsLabel].(map[string]string)
	}

	// We need to do a second call to get the metadata
	hashicorpSecretMetadata, err := s.client.Read(s.pathMetadata(id), nil)
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	metadata, err := formatHashicorpSecretMetadata(hashicorpSecretMetadata, version)
	if err != nil {
		return nil, err
	}

	return formatHashicorpSecret(id, value, tags, metadata), nil
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
