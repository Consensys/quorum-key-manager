package client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type Config struct {
	Endpoint  string
	Region    string
	AccessID  string
	SecretKey string
	Debug     bool
}

func NewConfig(region, accessID, secretKey string, debug bool) *Config {
	return &Config{
		Region:    region,
		AccessID:  accessID,
		SecretKey: secretKey,
		Debug:     debug,
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
