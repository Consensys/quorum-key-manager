package placeholder

import (
	"context"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

//go:generate mockgen -source=alias_parser.go -destination=mock/alias_parser.go -package=mock

// AliasParser parses and replace aliases.
type AliasParser interface {
	ParseAlias(alias string) (regName string, aliasKey string, isAlias bool)
	ReplaceAliases(ctx context.Context, aliasBackend aliasent.AliasBackend, addrs []string) ([]string, error)
}
