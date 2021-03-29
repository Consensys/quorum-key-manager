package hashicorp

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys"
	"path"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

const (
	endpoint        = "keys"
	idLabel         = "id"
	curveLabel      = "curve"
	algorithmLabel  = "algorithm"
	tagsLabel       = "tags"
	privateKeyLabel = "privateKey"
	dataLabel       = "data"
	signatureLabel  = "signatureLabel"
)

// Store is an implementation of key store relying on Hashicorp Vault ConsenSys secret engine
type hashicorpKeyStore struct {
	client     hashicorp.VaultClient
	mountPoint string
}

// New creates an HashiCorp key store
func New(client hashicorp.VaultClient, mountPoint string) keys.Store {
	return &hashicorpKeyStore{
		client:     client,
		mountPoint: mountPoint,
	}
}

func (s *hashicorpKeyStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.NotImplementedError
}

// Create a key
func (s *hashicorpKeyStore) Create(_ context.Context, id string, alg *entities.Algo, attr *entities.Attributes) (*entities.Key, error) {
	key := &entities.Key{}
	res, err := s.client.Write(s.pathKeys(""), map[string]interface{}{
		idLabel:        id,
		curveLabel:     alg.EllipticCurve,
		algorithmLabel: alg.Type,
		tagsLabel:      attr.Tags,
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	err = parseResponse(res.Data, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Import a key
func (s *hashicorpKeyStore) Import(_ context.Context, id string, privKey string, alg *entities.Algo, attr *entities.Attributes) (*entities.Key, error) {
	key := &entities.Key{}
	res, err := s.client.Write(path.Join(s.mountPoint, endpoint, "import"), map[string]interface{}{
		idLabel:         id,
		curveLabel:      alg.EllipticCurve,
		algorithmLabel:  alg.Type,
		tagsLabel:       attr.Tags,
		privateKeyLabel: privKey,
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	err = parseResponse(res.Data, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Get a key
func (s *hashicorpKeyStore) Get(_ context.Context, id string, version int) (*entities.Key, error) {
	key := &entities.Key{}

	// TODO: Versioning is not yet implemented on the plugin
	if version != 0 {
		return nil, errors.NotImplementedError
	}

	res, err := s.client.Read(s.pathKeys(id))
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	err = parseResponse(res.Data, key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// Get all secret ids
func (s *hashicorpKeyStore) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathKeys(""))
	if err != nil {
		return []string{}, parseErrorResponse(err)
	}

	ids, ok := res.Data["keys"].([]string)
	if !ok {
		return []string{}, nil
	}

	return ids, nil
}

// Update key tags
func (s *hashicorpKeyStore) Update(ctx context.Context, id string, tags map[string]string) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

// Refresh key (create new identical version with different TTL)
func (s *hashicorpKeyStore) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	return errors.NotImplementedError
}

// Delete a secret
func (s *hashicorpKeyStore) Delete(_ context.Context, id string, versions ...int) (*entities.Key, error) {
	return nil, errors.NotImplementedError
}

// Gets a deleted secret
func (s *hashicorpKeyStore) GetDeleted(_ context.Context, id string) (*entities.Key, error) {
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

func (s *hashicorpKeyStore) Sign(ctx context.Context, id string, data string, version int) (string, error) {
	// TODO: Versioning is not yet implemented on the plugin
	if version != 0 {
		return "", errors.NotImplementedError
	}

	res, err := s.client.Write(path.Join(s.pathKeys(id), "sign"), map[string]interface{}{
		dataLabel: data,
	})
	if err != nil {
		return "", parseErrorResponse(err)
	}

	return res.Data[signatureLabel].(string), nil
}

// Encrypt any arbitrary data using a specified key
func (s *hashicorpKeyStore) Encrypt(ctx context.Context, id string, data string) (string, error) {
	return "", errors.NotImplementedError

}

// Decrypt a single block of encrypted data.
func (s *hashicorpKeyStore) Decrypt(ctx context.Context, id string, data string) (string, error) {
	return "", errors.NotImplementedError
}

func (s *hashicorpKeyStore) pathKeys(id string) string {
	return path.Join(s.mountPoint, endpoint, id)
}
