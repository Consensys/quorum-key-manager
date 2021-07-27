package postgres

import (
	"context"
	"github.com/go-pg/pg/v10/orm"
)

//go:generate mockgen -source=postgres.go -destination=mocks/postgres.go -package=mocks

type Client interface {
	ModelContext(ctx context.Context, model ...interface{}) *orm.Query
	Insert(ctx context.Context, model ...interface{}) error
	SelectPK(ctx context.Context, model ...interface{}) error
	SelectDeletedPK(ctx context.Context, model ...interface{}) error
	SelectQuery(query *orm.Query) error
	Select(ctx context.Context, model ...interface{}) error
	SelectDeleted(ctx context.Context, model ...interface{}) error
	UpdatePK(ctx context.Context, model ...interface{}) error
	DeletePK(ctx context.Context, model ...interface{}) error
	ForceDeletePK(ctx context.Context, model ...interface{}) error
	RunInTransaction(ctx context.Context, persist func(client Client) error) error
}
