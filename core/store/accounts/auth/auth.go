package authorizedaccounts

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/auth"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
)

type Instrument struct{}

func NewInstrument() *Instrument {
	return &Instrument{}
}

func (i *Instrument) Apply(s accounts.Store) accounts.Store {
	return &store{
		accounts: s,
	}
}

// store instruments an account store with authorization capabilities
type store struct {
	accounts accounts.Store
}

// Create an account
func (s *store) Create(ctx context.Context, attr *models.Attributes) (*models.Account, error) {
	authInfo := auth.AuthFromContext(ctx)
	if err := authInfo.Policies.IsStoreAuthorized(s.accounts.Info(ctx)); err != nil {
		return nil, err
	}

	return s.accounts.Create(ctx, attr)
}

// Sign from a digest using the specified account
func (s *store) Sign(ctx context.Context, addr string, data []byte) (sig []byte, err error) {
	authInfo := auth.AuthFromContext(ctx)
	if err := authInfo.Policies.IsStoreAuthorized(s.accounts.Info(ctx)); err != nil {
		return nil, err
	}

	return s.accounts.Sign(ctx, addr, data)
}
