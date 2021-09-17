package placeholder

import (
	"context"
	"regexp"
	"sync"

	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

const aliasParseFormat = `{{(?m)(?P<registry>[a-zA-Z0-9-_+]+):(?P<alias>[a-zA-Z0-9-_+]+)}}$`

var (
	aliasParseRegexOnce sync.Once
	aliasParseRegex     *regexp.Regexp
)

func ParseAlias(alias string) (regName aliasent.RegistryName, aliasKey aliasent.AliasKey, isAlias bool, err error) {
	aliasParseRegexOnce.Do(func() {
		aliasParseRegex, err = regexp.Compile(aliasParseFormat)
	})
	if err != nil {
		return "", "", false, err
	}
	submatches := aliasParseRegex.FindStringSubmatch(alias)
	if len(submatches) < 3 {
		return "", "", false, nil
	}

	regName = aliasent.RegistryName(submatches[1])
	aliasKey = aliasent.AliasKey(submatches[2])

	return regName, aliasKey, true, nil
}

func ReplaceAliases(ctx context.Context, aliasBackend aliasent.AliasBackend, addrs []string) ([]string, error) {
	var values []string
	for _, addr := range addrs {
		regName, aliasKey, isAlias, err := ParseAlias(addr)
		if err != nil {
			return nil, err
		}

		// it is not an alias
		if !isAlias {
			values = append(values, addr)
			continue
		}

		alias, err := aliasBackend.GetAlias(ctx, regName, aliasKey)
		if err != nil {
			return nil, err
		}

		for _, v := range alias.Value {
			values = append(values, v)
		}
	}
	return values, nil
}
