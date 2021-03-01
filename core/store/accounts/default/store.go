package defaultaccounts

import (
	"context"
	"fmt"

	"github.com/ConsenSysQuorum/quorum-key-manager/common/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/keys"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
)

var NotImplementedError = fmt.Errorf("not implemented")

// [DRAFT] Store is an accounts store that relyies on an underlying secp256k1 compatible keys store
type Store struct {
	keys     keys.Store
	addrToID map[string]string
}

// NewStore creates an account store
func NewStore(keys keys.Store) *Store {
	return &Store{
		keys: keys,
	}
}

// Create an account
func (s *Store) Create(ctx context.Context, attr *types.Attributes) (*types.Account, error) {
	// TODO: to be implemented
	// TODO: problem account address can be computed only after key has been created in key store thus
	// it is not possible to index a key by account address making it impossible to query an account by address on keystore.
	return nil, NotImplementedError
}

// Import an externally created key and store account
func (s *Store) Import(ctx context.Context, privKey []byte, attr *types.Attributes) (*types.Account, error) {
	// TODO: to be implemented
	return nil, NotImplementedError
}

// Get account
func (s *Store) Get(ctx context.Context, addr string) (*types.Account, error) {
	id, err := s.getId(addr)
	if err != nil {
		return nil, err
	}

	key, err := s.keys.Get(ctx, id, 0)
	if err != nil {
		return nil, s.handleError(err, addr)
	}

	return types.KeyToAccount(key), nil
}

// List accounts
func (s *Store) List(ctx context.Context, count uint, skip string) (accounts []*types.Account, next string, err error) {
	keys, next, err := s.keys.List(ctx, count, skip)
	if err != nil {
		return nil, "", s.handleError(err, "")
	}

	return types.KeysToAccounts(keys), next, nil
}

// Update account attributes
func (s *Store) Update(ctx context.Context, addr string, attr *types.Attributes) (*types.Account, error) {
	id, err := s.getId(addr)
	if err != nil {
		return nil, err
	}

	key, err := s.keys.Update(ctx, id, attr)
	if err != nil {
		return nil, s.handleError(err, addr)
	}

	return types.KeyToAccount(key), nil
}

// Sign from a digest using the specified account
func (s *Store) Sign(ctx context.Context, addr string, data []byte) (sig []byte, err error) {
	id, err := s.getId(addr)
	if err != nil {
		return nil, err
	}

	sig, err = s.keys.Sign(ctx, id, data)
	if err != nil {
		return nil, s.handleError(err, addr)
	}

	return
}

// SignHomestead transaction
func (s *Store) SignHomestead(ctx context.Context, addr string, tx *ethereum.Transaction) (sig []byte, err error) {
	id, err := s.getId(addr)
	if err != nil {
		return nil, err
	}

	sig, err = ethereum.HomesteadSign(tx, func(data []byte) ([]byte, error) { return s.keys.Sign(ctx, id, data) })
	if err != nil {
		return nil, s.handleError(err, addr)
	}

	return
}

// TODO: implement all Store methods

func (s *Store) getId(addr string) (string, error) {
	// TODO: to be implemented
	return s.addrToID[addr], nil
}

func (s *Store) handleError(err error, addr string) error {
	// TODO: to be implemented
	return err
}
