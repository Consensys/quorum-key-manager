package aliases

import "context"

//go:generate mockgen -destination=mock/parser.go -package=mock . Parser

// Parser parses and replace aliases.
type Parser interface {
	ParseAlias(alias string) (regName string, aliasKey string, isAlias bool)
	ReplaceAliases(ctx context.Context, addrs []string) ([]string, error)
}
