package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

//go:generate mockgen -source=aws.go -destination=mocks/aws.go -package=mocks

type SecretsManagerClient interface {
	GetSecret(ctx context.Context, input *secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error)
	CreateSecret(ctx context.Context, input *secretsmanager.CreateSecretInput) (*secretsmanager.CreateSecretOutput, error)
	PutSecretValue(ctx context.Context, input *secretsmanager.PutSecretValueInput) (*secretsmanager.PutSecretValueOutput, error)
	TagSecretResource(ctx context.Context, input *secretsmanager.TagResourceInput) (*secretsmanager.TagResourceOutput, error)
	DescribeSecret(ctx context.Context, input *secretsmanager.DescribeSecretInput) (*secretsmanager.DescribeSecretOutput, error)
	ListSecrets(ctx context.Context, criteria *secretsmanager.ListSecretsInput) (*secretsmanager.ListSecretsOutput, error)
	UpdateSecret(ctx context.Context, input *secretsmanager.UpdateSecretInput) (*secretsmanager.UpdateSecretOutput, error)
	RestoreSecret(ctx context.Context, input *secretsmanager.RestoreSecretInput) (*secretsmanager.RestoreSecretOutput, error)
	DeleteSecret(ctx context.Context, input *secretsmanager.DeleteSecretInput) (*secretsmanager.DeleteSecretOutput, error)
}
