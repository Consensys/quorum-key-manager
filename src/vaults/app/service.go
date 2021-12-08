package app

import (
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/vaults/service/vaults"
)

func RegisterService(logger log.Logger, roles auth.Roles) *vaults.Vaults {
	// Business layer
	return vaults.New(roles, logger)
}
