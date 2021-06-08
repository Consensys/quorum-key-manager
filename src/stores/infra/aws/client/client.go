package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AwsSecretsClient struct {
	client secretsmanager.SecretsManager
}

type AwsKmsClient struct {
	client kms.KMS
}

func NewSecretsClient(cfg *Config) (*AwsSecretsClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create a Secrets Manager client
	client := secretsmanager.New(newSession,
		aws.NewConfig().WithRegion(cfg.Region).WithLogLevel(aws.LogDebug))

	return &AwsSecretsClient{*client}, nil
}

func NewKmsClient(cfg *Config) (*AwsKmsClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create a Secrets Manager client
	client := kms.New(newSession,
		aws.NewConfig().WithRegion(cfg.Region).WithLogLevel(aws.LogDebug))

	return &AwsKmsClient{*client}, nil
}
