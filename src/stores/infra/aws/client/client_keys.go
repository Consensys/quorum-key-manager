package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type AwsKmsClient struct {
	client kms.KMS
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
