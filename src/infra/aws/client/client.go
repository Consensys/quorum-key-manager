package client

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AwsVaultClient struct {
	client secretsmanager.SecretsManager
}

func NewClient(cfg *Config) (*AwsVaultClient, error) {
	//Create a Secrets Manager client
	client := secretsmanager.New(session.New(),
		aws.NewConfig().WithRegion(cfg.Region).WithLogLevel(aws.LogDebug))

	return &AwsVaultClient{*client}, nil
}

func NewClientWithEndpoint(cfg *Config) (*AwsVaultClient, error) {
	//Create a new session
	session, _ := session.NewSession()
	//Create a Secrets Manager client
	config := aws.NewConfig().
		WithRegion(cfg.Region).
		WithEndpoint(cfg.Endpoint)

	if isDebugOn() {
		config.WithLogLevel(aws.LogDebug)
	}
	client := secretsmanager.New(session, config)

	return &AwsVaultClient{*client}, nil

}

func (c *AwsVaultClient) GetSecret(ctx context.Context, id, version string) (*secretsmanager.GetSecretValueOutput, error) {
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
	return c.client.GetSecretValue(getSecretInput)
}
func (c *AwsVaultClient) CreateSecret(ctx context.Context, id, value string) (*secretsmanager.CreateSecretOutput, error) {
	return c.client.CreateSecret(&secretsmanager.CreateSecretInput{
		Name:         &id,
		SecretString: &value,
	})
}

func (c *AwsVaultClient) PutSecretValue(ctx context.Context, id, value string) (*secretsmanager.PutSecretValueOutput, error) {
	return c.client.PutSecretValue(&secretsmanager.PutSecretValueInput{
		SecretId:     &id,
		SecretString: &value,
	})
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
	return c.client.TagResource(&secretsmanager.TagResourceInput{
		SecretId: &id,
		Tags:     inputTags,
	})
}

func (c *AwsVaultClient) DescribeSecret(ctx context.Context, id string) (*secretsmanager.DescribeSecretOutput, error) {
	return c.client.DescribeSecret(&secretsmanager.DescribeSecretInput{
		SecretId: &id,
	})
}

func (c *AwsVaultClient) ListSecrets(ctx context.Context) (*secretsmanager.ListSecretsOutput, error) {
	return c.client.ListSecrets(&secretsmanager.ListSecretsInput{})

}
func (c *AwsVaultClient) UpdateSecret(ctx context.Context, id, value, keyID, desc string) (*secretsmanager.UpdateSecretOutput, error) {
	return c.client.UpdateSecret(&secretsmanager.UpdateSecretInput{
		SecretId:     &id,
		SecretString: &value,
		KmsKeyId:     &keyID,
		Description:  &desc,
	})
}

func (c *AwsVaultClient) RestoreSecret(ctx context.Context, id string) (*secretsmanager.RestoreSecretOutput, error) {
	return c.client.RestoreSecret(&secretsmanager.RestoreSecretInput{
		SecretId: &id,
	})
}
func (c *AwsVaultClient) DeleteSecret(ctx context.Context, id string, force bool) (*secretsmanager.DeleteSecretOutput, error) {

	return c.client.DeleteSecret(&secretsmanager.DeleteSecretInput{
		SecretId:                   &id,
		ForceDeleteWithoutRecovery: &force,
	})
}

func isDebugOn() bool {
	val, ok := os.LookupEnv("AWS_DEBUG")
	if !ok {
		return false
	} else {
		return strings.EqualFold("true", val) ||
			strings.EqualFold("on", val) ||
			strings.EqualFold("yes", val)
	}
}
