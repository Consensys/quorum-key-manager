package hashicorp

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	hashicorpclient "github.com/ConsenSysQuorum/quorum-key-manager/src/infra/hashicorp/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets"

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
type Store struct {
	client     hashicorp.VaultClient
	mountPoint string
	logger     *log.Logger
}

var _ secrets.Store = &Store{}

// New creates an Hashicorp secret store
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

// Set a secret
func (s *Store) Set(_ context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := s.logger.WithField("id", id)
	errMsg := "failed to set secret"
	data := map[string]interface{}{
		dataLabel: map[string]interface{}{
			valueLabel: value,
			tagsLabel:  attr.Tags,
		},
	}

	hashicorpSecret, err := s.client.Write(s.pathData(id), data)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	metadata, err := formatHashicorpSecretData(hashicorpSecret.Data)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return nil, err
	}

	s.logger.Info("secret was set successfully")
	return formatHashicorpSecret(id, value, attr.Tags, metadata), nil
}

// Get a secret
func (s *Store) Get(_ context.Context, id, version string) (*entities.Secret, error) {
	logger := s.logger.WithField("id", id)
	errMsg := "failed to get secret"
	var callData map[string][]string
	if version != "" {
		_, err := strconv.Atoi(version)
		if err != nil {
			err2 := errors.InvalidParameterError("version must be a number")
			logger.WithError(err2).Error(errMsg)
			return nil, err2
		}

		callData = map[string][]string{
			versionLabel: {version},
		}
	}

	hashicorpSecretData, err := s.client.Read(s.pathData(id), callData)
	if err != nil {
		logger.WithError(err).Error(errMsg)
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
		logger.WithError(err).Error(errMsg)
		return nil, hashicorpclient.ParseErrorResponse(err)
	}

	metadata, err := formatHashicorpSecretMetadata(hashicorpSecretMetadata, version)
	if err != nil {
		logger.WithError(err).Error(errMsg)
		return nil, err
	}

	logger.Debug("secret was retrieved successfully")
	return formatHashicorpSecret(id, value, tags, metadata), nil
}

// Get all secret ids
func (s *Store) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathMetadata(""))
	if err != nil {
		s.logger.WithError(err).Error("failed to list secrets")
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

	s.logger.Debug("secrets were listed successfully")
	return keysStr, nil
}

// Refresh an existing secret by extending its TTL
func (s *Store) Refresh(_ context.Context, id, _ string, expirationDate time.Time) error {
	logger := s.logger.WithField("id", id)

	data := make(map[string]interface{})
	if !expirationDate.IsZero() {
		data[deleteAfterLabel] = time.Until(expirationDate).String()
	}

	_, err := s.client.Write(s.pathMetadata(id), data)
	if err != nil {
		logger.WithError(err).Error("failed to refresh secret")
		return hashicorpclient.ParseErrorResponse(err)
	}

	logger.Info("secret was refreshed successfully")
	return nil
}

// Delete a secret
func (s *Store) Delete(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.ErrNotImplemented
}

// Gets a deleted secret
func (s *Store) GetDeleted(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.ErrNotImplemented
}

// Lists all deleted secrets
func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

// Undelete a previously deleted secret
func (s *Store) Undelete(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

// Destroy a secret permanently
func (s *Store) Destroy(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

// path compute path from hashicorp mount
func (s *Store) pathData(id string) string {
	return path.Join(s.mountPoint, dataLabel, id)
}

func (s *Store) pathMetadata(id string) string {
	return path.Join(s.mountPoint, metadataLabel, id)
}
