package akv

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
)

//go:generate mockgen -source=akv.go -destination=mocks/akv.go -package=mocks

type Client interface {
	SecretClient
	KeysClient
}

type SecretClient interface {
	SetSecret(ctx context.Context, secretName string, value string, tags map[string]string) (keyvault.SecretBundle, error)
	GetSecret(ctx context.Context, secretName, secretVersion string) (keyvault.SecretBundle, error)
	GetSecrets(ctx context.Context, maxResults int32) ([]keyvault.SecretItem, error)
	UpdateSecret(ctx context.Context, secretName string, secretVersion string, expireAt time.Time) (keyvault.SecretBundle, error)
	DeleteSecret(ctx context.Context, secretName string) (keyvault.DeletedSecretBundle, error)
	GetDeletedSecret(ctx context.Context, secretName string) (keyvault.DeletedSecretBundle, error)
	GetDeletedSecrets(ctx context.Context, maxResults int32) ([]keyvault.DeletedSecretItem, error)
	PurgeDeletedSecret(ctx context.Context, secretName string) (bool, error)
	RecoverSecret(ctx context.Context, secretName string) (keyvault.SecretBundle, error)
}

type KeysClient interface {
	CreateKey(ctx context.Context, keyName string, kty keyvault.JSONWebKeyType, crv keyvault.JSONWebKeyCurveName, attr *keyvault.KeyAttributes, ops []keyvault.JSONWebKeyOperation, tags map[string]string) (keyvault.KeyBundle, error)
	ImportKey(ctx context.Context, keyName string, k *keyvault.JSONWebKey, attr *keyvault.KeyAttributes, tags map[string]string) (keyvault.KeyBundle, error)
	GetKey(ctx context.Context, name string, version string) (keyvault.KeyBundle, error)
	GetKeys(ctx context.Context, maxResults int32) ([]keyvault.KeyItem, error)
	UpdateKey(ctx context.Context, keyName string, version string, attr *keyvault.KeyAttributes, ops []keyvault.JSONWebKeyOperation, tags map[string]string) (keyvault.KeyBundle, error)
	DeleteKey(ctx context.Context, keyName string) (keyvault.DeletedKeyBundle, error)
	GetDeletedKey(ctx context.Context, keyName string) (keyvault.DeletedKeyBundle, error)
	GetDeletedKeys(ctx context.Context, maxResults int32) ([]keyvault.DeletedKeyItem, error)
	PurgeDeletedKey(ctx context.Context, keyName string) (bool, error)
	RecoverDeletedKey(ctx context.Context, keyName string) (keyvault.KeyBundle, error)
	Sign(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeySignatureAlgorithm, payload string) (string, error)
	Encrypt(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, payload string) (string, error)
	Decrypt(ctx context.Context, keyName string, version string, alg keyvault.JSONWebKeyEncryptionAlgorithm, value string) (string, error)
}
