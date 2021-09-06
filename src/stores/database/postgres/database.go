package postgres

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/consensys/quorum-key-manager/src/stores/database"
)

type Database struct {
	logger log.Logger
	client postgres.Client
}

var _ database.Database = &Database{}

func New(logger log.Logger, client postgres.Client) *Database {
	return &Database{
		logger: logger,
		client: client,
	}
}

func (db *Database) Ping(ctx context.Context) error {
	return db.client.Ping(ctx)
}

func (db *Database) ETH1Accounts(storeID string) database.ETH1Accounts {
	return NewETH1Accounts(storeID, db.client, db.logger.With("store_id", storeID))
}

func (db *Database) Keys(storeID string) database.Keys {
	return NewKeys(storeID, db.client, db.logger.With("store_id", storeID))
}

func (db *Database) Secrets(storeID string) database.Secrets {
	return NewSecrets(storeID, db.client, db.logger.With("store_id", storeID))
}
