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

// [DRAFT] Store wraps an account store and make its auditable
type Store struct {
	accountsStore accounts.Store
	auditor       audit.Auditor
}

// Create an account
func (s *Store) Create(ctx context.Context, attr *types.Attributes) (*types.Account, error) {
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
	account, err := s.accountsStore.Create(ctx, attr)

	// enrich operation data with results
	op.EndTime = time.Now()
	op.Data["account"] = account
	op.Error = err

	// audit operation end
	// TODO: what to do in case of auditing error?
	_ = s.auditor.EndOperation(ctx, op)

	return account, err
}
