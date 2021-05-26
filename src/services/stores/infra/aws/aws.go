package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

//go:generate mockgen -source=aws.go -destination=mocks/aws.go -package=mocks

type SecretsManagerClient interface {
	GetSecret(ctx context.Context, id, version string) (*secretsmanager.GetSecretValueOutput, error)
	CreateSecret(ctx context.Context, id, value string) (*secretsmanager.CreateSecretOutput, error)
	PutSecretValue(ctx context.Context, id, value string) (*secretsmanager.PutSecretValueOutput, error)
	TagSecretResource(ctx context.Context, id string, tags map[string]string) (*secretsmanager.TagResourceOutput, error)
	DescribeSecret(ctx context.Context, id string) (*secretsmanager.DescribeSecretOutput, error)
	ListSecrets(ctx context.Context, maxResults int64, nextToken string) (*secretsmanager.ListSecretsOutput, error)
	UpdateSecret(ctx context.Context, id, value, keyID, desc string) (*secretsmanager.UpdateSecretOutput, error)
	RestoreSecret(ctx context.Context, id string) (*secretsmanager.RestoreSecretOutput, error)
	DeleteSecret(ctx context.Context, id string, force bool) (*secretsmanager.DeleteSecretOutput, error)
}
