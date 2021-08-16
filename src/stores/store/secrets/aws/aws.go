package aws

import (
	"context"
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
)

const (
	maxTagsAllowed = 50
)

type Store struct {
	client aws.SecretsManagerClient
	logger log.Logger
}

var _ stores.SecretStore = &Store{}

func New(client aws.SecretsManagerClient, logger log.Logger) *Store {
	return &Store{
		client: client,
		logger: logger,
	}
}

func (s *Store) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := s.logger.With("id", id)

	_, err := s.client.CreateSecret(ctx, id, value)
	if err != nil && errors.IsAlreadyExistsError(err) {
		_, err1 := s.client.PutSecretValue(ctx, id, value)
		if err1 != nil {
			errMessage := "failed to set existing AWS secret"
			logger.WithError(err).Error(errMessage)
			return nil, errors.FromError(err).SetMessage(errMessage)
		}
	} else if err != nil {
		errMessage := "failed to create AWS secret"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	// Tag secret resource when tags found
	if len(attr.Tags) > 0 {
		// check overall len must be limited to max according to doc
		if len(attr.Tags) > maxTagsAllowed {
			errMessage := fmt.Sprintf("resource may not be tagged with more than %d items", maxTagsAllowed)
			logger.WithError(err).Error(errMessage)
			return nil, errors.InvalidParameterError(errMessage)
		}

		_, err = s.client.TagSecretResource(ctx, id, attr.Tags)
		if err != nil {
			errMessage := "failed to set AWS secret tags"
			logger.WithError(err).Error(errMessage)
			return nil, errors.FromError(err).SetMessage(errMessage)
		}
	}

	tags, metadata, err := s.client.DescribeSecret(ctx, id)
	if err != nil {
		errMessage := "failed to get AWS secret after creation"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return formatAwsSecret(id, value, tags, metadata), nil
}

func (s *Store) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := s.logger.With("id", id)

	getSecretOutput, err := s.client.GetSecret(ctx, id, version)
	if err != nil {
		errMessage := "failed to get AWS secret"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	tags, metadata, err := s.client.DescribeSecret(ctx, id)
	if err != nil {
		errMessage := "failed to get AWS secret description"
		logger.WithError(err).Error(errMessage)
		return nil, errors.FromError(err).SetMessage(errMessage)
	}

	return formatAwsSecret(id, *getSecretOutput.SecretString, tags, metadata), nil
}

func (s *Store) List(ctx context.Context) ([]string, error) {
	result := []string{}
	nextToken := ""

	// Loop until the entire list is constituted
	for {
		ret, retToken, err := s.listPaginated(ctx, 0, nextToken)
		if err != nil {
			return nil, err
		}

		result = append(result, ret...)
		if retToken == nil {
			break
		}

		nextToken = *retToken
	}

	return result, nil
}

func (s *Store) Delete(ctx context.Context, id, _ string) error {
	_, err := s.client.DeleteSecret(ctx, id)
	if err != nil {
		errMessage := "failed to delete AWS secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) GetDeleted(_ context.Context, _, _ string) (*entities.Secret, error) {
	err := errors.NotSupportedError("get deleted secret is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) ListDeleted(_ context.Context) ([]string, error) {
	err := errors.NotSupportedError("list deleted secret is not supported")
	s.logger.Warn(err.Error())
	return nil, err
}

func (s *Store) Restore(ctx context.Context, id, version string) error {
	_, err := s.client.RestoreSecret(ctx, id)
	if err != nil {
		errMessage := "failed to restore AWS secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) Destroy(ctx context.Context, id, _ string) error {
	_, err := s.client.DestroySecret(ctx, id)
	if err != nil {
		errMessage := "failed to permanently delete AWS secret"
		s.logger.With("id", id).WithError(err).Error(errMessage)
		return errors.FromError(err).SetMessage(errMessage)
	}

	return nil
}

func (s *Store) listPaginated(ctx context.Context, maxResults int64, nextToken string) (resList []string, resNextToken *string, err error) {
	listOutput, err := s.client.ListSecrets(ctx, maxResults, nextToken)
	if err != nil {
		errMessage := "failed to list AWS secrets"
		s.logger.WithError(err).Error(errMessage)
		return nil, nil, errors.FromError(err).SetMessage(errMessage)
	}

	// return only a list of secret names (IDs)
	secretNamesList := []string{}
	for _, secret := range listOutput.SecretList {
		secretNamesList = append(secretNamesList, *secret.Name)
	}

	return secretNamesList, listOutput.NextToken, nil
}
