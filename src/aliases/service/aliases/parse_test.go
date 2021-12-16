package aliases

import (
	"testing"

	mock2 "github.com/consensys/quorum-key-manager/src/aliases/database/mock"
	"github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		input  string
		reg    string
		key    string
		parsed bool
	}{
		"single {":         {`{ok_registry:ok_key}`, "", "", false},
		"column missing":   {`{{ok_registry ok_key}}`, "", "", false},
		"too many columns": {`{{ok_registry:ok_key:}}`, "", "", false},
		"base 64 key":      {`ROAZBWtSacxXQrOe3FGAqJDyJjFePR5ce4TSIzmJ0Bc=`, "", "", false},
		"ok":               {`{{ok_registry:ok_key}}`, "ok_registry", "ok_key", true},
	}

	ctrl := gomock.NewController(t)
	loggerMock := testutils.NewMockLogger(ctrl)
	mockDB := mock2.NewMockAlias(ctrl)
	mockRegistryDB := mock2.NewMockRegistry(ctrl)
	mockRoles := mock.NewMockRoles(ctrl)

	aConn := New(mockDB, mockRegistryDB, mockRoles, loggerMock)

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			reg, key, parsed := aConn.Parse(c.input)
			assert.Equal(t, c.reg, reg)
			assert.Equal(t, c.key, key)
			assert.Equal(t, c.parsed, parsed)
		})

	}
}
