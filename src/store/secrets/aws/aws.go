package aws

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/infra/aws"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	sdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	dataLabel        = "data"
	metadataLabel    = "metadata"
	valueLabel       = "value"
	deleteAfterLabel = "delete_version_after"
	tagsLabel        = "tags"
	versionLabel     = "version"
	maxTagsAllowed   = 50
)

// Store is an implementation of secret store relying on Hashicorp Vault kv-v2 secret engine
type SecretStore struct {
	client aws.SecretsManagerClient
}

// New creates an HashiCorp secret store
func New(client aws.SecretsManagerClient) *SecretStore {
	return &SecretStore{
		client: client,
	}
}

func (s *SecretStore) Info(context.Context) (*entities.StoreInfo, error) {
	return nil, errors.ErrNotImplemented
}

//Set Set a secret and tag it when tags exist
func (s *SecretStore) Set(ctx context.Context, id, value string, attr *entities.Attributes) (*entities.Secret, error) {

	createSecretInput := &secretsmanager.CreateSecretInput{
		SecretString: &value,
		Name:         &id,
	}

	_, err := s.client.CreateSecret(ctx, createSecretInput)
	if err != nil {
		//TODO parse aws flavored errors
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case secretsmanager.ErrCodeResourceExistsException:
				putSecretInput := &secretsmanager.PutSecretValueInput{
					SecretId:     &id,
					SecretString: &value,
				}
				_, err1 := s.client.PutSecretValue(ctx, putSecretInput)
				if err1 != nil {
					log.Printf("PutSecretValue error for secret %v : %v", id, err1)
					return nil, err1
				}
			default:
				return nil, err
			}
		}

	}

	//Tag secret resource when tags found
	if len(attr.Tags) > 0 {
		//check overall len must be limited to max according to doc
		if len(attr.Tags) > maxTagsAllowed {
			return nil, fmt.Errorf("resource may not be tagged with more than %d items", maxTagsAllowed)
		}
		inputTags := []*secretsmanager.Tag{}

		for key, value := range attr.Tags {
			k, v := key, value
			var in secretsmanager.Tag = secretsmanager.Tag{
				Key:   &k,
				Value: &v,
			}
			inputTags = append(inputTags, &in)
		}

		tagResourceInput := &secretsmanager.TagResourceInput{
			SecretId: &id,
			Tags:     inputTags,
		}

		_, err = s.client.TagSecretResource(ctx, tagResourceInput)
		if err != nil {
			//TODO parse aws flavored errors
			log.Printf("Tagging error : %s", err.Error())
		}
	}

	//Now retrieve resource description for metadata
	describeInput := &secretsmanager.DescribeSecretInput{
		SecretId: &id,
	}

	tags := make(map[string]string)
	metadata := &entities.Metadata{}

	describeOutput, err := s.client.DescribeSecret(ctx, describeInput)

	if err != nil {
		log.Printf("Description error for secret %v", describeInput.SecretId)
	}

	if err == nil && describeOutput != nil {
		log.Printf("Description for secret %v", describeOutput)
		//Trick to help us getting the actual current version as there is no versionID metadata
		currentVersion := ""
		for version, stages := range describeOutput.VersionIdsToStages {
			for _, stage := range stages {
				if *stage == "AWSCURRENT" {
					currentVersion = version
				}
			}
		}

		metadata = &entities.Metadata{
			Version:   currentVersion,
			CreatedAt: sdk.TimeValue(describeOutput.CreatedDate),
			UpdatedAt: sdk.TimeValue(describeOutput.LastChangedDate),
			DeletedAt: sdk.TimeValue(describeOutput.DeletedDate),
		}

		for _, outTag := range describeOutput.Tags {
			tags[*outTag.Key] = *outTag.Value
		}

	}
	return formatAwsSecret(id, value, tags, metadata), nil
}

//Get Get a secret and its description
func (s *SecretStore) Get(ctx context.Context, id, version string) (*entities.Secret, error) {

	getSecretInput := &secretsmanager.GetSecretValueInput{
		SecretId:  &id,
		VersionId: &version,
	}

	if len(version) == 0 {
		//Get with secret-id only
		//Here adding version would cause a not found error
		getSecretInput = &secretsmanager.GetSecretValueInput{
			SecretId: &id,
		}
	}

	getSecretOutput, err := s.client.GetSecret(ctx, getSecretInput)
	if err != nil || getSecretOutput == nil {
		return nil, errors.NotFoundError("secret not found")
	}

	describeInput := &secretsmanager.DescribeSecretInput{
		SecretId: &id,
	}

	//Prepare to get tags and metadata via description
	tags := make(map[string]string)
	metadata := &entities.Metadata{}

	describeOutput, err := s.client.DescribeSecret(ctx, describeInput)

	if err == nil && describeOutput != nil {
		metadata = &entities.Metadata{
			Version:   *getSecretOutput.VersionId,
			CreatedAt: sdk.TimeValue(describeOutput.CreatedDate),
			UpdatedAt: sdk.TimeValue(describeOutput.LastChangedDate),
			DeletedAt: sdk.TimeValue(describeOutput.DeletedDate),
		}

		for _, outTag := range describeOutput.Tags {
			tags[*outTag.Key] = *outTag.Value
		}
	}

	return formatAwsSecret(id, *getSecretOutput.SecretString, tags, metadata), nil
}

//List Get all secret ids as a slice of arns
func (s *SecretStore) List(ctx context.Context) ([]string, error) {

	//Leaving criteria unchanged should return all the keys (full list)
	listInput := &secretsmanager.ListSecretsInput{}
	listOutput, err := s.client.ListSecrets(ctx, listInput)
	if err != nil {
		return nil, errors.NotFoundError("secret not found")
	}

	//return only a list of secret names (IDs)
	secretNamesList := make([]string, len(listOutput.SecretList))
	for _, secret := range listOutput.SecretList {
		secretNamesList = append(secretNamesList, *secret.Name)
	}
	return secretNamesList, nil
}

// Refresh an existing secret by extending its TTL
func (s *SecretStore) Refresh(_ context.Context, id, _ string, expirationDate time.Time) error {
	return errors.ErrNotImplemented
}

//Delete Delete a secret
func (s *SecretStore) Delete(ctx context.Context, id string, versions ...string) (*entities.Secret, error) {
	deleteInput := &secretsmanager.DeleteSecretInput{
		SecretId: &id,
	}
	deleteOutput, err := s.client.DeleteSecret(ctx, deleteInput)
	if err != nil {
		return nil, errors.NotFoundError("secret not found")
	}
	return formatAwsSecret(*deleteOutput.Name, "", nil, nil), nil
}

// Gets a deleted secret
func (s *SecretStore) GetDeleted(_ context.Context, id string) (*entities.Secret, error) {
	return nil, errors.ErrNotImplemented
}

// Lists all deleted secrets
func (s *SecretStore) ListDeleted(ctx context.Context) ([]string, error) {
	return nil, errors.ErrNotImplemented
}

// Undelete a previously deleted secret
func (s *SecretStore) Undelete(ctx context.Context, id string) error {
	restoreInput := &secretsmanager.RestoreSecretInput{
		SecretId: &id,
	}

	_, err := s.client.RestoreSecret(ctx, restoreInput)
	if err != nil {
		return errors.NotFoundError("secret not found")
	}
	return nil
}

// Destroy a secret permanently (force deletion, secret will be unrecoverable)
func (s *SecretStore) Destroy(ctx context.Context, id string, versions ...string) error {
	forceDeletion := true
	deleteInput := &secretsmanager.DeleteSecretInput{
		SecretId:                   &id,
		ForceDeleteWithoutRecovery: &forceDeletion,
	}
	_, err := s.client.DeleteSecret(ctx, deleteInput)
	if err != nil {
		return errors.NotFoundError("secret not found")
	}
	return nil
}
