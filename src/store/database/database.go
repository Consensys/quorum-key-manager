package database

import (
	"context"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	ETH1Accounts() ETH1Accounts
}

type ETH1Accounts interface {
	GetAccount(ctx context.Context, addr string) (*entities.ETH1Account, error)
	GetDeletedAccount(ctx context.Context, addr string) (*entities.ETH1Account, error)
	GetAllAccounts(ctx context.Context) ([]*entities.ETH1Account, error)
	GetAllDeletedAccounts(ctx context.Context) ([]*entities.ETH1Account, error)
	AddAccount(ctx context.Context, account *entities.ETH1Account) error
	AddDeletedAccount(ctx context.Context, account *entities.ETH1Account) error
	RemoveAccount(ctx context.Context, addr string) error
	RemoveDeletedAccount(ctx context.Context, addr string) error
}
