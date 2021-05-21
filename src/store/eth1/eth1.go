package eth1

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

//go:generate mockgen -source=eth1.go -destination=mock/eth1.go -package=mock

type Store interface {
	// Info returns store information
	Info(context.Context) (*entities.StoreInfo, error)

	// Create an account
	Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETH1Account, error)

	// Import an externally created key and store account
	Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETH1Account, error)

	// Get account
	Get(ctx context.Context, addr string) (*entities.ETH1Account, error)

	// List account ids
	List(ctx context.Context) ([]string, error)

	// GetAll accounts
	GetAll(ctx context.Context) ([]*entities.ETH1Account, error)

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
	Sign(ctx context.Context, addr string, data []byte) ([]byte, error)

	// Sign EIP-712 formatted data using the specified account
	SignTypedData(ctx context.Context, addr string, typedData *core.TypedData) ([]byte, error)

	// SignTransaction transaction
	SignTransaction(ctx context.Context, addr string, chainID *big.Int, tx *types.Transaction) ([]byte, error)

	// SignEEA transaction
	SignEEA(ctx context.Context, addr string, chainID *big.Int, tx *ethereum.EEATxData, args *ethereum.PrivateArgs) ([]byte, error)

	// SignPrivate transaction
	SignPrivate(ctx context.Context, addr string, tx *types.Transaction) ([]byte, error)

	// ECRevocer returns the address from a signature and data
	ECRevocer(ctx context.Context, data, sig []byte) (string, error)

	// Verify verifies that a signature belongs to a given address
	Verify(ctx context.Context, addr string, data, sig []byte) error

	// Verify verifies that a typed data signature belongs to a given address
	VerifyTypedData(ctx context.Context, addr string, sig []byte, typedData *core.TypedData) error

	// Encrypt any arbitrary data using a specified account
	Encrypt(ctx context.Context, addr string, data []byte) ([]byte, error)

	// Decrypt a single block of encrypted data.
	Decrypt(ctx context.Context, addr string, data []byte) ([]byte, error)
}
