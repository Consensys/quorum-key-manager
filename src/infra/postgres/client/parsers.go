package client

import (
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/go-pg/pg/v10"
)

func parseErrorResponse(err error) error {
	pgErr, ok := err.(pg.Error)
	if !ok {
		return errors.PostgresError(err.Error())
	}

	switch {
	case pg.ErrNoRows == err:
		return errors.NotFoundError(err.Error())
	case pgErr.IntegrityViolation():
		return errors.StatusConflictError(pgErr.Error())
	default:
		return errors.PostgresError(pgErr.Error())
	}
}
