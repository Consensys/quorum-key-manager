package eth1

import (
	"context"
	"math/big"

	quorumtypes "github.com/consensys/quorum/core/types"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

//go:generate mockgen -source=eth1.go -destination=mock/eth1.go -package=mock

type Store interface {
	// Info returns store information
	Info(context.Context) (*entities.StoreInfo, error)

	// Create creates an Ethereum account
	Create(ctx context.Context, id string, attr *entities.Attributes) (*entities.ETH1Account, error)

	// Import imports an externally created Ethereum account
	Import(ctx context.Context, id string, privKey []byte, attr *entities.Attributes) (*entities.ETH1Account, error)

	// Get gets an Ethereum account
	Get(ctx context.Context, addr string) (*entities.ETH1Account, error)

	// GetAll gets all Ethereum accounts
	GetAll(ctx context.Context) ([]*entities.ETH1Account, error)

	// List lists all Ethereum account addresses
	List(ctx context.Context) ([]string, error)

	// Update updates Ethereum account attributes
	Update(ctx context.Context, addr string, attr *entities.Attributes) (*entities.ETH1Account, error)

	// Delete deletes an account temporarily, by using Undelete the account can be restored
	Delete(ctx context.Context, addr string) error

	// GetDeleted Gets a deleted Ethereum accounts
	GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error)

	// ListDeleted lists all deleted Ethereum accounts
	ListDeleted(ctx context.Context) ([]string, error)

	// Undelete restores a previously deleted Ethereum account
	Undelete(ctx context.Context, addr string) error

	// Destroy destroys (purges) an Ethereum account permanently
	Destroy(ctx context.Context, addr string) error

	// Sign signs a payload using the specified Ethereum account
	Sign(ctx context.Context, addr string, data []byte) ([]byte, error)

	// Sign signs EIP-712 formatted data using the specified Ethereum account
	SignTypedData(ctx context.Context, addr string, typedData *core.TypedData) ([]byte, error)

	// SignTransaction signs a plublic Ethereum transaction
	SignTransaction(ctx context.Context, addr string, chainID *big.Int, tx *types.Transaction) ([]byte, error)

	// SignEEA signs an EEA transaction
	SignEEA(ctx context.Context, addr string, chainID *big.Int, tx *types.Transaction, args *ethereum.PrivateArgs) ([]byte, error)

	// SignPrivate signs a Quorum private transaction
	SignPrivate(ctx context.Context, addr string, tx *quorumtypes.Transaction) ([]byte, error)

	// ECRevocer returns the Ethereum address from a signature and data
	ECRevocer(ctx context.Context, data, sig []byte) (string, error)

	// Verify verifies that a signature belongs to a given address
	Verify(ctx context.Context, addr string, data, sig []byte) error

	// Verify verifies that a typed data signature belongs to a given address
	VerifyTypedData(ctx context.Context, addr string, sig []byte, typedData *core.TypedData) error

	// Encrypt encrypts any arbitrary data using a specified account
	Encrypt(ctx context.Context, addr string, data []byte) ([]byte, error)

	// Decrypt decrypts a single block of encrypted data.
	Decrypt(ctx context.Context, addr string, data []byte) ([]byte, error)
}
