package akv

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

//go:generate mockgen -source=akv.go -destination=mocks/akv.go -package=mocks

type Client interface {
	SecretClient
	KeysClient
}

type SecretClient interface {
	SetSecret(ctx context.Context, secretName string, value string, tags map[string]string) (*entities.Secret, error)
	GetSecret(ctx context.Context, secretName, secretVersion string) (*entities.Secret, error)
	GetSecrets(ctx context.Context, maxResults int32) ([]*entities.Secret, error)
	UpdateSecret(ctx context.Context, secretName string, secretVersion string, expireAt time.Time) (*entities.Secret, error)
	DeleteSecret(ctx context.Context, secretName string) (*entities.Secret, error)
	RecoverSecret(ctx context.Context, secretName string) (*entities.Secret, error)
	GetDeletedSecret(ctx context.Context, secretName string) (*entities.Secret, error)
	GetDeletedSecrets(ctx context.Context, maxResults int32) ([]*entities.Secret, error)
	PurgeDeletedSecret(ctx context.Context, secretName string) (bool, error)
}

type KeysClient interface {
	CreateKey(ctx context.Context, keyName string, kty keyvault.JSONWebKeyType, crv keyvault.JSONWebKeyCurveName, attr *entities.Attributes, ops []entities.CryptoOperation, tags map[string]string) (*entities.Key, error)
	ImportKey(ctx context.Context, keyName string, k *keyvault.JSONWebKey, attr *entities.Attributes, tags map[string]string) (*entities.Key, error)
	GetKey(ctx context.Context, name string, version string) (*entities.Key, error)
	GetKeys(ctx context.Context, maxResults int32) ([]*entities.Key, error)
	UpdateKey(ctx context.Context, keyName string, version string, attr *keyvault.KeyAttributes, ops []entities.CryptoOperation, tags map[string]string) (*entities.Key, error)
	DeleteKey(ctx context.Context, keyName string) (*entities.Key, error)
	GetDeletedKey(ctx context.Context, keyName string) (*entities.Key, error)
	GetDeletedKeys(ctx context.Context, maxResults int32) ([]*entities.Key, error)
	PurgeDeletedKey(ctx context.Context, keyName string) (bool, error)
	RecoverDeletedKey(ctx context.Context, keyName string) (*entities.Key, error)
	Sign(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeySignatureAlgorithm, payload string) (string, error)
	Encrypt(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, payload string) (string, error)
	Decrypt(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, value string) (string, error)
}
