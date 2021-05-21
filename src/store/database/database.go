package database

import "context"

//go:generate mockgen -source=database.go -destination=mock/database.go -package=mock

type Database interface {
	GetID(ctx context.Context, addr string) (string, error)
	GetDeletedID(ctx context.Context, addr string) (string, error)
	GetAll(ctx context.Context) ([]string, error)
	GetAllDeleted(ctx context.Context) ([]string, error)
	GetAllIDs(ctx context.Context) ([]string, error)
	AddID(ctx context.Context, addr, id string) error
	AddDeletedID(ctx context.Context, addr, id string) error
	RemoveID(ctx context.Context, addr string) error
	RemoveDeletedID(ctx context.Context, addr string) error
}
