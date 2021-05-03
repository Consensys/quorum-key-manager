package hashicorp

import (
	"context"
	"path"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp"
	hashicorpclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

const (
	urlPath         = "keys"
	idLabel         = "id"
	curveLabel      = "curve"
	algorithmLabel  = "algorithm"
	tagsLabel       = "tags"
	publicKeyLabel  = "public_key"
	privateKeyLabel = "privateKey"
	dataLabel       = "data"
	signatureLabel  = "signature"
	versionLabel    = "version"
	createdAtLabel  = "created_at"
	updatedAtLabel  = "updated_at"
)

// Store is an implementation of key store relying on Hashicorp Vault ConsenSys secret engine
type Store struct {
	client     hashicorp.VaultClient
	mountPoint string
}

var _ keys.Store = &Store{}

// New creates an HashiCorp key store
func New(client hashicorp.VaultClient, mountPoint string) *Store {
	return &Store{
		client:     client,
		mountPoint: mountPoint,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

// Create a key
func (s *Store) Create(_ context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	res, err := s.client.Write(s.pathKeys(""), map[string]interface{}{
		idLabel:        id,
		curveLabel:     alg.EllipticCurve,
		algorithmLabel: alg.Type,
		tagsLabel:      attr.Tags,
	})
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	return parseResponse(res), nil
}

// Import a key
func (s *Store) Import(_ context.Context, id, privKey string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	res, err := s.client.Write(s.pathKeys("import"), map[string]interface{}{
		idLabel:         id,
		curveLabel:      alg.EllipticCurve,
		algorithmLabel:  alg.Type,
		tagsLabel:       attr.Tags,
		privateKeyLabel: privKey,
	})
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	return parseResponse(res), nil
}

// Get a key
func (s *Store) Get(_ context.Context, id, version string) (*entities.Key, error) {
	// TODO: Versioning is not yet implemented on the plugin
	if version != "" {
		return nil, errors.NotImplementedError
	}

	res, err := s.client.Read(s.pathKeys(id), nil)
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	if res.Data["error"] != nil {
		return nil, errors.NotFoundError("could not find key pair")
	}

	return parseResponse(res), nil
}

// Get all key ids
func (s *Store) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathKeys(""))
	if err != nil {
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	keys, ok := res.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	var ids []string
	for _, id := range keys {
		ids = append(ids, id.(string))
	}

	return ids, nil
}

// Update key tags
func (s *Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

// Refresh key (create new identical version with different TTL)
func (s *Store) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	return errors.NotImplementedError
}

// Delete a key
func (s *Store) Delete(_ context.Context, id string) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

// Gets a deleted key
func (s *Store) GetDeleted(_ context.Context, id string) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

// Lists all deleted keys
func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.NotImplementedError
}

// Undelete a previously deleted key
func (s *Store) Undelete(ctx context.Context, id string) error {
	return errors.NotImplementedError
}

// Destroy a key permanently
func (s *Store) Destroy(ctx context.Context, id string) error {
	return errors.NotImplementedError
}

// Sign any arbitrary data
func (s *Store) Sign(_ context.Context, id, data, version string) (string, error) {
	// TODO: Versioning is not yet implemented on the plugin
	if version != "" {
		return "", errors.NotImplementedError
	}

	res, err := s.client.Write(path.Join(s.pathKeys(id), "sign"), map[string]interface{}{
		dataLabel: data,
	})
	if err != nil {
		return "", hashicorpclient.ParseErrorResponse(err)
	}

	return res.Data[signatureLabel].(string), nil
}

// Encrypt any arbitrary data using a specified key
func (s *Store) Encrypt(ctx context.Context, id, version, data string) (string, error) {
	return "", errors.NotImplementedError

}

// Decrypt a single block of encrypted data.
func (s *Store) Decrypt(ctx context.Context, id, version, data string) (string, error) {
	return "", errors.NotImplementedError
}

func (s *Store) pathKeys(suffix string) string {
	return path.Join(s.mountPoint, urlPath, suffix)
}
