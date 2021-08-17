package client

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type SecretsClient interface {
	SetSecret(ctx context.Context, storeName, id string, request *types.SetSecretRequest) (*types.SecretResponse, error)
	GetSecret(ctx context.Context, storeName, id, version string) (*types.SecretResponse, error)
	GetDeletedSecret(ctx context.Context, storeName, id string) (*types.SecretResponse, error)
	DeleteSecret(ctx context.Context, storeName, id string) error
	RestoreSecret(ctx context.Context, storeName, id string) error
	DestroySecret(ctx context.Context, storeName, id string) error
	ListSecrets(ctx context.Context, storeName string) ([]string, error)
	ListDeletedSecrets(ctx context.Context, storeName string) ([]string, error)
}

type KeysClient interface {
	CreateKey(ctx context.Context, storeName, id string, request *types.CreateKeyRequest) (*types.KeyResponse, error)
	ImportKey(ctx context.Context, storeName, id string, request *types.ImportKeyRequest) (*types.KeyResponse, error)
	SignKey(ctx context.Context, storeName, id string, request *types.SignBase64PayloadRequest) (string, error)
	VerifyKeySignature(ctx context.Context, storeName string, request *types.VerifyKeySignatureRequest) error
	GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error)
	ListKeys(ctx context.Context, storeName string) ([]string, error)
	DestroyKey(ctx context.Context, storeName, id string) error
}

type Eth1Client interface {
	CreateEth1Account(ctx context.Context, storeName string, request *types.CreateEth1AccountRequest) (*types.Eth1AccountResponse, error)
	ImportEth1Account(ctx context.Context, storeName string, request *types.ImportEth1AccountRequest) (*types.Eth1AccountResponse, error)
	UpdateEth1Account(ctx context.Context, storeName, address string, request *types.UpdateEth1AccountRequest) (*types.Eth1AccountResponse, error)
	SignMessage(ctx context.Context, storeName, account string, request *types.SignMessageRequest) (string, error)
	SignTypedData(ctx context.Context, storeName, address string, request *types.SignTypedDataRequest) (string, error)
	SignTransaction(ctx context.Context, storeName, address string, request *types.SignETHTransactionRequest) (string, error)
	SignQuorumPrivateTransaction(ctx context.Context, storeName, address string, request *types.SignQuorumPrivateTransactionRequest) (string, error)
	SignEEATransaction(ctx context.Context, storeName, address string, request *types.SignEEATransactionRequest) (string, error)
	GetEth1Account(ctx context.Context, storeName, address string) (*types.Eth1AccountResponse, error)
	ListEth1Accounts(ctx context.Context, storeName string) ([]string, error)
	DeleteEth1Account(ctx context.Context, storeName, address string) error
	DestroyEth1Account(ctx context.Context, storeName, address string) error
	RestoreEth1Account(ctx context.Context, storeName, address string) error
	ECRecover(ctx context.Context, storeName string, request *types.ECRecoverRequest) (string, error)
	VerifyEth1Signature(ctx context.Context, storeName string, request *types.VerifyEth1SignatureRequest) error
	VerifyTypedDataSignature(ctx context.Context, storeName string, request *types.VerifyTypedDataRequest) error
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
