package client

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

const (
	CurrentVersionMark = "AWSCURRENT"
)

type AwsSecretsClient struct {
	client *secretsmanager.SecretsManager
	cfg    *Config
}

func NewSecretsClient(cfg *Config) (*AwsSecretsClient, error) {
	newSession, err := session.NewSession(cfg.ToAWSConfig())
	if err != nil {
		return nil, err
	}

	return &AwsSecretsClient{
		client: secretsmanager.New(newSession),
		cfg:    cfg,
	}, nil
}

func (c *AwsSecretsClient) GetSecret(_ context.Context, id, version string) (*secretsmanager.GetSecretValueOutput, error) {
	getSecretInput := &secretsmanager.GetSecretValueInput{
		SecretId:  &id,
		VersionId: nil,
	}

	if version != "" {
		getSecretInput.VersionId = &version
	}

	output, err := c.client.GetSecretValue(getSecretInput)
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}
func (c *AwsSecretsClient) CreateSecret(_ context.Context, id, value string) (*secretsmanager.CreateSecretOutput, error) {
	output, err := c.client.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         &id,
		SecretString: &value,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}

func (c *AwsSecretsClient) PutSecretValue(_ context.Context, id, value string) (*secretsmanager.PutSecretValueOutput, error) {
	output, err := c.client.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretId:     &id,
		SecretString: &value,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}

func (c *AwsSecretsClient) TagSecretResource(_ context.Context, id string, tags map[string]string) (*secretsmanager.TagResourceOutput, error) {

	var inputTags []*secretsmanager.Tag

	for key, value := range tags {
		k, v := key, value
		var inTag = secretsmanager.Tag{
			Key:   &k,
			Value: &v,
		}
		inputTags = append(inputTags, &inTag)
	}
	output, err := c.client.TagResource(&secretsmanager.TagResourceInput{
		SecretId: &id,
		Tags:     inputTags,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}

func (c *AwsSecretsClient) DescribeSecret(_ context.Context, id string) (tags map[string]string, metadata *entities.Metadata, err error) {
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
	if err != nil {
		return nil, nil, parseSecretsManagerErrorResponse(err)
	}

	return outTags, outMeta, nil
}

func (c *AwsSecretsClient) ListSecrets(_ context.Context, maxResults int64, nextToken string) (*secretsmanager.ListSecretsOutput, error) {
	listInput := &secretsmanager.ListSecretsInput{}
	if len(nextToken) > 0 {
		listInput.NextToken = &nextToken
	}
	if maxResults > 0 {
		listInput.MaxResults = &maxResults
	}
	output, err := c.client.ListSecrets(listInput)
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil

}
func (c *AwsSecretsClient) UpdateSecret(_ context.Context, id, value, keyID, desc string) (*secretsmanager.UpdateSecretOutput, error) {
	output, err := c.client.UpdateSecret(&secretsmanager.UpdateSecretInput{
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

func (c *AwsSecretsClient) RestoreSecret(_ context.Context, id string) (*secretsmanager.RestoreSecretOutput, error) {
	output, err := c.client.RestoreSecret(&secretsmanager.RestoreSecretInput{
		SecretId: &id,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}
func (c *AwsSecretsClient) DeleteSecret(_ context.Context, id string, force bool) (*secretsmanager.DeleteSecretOutput, error) {

	if force {
		// check appropriate state with description
		desc, err := c.client.DescribeSecret(&secretsmanager.DescribeSecretInput{
			SecretId: &id,
		})
		if err != nil {
			return nil, parseSecretsManagerErrorResponse(err)
		}
		if desc.DeletedDate != nil {
			return nil, errors.InvalidParameterError("failed to destroy, must be deleted first")
		}
	}
	output, err := c.client.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   &id,
		ForceDeleteWithoutRecovery: &force,
	})
	if err != nil {
		return nil, parseSecretsManagerErrorResponse(err)
	}

	return output, nil
}
