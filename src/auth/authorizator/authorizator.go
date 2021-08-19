package authorizator

import (
	"fmt"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/auth/types"
)

type Authorizator struct {
	logger      log.Logger
	permissions map[types.Permission]bool // We use a map to avoid iterating an array, the boolean is irrelevant and always true
}

func New(userInfo *types.UserInfo, logger log.Logger) *Authorizator {
	pMap := map[types.Permission]bool{}
	for _, p := range userInfo.Permissions {
		pMap[p] = true
	}

	return &Authorizator{
		permissions: pMap,
		logger:      logger,
	}
}

// Check checks whether an operation is authorized or not
func (auth *Authorizator) Check(ops ...*types.Operation) error {
	for _, op := range ops {
		permission := buildPermission(op.Action, op.Resource)
		if _, ok := auth.permissions[permission]; !ok {
			errMessage := "user is not authorized to perform this operation"
			auth.logger.With("permission", permission).Error(errMessage)
			return errors.ForbiddenError(errMessage)
		}
	}

	return nil
}

func buildPermission(action types.OpAction, resource types.OpResource) types.Permission {
	return types.Permission(fmt.Sprintf("%s:%s", action, resource))
}
