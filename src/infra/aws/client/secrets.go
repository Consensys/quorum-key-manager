package client

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/entities"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const CurrentVersionMark = "AWSCURRENT"

func (c *AWSClient) GetSecret(_ context.Context, id, version string) (*secretsmanager.GetSecretValueOutput, error) {
	getSecretInput := &secretsmanager.GetSecretValueInput{
		SecretId:  &id,
		VersionId: nil,
	}

	if version != "" {
		getSecretInput.VersionId = &version
	}

	output, err := c.secretsClient.GetSecretValue(getSecretInput)
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}
func (c *AWSClient) CreateSecret(_ context.Context, id, value string) (*secretsmanager.CreateSecretOutput, error) {
	output, err := c.secretsClient.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         &id,
		SecretString: &value,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}

func (c *AWSClient) PutSecretValue(_ context.Context, id, value string) (*secretsmanager.PutSecretValueOutput, error) {
	output, err := c.secretsClient.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretId:     &id,
		SecretString: &value,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}

func (c *AWSClient) TagSecretResource(_ context.Context, id string, tags map[string]string) (*secretsmanager.TagResourceOutput, error) {
	var inputTags []*secretsmanager.Tag
	for key, value := range tags {
		k, v := key, value
		var inTag = secretsmanager.Tag{
			Key:   &k,
			Value: &v,
		}
		inputTags = append(inputTags, &inTag)
	}
	output, err := c.secretsClient.TagResource(&secretsmanager.TagResourceInput{
		SecretId: &id,
		Tags:     inputTags,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}

func (c *AWSClient) DescribeSecret(_ context.Context, id string) (tags map[string]string, metadata *entities.Metadata, err error) {
	output, err := c.secretsClient.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: &id,
	})

	outTags := make(map[string]string)
	outMeta := &entities.Metadata{}

	if output != nil {

		// Trick to help us to get the actual current version as there is no versionID metadata
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
	if err != nil {
		return nil, nil, parseSecretsManagerErrorResponse(err)
	}

	return outTags, outMeta, nil
}

func (c *AWSClient) ListSecrets(_ context.Context, maxResults int64, nextToken string) (*secretsmanager.ListSecretsOutput, error) {
	listInput := &secretsmanager.ListSecretsInput{}
	if len(nextToken) > 0 {
		listInput.NextToken = &nextToken
	}
	if maxResults > 0 {
		listInput.MaxResults = &maxResults
	}
	output, err := c.secretsClient.ListSecrets(listInput)
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil

}
func (c *AWSClient) UpdateSecret(_ context.Context, id, value, keyID, desc string) (*secretsmanager.UpdateSecretOutput, error) {
	output, err := c.secretsClient.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     &id,
		SecretString: &value,
		KmsKeyId:     &keyID,
		Description:  &desc,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}

func (c *AWSClient) RestoreSecret(_ context.Context, id string) (*secretsmanager.RestoreSecretOutput, error) {
	output, err := c.secretsClient.RestoreSecret(&secretsmanager.RestoreSecretInput{
		SecretId: &id,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}
func (c *AWSClient) DeleteSecret(_ context.Context, id string) (*secretsmanager.DeleteSecretOutput, error) {
	output, err := c.secretsClient.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   &id,
		ForceDeleteWithoutRecovery: common.ToPtr(false).(*bool),
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}

func (c *AWSClient) DestroySecret(ctx context.Context, id string) (*secretsmanager.DeleteSecretOutput, error) {
	// check appropriate state with description
	desc, err := c.secretsClient.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: &id,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}
	if desc.DeletedDate == nil {
		return nil, errors.InvalidParameterError("failed to destroy, must be deleted first")
	}

	// We need to restore before we can destroy the key
	_, err = c.RestoreSecret(ctx, id)
	if err != nil {
		return nil, err
	}

	output, err := c.secretsClient.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   &id,
		ForceDeleteWithoutRecovery: common.ToPtr(true).(*bool),
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}
