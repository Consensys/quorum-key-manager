package client

import (
	"context"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/sdk/types"
)

//go:generate mockgen -source=client.go -destination=mocks/mock.go -package=mocks

type KeyManagerClient interface {
	SecretsClient
}

type SecretsClient interface {
	CreateSecret(ctx context.Context, request *types.CreateSecretRequest) (*types.Secret, error)
}
