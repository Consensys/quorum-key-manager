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
	ListSecrets(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	ListDeletedSecrets(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
}

type KeysClient interface {
	CreateKey(ctx context.Context, storeName, id string, request *types.CreateKeyRequest) (*types.KeyResponse, error)
	ImportKey(ctx context.Context, storeName, id string, request *types.ImportKeyRequest) (*types.KeyResponse, error)
	SignKey(ctx context.Context, storeName, id string, request *types.SignBase64PayloadRequest) (string, error)
	GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error)
	ListKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	DeleteKey(ctx context.Context, storeName, id string) error
	GetDeletedKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error)
	ListDeletedKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	RestoreKey(ctx context.Context, storeName, id string) error
	DestroyKey(ctx context.Context, storeName, id string) error
}

type EthClient interface {
	CreateEthAccount(ctx context.Context, storeName string, request *types.CreateEthAccountRequest) (*types.EthAccountResponse, error)
	ImportEthAccount(ctx context.Context, storeName string, request *types.ImportEthAccountRequest) (*types.EthAccountResponse, error)
	UpdateEthAccount(ctx context.Context, storeName, address string, request *types.UpdateEthAccountRequest) (*types.EthAccountResponse, error)
	SignMessage(ctx context.Context, storeName, account string, request *types.SignMessageRequest) (string, error)
	SignTypedData(ctx context.Context, storeName, address string, request *types.SignTypedDataRequest) (string, error)
	SignTransaction(ctx context.Context, storeName, address string, request *types.SignETHTransactionRequest) (string, error)
	SignQuorumPrivateTransaction(ctx context.Context, storeName, address string, request *types.SignQuorumPrivateTransactionRequest) (string, error)
	SignEEATransaction(ctx context.Context, storeName, address string, request *types.SignEEATransactionRequest) (string, error)
	GetEthAccount(ctx context.Context, storeName, address string) (*types.EthAccountResponse, error)
	ListEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	ListDeletedEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	DeleteEthAccount(ctx context.Context, storeName, address string) error
	DestroyEthAccount(ctx context.Context, storeName, address string) error
	RestoreEthAccount(ctx context.Context, storeName, address string) error
}

type UtilsClient interface {
	VerifyKeySignature(ctx context.Context, request *types.VerifyKeySignatureRequest) error
	ECRecover(ctx context.Context, request *types.ECRecoverRequest) (string, error)
	VerifyMessage(ctx context.Context, request *types.VerifyRequest) error
	VerifyTypedData(ctx context.Context, request *types.VerifyTypedDataRequest) error
}

type JSONRPC interface {
	Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error)
}

type KeyManagerClient interface {
	SecretsClient
	KeysClient
	EthClient
	UtilsClient
	JSONRPC
}
