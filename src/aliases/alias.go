package aliases

import "context"

type API interface {
	Aliaser
	Registrer
}

type AliasID string

type Alias struct {
	tableName struct{} `pg:"aliases" json:"-"`

	ID         AliasID
	RegistryID RegistryID
}

type Aliaser interface {
	CreateAlias(ctx context.Context, registry RegistryID, alias AliasID) error
	GetAlias(ctx context.Context, registry RegistryID, alias AliasID) (*Alias, error)
	UpdateAlias(ctx context.Context, registry RegistryID, alias AliasID) error
	DeleteAlias(ctx context.Context, registry RegistryID, alias AliasID) error
}

type RegistryID string

type Registry struct {
	ID RegistryID
}

type Registrer interface {
	GetRegistry(ctx context.Context, registry RegistryID) (*Registry, error)
	DeleteRegistry(ctx context.Context, registry RegistryID) error
}
