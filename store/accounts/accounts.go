package accounts

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/store/types"
)

type Store interface {
	// Info returns store information
	Info(context.Context) *types.StoreInfo

	// Create an account
	Create(ctx context.Context, attr *types.Attributes) (*types.Account, error)

	// Import an externally created key and store account
	Import(ctx context.Context, privKey []byte, attr *types.Attributes) (*types.Account, error)

	// Get account
	Get(ctx context.Context, addr string) (*types.Account, error)

	// List accounts
	List(ctx context.Context, count uint, skip string) (accounts []*types.Account, next string, err error)

	// Update account attributes
	Update(ctx context.Context, addr string, attr *types.Attributes) (*types.Account, error)

	// Delete account not parmently, by using Undelete the account can be retrieve
	Delete(ctx context.Context, addrs ...string) (*types.Account, error)

	// GetDeleted accounts
	GetDeleted(ctx context.Context, addr string)

	// ListDeleted accounts
	ListDeleted(ctx context.Context, count uint, skip string) (keys []*types.Account, next string, err error)

	// Undelete a previously deleted account
	Undelete(ctx context.Context, addr string) error

	// Destroy account permanenty
	Destroy(ctx context.Context, addrs ...string) error

	// Sign from a digest using the specified account
	Sign(ctx context.Context, addr string, data []byte) (sig []byte, err error)

	// SignHomestead transaction
	SignHomestead(ctx context.Context, addr string, tx *ethereum.Transaction) (sig []byte, err error)

	// SignEIP155 transaction
	SignEIP155(ctx context.Context, addr string, chainID string, tx *ethereum.Transaction) (sig []byte, err error)

	// SignEEA transaction
	SignEEA(ctx context.Context, addr string, chainID string, tx *ethereum.Transaction, args *ethereum.EEAPrivateArgs) (sig []byte, err error)

	// SignPrivate transaction
	SignPrivate(ctx context.Context, addr string, tx *ethereum.Transaction) (sig []byte, err error)

	// Verify a signature using a specified key
	ECRevocer(ctx context.Context, addr string, data []byte, sig []byte) (*types.Account, error)
}
