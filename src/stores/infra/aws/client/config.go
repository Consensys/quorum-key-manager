package client

import (
	"github.com/aws/aws-sdk-go/aws"
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
	awsConfig := aws.NewConfig().WithRegion(c.Region)

	if c.Debug {
		awsConfig.WithLogLevel(aws.LogDebug)
	}

	return awsConfig
}
