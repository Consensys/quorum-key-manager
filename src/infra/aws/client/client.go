package client

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/cenkalti/backoff/v4"
	awsinfra "github.com/consensys/quorum-key-manager/src/infra/aws"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type AWSClient struct {
	secretsClient *secretsmanager.SecretsManager
	kmsClient     *kms.KMS
	cfg           *Config
	backOff       backoff.BackOff
	logger        log.Logger
}

var _ awsinfra.Client = &AWSClient{}

func New(cfg *Config, logger log.Logger) (*AWSClient, error) {
	sess, err := session.NewSession(cfg.ToAWSConfig())
	if err != nil {
		return nil, err
	}

	return &AWSClient{
		kmsClient:     kms.New(sess),
		secretsClient: secretsmanager.New(sess),
		// Max wait of 5 seconds to wait for KMS to transition the state of assets
		backOff: backoff.WithMaxRetries(backoff.NewConstantBackOff(time.Second), 5),
		cfg:     cfg,
		logger:  logger,
	}, nil
}
