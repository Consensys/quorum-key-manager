package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"os"
	"strings"
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
	// Create config
	config := aws.NewConfig().WithRegion(cfg.Region)

	if isDebugOn() {
		config.WithLogLevel(aws.LogDebug)
	}
	// Create a Secrets Manager client
	client := secretsmanager.New(newSession, config)

	return &AwsSecretsClient{*client}, nil
}

func NewKmsClient(cfg *Config) (*AwsKmsClient, error) {
	// Create a new newSession
	newSession, _ := session.NewSession()
	// Create config
	config := aws.NewConfig().WithRegion(cfg.Region)

	if isDebugOn() {
		config.WithLogLevel(aws.LogDebug)
	}
	// Create a Secrets Manager client
	client := kms.New(newSession, config)

	return &AwsKmsClient{*client}, nil
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
