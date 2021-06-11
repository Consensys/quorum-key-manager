package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/service/kms"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"

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
	CreateKey(ctx context.Context, id string, alg *entities.Algorithm, attr *entities.Attributes) (*kms.CreateKeyOutput, *string, error)
	// ImportKey(ctx context.Context, input *kms.ImportKeyMaterialInput, tags map[string]string) (*kms.ImportKeyMaterialOutput, error)
	GetPublicKey(ctx context.Context, name string) (*kms.GetPublicKeyOutput, error)
	ListKeys(ctx context.Context, limit int64, marker string) (*kms.ListKeysOutput, error)
	ListTags(ctx context.Context, id, marker string) (*kms.ListResourceTagsOutput, error)
	ListAliases(ctx context.Context, id, marker string) (*kms.ListAliasesOutput, []string, error)
	DescribeKey(ctx context.Context, id string) (*kms.DescribeKeyOutput, error)
	UpdateKey(ctx context.Context, id string, tags map[string]string) (*kms.UpdateCustomKeyStoreOutput, error)
	/*
		GetDeletedKey(ctx context.Context, keyName string) (keyvault.DeletedKeyBundle, error)
		GetDeletedKeys(ctx context.Context, maxResults int32) ([]keyvault.DeletedKeyItem, error)
		PurgeDeletedKey(ctx context.Context, keyName string) (bool, error)
		RecoverDeletedKey(ctx context.Context, keyName string) (keyvault.KeyBundle, error)*/
	Sign(ctx context.Context, id string, msg []byte) (*kms.SignOutput, error)
	Verify(ctx context.Context, id string, msg, signature []byte) (*kms.VerifyOutput, error)
	DeleteKey(ctx context.Context, id string) (*kms.DisableKeyOutput, error)
	/*Encrypt(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, payload string) (string, error)
	Decrypt(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, value string) (string, error)
	*/
}
