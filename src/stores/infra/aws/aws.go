package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/kms"

	"github.com/consensysquorum/quorum-key-manager/src/stores/store/entities"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

//go:generate mockgen -source=aws.go -destination=mocks/aws.go -package=mocks

type SecretsManagerClient interface {
	GetSecret(ctx context.Context, id, version string) (*secretsmanager.GetSecretValueOutput, error)
	CreateSecret(ctx context.Context, id, value string) (*secretsmanager.CreateSecretOutput, error)
	PutSecretValue(ctx context.Context, id, value string) (*secretsmanager.PutSecretValueOutput, error)
	TagSecretResource(ctx context.Context, id string, tags map[string]string) (*secretsmanager.TagResourceOutput, error)
	DescribeSecret(ctx context.Context, id string) (tags map[string]string, metadata *entities.Metadata, err error)
	ListSecrets(ctx context.Context, maxResults int64, nextToken string) (*secretsmanager.ListSecretsOutput, error)
	UpdateSecret(ctx context.Context, id, value, keyID, desc string) (*secretsmanager.UpdateSecretOutput, error)
	RestoreSecret(ctx context.Context, id string) (*secretsmanager.RestoreSecretOutput, error)
	DeleteSecret(ctx context.Context, id string, force bool) (*secretsmanager.DeleteSecretOutput, error)
}

type KmsClient interface {
	CreateKey(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*kms.CreateKeyOutput, error)
	GetPublicKey(ctx context.Context, keyID string) (*kms.GetPublicKeyOutput, error)
	ListKeys(ctx context.Context) ([]string, error)
	ListTags(ctx context.Context, keyID, marker string) (*kms.ListResourceTagsOutput, error)
	DescribeKey(ctx context.Context, id string) (*kms.DescribeKeyOutput, error)
	UpdateKey(ctx context.Context, keyID string, tags map[string]string) (*kms.TagResourceOutput, error)
	Sign(ctx context.Context, keyID string, msg []byte, signingAlgorithm string) (*kms.SignOutput, error)
	DeleteKey(ctx context.Context, keyID string) (*kms.ScheduleKeyDeletionOutput, error)
	RestoreKey(ctx context.Context, keyID string) (*kms.CancelKeyDeletionOutput, error)
	GetAlias(ctx context.Context, keyID string) (string, error)
}
