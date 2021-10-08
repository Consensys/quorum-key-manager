package manager

import (
	"github.com/consensys/quorum-key-manager/src/auth/types"
)

type RoleSpecs struct {
	Permissions []types.Permission `json:"permission"`
}
