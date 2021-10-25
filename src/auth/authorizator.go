package auth

import "github.com/consensys/quorum-key-manager/src/auth/entities"

//go:generate mockgen -source=authorizator.go -destination=mock/authorizator.go -package=mock

// Authorizator allows managing authorizations given a set of permissions
type Authorizator interface {
	CheckPermission(ops ...*entities.Operation) error
	CheckAccess(allowedTenants []string) error
}
