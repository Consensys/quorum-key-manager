package client

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	CurrentVersionMark = "AWSCURRENT"
)

func (c *AwsSecretsClient) GetSecret(ctx context.Context, id, version string) (*secretsmanager.GetSecretValueOutput, error) {
	getSecretInput := &secretsmanager.GetSecretValueInput{
		SecretId:  &id,
		VersionId: &version,
	}

	if version == "" {
		// Get with secret-id only
		// Here adding version would cause a not found error
		getSecretInput = &secretsmanager.GetSecretValueInput{
			SecretId: &id,
		}
	}
	return c.client.GetSecretValue(getSecretInput)
}
func (c *AwsSecretsClient) CreateSecret(ctx context.Context, id, value string) (*secretsmanager.CreateSecretOutput, error) {
	return c.client.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         &id,
		SecretString: &value,
	})
}

func (c *AwsSecretsClient) PutSecretValue(ctx context.Context, id, value string) (*secretsmanager.PutSecretValueOutput, error) {
	return c.client.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretId:     &id,
		SecretString: &value,
	})
}

func (c *AwsSecretsClient) TagSecretResource(ctx context.Context, id string, tags map[string]string) (*secretsmanager.TagResourceOutput, error) {

	var inputTags []*secretsmanager.Tag

	for key, value := range tags {
		k, v := key, value
		var inTag = secretsmanager.Tag{
			Key:   &k,
			Value: &v,
		}
		inputTags = append(inputTags, &inTag)
	}
	return c.client.TagResource(&secretsmanager.TagResourceInput{
		SecretId: &id,
		Tags:     inputTags,
	})
}

func (c *AwsSecretsClient) DescribeSecret(ctx context.Context, id string) (tags map[string]string, metadata *entities.Metadata, err error) {
	output, err := c.client.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: &id,
	})

	outTags := make(map[string]string)
	outMeta := &entities.Metadata{}

	if output != nil {

		// Trick to help us getting the actual current version as there is no versionID metadata
		currentVersion := ""
		for version, stages := range output.VersionIdsToStages {
			for _, stage := range stages {
				if *stage == CurrentVersionMark {
					currentVersion = version
				}
			}
		}

		outMeta = &entities.Metadata{
			Version:   currentVersion,
			CreatedAt: aws.TimeValue(output.CreatedDate),
			UpdatedAt: aws.TimeValue(output.LastChangedDate),
			DeletedAt: aws.TimeValue(output.DeletedDate),
		}

		for _, outTag := range output.Tags {
			outTags[*outTag.Key] = *outTag.Value
		}

	}
	return outTags, outMeta, translateAwsError(err)
}

func (c *AwsSecretsClient) ListSecrets(ctx context.Context, maxResults int64, nextToken string) (*secretsmanager.ListSecretsOutput, error) {
	listInput := &secretsmanager.ListSecretsInput{}
	if len(nextToken) > 0 {
		listInput.NextToken = &nextToken
	}
	if maxResults > 0 {
		listInput.MaxResults = &maxResults
	}
	return c.client.ListSecrets(listInput)

}
func (c *AwsSecretsClient) UpdateSecret(ctx context.Context, id, value, keyID, desc string) (*secretsmanager.UpdateSecretOutput, error) {
	return c.client.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     &id,
		SecretString: &value,
		KmsKeyId:     &keyID,
		Description:  &desc,
	})
}

func (c *AwsSecretsClient) RestoreSecret(ctx context.Context, id string) (*secretsmanager.RestoreSecretOutput, error) {
	return c.client.RestoreSecret(&secretsmanager.RestoreSecretInput{
		SecretId: &id,
	})
}
func (c *AwsSecretsClient) DeleteSecret(ctx context.Context, id string, force bool) (*secretsmanager.DeleteSecretOutput, error) {

	return c.client.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   &id,
		ForceDeleteWithoutRecovery: &force,
	})
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
