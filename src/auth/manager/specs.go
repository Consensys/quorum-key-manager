package manager

import (
	"github.com/consensys/quorum-key-manager/src/auth/types"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
)

var RoleKind manifest.Kind = "Role"

type RoleSpecs struct {
	Permissions []types.Permission `json:"permission"`
}
