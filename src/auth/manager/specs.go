package manager

import (
	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

type RoleSpecs struct {
	Permissions []entities.Permission `json:"permission"`
}
