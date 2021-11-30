package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/consensys/quorum-key-manager/src/entities"
)

type Config struct {
	Region    string
	AccessID  string
	SecretKey string
	Debug     bool
}

func NewConfig(cfg *entities.AWSConfig) *Config {
	return &Config{
		Region:    cfg.Region,
		AccessID:  cfg.AccessID,
		SecretKey: cfg.SecretKey,
		Debug:     cfg.Debug,
	}
}

func (c *Config) ToAWSConfig() *aws.Config {
	awsConfig := &aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessID, c.SecretKey, ""),
		Region:      aws.String(c.Region),
	}

	if c.Debug {
		awsConfig.WithLogLevel(aws.LogDebug)
		awsConfig.CredentialsChainVerboseErrors = aws.Bool(true)
	}

	return awsConfig
}
