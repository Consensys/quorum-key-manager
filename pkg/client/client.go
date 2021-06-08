package client

import (
	"context"
	types2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type SecretsClient interface {
	SetSecret(ctx context.Context, storeName string, request *types2.SetSecretRequest) (*types2.SecretResponse, error)
	GetSecret(ctx context.Context, storeName, id, version string) (*types2.SecretResponse, error)
	ListSecrets(ctx context.Context, storeName string) ([]string, error)
}

type KeysClient interface {
	CreateKey(ctx context.Context, storeName string, request *types2.CreateKeyRequest) (*types2.KeyResponse, error)
	ImportKey(ctx context.Context, storeName string, request *types2.ImportKeyRequest) (*types2.KeyResponse, error)
	Sign(ctx context.Context, storeName, id string, request *types2.SignBase64PayloadRequest) (string, error)
	GetKey(ctx context.Context, storeName, id string) (*types2.KeyResponse, error)
	ListKeys(ctx context.Context, storeName string) ([]string, error)
	DestroyKey(ctx context.Context, storeName, id string) error
}

type Eth1Client interface {
	CreateEth1Account(ctx context.Context, storeName string, request *types2.CreateEth1AccountRequest) (*types2.Eth1AccountResponse, error)
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
