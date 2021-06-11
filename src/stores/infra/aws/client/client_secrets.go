package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type AwsSecretsClient struct {
	client secretsmanager.SecretsManager
}

func NewSecretsClient(cfg *Config) (*AwsSecretsClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create config
	config := aws.NewConfig().WithRegion(cfg.Region)

	if isDebugOn() {
		config.WithLogLevel(aws.LogDebug)
	}
	// Create a Secrets Manager client
	client := secretsmanager.New(newSession, config)

	return &AwsSecretsClient{*client}, nil
}
