package placeholder

import (
	"context"
	"regexp"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/aliases"
)

// make sure Parser implements AliasParser
var _ AliasParser = &Parser{}

// Parser parses alias strings in the form of {{registry_name:alias_key}}.
type Parser struct {
	regex *regexp.Regexp
}

// New creates a new Parser or fails if the regexp is invalid.
func New() (*Parser, error) {
	const aliasParseFormat = `{{(?m)(?P<registry>[a-zA-Z0-9-_+]+):(?P<alias>[a-zA-Z0-9-_+]+)}}$`
	regex, err := regexp.Compile(aliasParseFormat)
	if err != nil {
		return nil, errors.ConfigError("bad regexp format '%v': %v", aliasParseFormat, err)
	}
	return &Parser{
		regex: regex,
	}, nil
}

// ParseAlias parses an alias string and returns the registryName and the aliasKey
// as well as if the string isAlias. If the string is not isAlias, we'll consider it
// as a valid key.
func (p *Parser) ParseAlias(alias string) (regName, aliasKey string, isAlias bool) {
	submatches := p.regex.FindStringSubmatch(alias)
	if len(submatches) < 3 {
		return "", "", false
	}

	regName = submatches[1]
	aliasKey = submatches[2]

	return regName, aliasKey, true
}

// ReplaceAliases replace a slice of potential aliases with a slice having all the aliases replaced by their value.
// It will fail if no aliases can be found.
func (p *Parser) ReplaceAliases(ctx context.Context, aliasBackend aliases.AliasBackend, addrs []string) ([]string, error) {
	var values []string
	for _, addr := range addrs {
		regName, aliasKey, isAlias := p.ParseAlias(addr)

		// it is not an alias
		if !isAlias {
			values = append(values, addr)
			continue
		}

		alias, err := aliasBackend.GetAlias(ctx, regName, aliasKey)
		if err != nil {
			return nil, err
		}

		values = append(values, alias.Value...)
	}
	return values, nil
}
