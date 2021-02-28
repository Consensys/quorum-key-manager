package audit

import (
	"context"
)

// Auditor allows to audit
type Auditor interface {
	// StartOperation synchronously persists a starting operation, this is done after authentication succeeded but before the operation is executed
	StartOperation(context.Context, *Operation) error

	// EndOperation persists an ending operation this is done after the operation ended but before returning
	EndOperation(context.Context, *Operation) error

	// Get return an operation by id
	Get(ctx context.Context, id string) (*Operation, error)

	// Lookup queries operations

	// There should
	Lookup(ctx context.Context, sel *Selector, count int, skip int) ([]*Operation, error)
}
