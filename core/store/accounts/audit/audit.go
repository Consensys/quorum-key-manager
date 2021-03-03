package auditedaccount

import (
	"context"
	"fmt"
	"time"

	"github.com/ConsenSysQuorum/quorum-key-manager/core/audit"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/accounts"
	"github.com/ConsenSysQuorum/quorum-key-manager/core/store/types"
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
func (s *store) Create(ctx context.Context, attr *models.Attributes) (*models.Account, error) {
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
