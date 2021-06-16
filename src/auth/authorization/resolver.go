package authorization

import (
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

type Operation struct {
	Action string
	Path   string
}

type Result struct {
	isAllowed, isRoot bool
	err               error
}

func (res *Result) Allowed() bool {
	return res.isAllowed
}

func (res *Result) IsRoot() bool {
	return res.isRoot
}

func (res *Result) Error() error {
	return res.err
}

type Resolver struct {
}

func NewResolver(policies []*types.Policy) (*Resolver, error) {
	return &Resolver{}, nil
}

func (r *Resolver) IsAuthorized(op *Operation) *Result {
	return &Result{
		isAllowed: true,
	}
}
