package authorizator

import (
	"fmt"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
)

type Authorizator struct {
	logger      log.Logger
	permissions map[entities.Permission]bool // We use a map to avoid iterating an array, the boolean is irrelevant and always true
	tenant      string
}

var _ auth.Authorizator = &Authorizator{}

func New(permissions []entities.Permission, tenant string, logger log.Logger) *Authorizator {
	pMap := map[entities.Permission]bool{}
	for _, p := range permissions {
		pMap[p] = true
	}

	return &Authorizator{
		permissions: pMap,
		tenant:      tenant,
		logger:      logger,
	}
}

func (author *Authorizator) CheckPermission(ops ...*entities.Operation) error {
	for _, op := range ops {
		permission := buildPermission(op.Action, op.Resource)
		if _, ok := author.permissions[permission]; !ok {
			errMessage := "user is not authorized to perform this operation"
			author.logger.With("permission", permission).Error(errMessage)
			return errors.ForbiddenError(errMessage)
		}
	}

	return nil
}

func (author *Authorizator) CheckAccess(allowedTenants []string) error {
	if len(allowedTenants) == 0 {
		return nil
	}

	if author.tenant == "" {
		errMessage := "missing tenant in credentials"
		author.logger.Error(errMessage)
		return errors.UnauthorizedError(errMessage)
	}

	for _, t := range allowedTenants {
		if t == author.tenant {
			return nil
		}
	}

	errMessage := "resource not found"
	author.logger.With("tenant", author.tenant, "allowed_tenants", allowedTenants).Error(errMessage)
	return errors.NotFoundError(errMessage)
}

func buildPermission(action entities.OpAction, resource entities.OpResource) entities.Permission {
	return entities.Permission(fmt.Sprintf("%s:%s", action, resource))
}
