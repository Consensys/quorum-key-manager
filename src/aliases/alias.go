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
	Kind       AliasKind
	Value      AliasValue
}

type AliasValue string

type AliasKind string

const (
	AliasKindUnknown = ""
	AliasKindString  = "string"
	AliasKindArray   = "array"
)

type Aliaser interface {
	CreateAlias(ctx context.Context, alias Alias) error
	GetAlias(ctx context.Context, registry RegistryID, alias AliasID) (*Alias, error)
	UpdateAlias(ctx context.Context, alias Alias) error
	DeleteAlias(ctx context.Context, registry RegistryID, alias AliasID) error

	ListAliases(ctx context.Context, registry RegistryID) ([]Alias, error)
}

type RegistryID string

type Registry struct {
	tableName struct{} `pg:"registries" json:"-"`

	ID RegistryID
}

type Registrer interface {
	DeleteRegistry(ctx context.Context, registry RegistryID) error
}
