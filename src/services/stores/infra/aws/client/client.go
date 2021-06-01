package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AwsVaultClient struct {
	client secretsmanager.SecretsManager
}

type AwsKmsClient struct {
	client kms.KMS
}

func NewSecretsClient(cfg *Config) (*AwsVaultClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create a Secrets Manager client
	client := secretsmanager.New(newSession,
		aws.NewConfig().WithRegion(cfg.Region).WithLogLevel(aws.LogDebug))

	return &AwsVaultClient{*client}, nil
}

func NewSecretsClientWithEndpoint(cfg *Config) (*AwsVaultClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create a Secrets Manager client
	config := aws.NewConfig().
		WithRegion(cfg.Region).
		WithEndpoint(cfg.Endpoint)

	if isDebugOn() {
		config.WithLogLevel(aws.LogDebug)
	}
	client := secretsmanager.New(newSession, config)

	return &AwsVaultClient{*client}, nil

}

func NewKmsClient(cfg *Config) (*AwsKmsClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create a Secrets Manager client
	client := kms.New(newSession,
		aws.NewConfig().WithRegion(cfg.Region).WithLogLevel(aws.LogDebug))

	return &AwsKmsClient{*client}, nil
}

func NewKmsClientWithEndpoint(cfg *Config) (*AwsKmsClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create a Secrets Manager client
	config := aws.NewConfig().
		WithRegion(cfg.Region).
		WithEndpoint(cfg.Endpoint)

	if isDebugOn() {
		config.WithLogLevel(aws.LogDebug)
	}
	client := kms.New(newSession, config)

	return &AwsKmsClient{*client}, nil

}
