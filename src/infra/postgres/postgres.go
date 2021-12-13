package postgres

import (
	"context"
)

//go:generate mockgen -source=postgres.go -destination=mocks/postgres.go -package=mocks

type Client interface {
	QueryOne(ctx context.Context, result, query interface{}, params ...interface{}) error
	Query(ctx context.Context, result, query interface{}, params ...interface{}) error
	Insert(ctx context.Context, model ...interface{}) error
	SelectPK(ctx context.Context, model ...interface{}) error
	SelectDeletedPK(ctx context.Context, model ...interface{}) error
	Select(ctx context.Context, model ...interface{}) error
	SelectDeleted(ctx context.Context, model ...interface{}) error
	SelectWhere(ctx context.Context, model interface{}, where string, relations []string, args ...interface{}) error
	SelectDeletedWhere(ctx context.Context, model interface{}, where string, args ...interface{}) error
	UpdatePK(ctx context.Context, model interface{}) error
	UpdateWhere(ctx context.Context, model interface{}, where string, params ...interface{}) error
	DeletePK(ctx context.Context, model ...interface{}) error
	DeleteWhere(ctx context.Context, model interface{}, where string, params ...interface{}) error
	UndeletePK(ctx context.Context, model ...interface{}) error
	UndeleteWhere(ctx context.Context, model interface{}, where string, params ...interface{}) error
	ForceDeletePK(ctx context.Context, model ...interface{}) error
	ForceDeleteWhere(ctx context.Context, model interface{}, where string, params ...interface{}) error
	RunInTransaction(ctx context.Context, persist func(client Client) error) error
	Ping(ctx context.Context) error
}
