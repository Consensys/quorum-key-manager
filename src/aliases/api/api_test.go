package aliasapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aliasapi "github.com/consensys/quorum-key-manager/src/aliases/api"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/aliases/mock"
)

type apiHelper struct {
	ctx     context.Context
	mock    *mock.MockAlias
	rec     *httptest.ResponseRecorder
	router  *mux.Router
	handler *aliasapi.AliasAPI
}

func NewAPIHelper(t *testing.T) *apiHelper {
	ctrl := gomock.NewController(t)
	store := mock.NewMockAlias(ctrl)
	handler := aliasapi.New(store)
	router := mux.NewRouter()
	handler.Register(router)

	return &apiHelper{
		ctx:     context.Background(),
		mock:    store,
		rec:     httptest.NewRecorder(),
		router:  router,
		handler: handler,
	}
}

func TestCreateAlias(t *testing.T) {
	helper := NewAPIHelper(t)

	cases := map[string]struct {
		reg   string
		key   string
		value string

		status int
	}{
		"empty array": {"testr", "akey", `[]`, http.StatusOK},
		"one element": {"testr", "akey2", `[ "0123" ]`, http.StatusOK},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			alias := types.CreateAliasRequest{
				Alias: newAPIAlias(c.key, c.value),
			}
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(alias)
			require.NoError(t, err)

			ent := newEntAlias(c.reg, c.key, c.value)
			helper.mock.EXPECT().CreateAlias(gomock.Any(), aliasent.RegistryName(c.reg), ent)

			path := fmt.Sprintf("/aliases/registries/%s/aliases", c.reg)
			r, err := http.NewRequestWithContext(helper.ctx, "POST", path, &b)
			require.NoError(t, err)

			helper.router.ServeHTTP(helper.rec, r)
			assert.Equal(t, helper.rec.Code, c.status)
			assert.Contains(t, helper.rec.HeaderMap["Content-Type"], "application/json; charset=UTF-8;")
		})
	}
}

func newEntAlias(registry, key, value string) aliasent.Alias {
	return aliasent.Alias{
		RegistryName: aliasent.RegistryName(registry),
		Key:          aliasent.AliasKey(key),
		Value:        aliasent.AliasValue(value),
	}
}

func newAPIAlias(key, value string) types.Alias {
	return types.Alias{
		Key:   types.AliasKey(key),
		Value: types.AliasValue(value),
	}
}

func TestJSONHeader(t *testing.T) {
	//assert.Contains(t, rr.HeaderMap["Content-Type"], "application/json; charset=UTF-8;")
}

func TestFullScenario(t *testing.T) {
}
