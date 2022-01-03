package client

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	aliastypes "github.com/consensys/quorum-key-manager/src/aliases/api/types"
	storestypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
	utilstypes "github.com/consensys/quorum-key-manager/src/utils/api/types"
)

//go:generate mockgen -source=client.go -destination=mock/mock.go -package=mock

type SecretsClient interface {
	SetSecret(ctx context.Context, storeName, id string, request *storestypes.SetSecretRequest) (*storestypes.SecretResponse, error)
	GetSecret(ctx context.Context, storeName, id, version string) (*storestypes.SecretResponse, error)
	GetDeletedSecret(ctx context.Context, storeName, id string) (*storestypes.SecretResponse, error)
	DeleteSecret(ctx context.Context, storeName, id string) error
	RestoreSecret(ctx context.Context, storeName, id string) error
	DestroySecret(ctx context.Context, storeName, id string) error
	ListSecrets(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	ListDeletedSecrets(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
}

type KeysClient interface {
	CreateKey(ctx context.Context, storeName, id string, request *storestypes.CreateKeyRequest) (*storestypes.KeyResponse, error)
	ImportKey(ctx context.Context, storeName, id string, request *storestypes.ImportKeyRequest) (*storestypes.KeyResponse, error)
	SignKey(ctx context.Context, storeName, id string, request *storestypes.SignBase64PayloadRequest) (string, error)
	GetKey(ctx context.Context, storeName, id string) (*storestypes.KeyResponse, error)
	ListKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	DeleteKey(ctx context.Context, storeName, id string) error
	GetDeletedKey(ctx context.Context, storeName, id string) (*storestypes.KeyResponse, error)
	ListDeletedKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	RestoreKey(ctx context.Context, storeName, id string) error
	DestroyKey(ctx context.Context, storeName, id string) error
}

type EthClient interface {
	CreateEthAccount(ctx context.Context, storeName string, request *storestypes.CreateEthAccountRequest) (*storestypes.EthAccountResponse, error)
	ImportEthAccount(ctx context.Context, storeName string, request *storestypes.ImportEthAccountRequest) (*storestypes.EthAccountResponse, error)
	UpdateEthAccount(ctx context.Context, storeName, address string, request *storestypes.UpdateEthAccountRequest) (*storestypes.EthAccountResponse, error)
	SignMessage(ctx context.Context, storeName, account string, request *storestypes.SignMessageRequest) (string, error)
	SignTypedData(ctx context.Context, storeName, address string, request *storestypes.SignTypedDataRequest) (string, error)
	SignTransaction(ctx context.Context, storeName, address string, request *storestypes.SignETHTransactionRequest) (string, error)
	SignQuorumPrivateTransaction(ctx context.Context, storeName, address string, request *storestypes.SignQuorumPrivateTransactionRequest) (string, error)
	SignEEATransaction(ctx context.Context, storeName, address string, request *storestypes.SignEEATransactionRequest) (string, error)
	GetEthAccount(ctx context.Context, storeName, address string) (*storestypes.EthAccountResponse, error)
	ListEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	ListDeletedEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error)
	DeleteEthAccount(ctx context.Context, storeName, address string) error
	DestroyEthAccount(ctx context.Context, storeName, address string) error
	RestoreEthAccount(ctx context.Context, storeName, address string) error
}

type UtilsClient interface {
	VerifyKeySignature(ctx context.Context, request *utilstypes.VerifyKeySignatureRequest) error
	ECRecover(ctx context.Context, request *utilstypes.ECRecoverRequest) (string, error)
	VerifyMessage(ctx context.Context, request *utilstypes.VerifyRequest) error
	VerifyTypedData(ctx context.Context, request *utilstypes.VerifyTypedDataRequest) error
}

type AliasClient interface {
	CreateAlias(ctx context.Context, registry, aliasKey string, req *aliastypes.AliasRequest) (*aliastypes.AliasResponse, error)
	GetAlias(ctx context.Context, registry, aliasKey string) (*aliastypes.AliasResponse, error)
	UpdateAlias(ctx context.Context, registry, aliasKey string, req *aliastypes.AliasRequest) (*aliastypes.AliasResponse, error)
	DeleteAlias(ctx context.Context, registry, aliasKey string) error
}

type AliasRegistryClient interface {
	CreateRegistry(ctx context.Context, registry string, req *aliastypes.CreateRegistryRequest) (*aliastypes.RegistryResponse, error)
	GetRegistry(ctx context.Context, registry string) (*aliastypes.RegistryResponse, error)
	DeleteRegistry(ctx context.Context, registry string) error
}

type JSONRPC interface {
	Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error)
}

type KeyManagerClient interface {
	SecretsClient
	KeysClient
	EthClient
	UtilsClient
	AliasRegistryClient
	AliasClient
	JSONRPC
}
