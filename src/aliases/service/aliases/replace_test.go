package aliases

import (
	"context"
	"testing"

	mock2 "github.com/consensys/quorum-key-manager/src/aliases/database/mock"
	"github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	auth "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/entities"
)

type backendCall struct {
	reg   string
	key   string
	kind  string
	value interface{}
	err   error
}

func TestReplace(t *testing.T) {
	groupACall := backendCall{"my-registry", "group-A", entities.AliasKindArray, []interface{}{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}, nil}
	JPMCall := backendCall{"my-registry", "JPM", entities.AliasKindArray, []interface{}{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="}, nil}
	GSCall := backendCall{"my-registry", "GS", entities.AliasKindArray, []interface{}{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}, nil}
	errCall := backendCall{"unknown-registry", "unknown-key", entities.AliasKindArray, []interface{}{""}, errors.InvalidFormatError("bad format")}
	user := auth.NewWildcardUser()

	cases := map[string]struct {
		addrs  []string
		calls  []backendCall
		expLen int
	}{
		"unknown registry": {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{unknown-registry:unknown-key}}"}, []backendCall{errCall}, 2},
		"ok without alias": {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}, nil, 1},
		"ok 1":             {[]string{"{{my-registry:group-A}}"}, []backendCall{groupACall}, 2},
		"ok 2":             {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{my-registry:JPM}}"}, []backendCall{JPMCall}, 2},
		"ok 3":             {[]string{"2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0=", "{{my-registry:GS}}", "{{my-registry:JPM}}"}, []backendCall{JPMCall, GSCall}, 3},
	}

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := mock2.NewMockAlias(ctrl)
	mockRegistryDB := mock2.NewMockRegistry(ctrl)
	loggerMock := testutils.NewMockLogger(ctrl)
	mockRoles := mock.NewMockRoles(ctrl)
	mockRoles.EXPECT().UserPermissions(gomock.Any(), user).Return(auth.ListPermissions()).AnyTimes()

	aConn := New(mockDB, mockRegistryDB, mockRoles, loggerMock)

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			for _, call := range c.calls {
				mockDB.EXPECT().FindOne(gomock.Any(), call.reg, call.key, user.Tenant).Return(&entities.Alias{Kind: call.kind, Value: call.value}, call.err)
			}

			addrs, err := aConn.Replace(ctx, c.addrs, user)
			if err != nil {
				require.True(t, errors.IsInvalidFormatError(err))
				return
			}

			assert.Len(t, addrs, c.expLen)
			// we check the aliases have been extracted to the results
			for _, call := range c.calls {
				present := false
				for _, addr := range addrs {
					for _, v := range call.value.([]interface{}) {
						if addr == v {
							present = true
							break
						}
					}
				}
				assert.True(t, present)
			}
		})
	}
}

func TestReplaceSimple(t *testing.T) {
	groupACall := backendCall{"my-registry", "group-A", entities.AliasKindArray, []interface{}{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=", "2T7xkjblN568N1QmPeElTjoeoNT4tkWYOJYxSMDO5i0="}, nil}
	JPMCall := backendCall{"my-registry", "JPM", entities.AliasKindArray, []interface{}{"ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc="}, nil}
	user := auth.NewWildcardUser()

	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockDB := mock2.NewMockAlias(ctrl)
	loggerMock := testutils.NewMockLogger(ctrl)
	mockRegistryDB := mock2.NewMockRegistry(ctrl)
	mockRoles := mock.NewMockRoles(ctrl)
	mockRoles.EXPECT().UserPermissions(gomock.Any(), user).Return(auth.ListPermissions()).AnyTimes()

	aConn := New(mockDB, mockRegistryDB, mockRoles, loggerMock)

	t.Run("no alias found", func(t *testing.T) {
		mockDB.EXPECT().FindOne(gomock.Any(), groupACall.reg, groupACall.key, user.Tenant).Return(nil, errors.NotFoundError("resource not found"))
		_, err := aConn.ReplaceSimple(ctx, "{{my-registry:group-A}}", user)
		require.Error(t, err)
		assert.True(t, errors.IsNotFoundError(err))
	})

	t.Run("more than 1 alias value", func(t *testing.T) {
		mockDB.EXPECT().FindOne(gomock.Any(), groupACall.reg, groupACall.key, user.Tenant).Return(&entities.Alias{Kind: groupACall.kind, Value: groupACall.value}, nil)
		_, err := aConn.ReplaceSimple(ctx, "{{my-registry:group-A}}", user)
		require.Error(t, err)
		assert.True(t, errors.IsEncodingError(err))
	})

	t.Run("1 alias value", func(t *testing.T) {
		mockDB.EXPECT().FindOne(gomock.Any(), JPMCall.reg, JPMCall.key, user.Tenant).Return(&entities.Alias{Kind: JPMCall.kind, Value: JPMCall.value}, nil)
		addr, err := aConn.ReplaceSimple(ctx, "{{my-registry:JPM}}", user)
		require.NoError(t, err)
		assert.Equal(t, groupACall.value.([]interface{})[0], addr)
	})
}
