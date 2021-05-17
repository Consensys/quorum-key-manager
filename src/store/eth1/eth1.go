package eth1

import (
	"context"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

//go:generate mockgen -source=eth1.go -destination=mocks/eth1.go -package=mocks

type Store interface {
	// Info returns store information
	Info(context.Context) (*entities.StoreInfo, error)

	// Create an account
	Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETH1Account, error)

	// Import an externally created key and store account
	Import(ctx context.Context, id, privKey string, attr *entities.Attributes) (*entities.ETH1Account, error)

	// Get account
	Get(ctx context.Context, addr string) (*entities.ETH1Account, error)

	// List accounts
	List(ctx context.Context) ([]string, error)

	// Update account attributes
	Update(ctx context.Context, addr string, attr *entities.Attributes) (*entities.ETH1Account, error)

	// Delete account not permanently, by using Undelete the account can be retrieved
	Delete(ctx context.Context, addr string) error

	// GetDeleted accounts
	GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error)

	// ListDeleted accounts
	ListDeleted(ctx context.Context) ([]string, error)

	// Undelete a previously deleted account
	Undelete(ctx context.Context, addr string) error

	// Destroy account permanently
	Destroy(ctx context.Context, addr string) error

	// Sign from a digest using the specified account
	Sign(ctx context.Context, addr, data string) (string, error)

	// Sign EIP-712 formatted data using the specified account
	SignTypedData(ctx context.Context, addr string, typedData *core.TypedData) (string, error)

	// SignTransaction transaction
	SignTransaction(ctx context.Context, addr, chainID string, tx *types.Transaction) (string, error)

	// SignEEA transaction
	SignEEA(ctx context.Context, addr, chainID string, tx *ethereum.EEATxData, args *ethereum.PrivateArgs) (string, error)

	// SignPrivate transaction
	SignPrivate(ctx context.Context, addr string, tx *types.Transaction) (string, error)

	// ECRevocer returns the address from a signature and data
	ECRevocer(ctx context.Context, data, sig string) (string, error)

	// Verify verifies that a signature belongs to a given address
	Verify(ctx context.Context, addr, sig, payload string) error

	// Verify verifies that a typed data signature belongs to a given address
	VerifyTypedData(ctx context.Context, addr, sig string, typedData *core.TypedData) error

	// Encrypt any arbitrary data using a specified account
	Encrypt(ctx context.Context, addr, data string) (string, error)

	// Decrypt a single block of encrypted data.
	Decrypt(ctx context.Context, addr, data string) (string, error)
}
