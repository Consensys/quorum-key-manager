package database

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	ETH1Accounts() ETH1Accounts
	Keys() Keys
}

type ETH1Accounts interface {
	Get(ctx context.Context, addr string) (*entities.ETH1Account, error)
	GetDeleted(ctx context.Context, addr string) (*entities.ETH1Account, error)
	GetAll(ctx context.Context) ([]*entities.ETH1Account, error)
	GetAllDeleted(ctx context.Context) ([]*entities.ETH1Account, error)
	Add(ctx context.Context, account *entities.ETH1Account) error
	Update(ctx context.Context, account *entities.ETH1Account) error
	AddDeleted(ctx context.Context, account *entities.ETH1Account) error
	Remove(ctx context.Context, addr string) error
	RemoveDeleted(ctx context.Context, addr string) error
}

type Keys interface {
	Get(ctx context.Context, id string) (*entities.Key, error)
	GetDeleted(ctx context.Context, id string) (*entities.Key, error)
	GetAll(ctx context.Context) ([]*entities.Key, error)
	GetAllDeleted(ctx context.Context) ([]*entities.Key, error)
	Add(ctx context.Context, key *entities.Key) error
	Update(ctx context.Context, key *entities.Key) error
	AddDeleted(ctx context.Context, key *entities.Key) error
	Remove(ctx context.Context, addr string) error
	RemoveDeleted(ctx context.Context, addr string) error
}
