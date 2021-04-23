package client

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type SecretsClient interface {
	Set(ctx context.Context, request *types.SetSecretRequest) (*types.SecretResponse, error)
}

type KeyManagerClient interface {
	SecretsClient
}
