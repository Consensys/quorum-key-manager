package testdata

import (
	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

func FakeUserClaims() *entities.UserClaims {
	return &entities.UserClaims{
		Tenant:      "TenantOne|Alice",
		Permissions: "read:key write:key",
		Roles:       "guest admin",
	}
}
