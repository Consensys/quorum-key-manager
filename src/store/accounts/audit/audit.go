package auditedaccount

import (
	"context"
	"fmt"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/audit"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

var opPrefix = "accounts."

type Instrument struct {
	auditor audit.Auditor
}

func NewInstrument(auditor audit.Auditor) *Instrument {
	return &Instrument{
		auditor: auditor,
	}
}

func (i *Instrument) Apply(s accounts.Store) accounts.Store {
	return &store{
		accounts: s,
		auditor:  i.auditor,
	}
}

// [DRAFT] store instruments an account store with audit capabilities
type store struct {
	accounts accounts.Store
	auditor  audit.Auditor
}

// Create an account
func (s *store) Create(ctx context.Context, attr *entities.Attributes) (*entities.Account, error) {
	// create operation object
	// TODO: Can probably be improved by relying extracting already existing operation from context
	// TODO: Auth should be extracted from context
	op := &audit.Operation{
		Type:      fmt.Sprintf("%v.create", opPrefix),
		StartTime: time.Now(),
		Data: map[string]interface{}{
			"attr": attr,
		},
	}

	// audit operation start
	// TODO: what to do in case of auditing error?
	_ = s.auditor.StartOperation(ctx, op)

	// execute operation
	account, err := s.accounts.Create(ctx, attr)

	// enrich operation data with results
	op.EndTime = time.Now()
	op.Data["account"] = account
	op.Error = err

	// audit operation end
	// TODO: what to do in case of auditing error?
	_ = s.auditor.EndOperation(ctx, op)

	return account, err
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

func (s *store) Sign(ctx context.Context, addr string, data []byte) (sig []byte, err error) {
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
