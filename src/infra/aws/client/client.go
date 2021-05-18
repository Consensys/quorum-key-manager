package client

import (
	"context"

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
	client := secretsmanager.New(session, aws.NewConfig().
		WithRegion(cfg.Region).
		WithLogLevel(aws.LogDebug).
		WithEndpoint(cfg.Endpoint))

	return &AwsVaultClient{*client}, nil

}

func (c *AwsVaultClient) GetSecret(ctx context.Context, input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	return c.client.GetSecretValue(input)
}
func (c *AwsVaultClient) CreateSecret(ctx context.Context, input *secretsmanager.CreateSecretInput) (*secretsmanager.CreateSecretOutput, error) {
	return c.client.CreateSecret(input)
}

func (c *AwsVaultClient) PutSecretValue(ctx context.Context, input *secretsmanager.PutSecretValueInput) (*secretsmanager.PutSecretValueOutput, error) {
	return c.client.PutSecretValue(input)
}

func (c *AwsVaultClient) TagSecretResource(ctx context.Context, input *secretsmanager.TagResourceInput) (*secretsmanager.TagResourceOutput, error) {
	return c.client.TagResource(input)
}

func (c *AwsVaultClient) DescribeSecret(ctx context.Context, input *secretsmanager.DescribeSecretInput) (*secretsmanager.DescribeSecretOutput, error) {
	return c.client.DescribeSecret(input)
}

func (c *AwsVaultClient) ListSecrets(ctx context.Context, criteria *secretsmanager.ListSecretsInput) (*secretsmanager.ListSecretsOutput, error) {
	return c.client.ListSecrets(criteria)

}
func (c *AwsVaultClient) UpdateSecret(ctx context.Context, input *secretsmanager.UpdateSecretInput) (*secretsmanager.UpdateSecretOutput, error) {
	return c.client.UpdateSecret(input)
}

func (c *AwsVaultClient) RestoreSecret(ctx context.Context, input *secretsmanager.RestoreSecretInput) (*secretsmanager.RestoreSecretOutput, error) {
	return c.client.RestoreSecret(input)
}
func (c *AwsVaultClient) DeleteSecret(ctx context.Context, input *secretsmanager.DeleteSecretInput) (*secretsmanager.DeleteSecretOutput, error) {
	return c.client.DeleteSecret(input)
}
