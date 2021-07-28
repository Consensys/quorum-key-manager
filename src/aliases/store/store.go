package aliasstore

import "context"

//go:generate mockgen -source=store.go -destination=mock/store.go -package=mock

// Store handles the alias storing.
type Store interface {
	// CreateAlias creates an alias in the registry.
	CreateAlias(ctx context.Context, registry RegistryName, alias Alias) error
	// GetAlias gets an alias from the registry.
	GetAlias(ctx context.Context, registry RegistryName, aliasKey AliasKey) (*Alias, error)
	// UpdateAlias updates an alias in the registry.
	UpdateAlias(ctx context.Context, registry RegistryName, alias Alias) error
	// GetAlias deletes an alias from the registry.
	DeleteAlias(ctx context.Context, registry RegistryName, aliasKey AliasKey) error

	// ListAlias lists all aliases from a registry.
	ListAliases(ctx context.Context, registry RegistryName) ([]Alias, error)

	// DeleteRegistry deletes a registry, with all the aliases it contained.
	DeleteRegistry(ctx context.Context, registry RegistryName) error
}

// Alias allows the user to associates a RegistryName + a Key to 1 or more public keys stored
// in Value. The Value has 2 formats:
// - a JSON string if AliasKind is an AliasKindString.
// - a JSON array of strings if AliasKind is an AliasKindArray.
type Alias struct {
	tableName struct{} `pg:"aliases"` // nolint:unused,structcheck // reason

	Key          AliasKey     `pg:",pk"`
	RegistryName RegistryName `pg:",pk"`
	Kind         AliasKind
	Value        AliasValue
}

type AliasKey string

type AliasValue string

type AliasKind string

const (
	AliasKindUnknown = ""
	AliasKindString  = "string"
	AliasKindArray   = "array"
)

type RegistryName string
