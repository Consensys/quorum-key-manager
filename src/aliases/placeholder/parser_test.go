package placeholder_test

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/aliases/entities/mock"
	"github.com/consensys/quorum-key-manager/src/aliases/placeholder"
	"github.com/golang/mock/gomock"
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
		})

	}
}

func TestReplaceAliases(t *testing.T) {
	type backendCall struct {
		reg   aliasent.RegistryName
		key   aliasent.AliasKey
		value aliasent.AliasValue
		err   error
	}

	groupACall := backendCall{"my-registry", "group-A", `["ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=","2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="]`, nil}
	JPMCall := backendCall{"my-registry", "JPM", `["ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="]`, nil}
	GSCall := backendCall{"my-registry", "GS", `["2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="]`, nil}
	errCall := backendCall{"unknown-registry", "unknown-key", "", errors.InvalidFormatError("bad format")}

	cases := map[string]struct {
		addrs []string
		calls []backendCall
	}{
		"bad key":          {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{unknown-registry:bad/key}}"}, nil},
		"unknown registry": {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{unknown-registry:unknown-key}}"}, []backendCall{errCall}},
		"bad registry":     {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{bad#registry:unknown-key}}"}, nil},
		"ok without alias": {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}, nil},
		"ok 1":             {[]string{"{{my-registry:group-A}}"}, []backendCall{groupACall}},
		"ok 2":             {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{my-registry:JPM}}"}, []backendCall{JPMCall}},
		"ok 3":             {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{my-registry:GS}}", "{{my-registry:JPM}}"}, []backendCall{JPMCall, GSCall}},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	aliasBackend := mock.NewMockAliasBackend(ctrl)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			for _, call := range c.calls {
				aliasBackend.EXPECT().GetAlias(gomock.Any(), call.reg, call.key).Return(&aliasent.Alias{Value: call.value}, call.err)
			}

			addrs, err := placeholder.ReplaceAliases(ctx, aliasBackend, c.addrs)
			if err != nil {
				require.True(t, errors.IsInvalidFormatError(err))
				return
			}
			t.Log("err", err)

			assert.Len(t, addrs, len(c.addrs))
			// we check the aliases have been extracted to the results
			for _, call := range c.calls {
				present := false
				for _, addr := range addrs {
					if addr == string(call.value) {
						present = true
						break
					}
				}
				assert.True(t, present)
			}
			t.Log(c.addrs, addrs)
		})
	}
}
