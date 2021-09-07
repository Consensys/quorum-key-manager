package placeholder

import (
	"context"
	"regexp"
	"sync"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
)

const aliasParseFormat = `{{(?m)(?P<registry>[a-zA-Z0-9-_+]+):(?P<alias>[a-zA-Z0-9-_+]+)}}$`

var (
	aliasParseRegexOnce sync.Once
	aliasParseRegex     *regexp.Regexp
)

func ParseAlias(alias string) (regName aliasent.RegistryName, aliasKey aliasent.AliasKey, parsed bool, err error) {
	formatError := errors.InvalidFormatError(`alias not in the format "{{registry_name:alias_key}}"`)
	aliasParseRegexOnce.Do(func() {
		aliasParseRegex, err = regexp.Compile(aliasParseFormat)
	})
	if err != nil {
		return "", "", false, err
	}
	submatches := aliasParseRegex.FindStringSubmatch(alias)
	if len(submatches) < 3 {
		return "", "", false, formatError
	}

	regName = aliasent.RegistryName(submatches[1])
	aliasKey = aliasent.AliasKey(submatches[2])

	return regName, aliasKey, true, nil
}

func ReplaceAliases(ctx context.Context, aliasBackend aliasent.AliasBackend, addrs []string) ([]string, error) {
	var values []string
	for _, v := range addrs {
		regName, aliasKey, parsed, err := ParseAlias(v)
		if err != nil {
			values = append(values, v)
			continue
		}
		if parsed {
		}

		if err != nil {
			return nil, err
		}

		alias, err := aliasBackend.GetAlias(ctx, regName, aliasKey)
		if err != nil {
			return nil, err
		}
		values = append(values, string(alias.Value))
	}
	return values, nil
}
