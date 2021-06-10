package aws

import (
	"context"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/aws"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
)

const (
	maxTagsAllowed = 50
)

// SecretStore is an implementation of secret store relying on AWS secretsmanager
type SecretStore struct {
	client aws.SecretsManagerClient
	logger *log.Logger
}

// New creates an AWS secret store
func New(client aws.SecretsManagerClient, logger *log.Logger) *SecretStore {
	return &SecretStore{
		client: client,
		logger: logger,
	}
}

func (s *SecretStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

// Set Set a secret and tag it when tags exist
func (s *SecretStore) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {
	logger := s.logger.WithField("id", id)

	_, err := s.client.CreateSecret(ctx, id, value)

	if err != nil && errors.IsAlreadyExistsError(err) {
		_, err1 := s.client.PutSecretValue(ctx, id, value)
		if err1 != nil {
			logger.WithError(err).Error("failed to update secret")
			return nil, err1
		}
	} else if err != nil && !errors.IsAlreadyExistsError(err) {
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

// Get Gets a secret and its description
func (s *SecretStore) Get(ctx context.Context, id, version string) (*entities.Secret, error) {
	logger := s.logger.WithField("id", id)

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

	logger.Info("secret was retrieved successfully")
	return formatAwsSecret(id, *getSecretOutput.SecretString, tags, metadata), nil
}

// List Gets all secret ids as a slice of names
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
	return secrets, nil
}

// ListPaginated Gets all secret ids as a slice of names
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

// Refresh Updates an existing secret by extending its TTL
func (s *SecretStore) Refresh(_ context.Context, id, _ string, expirationDate time.Time) error {
	return errors.ErrNotImplemented
}

// Delete Deletes a secret
func (s *SecretStore) Delete(ctx context.Context, id string) error {
	logger := s.logger.WithField("id", id)

	_, err := s.client.DeleteSecret(ctx, id, false)
	if err != nil {
		logger.WithError(err).Error("failed to delete secret")
		return err
	}

	logger.Info("secret was deleted successfully")
	return nil
}

// GetDeleted Gets a deleted secret
func (s *SecretStore) GetDeleted(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.ErrNotImplemented
}

// ListDeleted Lists all deleted secrets
func (s *SecretStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

// Undelete Restores a previously deleted secret
func (s *SecretStore) Undelete(ctx context.Context, id string) error {
	logger := s.logger.WithField("id", id)

	_, err := s.client.RestoreSecret(ctx, id)
	if err != nil {
		logger.WithError(err).Error("failed to restore secret")
		return err
	}

	logger.Info("secret has been restored successfully")
	return nil
}

// Destroy Deletes a secret permanently (force deletion, secret will be unrecoverable)
func (s *SecretStore) Destroy(ctx context.Context, id string) error {
	logger := s.logger.WithField("id", id)

	_, err := s.client.DeleteSecret(ctx, id, true)
	if err != nil {
		logger.WithError(err).Error("failed to destroy secret")
		return err
	}

	logger.Info("secret has been destroyed successfully")
	return nil
}
