package aliases

import "context"

//TODO the: rename, API is normally all the public exposed HTTP services
type Backend interface {
	Aliaser
	Registrer
}

type AliasKey string

type Alias struct {
	tableName struct{} `pg:"aliases"` // nolint:unused,structcheck // reason

	Key          AliasKey     `pg:",pk"`
	RegistryName RegistryName `pg:",pk"`
	Kind         AliasKind
	Value        AliasValue
}

type AliasValue string

type AliasKind string

const (
	AliasKindUnknown = ""
	AliasKindString  = "string"
	AliasKindArray   = "array"
)

type Aliaser interface {
	CreateAlias(ctx context.Context, registry RegistryName, alias Alias) error
	GetAlias(ctx context.Context, registry RegistryName, aliasKey AliasKey) (*Alias, error)
	UpdateAlias(ctx context.Context, registry RegistryName, alias Alias) error
	DeleteAlias(ctx context.Context, registry RegistryName, aliasKey AliasKey) error

	ListAliases(ctx context.Context, registry RegistryName) ([]Alias, error)

	DeleteRegistry(ctx context.Context, registry RegistryName) error
}

type RegistryName string

type Registry struct {
	tableName struct{} `pg:"registries"` // nolint:unused,structcheck // reason

	Name RegistryName `pg:",pk"`
}

type Registrer interface {
	DeleteRegistry(ctx context.Context, registry RegistryName) error
}
