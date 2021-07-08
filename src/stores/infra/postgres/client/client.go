package client

import (
	"github.com/consensys/quorum-key-manager/src/stores/infra/postgres"
)

type PostgresClient struct {
	cfg *Config
}

var _ postgres.Client = &PostgresClient{}

func NewClient(cfg *Config) (*PostgresClient, error) {
	return &PostgresClient{
		cfg: cfg,
	}, nil
}
