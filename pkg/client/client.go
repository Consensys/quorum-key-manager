package client

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/api/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type SecretsClient interface {
	SetSecret(ctx context.Context, storeName string, request *types.SetSecretRequest) (*types.SecretResponse, error)
	GetSecret(ctx context.Context, storeName, id, version string) (*types.SecretResponse, error)
	ListSecrets(ctx context.Context, storeName string) ([]string, error)
}

type KeysClient interface {
	CreateKey(ctx context.Context, storeName string, request *types.CreateKeyRequest) (*types.KeyResponse, error)
	ImportKey(ctx context.Context, storeName string, request *types.ImportKeyRequest) (*types.KeyResponse, error)
	Sign(ctx context.Context, storeName, id string, request *types.SignBase64PayloadRequest) (string, error)
	GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error)
	ListKeys(ctx context.Context, storeName string) ([]string, error)
	DestroyKey(ctx context.Context, storeName, id string) error
}

type Eth1Client interface {
	CreateEth1Account(ctx context.Context, storeName string, request *types.CreateEth1AccountRequest) (*types.Eth1AccountResponse, error)
}

type JSONRPC interface {
	Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error)
}

type KeyManagerClient interface {
	SecretsClient
	KeysClient
	Eth1Client
	JSONRPC
}
