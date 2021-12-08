package client

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	awsinfra "github.com/consensys/quorum-key-manager/src/infra/aws"
)

type AWSClient struct {
	secretsClient *secretsmanager.SecretsManager
	kmsClient     *kms.KMS
	cfg           *Config
}

var _ awsinfra.Client = &AWSClient{}

func New(cfg *Config) (*AWSClient, error) {
	sess, err := session.NewSession(cfg.ToAWSConfig())
	if err != nil {
		return nil, err
	}

	return &AWSClient{
		kmsClient:     kms.New(sess),
		secretsClient: secretsmanager.New(sess),
		cfg:           cfg,
	}, nil
}
