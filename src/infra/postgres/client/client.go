package client

import (
	"context"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/lib/pq"

	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

type PostgresClient struct {
	cfg *Config
	db  orm.DB
}

var _ postgres.Client = &PostgresClient{}

func New(cfg *Config) (*PostgresClient, error) {
	pgOptions, err := cfg.ToPGOptions()
	if err != nil {
		return nil, err
	}

	db := pg.Connect(pgOptions)

	return &PostgresClient{
		cfg: cfg,
		db:  db,
	}, nil
}

func (c *PostgresClient) Config() *Config {
	return c.cfg
}

func (c *PostgresClient) QueryOne(ctx context.Context, result, query interface{}, params ...interface{}) error {
	_, err := c.db.QueryOneContext(ctx, pg.Scan(result), query, params...)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) Query(ctx context.Context, result, query interface{}, params ...interface{}) error {
	_, err := c.db.QueryContext(ctx, pg.Scan(pq.Array(result)), query, params...)
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) Insert(ctx context.Context, model ...interface{}) error {
	_, err := c.db.ModelContext(ctx, model...).Insert()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) SelectPK(ctx context.Context, model ...interface{}) error {
	err := c.db.ModelContext(ctx, model...).WherePK().Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) SelectDeletedPK(ctx context.Context, model ...interface{}) error {
	err := c.db.ModelContext(ctx, model...).WherePK().Deleted().Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) Select(ctx context.Context, model ...interface{}) error {
	err := c.db.ModelContext(ctx, model...).Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) SelectDeleted(ctx context.Context, model ...interface{}) error {
	err := c.db.ModelContext(ctx, model...).Deleted().Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) SelectDeletedWhere(ctx context.Context, model interface{}, where string, args ...interface{}) error {
	err := c.db.ModelContext(ctx, model).Deleted().Where(where, args...).Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) SelectWhere(ctx context.Context, model interface{}, where string, args ...interface{}) error {
	err := c.db.ModelContext(ctx, model).Where(where, args...).Select()
	if err != nil {
		return parseErrorResponse(err)
	}
	return nil
}

func (c *PostgresClient) UpdatePK(ctx context.Context, model interface{}) error {
	q := c.db.ModelContext(ctx, model)
	r, err := q.WherePK().UpdateNotZero()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no rows were updated")
	}

	return nil
}

func (c *PostgresClient) UpdateWhere(ctx context.Context, model interface{}, where string, params ...interface{}) error {
	q := c.db.ModelContext(ctx, model)
	r, err := q.Where(where, params...).UpdateNotZero()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no matched rows were updated")
	}

	return nil
}

func (c *PostgresClient) Delete(ctx context.Context, model ...interface{}) error {
	r, err := c.db.ModelContext(ctx, model...).Delete()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no rows were deleted")
	}

	return nil
}

func (c *PostgresClient) DeletePK(ctx context.Context, model ...interface{}) error {
	r, err := c.db.ModelContext(ctx, model...).WherePK().Delete()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no rows were deleted with PK")
	}

	return nil
}

func (c *PostgresClient) DeleteWhere(ctx context.Context, model interface{}, where string, params ...interface{}) error {
	r, err := c.db.ModelContext(ctx, model).Where(where, params...).Delete()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no matched rows were deleted using where condition")
	}

	return nil
}

func (c *PostgresClient) UndeletePK(ctx context.Context, model ...interface{}) error {
	q := c.db.ModelContext(ctx, model...).AllWithDeleted().WherePK()
	if q.TableModel().Table().SoftDeleteField == nil {
		return errors.PostgresError("models does not support soft-delete")
	}

	r, err := q.Set("? = ?", q.TableModel().Table().SoftDeleteField.Column, nil).UpdateNotZero()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no rows were undeleted using PK")
	}

	return nil
}

func (c *PostgresClient) UndeleteWhere(ctx context.Context, model interface{}, where string, params ...interface{}) error {
	q := c.db.ModelContext(ctx, model).AllWithDeleted().Where(where, params...)
	if q.TableModel().Table().SoftDeleteField == nil {
		return errors.PostgresError("model does not support soft-delete")
	}

	r, err := q.Set("? = ?", q.TableModel().Table().SoftDeleteField.Column, nil).UpdateNotZero()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no rows were undeleted using where condition")
	}

	return nil
}

func (c *PostgresClient) ForceDeletePK(ctx context.Context, model ...interface{}) error {
	r, err := c.db.ModelContext(ctx, model...).WherePK().ForceDelete()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no rows were force deleted")
	}

	return nil
}

func (c *PostgresClient) ForceDeleteWhere(ctx context.Context, model interface{}, where string, params ...interface{}) error {
	r, err := c.db.ModelContext(ctx, model).Where(where, params...).ForceDelete()
	if err != nil {
		return parseErrorResponse(err)
	}
	if r.RowsAffected() == 0 {
		return errors.NotFoundError("no matched  rows were force")
	}

	return nil
}

func (c PostgresClient) RunInTransaction(ctx context.Context, persist func(client postgres.Client) error) (err error) {
	persistFunc := func(tx *pg.Tx) error {
		c.db = tx
		return persist(&c)
	}

	// CheckPermission whether we already are in a tx or not to allow for nested DB transactions
	dbtx, isTx := c.db.(*pg.Tx)
	if isTx {
		return dbtx.RunInTransaction(ctx, persistFunc)
	}

	return c.db.(*pg.DB).RunInTransaction(ctx, persistFunc)
}

func (c PostgresClient) Ping(ctx context.Context) error {
	_, err := c.db.ExecContext(ctx, "SELECT 1")
	return err
}
