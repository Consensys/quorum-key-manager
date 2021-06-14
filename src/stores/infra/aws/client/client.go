package client

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/consensysquorum/quorum-key-manager/pkg/errors"
	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"
)

const (
	CurrentVersionMark = "AWSCURRENT"
)

type AwsVaultClient struct {
	client secretsmanager.SecretsManager
}

func NewClient(cfg *Config) (*AwsVaultClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create a Secrets Manager client
	client := secretsmanager.New(newSession,
		aws.NewConfig().WithRegion(cfg.Region).WithLogLevel(aws.LogDebug))

	return &AwsVaultClient{*client}, nil
}

func NewClientWithEndpoint(cfg *Config) (*AwsVaultClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create a Secrets Manager client
	config := aws.NewConfig().
		WithRegion(cfg.Region).
		WithEndpoint(cfg.Endpoint)

	// TODO: Use field in config to get this param instead of checking the ENV var directly
	if isDebugOn() {
		config.WithLogLevel(aws.LogDebug)
	}
	client := secretsmanager.New(newSession, config)

	return &AwsVaultClient{*client}, nil

}

func (c *AwsVaultClient) GetSecret(ctx context.Context, id, version string) (*secretsmanager.GetSecretValueOutput, error) {
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
	output, err := c.client.GetSecretValue(getSecretInput)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return output, nil
}

func (c *AwsVaultClient) CreateSecret(ctx context.Context, id, value string) (*secretsmanager.CreateSecretOutput, error) {
	output, err := c.client.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         &id,
		SecretString: &value,
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return output, nil
}

func (c *AwsVaultClient) PutSecretValue(ctx context.Context, id, value string) (*secretsmanager.PutSecretValueOutput, error) {
	output, err := c.client.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretId:     &id,
		SecretString: &value,
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return output, nil
}

func (c *AwsVaultClient) TagSecretResource(ctx context.Context, id string, tags map[string]string) (*secretsmanager.TagResourceOutput, error) {

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
		return nil, parseErrorResponse(err)
	}

	return output, nil
}

func (c *AwsVaultClient) DescribeSecret(ctx context.Context, id string) (tags map[string]string, metadata *entities.Metadata, err error) {
	output, err := c.client.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: &id,
	})
	if err != nil {
		return nil, nil, parseErrorResponse(err)
	}

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

	return outTags, outMeta, nil
}

func (c *AwsVaultClient) ListSecrets(ctx context.Context, maxResults int64, nextToken string) (*secretsmanager.ListSecretsOutput, error) {
	listInput := &secretsmanager.ListSecretsInput{}
	if len(nextToken) > 0 {
		listInput.NextToken = &nextToken
	}
	if maxResults > 0 {
		listInput.MaxResults = &maxResults
	}
	output, err := c.client.ListSecrets(listInput)
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return output, nil

}
func (c *AwsVaultClient) UpdateSecret(ctx context.Context, id, value, keyID, desc string) (*secretsmanager.UpdateSecretOutput, error) {
	output, err := c.client.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     &id,
		SecretString: &value,
		KmsKeyId:     &keyID,
		Description:  &desc,
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return output, nil
}

func (c *AwsVaultClient) RestoreSecret(ctx context.Context, id string) (*secretsmanager.RestoreSecretOutput, error) {
	output, err := c.client.RestoreSecret(&secretsmanager.RestoreSecretInput{
		SecretId: &id,
	})
	if err != nil {
		return nil, parseErrorResponse(err)
	}

	return output, nil
}

func (c *AwsVaultClient) DeleteSecret(ctx context.Context, id string, force bool) (*secretsmanager.DeleteSecretOutput, error) {
	if force {
		// check appropriate state with description
		desc, err := c.client.DescribeSecret(&secretsmanager.DescribeSecretInput{
			SecretId: &id,
		})
		if err != nil {
			return nil, parseErrorResponse(err)
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
		return nil, parseErrorResponse(err)
	}

	return output, nil
}

// TODO: Use field in config to get this param instead of checking the ENV var directly
func isDebugOn() bool {
	val, ok := os.LookupEnv("AWS_DEBUG")
	if !ok {
		return false
	}
	return strings.EqualFold("true", val) || strings.EqualFold("on", val) || strings.EqualFold("yes", val)
}
