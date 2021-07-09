package aws

import (
	"context"

	aws2 "github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

const (
	maxTagsAllowed = 50
)

type SecretStore struct {
	client aws2.SecretsManagerClient
	logger log.Logger
}

func New(client aws2.SecretsManagerClient, logger log.Logger) *SecretStore {
	return &SecretStore{
		client: client,
		logger: logger,
	}
}

func (s *SecretStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *SecretStore) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := s.logger.With("id", id)
	logger.Debug("creating secret")

	_, err := s.client.CreateSecret(ctx, id, value)
	if err != nil && errors.IsAlreadyExistsError(err) {
		_, err1 := s.client.PutSecretValue(ctx, id, value)
		if err1 != nil {
			logger.WithError(err).Error("failed to replace secret")
			return nil, err1
		}
	} else if err != nil {
		logger.WithError(err).Error("failed to create aws secret")
		return nil, err
	}

	// Tag secret resource when tags found
	if len(attr.Tags) > 0 {
		// check overall len must be limited to max according to doc
		if len(attr.Tags) > maxTagsAllowed {
			return nil, errors.InvalidParameterError("resource may not be tagged with more than %d items", maxTagsAllowed)
		}

		_, err = s.client.TagSecretResource(ctx, id, attr.Tags)
		if err != nil {
			logger.WithError(err).Error("failed to tag secret")
			return nil, err
		}
	}
	tags, metadata, err := s.client.DescribeSecret(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to describe secret")
		return nil, err
	}

	logger.Info("secret set successfully")
	return formatAwsSecret(id, value, tags, metadata), nil
}

func (s *SecretStore) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := s.logger.With("id", id)

	getSecretOutput, err := s.client.GetSecret(ctx, id, version)
	if err != nil {
		logger.WithError(err).Error("failed to get secret")
		return nil, err
	}

	tags, metadata, err := s.client.DescribeSecret(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to describe secret")
		return nil, err
	}

	logger.Debug("secret retrieved successfully")
	return formatAwsSecret(id, *getSecretOutput.SecretString, tags, metadata), nil
}

func (s *SecretStore) List(ctx context.Context) ([]string, error) {
	secrets := []string{}
	nextToken := ""

	// Loop until the entire list is constituted
	for {
		ret, retToken, err := s.listPaginated(ctx, 0, nextToken)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, ret...)
		if retToken == nil {
			break
		}
		nextToken = *retToken

	}

	s.logger.Debug("secrets listed successfully")
	return secrets, nil
}

func (s *SecretStore) listPaginated(ctx context.Context, maxResults int64, nextToken string) (resList []string, resNextToken *string, err error) {
	listOutput, err := s.client.ListSecrets(ctx, maxResults, nextToken)
	if err != nil {
		s.logger.WithError(err).Error("failed to list secrets")
		return nil, nil, err
	}

	// return only a list of secret names (IDs)
	secretNamesList := []string{}
	for _, secret := range listOutput.SecretList {
		secretNamesList = append(secretNamesList, *secret.Name)
	}

	s.logger.Info("secrets were listed successfully")
	return secretNamesList, listOutput.NextToken, nil
}

func (s *SecretStore) Delete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("deleting secret")

	_, err := s.client.DeleteSecret(ctx, id, false)
	if err != nil {
		logger.WithError(err).Error("failed to delete secret")
		return err
	}

	logger.Info("secret deleted successfully")
	return nil
}

func (s *SecretStore) GetDeleted(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.ErrNotImplemented
}

func (s *SecretStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

func (s *SecretStore) Undelete(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)

	_, err := s.client.RestoreSecret(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to restore secret")
		return err
	}

	logger.Info("secret has been restored successfully")
	return nil
}

func (s *SecretStore) Destroy(ctx context.Context, id string) error {
	logger := s.logger.With("id", id)
	logger.Debug("destroying key")

	_, err := s.client.DeleteSecret(ctx, id, true)
	if err != nil {
		logger.Error("failed to permanently delete secret")
		return err
	}

	logger.Info("secret permanently deleted")
	return nil
}
