package hashicorp

import (
	"context"
	"path"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
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
	logger     *log.Logger
}

var _ keys.Store = &Store{}

// New creates an Hashicorp key store
func New(client hashicorp.VaultClient, mountPoint string, logger *log.Logger) *Store {
	return &Store{
		client:     client,
		mountPoint: mountPoint,
		logger:     logger,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

// Create a key
func (s *Store) Create(_ context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.WithField("id", id)
	res, err := s.client.Write(s.pathKeys(""), map[string]interface{}{
		idLabel:        id,
		curveLabel:     alg.EllipticCurve,
		algorithmLabel: alg.Type,
		tagsLabel:      attr.Tags,
	})
	if err != nil {
		logger.WithError(err).Error("failed to create key")
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	logger.Info("key was created successfully")
	return parseResponse(res), nil
}

// Import a key
func (s *Store) Import(_ context.Context, id, privKey string, alg *entities.Algorithm, attr *entities.Attributes) (*entities.Key, error) {
	logger := s.logger.WithField("id", id)
	res, err := s.client.Write(s.pathKeys("import"), map[string]interface{}{
		idLabel:         id,
		curveLabel:      alg.EllipticCurve,
		algorithmLabel:  alg.Type,
		tagsLabel:       attr.Tags,
		privateKeyLabel: privKey,
	})

	if err != nil {
		logger.WithError(err).Error("failed to import key")
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	logger.Info("key was imported successfully")
	return parseResponse(res), nil
}

// Get a key
func (s *Store) Get(_ context.Context, id string) (*entities.Key, error) {
	logger := s.logger.WithField("id", id)

	res, err := s.client.Read(s.pathKeys(id), nil)
	if err != nil {
		logger.WithError(err).Error("failed to get key")
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	if res.Data["error"] != nil {
		return nil, errors.NotFoundError("could not find key pair")
	}

	logger.Debug("key was retrieved successfully")
	return parseResponse(res), nil
}

// Get all key ids
func (s *Store) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathKeys(""))
	if err != nil {
		s.logger.WithError(err).Error("failed to list keys")
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	keyIds, ok := res.Data["keys"].([]interface{})
	if !ok {
		return []string{}, nil
	}

	var ids []string
	for _, id := range keyIds {
		ids = append(ids, id.(string))
	}

	s.logger.Debug("keys were listed successfully")
	return ids, nil
}

// Update key tags
func (s *Store) Update(ctx context.Context, id string, attr *entities.Attributes) (*entities.Key, error) {
	return nil, errors.ErrNotImplemented
}

// Refresh key (create new identical version with different TTL)
func (s *Store) Refresh(ctx context.Context, id string, expirationDate time.Time) error {
	return errors.ErrNotImplemented
}

// Delete a key
func (s *Store) Delete(_ context.Context, id string) (*entities.Key, error) {
	return nil, errors.ErrNotImplemented
}

// Gets a deleted key
func (s *Store) GetDeleted(_ context.Context, id string) (*entities.Key, error) {
	return nil, errors.ErrNotImplemented
}

// Lists all deleted keys
func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

// Undelete a previously deleted key
func (s *Store) Undelete(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

// Destroy a key permanently
func (s *Store) Destroy(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

// Sign any arbitrary data
func (s *Store) Sign(_ context.Context, id, data string) (string, error) {
	logger := s.logger.WithField("id", id).WithField("data", common.ShortString(data, 5))

	res, err := s.client.Write(path.Join(s.pathKeys(id), "sign"), map[string]interface{}{
		dataLabel: data,
	})
	if err != nil {
		logger.WithError(err).Error("failed to sign data")
		return "", hashicorpclient.ParseErrorResponse(err)
	}

	logger.Debug("data was signed successfully")
	return res.Data[signatureLabel].(string), nil
}

// Encrypt any arbitrary data using a specified key
func (s *Store) Encrypt(ctx context.Context, id, data string) (string, error) {
	return "", errors.ErrNotImplemented

}

// Decrypt a single block of encrypted data.
func (s *Store) Decrypt(ctx context.Context, id, data string) (string, error) {
	return "", errors.ErrNotImplemented
}

func (s *Store) pathKeys(suffix string) string {
	return path.Join(s.mountPoint, urlPath, suffix)
}
