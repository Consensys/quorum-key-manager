package client

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/go-pg/pg/v10"
)

func parseErrorResponse(err error) error {
	if pg.ErrNoRows == err {
		return errors.NotFoundError("resource not found")
	}
	if pg.ErrMultiRows == err {
		return errors.StatusConflictError("multiple resources found, only expected one")
	}

	pgErr, ok := err.(pg.Error)
	if !ok {
		return errors.PostgresError(err.Error())
	}

	switch {
	case pgErr.IntegrityViolation():
		return errors.StatusConflictError(pgErr.Error())
	default:
		return errors.PostgresError(pgErr.Error())
	}
}
