package aws

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

const (
	maxTagsAllowed = 50
)

type Store struct {
	client aws.SecretsManagerClient
	logger log.Logger
}

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
	_, err := s.client.CreateSecret(ctx, id, value)
	if err != nil && errors.IsAlreadyExistsError(err) {
		_, err1 := s.client.PutSecretValue(ctx, id, value)
		if err1 != nil {
			return nil, err1
		}
	} else if err != nil {
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
			return nil, err
		}
	}
	tags, metadata, err := s.client.DescribeSecret(ctx, id)
	if err != nil {
		return nil, err
	}

	return formatAwsSecret(id, value, tags, metadata), nil
}

func (s *Store) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	getSecretOutput, err := s.client.GetSecret(ctx, id, version)
	if err != nil {
		return nil, err
	}

	tags, metadata, err := s.client.DescribeSecret(ctx, id)
	if err != nil {
		return nil, err
	}

	return formatAwsSecret(id, *getSecretOutput.SecretString, tags, metadata), nil
}

func (s *Store) List(ctx context.Context) ([]string, error) {
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

	return secrets, nil
}

func (s *Store) listPaginated(ctx context.Context, maxResults int64, nextToken string) (resList []string, resNextToken *string, err error) {
	listOutput, err := s.client.ListSecrets(ctx, maxResults, nextToken)
	if err != nil {
		return nil, nil, err
	}

	// return only a list of secret names (IDs)
	secretNamesList := []string{}
	for _, secret := range listOutput.SecretList {
		secretNamesList = append(secretNamesList, *secret.Name)
	}

	return secretNamesList, listOutput.NextToken, nil
}

func (s *Store) Delete(ctx context.Context, id string) error {
	_, err := s.client.DeleteSecret(ctx, id, false)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetDeleted(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Store) Undelete(ctx context.Context, id string) error {
	_, err := s.client.RestoreSecret(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Destroy(ctx context.Context, id string) error {
	_, err := s.client.DeleteSecret(ctx, id, true)
	if err != nil {
		return err
	}

	return nil
}
