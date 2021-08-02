package postgres

import (
	"context"
)

//go:generate mockgen -source=postgres.go -destination=mocks/postgres.go -package=mocks

type Client interface {
	Insert(ctx context.Context, model ...interface{}) error
	SelectPK(ctx context.Context, model ...interface{}) error
	SelectDeletedPK(ctx context.Context, model ...interface{}) error
	Select(ctx context.Context, model ...interface{}) error
	SelectDeleted(ctx context.Context, model ...interface{}) error
	SelectWhere(ctx context.Context, model interface{}, condition string, params ...interface{}) error
	UpdatePK(ctx context.Context, model ...interface{}) error
	DeletePK(ctx context.Context, model ...interface{}) error
	ForceDeletePK(ctx context.Context, model ...interface{}) error
	RunInTransaction(ctx context.Context, persist func(client Client) error) error
}
