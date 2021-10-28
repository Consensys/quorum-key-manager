package testutils

import (
	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

func FakeUserClaims() *entities.UserClaims {
	return &entities.UserClaims{
		Subject: "TenantOne|Alice",
		Scope:   "read:key write:key",
		Roles:   "guest admin",
	}
}
