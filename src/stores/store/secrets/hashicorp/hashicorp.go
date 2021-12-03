package hashicorp

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/hashicorp"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

const (
	valueLabel   = "value"
	tagsLabel    = "tags"
	versionLabel = "version"
)

type Store struct {
	client hashicorp.Kvv2Client
	db     database.Secrets
	logger log.Logger
}

var _ stores.SecretStore = &Store{}

func New(client hashicorp.Kvv2Client, db database.Secrets, logger log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
		db:     db,
	}
}

func (s *Store) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := s.logger.With("id", id)

	secretItem, err := s.client.SetSecret(id, map[string]interface{}{
		valueLabel: value,
		tagsLabel:  attr.Tags,
	})
	if err != nil {
		errMessage := "failed to create Hashicorp secret"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return s.Get(ctx, id, string(secretItem.Data[versionLabel].(json.Number)))
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

	hashicorpSecretData, err := s.client.ReadData(id, callData)
	if err != nil {
		errMessage := "failed to get Hashicorp secret data"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	} else if hashicorpSecretData == nil {
		errMessage := "Hashicorp secret not found"
		logger.WithError(err).Error(errMessage)
		return nil, errors.NotFoundError(errMessage)
	}

	data := hashicorpSecretData.Data["data"].(map[string]interface{})
	value := data[valueLabel].(string)

	// We need to do a second call to get the metadata
	hashicorpSecretMetadata, err := s.client.ReadMetadata(id)
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

func (s *Store) List(_ context.Context, _, _ uint64) ([]string, error) {
	res, err := s.client.ListSecrets()
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

func (s *Store) Delete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)

	versions, err := s.listVersions(ctx, id, false)
	if err != nil {
		return err
	}
	hashicorpSecretData, err := s.client.ReadData(id, map[string][]string{
		"versions": versions,
	})
	if err != nil {
		errMessage := "failed to get Hashicorp secret data for deletion"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}
	if hashicorpSecretData == nil {
		errMessage := "Hashicorp secret not found for deletion"
		logger.WithError(err).Error(errMessage)
		return errors.NotFoundError(errMessage)
	}

	err = s.client.DeleteSecret(id, map[string][]string{
		"versions": versions,
	})
	if err != nil {
		errMessage := "failed to delete Hashicorp secret"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) GetDeleted(_ context.Context, _ string) (*entities.Secret, error) {
	err := errors.NotSupportedError("get deleted secret is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) ListDeleted(_ context.Context, _, _ uint64) ([]string, error) {
	err := errors.NotSupportedError("list deleted secret is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) Restore(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)

	versions, err := s.listVersions(ctx, id, true)
	if err != nil {
		return err
	}

	err = s.client.RestoreSecret(id, map[string][]string{
		"versions": versions,
	})
	if err != nil {
		errMessage := "failed to restore Hashicorp secret"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)

	versions, err := s.listVersions(ctx, id, true)
	if err != nil {
		return err
	}

	err = s.client.DestroySecret(id, map[string][]string{
		"versions": versions,
	})
	if err != nil {
		errMessage := "failed to destroy Hashicorp secret"
		logger.WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) listVersions(ctx context.Context, id string, isDeleted bool) ([]string, error) {
	versionList, err := s.db.ListVersions(ctx, id, isDeleted)
	if err != nil {
		errMessage := "failed to list secret versions"
		s.logger.WithError(err).Error(errMessage, "id", id)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	if len(versionList) == 0 {
		errMsg := "no versions were found for secret"
		err := errors.NotFoundError("unexpected empty list of secret versions")
		s.logger.WithError(err).Error(errMsg, "id", id)
		return nil, errors.FromError(err).SetMessage(errMsg)
	}

	return versionList, nil
}
