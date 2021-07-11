package policy

import (
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

type Operation struct {
	Action string
	Path   string
}

type Result struct {
	isAllowed bool
	err       error
}

func (res *Result) Allowed() bool {
	return res.isAllowed
}

func (res *Result) Error() error {
	return res.err
}

// Resolver is responsible to control whether an operation is authorized or not
// depending on the set of policies attached to the resolver
type Resolver struct {
	// TODO: to be implemented
}

func NewResolver(policies []types.Policy) (*Resolver, error) {
	return &Resolver{}, nil
}

// IsAuthorized controls whether an operation is authorized or not
func (r *Resolver) IsAuthorized(op ...*Operation) *Result {
	// TODO: to be implemented
	return &Result{
		isAllowed: true,
	}
}
