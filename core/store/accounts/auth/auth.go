package authorizedaccounts

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/auth"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
)

// Store instruments an account store with authorization capabilities
type Store struct {
	accounts accounts.Store
}

// Create an account
func (s *Store) Create(ctx context.Context, attr *types.Attributes) (*types.Account, error) {
	authInfo := auth.AuthFromContext(ctx)
	if err := authInfo.Policies.IsStoreAuthorized(s.accounts.Info(ctx)); err != nil {
		return nil, err
	}

	return s.accounts.Create(ctx, attr)
}

// Sign from a digest using the specified account
func (s *Store) Sign(ctx context.Context, addr string, data []byte) (sig []byte, err error) {
	authInfo := auth.AuthFromContext(ctx)
	if err := authInfo.Policies.IsStoreAuthorized(s.accounts.Info(ctx)); err != nil {
		return nil, err
	}

	return s.accounts.Sign(ctx, addr, data)
}
