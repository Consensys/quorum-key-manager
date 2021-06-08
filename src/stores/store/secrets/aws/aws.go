package aws

import (
	"context"
	aws2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/infra/aws"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	sdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	CurrentVersionMark = "AWSCURRENT"
	maxTagsAllowed     = 50
)

// Store is an implementation of secret store relying on AWS secretsmanager
type SecretStore struct {
	client aws2.SecretsManagerClient
	logger *log.Logger
}

// New creates an AWS secret store
func New(client aws2.SecretsManagerClient, logger *log.Logger) *SecretStore {
	return &SecretStore{
		client: client,
		logger: logger,
	}
}

func (s *SecretStore) Info(context.Context) (*entities2.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

// Set Set a secret and tag it when tags exist
func (s *SecretStore) Set(ctx context.Context, id, value string, attr *entities2.Attributes) (*entities2.Secret, error) {
	logger := s.logger.WithField("id", id)

	_, err := s.client.CreateSecret(ctx, id, value)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceExistsException:
				_, err1 := s.client.PutSecretValue(ctx, id, value)
				if err1 != nil {
					logger.Error("failed to update secret")
					return nil, translateAwsError(err1)
				}
			default:
				logger.Error("failed to create secret")
				return nil, translateAwsError(err)
			}
		} else {
			logger.Error("failed to create secret")
			return nil, translateAwsError(err)
		}
	}

	// Tag secret resource when tags found
	if len(attr.Tags) > 0 {
		// check overall len must be limited to max according to doc
		if len(attr.Tags) > maxTagsAllowed {
			return nil, errors.InvalidParameterError("resource may not be tagged with more than %d items", maxTagsAllowed)
		}

		_, err = s.client.TagSecretResource(ctx, id, attr.Tags)
		if err != nil {
			logger.Error("failed to tag secret")
			return nil, translateAwsError(err)
		}
	}

	tags := make(map[string]string)
	metadata := &entities2.Metadata{}

	describeOutput, err := s.client.DescribeSecret(ctx, id)

	if err != nil {
		logger.Error("failed to describe secret")
		return nil, translateAwsError(err)
	}

	if describeOutput != nil {
		// Trick to help us getting the actual current version as there is no versionID metadata
		currentVersion := ""
		for version, stages := range describeOutput.VersionIdsToStages {
			for _, stage := range stages {
				if *stage == CurrentVersionMark {
					currentVersion = version
				}
			}
		}

		metadata = &entities2.Metadata{
			Version:   currentVersion,
			CreatedAt: sdk.TimeValue(describeOutput.CreatedDate),
			UpdatedAt: sdk.TimeValue(describeOutput.LastChangedDate),
			DeletedAt: sdk.TimeValue(describeOutput.DeletedDate),
		}

		for _, outTag := range describeOutput.Tags {
			tags[*outTag.Key] = *outTag.Value
		}

	}
	logger.Info("secret set successfully")
	return formatAwsSecret(id, value, tags, metadata), nil
}

// Get Gets a secret and its description
func (s *SecretStore) Get(ctx context.Context, id, version string) (*entities2.Secret, error) {
	logger := s.logger.WithField("id", id)

	getSecretOutput, err := s.client.GetSecret(ctx, id, version)
	if err != nil {
		logger.Error("secret not found")
		return nil, errors.NotFoundError("secret not found")
	}

	// Prepare to get tags and metadata via description
	tags := make(map[string]string)
	metadata := &entities2.Metadata{}

	describeOutput, err := s.client.DescribeSecret(ctx, id)

	if err != nil {
		logger.Error("failed to describe secret")
		return nil, translateAwsError(err)
	}

	if describeOutput != nil {
		metadata = &entities2.Metadata{
			Version:   *getSecretOutput.VersionId,
			CreatedAt: sdk.TimeValue(describeOutput.CreatedDate),
			UpdatedAt: sdk.TimeValue(describeOutput.LastChangedDate),
			DeletedAt: sdk.TimeValue(describeOutput.DeletedDate),
		}

		for _, outTag := range describeOutput.Tags {
			tags[*outTag.Key] = *outTag.Value
		}
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
		ret, retToken, err := s.ListPaginated(ctx, 0, nextToken)
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
func (s *SecretStore) ListPaginated(ctx context.Context, maxResults int64, nextToken string) (resList []string, resNextToken *string, err error) {

	listOutput, err := s.client.ListSecrets(ctx, maxResults, nextToken)
	if err != nil {
		s.logger.Error("failed to list secrets")
		return nil, nil, translateAwsError(err)
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
	destroy := false
	_, err := s.client.DeleteSecret(ctx, id, destroy)
	if err != nil {
		logger.Error("failed to delete secret")
		return translateAwsError(err)
	}

	logger.Info("secret was deleted successfully")
	return nil
}

// GetDeleted Gets a deleted secret
func (s *SecretStore) GetDeleted(_ context.Context, id string) (*entities2.Secret, error) {
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
		logger.Error("failed to restore secret")
		return translateAwsError(err)
	}
	logger.Info("secret has been restored successfully")
	return nil
}

// Destroy Deletes a secret permanently (force deletion, secret will be unrecoverable)
func (s *SecretStore) Destroy(ctx context.Context, id string) error {
	logger := s.logger.WithField("id", id)
	destroy := true

	_, err := s.client.DeleteSecret(ctx, id, destroy)
	if err != nil {
		logger.Error("failed to destroy secret")
		return translateAwsError(err)
	}
	logger.Info("secret has been destroyed successfully")
	return nil
}

func translateAwsError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case secretsmanager.ErrCodeResourceExistsException:
			return errors.AlreadyExistsError("resource already exists")
		case secretsmanager.ErrCodeInternalServiceError:
			return errors.InternalError("internal error")
		case secretsmanager.ErrCodeInvalidParameterException:
			return errors.InvalidParameterError("invalid parameter")
		case secretsmanager.ErrCodeInvalidRequestException:
			return errors.InvalidRequestError("invalid request")
		case secretsmanager.ErrCodeResourceNotFoundException:
			return errors.NotFoundError("resource was not found")
		case secretsmanager.ErrCodeInvalidNextTokenException:
			return errors.InvalidParameterError("invalid parameter, next token")
		case secretsmanager.ErrCodeLimitExceededException:
			return errors.InternalError("internal error, limit exceeded")
		case secretsmanager.ErrCodePreconditionNotMetException:
			return errors.InternalError("internal error, preconditions not met")
		case secretsmanager.ErrCodeEncryptionFailure:
			return errors.InternalError("internal error, encryption failed")
		case secretsmanager.ErrCodeDecryptionFailure:
			return errors.InternalError("internal error, decryption failed")
		case secretsmanager.ErrCodeMalformedPolicyDocumentException:
			return errors.InvalidParameterError("invalid policy documentation parameter")

		}
	}
	return err
}
