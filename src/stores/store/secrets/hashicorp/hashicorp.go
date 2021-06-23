package hashicorp

import (
	"context"
	"path"
	"strconv"

	"github.com/consensysquorum/quorum-key-manager/pkg/log"

	"github.com/consensysquorum/quorum-key-manager/src/stores/store/secrets"

	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/src/stores/infra/hashicorp"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
)

const (
	dataLabel     = "data"
	metadataLabel = "metadata"
	valueLabel    = "value"
	tagsLabel     = "tags"
	versionLabel  = "version"
)

type Store struct {
	client     hashicorp.VaultClient
	mountPoint string
	logger     log.Logger
}

var _ secrets.Store = &Store{}

func New(client hashicorp.VaultClient, mountPoint string, logger log.Logger) *Store {
	return &Store{
		client:     client,
		mountPoint: mountPoint,
		logger:     logger,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Set(_ context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := s.logger.With("id", id)
	logger.Debug("creating secret")

	data := map[string]interface{}{
		dataLabel: map[string]interface{}{
			valueLabel: value,
			tagsLabel:  attr.Tags,
		},
	}

	hashicorpSecret, err := s.client.Write(s.pathData(id), data)
	if err != nil {
		logger.WithError(err).Error("failed to set secret")
		return nil, err
	}

	metadata, err := formatHashicorpSecretData(hashicorpSecret.Data)
	if err != nil {
		logger.WithError(err).Error("failed to parse secret data")
		return nil, err
	}

	logger.Info("secret set successfully")
	return formatHashicorpSecret(id, value, attr.Tags, metadata), nil
}

func (s *Store) Get(_ context.Context, id, version string) (*entities.Secret, error) {
	logger := s.logger.With("id", id)

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
		logger.WithError(err).Error("failed to get secret data")
		return nil, err
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
		logger.WithError(err).Error("failed to get secret metadata")
		return nil, err
	}

	metadata, err := formatHashicorpSecretMetadata(hashicorpSecretMetadata, version)
	if err != nil {
		logger.WithError(err).Error("failed to format secret metadata")
		return nil, err
	}

	logger.Debug("secret retrieved successfully")
	return formatHashicorpSecret(id, value, tags, metadata), nil
}

func (s *Store) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathMetadata(""))
	if err != nil {
		s.logger.WithError(err).Error("failed to list secrets")
		return nil, err
	}

	if res == nil {
		return []string{}, nil
	}

	keysInterface := res.Data["keys"].([]interface{})
	keysStr := make([]string, len(keysInterface))
	for i, key := range keysInterface {
		keysStr[i] = key.(string)
	}

	s.logger.Debug("secrets listed successfully")
	return keysStr, nil
}

func (s *Store) Delete(_ context.Context, id string) error {
	return errors.ErrNotImplemented
}

func (s *Store) GetDeleted(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	return errors.ErrNotImplemented
}

func (s *Store) pathData(id string) string {
	return path.Join(s.mountPoint, dataLabel, id)
}

func (s *Store) pathMetadata(id string) string {
	return path.Join(s.mountPoint, metadataLabel, id)
}
