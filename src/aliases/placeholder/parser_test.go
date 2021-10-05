package placeholder_test

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/aliases/mock"
	"github.com/consensys/quorum-key-manager/src/aliases/placeholder"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAlias(t *testing.T) {
	cases := map[string]struct {
		input string

		reg    string
		key    string
		parsed bool
	}{
		"bad registry format": {`{{bad#registry:ok_key}}`, "", "", false},
		"bad key format":      {`{{ok_registry:bad>key}}`, "", "", false},
		"single {":            {`{ok_registry:ok_key}`, "", "", false},
		"column missing":      {`{{ok_registry ok_key}}`, "", "", false},
		"too many columns":    {`{{ok_registry:ok_key:}}`, "", "", false},
		"base 64 key":         {`ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=`, "", "", false},
		"ok":                  {`{{ok_registry:ok_key}}`, "ok_registry", "ok_key", true},
	}
	p, err := placeholder.New()
	require.NoError(t, err)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			reg, key, parsed := p.ParseAlias(c.input)
			assert.Equal(t, c.reg, reg)
			assert.Equal(t, c.key, key)
			assert.Equal(t, c.parsed, parsed)
		})

	}
}

func TestReplaceAliases(t *testing.T) {
	type backendCall struct {
		reg   string
		key   string
		value []string
		err   error
	}

	groupACall := backendCall{"my-registry", "group-A", []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}, nil}
	JPMCall := backendCall{"my-registry", "JPM", []string{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="}, nil}
	GSCall := backendCall{"my-registry", "GS", []string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}, nil}
	errCall := backendCall{"unknown-registry", "unknown-key", []string{""}, errors.InvalidFormatError("bad format")}

	cases := map[string]struct {
		addrs  []string
		calls  []backendCall
		expLen int
	}{
		"unknown registry": {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{unknown-registry:unknown-key}}"}, []backendCall{errCall}, 2},
		"bad key":          {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{unknown-registry:bad/key}}"}, nil, 2},
		"bad registry":     {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{bad#registry:unknown-key}}"}, nil, 2},
		"ok without alias": {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}, nil, 1},
		"ok 1":             {[]string{"{{my-registry:group-A}}"}, []backendCall{groupACall}, 2},
		"ok 2":             {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{my-registry:JPM}}"}, []backendCall{JPMCall}, 2},
		"ok 3":             {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{my-registry:GS}}", "{{my-registry:JPM}}"}, []backendCall{JPMCall, GSCall}, 3},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	aliasBackend := mock.NewMockAliasBackend(ctrl)

	p, err := placeholder.New()
	require.NoError(t, err)

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			for _, call := range c.calls {
				aliasBackend.EXPECT().GetAlias(gomock.Any(), call.reg, call.key).Return(&aliasent.Alias{Value: call.value}, call.err)
			}

			addrs, err := p.ReplaceAliases(ctx, aliasBackend, c.addrs)
			if err != nil {
				require.True(t, errors.IsInvalidFormatError(err))
				return
			}

			assert.Len(t, addrs, c.expLen)
			// we check the aliases have been extracted to the results
			for _, call := range c.calls {
				present := false
				for _, addr := range addrs {
					for _, v := range call.value {
						if addr == v {
							present = true
							break
						}
					}
				}
				assert.True(t, present)
			}
			t.Log(c.addrs, addrs)
		})
	}
}
