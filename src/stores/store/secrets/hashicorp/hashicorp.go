package hashicorp

import (
	"context"
	"path"
	"strconv"

	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/secrets"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
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

func (s *Store) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := s.logger.With("id", id)

	_, err := s.client.Write(s.pathData(id), map[string]interface{}{
		dataLabel: map[string]interface{}{
			valueLabel: value,
			tagsLabel:  attr.Tags,
		},
	})
	if err != nil {
		errMessage := "failed to create Hashicorp secret"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return s.Get(ctx, id, "")
}

func (s *Store) Get(_ context.Context, id, version string) (*entities.Secret, error) {
	logger := s.logger.With("id", id, "version", version)

	var callData map[string][]string
	if version != "" {
		_, err := strconv.Atoi(version)
		if err != nil {
			errMessage := "version must be a number"
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		callData = map[string][]string{
			versionLabel: {version},
		}
	}

	hashicorpSecretData, err := s.client.Read(s.pathData(id), callData)
	if err != nil {
		errMessage := "failed to get Hashicorp secret data"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	} else if hashicorpSecretData == nil {
		errMessage := "Hashicorp secret not found"
		logger.WithError(err).Error(errMessage)
		return nil, errors.NotFoundError(errMessage)
	}

	data := hashicorpSecretData.Data[dataLabel].(map[string]interface{})
	value := data[valueLabel].(string)

	// We need to do a second call to get the metadata
	hashicorpSecretMetadata, err := s.client.Read(s.pathMetadata(id), nil)
	if err != nil {
		errMessage := "failed to get Hashicorp secret metadata"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	metadata, err := formatHashicorpSecretMetadata(hashicorpSecretMetadata, version)
	if err != nil {
		errMessage := "failed to parse Hashicorp secret"
		logger.WithError(err).Error(errMessage)
		return nil, errors.HashicorpVaultError(errMessage)
	}

	var tags map[string]string
	if data[tagsLabel] != nil {
		tags = formatTags(data[tagsLabel].(map[string]interface{}))
	}

	return formatHashicorpSecret(id, value, tags, metadata), nil
}

func (s *Store) List(_ context.Context) ([]string, error) {
	res, err := s.client.List(s.pathMetadata(""))
	if err != nil {
		errMessage := "failed to list Hashicorp secrets"
		s.logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
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
