package authorizedaccounts

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/auth"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
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
func (s *store) Create(ctx context.Context, attr *entities.Attributes) (*entities.Account, error) {
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

func (s *store) Info(context.Context) *entities.StoreInfo {
	panic("implement me")
}

func (s *store) Import(ctx context.Context, privKey []byte, attr *entities.Attributes) (*entities.Account, error) {
	panic("implement me")
}

func (s *store) Get(ctx context.Context, addr string) (*entities.Account, error) {
	panic("implement me")
}

func (s *store) List(ctx context.Context, count uint, skip string) (accounts []*entities.Account, next string, err error) {
	panic("implement me")
}

func (s *store) Update(ctx context.Context, addr string, attr *entities.Attributes) (*entities.Account, error) {
	panic("implement me")
}

func (s *store) Delete(ctx context.Context, addrs ...string) (*entities.Account, error) {
	panic("implement me")
}

func (s *store) GetDeleted(ctx context.Context, addr string) {
	panic("implement me")
}

func (s *store) ListDeleted(ctx context.Context, count uint, skip string) (keys []*entities.Account, next string, err error) {
	panic("implement me")
}

func (s *store) Undelete(ctx context.Context, addr string) error {
	panic("implement me")
}

func (s *store) Destroy(ctx context.Context, addrs ...string) error {
	panic("implement me")
}

func (s *store) SignHomestead(ctx context.Context, addr string, tx *ethereum.Transaction) (sig []byte, err error) {
	panic("implement me")
}

func (s *store) SignEIP155(ctx context.Context, addr string, chainID string, tx *ethereum.Transaction) (sig []byte, err error) {
	panic("implement me")
}

func (s *store) SignEEA(ctx context.Context, addr string, chainID string, tx *ethereum.Transaction, args *ethereum.EEAPrivateArgs) (sig []byte, err error) {
	panic("implement me")
}

func (s *store) SignPrivate(ctx context.Context, addr string, tx *ethereum.Transaction) (sig []byte, err error) {
	panic("implement me")
}

func (s *store) ECRevocer(ctx context.Context, addr string, data []byte, sig []byte) (*entities.Account, error) {
	panic("implement me")
}

