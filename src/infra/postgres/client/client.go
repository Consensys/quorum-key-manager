package client

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/consensys/quorum-key-manager/src/infra/postgres"
)

type PostgresClient struct {
	cfg *Config
	db  orm.DB
}

var _ postgres.Client = &PostgresClient{}

func NewClient(cfg *Config) (*PostgresClient, error) {
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

func (c *PostgresClient) UpdatePK(ctx context.Context, model ...interface{}) error {
	_, err := c.db.ModelContext(ctx, model...).WherePK().Update()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) DeletePK(ctx context.Context, model ...interface{}) error {
	_, err := c.db.ModelContext(ctx, model...).WherePK().Delete()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) ForceDeletePK(ctx context.Context, model ...interface{}) error {
	_, err := c.db.ModelContext(ctx, model...).WherePK().ForceDelete()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c PostgresClient) RunInTransaction(ctx context.Context, persist func(client postgres.Client) error) (err error) {
	persistFunc := func(tx *pg.Tx) error {
		c.db = tx
		return persist(&c)
	}

	// Check whether we already are in a tx or not to allow for nested DB transactions
	dbtx, isTx := c.db.(*pg.Tx)
	if isTx {
		return dbtx.RunInTransaction(ctx, persistFunc)
	}

	return c.db.(*pg.DB).RunInTransaction(ctx, persistFunc)
}

func (c *PostgresClient) SelectMany(ctx context.Context, model, dst interface{}, condition string, params ...interface{}) error {
	q := c.db.ModelContext(ctx, model)
	err := q.Where(condition, params...).Select(dst)
	if err != nil {
		return err
	}
	return nil
}

func (c *PostgresClient) CreateTable(ctx context.Context, model interface{}) error {
	return c.db.ModelContext(ctx, model).CreateTable(nil)
}

func (c *PostgresClient) DropTable(ctx context.Context, model interface{}) error {
	return c.db.ModelContext(ctx, model).DropTable(nil)
}
