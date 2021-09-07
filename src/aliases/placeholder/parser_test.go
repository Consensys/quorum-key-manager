package placeholder_test

import (
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/aliases/placeholder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAlias(t *testing.T) {
	cases := map[string]struct {
		input string

		reg          aliasent.RegistryName
		key          aliasent.AliasKey
		parsed       bool
		errCompareFn func(error) bool
	}{
		"bad registry format": {`{{bad#registry:ok_key}}`, "", "", false, errors.IsInvalidFormatError},
		"bad key format":      {`{{ok_registry:bad>key}}`, "", "", false, errors.IsInvalidFormatError},
		"single {":            {`{ok_registry:ok_key}`, "", "", false, errors.IsInvalidFormatError},
		"column missing":      {`{{ok_registry ok_key}}`, "", "", false, errors.IsInvalidFormatError},
		"too many columns":    {`{{ok_registry:ok_key:}}`, "", "", false, errors.IsInvalidFormatError},
		"ok":                  {`{{ok_registry:ok_key}}`, "ok_registry", "ok_key", true, nil},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			reg, key, parsed, err := placeholder.ParseAlias(c.input)
			assert.Equal(t, c.reg, reg)
			assert.Equal(t, c.key, key)
			assert.Equal(t, c.parsed, parsed)
			if c.errCompareFn != nil {
				require.Error(t, err)
				assert.True(t, c.errCompareFn(err))
			}
			t.Log(reg, key, parsed, err)
		})

	}
}
