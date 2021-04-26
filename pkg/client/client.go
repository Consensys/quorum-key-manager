package client

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type SecretsClient interface {
	SetSecret(ctx context.Context, storeName string, request *types.SetSecretRequest) (*types.SecretResponse, error)
	GetSecret(ctx context.Context, storeName, id, version string) (*types.SecretResponse, error)
	ListSecrets(ctx context.Context, storeName string) ([]string, error)
}

type KeyManagerClient interface {
	SecretsClient
}
