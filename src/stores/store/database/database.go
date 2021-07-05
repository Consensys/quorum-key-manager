package database

import (
	"context"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
)

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	ETH1Accounts() ETH1Accounts
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
