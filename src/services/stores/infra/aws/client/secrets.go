package client

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
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

func (c *AwsSecretsClient) DescribeSecret(ctx context.Context, id string) (*secretsmanager.DescribeSecretOutput, error) {
	return c.client.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: &id,
	})
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

func isDebugOn() bool {
	val, ok := os.LookupEnv("AWS_DEBUG")
	if !ok {
		return false
	}
	return strings.EqualFold("true", val) ||
		strings.EqualFold("on", val) ||
		strings.EqualFold("yes", val)

}
