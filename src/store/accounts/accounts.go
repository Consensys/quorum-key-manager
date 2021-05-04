package accounts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

var (
	ErrorNotfound = fmt.Errorf("account not found")
)

type Store interface {
	// Info returns store information
	Info(context.Context) *entities.StoreInfo

	// Create an account
	Create(ctx context.Context, attr *entities.Attributes) (*entities.Account, error)

	// Import an externally created key and store account
	Import(ctx context.Context, privKey []byte, attr *entities.Attributes) (*entities.Account, error)

	// Get account
	Get(ctx context.Context, addr ethcommon.Address) (*entities.Account, error)

	// List accounts
	List(ctx context.Context, count uint, skip string) (accounts []*entities.Account, next string, err error)

	// Update account attributes
	Update(ctx context.Context, addr ethcommon.Address, attr *entities.Attributes) (*entities.Account, error)

	// Delete account not parmently, by using Undelete the account can be retrieve
	Delete(ctx context.Context, addrs ...string) (*entities.Account, error)

	// GetDeleted accounts
	GetDeleted(ctx context.Context, addr ethcommon.Address)

	// ListDeleted accounts
	ListDeleted(ctx context.Context, count uint, skip string) (keys []*entities.Account, next string, err error)

	// Undelete a previously deleted account
	Undelete(ctx context.Context, addr ethcommon.Address) error

	// Destroy account permanently
	Destroy(ctx context.Context, addrs ...string) error

	// Sign from a digest using the specified account
	Sign(ctx context.Context, addr ethcommon.Address, data []byte) (sig []byte, err error)

	// SignHomestead transaction
	SignHomestead(ctx context.Context, addr ethcommon.Address, tx *ethereum.TxData) (sig []byte, err error)

	// SignEIP155 transaction
	SignEIP155(ctx context.Context, chainID *big.Int, addr ethcommon.Address, tx *ethereum.TxData) (sig []byte, err error)

	// SignEEA transaction
	SignEEA(ctx context.Context, chainID *big.Int, addr ethcommon.Address, tx *ethereum.EEATxData, args *ethereum.PrivateArgs) (sig []byte, err error)

	// SignPrivate transaction
	SignPrivate(ctx context.Context, addr ethcommon.Address, tx *ethereum.TxData) (sig []byte, err error)

	// Verify a signature using a specified key
	ECRevocer(ctx context.Context, addr ethcommon.Address, data []byte, sig []byte) (*entities.Account, error)
}
