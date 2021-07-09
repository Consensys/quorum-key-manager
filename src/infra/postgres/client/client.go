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
