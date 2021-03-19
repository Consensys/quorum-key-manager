package hashicorp

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"path"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

const (
	endpoint = "keys"

	dataLabel     = "data"
	metadataLabel = "metadata"
)

// Store is an implementation of key store relying on Hashicorp Vault ConsenSys secret engine
type hashicorpKeyStore struct {
	client     hashicorp.VaultClient
	mountPoint string
}

// New creates an HashiCorp secret store
func New(client hashicorp.VaultClient, mountPoint string) keys.Store {
	return &hashicorpKeyStore{
		client:     client,
		mountPoint: mountPoint,
	}
}

func (s *hashicorpKeyStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

// Set a secret
func (s *hashicorpKeyStore) Create(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Key, error) {
	key := &entities.Key{}
	res, err := s.client.Write(path.Join(s.mountPoint, endpoint), map[string]interface{}{

	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	if res == nil || res.Data == nil {
		return nil, nil
	}

	err = parseResponse(res.Data, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Get a secret
func (s *hashicorpKeyStore) Get(_ context.Context, id string, version int) (*entities.Secret, error) {
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
func (s *hashicorpKeyStore) List(_ context.Context) ([]string, error) {
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
func (s *hashicorpKeyStore) Refresh(_ context.Context, id string, expirationDate time.Time) error {
	data := make(map[string]interface{})
	if !expirationDate.IsZero() {
		data[deleteAfterLabel] = time.Until(expirationDate).String()
	}

	_, err := s.client.Write(s.pathMetadata(id), data)
	if err != nil {
		return err
	}

	return nil
}

// Delete a secret
func (s *hashicorpKeyStore) Delete(_ context.Context, id string, versions ...int) (*entities.Secret, error) {
	return nil, errors.NotImplementedError
}

// Gets a deleted secret
func (s *hashicorpKeyStore) GetDeleted(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.NotImplementedError
}

// Lists all deleted secrets
func (s *hashicorpKeyStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.NotImplementedError
}

// Undelete a previously deleted secret
func (s *hashicorpKeyStore) Undelete(ctx context.Context, id string) error {
	return errors.NotImplementedError
}

// Destroy a secret permanently
func (s *hashicorpKeyStore) Destroy(ctx context.Context, id string, versions ...int) error {
	return errors.NotImplementedError
}

// path compute path from hashicorp mount
func (s *hashicorpKeyStore) pathData(id string) string {
	return path.Join(s.mountPoint, dataLabel, id)
}