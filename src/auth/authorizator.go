package auth

import "github.com/consensys/quorum-key-manager/src/auth/types"

//go:generate mockgen -source=authorizator.go -destination=mock/authorizator.go -package=mock

// Authorizator allows managing authorizations given a set of permissions
type Authorizator interface {
	CheckPermission(ops ...*types.Operation) error
	CheckAccess(allowedTenants []string) error
}
