package manager

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

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
	permissions map[types.Permission]bool
}

func NewResolver(permissions []types.Permission) *Resolver {
	pMap := map[types.Permission]bool{}
	for _, p := range permissions {
		pMap[p] = true
	}

	return &Resolver{
		permissions: pMap,
	}
}

// IsAuthorized controls whether an operation is authorized or not
func (r *Resolver) IsAuthorized(ops ...*Operation) *Result {
	for _, op := range ops {
		reqP := buildPermission(op.Action, op.Resource)
		if _, ok := r.permissions[reqP]; !ok {
			return &Result{
				isAllowed: false,
				err:       errors.UnauthorizedError("unauthorized operation. Required permission %s", reqP),
			}
		}
	}

	return &Result{
		isAllowed: true,
	}
}

func buildPermission(action OpAction, resource OpResource) types.Permission {
	return types.Permission(fmt.Sprintf("%s:%s", action, resource))
}
