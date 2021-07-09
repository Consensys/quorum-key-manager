package client

import (
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/go-pg/pg/v10"
)

type PostgresClient struct {
	cfg *Config
	db  *pg.DB
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

func (c *PostgresClient) Insert(model ...interface{}) error {
	_, err := c.db.Model(model...).Insert()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) SelectPK(model ...interface{}) error {
	err := c.db.Model(model...).WherePK().Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) SelectDeletedPK(model ...interface{}) error {
	err := c.db.Model(model...).WherePK().Deleted().Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) Select(model ...interface{}) error {
	err := c.db.Model(model...).Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) SelectDeleted(model ...interface{}) error {
	err := c.db.Model(model...).Deleted().Select()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) UpdatePK(model ...interface{}) error {
	_, err := c.db.Model(model...).WherePK().Update()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) DeletePK(model ...interface{}) error {
	_, err := c.db.Model(model...).WherePK().Delete()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}

func (c *PostgresClient) ForceDeletePK(model ...interface{}) error {
	_, err := c.db.Model(model...).WherePK().ForceDelete()
	if err != nil {
		return parseErrorResponse(err)
	}

	return nil
}
