package noop

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/audit"
)

// Auditor is a no operation auditor
type Auditor struct{}

func New() *Auditor {
	return &Auditor{}
}

// StartOperation synchronously persists a starting operation, this is done after authentication succeeded but before the operation is executed
func (a *Auditor) StartOperation(_ context.Context, _ *audit.Operation) error {
	return nil
}

// EndOperation persists an ending operation this is done after the operation ended but before returning
func (a *Auditor) EndOperation(_ context.Context, _ *audit.Operation) error {
	return nil
}

// Get return an operation by id
func (a *Auditor) Get(_ context.Context, _ string) (*audit.Operation, error) {
	return nil, nil
}

// Lookup queries operations
func (a *Auditor) Lookup(_ context.Context, _ *audit.Selector, _ int, _ int) ([]*audit.Operation, error) {
	return nil, nil
}
