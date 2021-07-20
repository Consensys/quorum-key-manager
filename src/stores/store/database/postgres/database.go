package postgres

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/postgres"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
)

type Database struct {
	logger       log.Logger
	client       postgres.Client
	eth1Accounts *ETH1Accounts
	keys         *Keys
}

var _ database.Database = &Database{}

func New(logger log.Logger, client postgres.Client) *Database {
	return &Database{
		logger:       logger,
		client:       client,
		eth1Accounts: NewETH1Accounts(logger), // TODO: Implement ETH1Accounts using Postgres and not in-memory
		keys:         NewKeys(logger, client),
	}
}

func (db *Database) ETH1Accounts() database.ETH1Accounts {
	return db.eth1Accounts
}

func (db *Database) Keys() database.Keys {
	return db.keys
}

func (db Database) RunInTransaction(ctx context.Context, persist func(dbtx database.Database) error) error {
	return db.client.RunInTransaction(ctx, func(newClient postgres.Client) error {
		db.client = newClient
		db.keys = NewKeys(db.logger, newClient)
		// TODO: Pass db.client to Eth1Accounts
		db.eth1Accounts = NewETH1Accounts(db.logger)

		return persist(&db)
	})
}
