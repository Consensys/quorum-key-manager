package akv

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.1/keyvault"
)

//go:generate mockgen -source=akv.go -destination=mocks/akv.go -package=mocks

type Client interface {
	SetSecret(ctx context.Context, secretName string, parameters keyvault.SecretSetParameters) (result keyvault.SecretBundle, err error)
	GetSecret(ctx context.Context, secretName, secretVersion string) (result keyvault.SecretBundle, err error)
	GetSecrets(ctx context.Context, maxResults *int32) (result keyvault.SecretListResultPage, err error)
	UpdateSecret(ctx context.Context, secretName string, secretVersion string, parameters keyvault.SecretUpdateParameters) (result keyvault.SecretBundle, err error)
	DeleteSecret(ctx context.Context, secretName string) (result keyvault.DeletedSecretBundle, err error)
}
