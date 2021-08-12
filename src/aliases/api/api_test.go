package aliasapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/errors"
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

type Case struct {
	reg   string
	key   string
	value string

	status int
}

func TestCreateAlias(t *testing.T) {
	t.Run("one element", func(t *testing.T) {
		c := Case{"testr", "akey2", `[ "0123" ]`, http.StatusOK}
		t.Parallel()
		helper := NewAPIHelper(t)
		req := types.CreateAliasRequest{
			Alias: newAPIAlias(c.key, c.value),
		}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(req)
		require.NoError(t, err)

		ent := newEntAlias(c.reg, c.key, c.value)
		helper.mock.EXPECT().CreateAlias(gomock.Any(), ent.RegistryName, ent).Return(&ent, nil)

		path := fmt.Sprintf("/aliases/registries/%s/aliases", c.reg)
		r, err := http.NewRequestWithContext(helper.ctx, "POST", path, &b)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, c.status)
		res, err := io.ReadAll(helper.rec.Body)
		require.NoError(t, err)

		var resp types.CreateAliasResponse
		err = json.Unmarshal(res, &resp)
		require.NoError(t, err)

		assert.Equal(t, types.CreateAliasResponse{Alias: types.Alias{Key: types.AliasKey(c.key), Value: types.AliasValue(c.value)}}, resp)
	})

	t.Run("already existing alias", func(t *testing.T) {
		t.Parallel()
		helper := NewAPIHelper(t)
		reg, key, value, status := "testr", "existing_key", "[]", http.StatusConflict
		req := types.CreateAliasRequest{
			Alias: newAPIAlias(key, value),
		}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(req)
		require.NoError(t, err)

		ent := newEntAlias(reg, key, value)
		helper.mock.EXPECT().CreateAlias(gomock.Any(), ent.RegistryName, ent).Return(nil, errors.AlreadyExistsError(""))

		path := fmt.Sprintf("/aliases/registries/%s/aliases", reg)
		r, err := http.NewRequestWithContext(helper.ctx, "POST", path, &b)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, status)
		res, err := io.ReadAll(helper.rec.Body)
		require.NoError(t, err)
		assert.Contains(t, string(res), `"code":"ST200"`)
	})

}

func TestUpdateAlias(t *testing.T) {
	t.Run("one element", func(t *testing.T) {
		c := Case{"testr", "akey2", `[ "0123" ]`, http.StatusOK}
		t.Parallel()
		helper := NewAPIHelper(t)
		req := types.UpdateAliasRequest{
			Value: types.AliasValue(c.value),
		}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(req)
		require.NoError(t, err)

		ent := newEntAlias(c.reg, c.key, c.value)
		helper.mock.EXPECT().UpdateAlias(gomock.Any(), ent.RegistryName, ent).Return(&ent, nil)

		path := fmt.Sprintf("/aliases/registries/%s/aliases/%s", c.reg, c.key)
		r, err := http.NewRequestWithContext(helper.ctx, "PUT", path, &b)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, c.status)
		res, err := io.ReadAll(helper.rec.Body)
		require.NoError(t, err)

		var resp types.UpdateAliasResponse
		err = json.Unmarshal(res, &resp)
		require.NoError(t, err)

		assert.Equal(t, types.UpdateAliasResponse{
			Alias: types.Alias{
				Key:   types.AliasKey(c.key),
				Value: types.AliasValue(c.value),
			},
		},
			resp)
	})

	t.Run("non-existing alias", func(t *testing.T) {
		t.Parallel()
		helper := NewAPIHelper(t)
		reg, key, value, status := "testr", "non_existing_key", "[]", http.StatusNotFound
		alias := types.UpdateAliasRequest{
			Value: types.AliasValue(value),
		}
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(alias)
		require.NoError(t, err)

		ent := newEntAlias(reg, key, value)
		helper.mock.EXPECT().UpdateAlias(gomock.Any(), ent.RegistryName, ent).Return(nil, errors.NotFoundError(""))

		path := fmt.Sprintf("/aliases/registries/%s/aliases/%s", reg, key)
		r, err := http.NewRequestWithContext(helper.ctx, "PUT", path, &b)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, status)
		res, err := io.ReadAll(helper.rec.Body)
		require.NoError(t, err)
		assert.Contains(t, string(res), errors.NotFound)
	})
}

func TestGetAlias(t *testing.T) {
	t.Run("one element", func(t *testing.T) {
		c := Case{"testr", "akey2", `[ "0123" ]`, http.StatusOK}
		t.Parallel()
		helper := NewAPIHelper(t)

		ent := newEntAlias(c.reg, c.key, c.value)
		helper.mock.EXPECT().GetAlias(gomock.Any(), ent.RegistryName, ent.Key).Return(&ent, nil)

		path := fmt.Sprintf("/aliases/registries/%s/aliases/%s", c.reg, c.key)
		r, err := http.NewRequestWithContext(helper.ctx, "GET", path, nil)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, c.status)
	})

	t.Run("non-existing alias", func(t *testing.T) {
		reg, key, status := "testr", "non_existing_key", http.StatusNotFound
		t.Parallel()
		helper := NewAPIHelper(t)

		ent := newEntAlias(reg, key, "")
		helper.mock.EXPECT().GetAlias(gomock.Any(), ent.RegistryName, ent.Key).Return(nil, errors.NotFoundError(""))

		path := fmt.Sprintf("/aliases/registries/%s/aliases/%s", reg, key)
		r, err := http.NewRequestWithContext(helper.ctx, "GET", path, nil)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, status)
		res, err := io.ReadAll(helper.rec.Body)
		require.NoError(t, err)
		assert.Contains(t, string(res), errors.NotFound)
	})
}

func TestDeleteAlias(t *testing.T) {
	t.Run("one element", func(t *testing.T) {
		c := Case{"testr", "akey2", `[ "0123" ]`, http.StatusNoContent}
		t.Parallel()
		helper := NewAPIHelper(t)

		ent := newEntAlias(c.reg, c.key, c.value)
		helper.mock.EXPECT().DeleteAlias(gomock.Any(), ent.RegistryName, ent.Key).Return(nil)

		path := fmt.Sprintf("/aliases/registries/%s/aliases/%s", c.reg, c.key)
		r, err := http.NewRequestWithContext(helper.ctx, "DELETE", path, nil)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, c.status)
	})
	t.Run("non-existing alias", func(t *testing.T) {
		reg, key, status := "testr", "non_existing_key", http.StatusNotFound
		t.Parallel()
		helper := NewAPIHelper(t)

		ent := newEntAlias(reg, key, "")
		helper.mock.EXPECT().DeleteAlias(gomock.Any(), ent.RegistryName, ent.Key).Return(errors.NotFoundError(""))

		path := fmt.Sprintf("/aliases/registries/%s/aliases/%s", reg, key)
		r, err := http.NewRequestWithContext(helper.ctx, "DELETE", path, nil)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, status)
		res, err := io.ReadAll(helper.rec.Body)
		require.NoError(t, err)
		assert.Contains(t, string(res), errors.NotFound)
	})
}

func TestListAliases(t *testing.T) {
	t.Run("non-existing registry", func(t *testing.T) {
		reg, status := "non_existing_registry", http.StatusNotFound
		t.Parallel()
		helper := NewAPIHelper(t)

		helper.mock.EXPECT().ListAliases(gomock.Any(), aliasent.RegistryName(reg)).Return(nil, errors.NotFoundError(""))

		path := fmt.Sprintf("/aliases/registries/%s/aliases", reg)
		r, err := http.NewRequestWithContext(helper.ctx, "GET", path, nil)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, status)
		res, err := io.ReadAll(helper.rec.Body)
		require.NoError(t, err)
		assert.Contains(t, string(res), errors.NotFound)
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		helper := NewAPIHelper(t)

		reg, status := "testr", http.StatusOK
		var ents []aliasent.Alias
		helper.mock.EXPECT().ListAliases(gomock.Any(), aliasent.RegistryName(reg)).Return(ents, nil)

		path := fmt.Sprintf("/aliases/registries/%s/aliases", reg)
		r, err := http.NewRequestWithContext(helper.ctx, "GET", path, nil)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, status)

		var als []types.Alias
		err = json.Unmarshal(helper.rec.Body.Bytes(), &als)
		require.NoError(t, err)

		alses := types.FromEntityAliases(ents)
		assert.Equal(t, als, alses)

	})

	t.Run("list of 2 elements", func(t *testing.T) {
		t.Parallel()
		helper := NewAPIHelper(t)

		reg, status := "testr", http.StatusOK
		ents := []aliasent.Alias{
			{
				Key:   "key_1",
				Value: `[ "value_1_1", "value_2_1" ]`,
			},
			{
				Key:   "key_2",
				Value: `[ "value_2_1", "value_2_2" ]`,
			},
		}
		helper.mock.EXPECT().ListAliases(gomock.Any(), aliasent.RegistryName(reg)).Return(ents, nil)

		path := fmt.Sprintf("/aliases/registries/%s/aliases", reg)
		r, err := http.NewRequestWithContext(helper.ctx, "GET", path, nil)
		require.NoError(t, err)

		helper.router.ServeHTTP(helper.rec, r)
		assert.Equal(t, helper.rec.Code, status)

		var als []types.Alias
		err = json.Unmarshal(helper.rec.Body.Bytes(), &als)
		require.NoError(t, err)
		assert.Equal(t, als, types.FromEntityAliases(ents))
	})
}

func TestJSONHeader(t *testing.T) {
	//assert.Contains(t, rr.HeaderMap["Content-Type"], "application/json; charset=UTF-8;")
	cases := []struct {
		method string
		path   string
		input  string
	}{
		{"POST", "/aliases/registries/1/aliases", `{"key": "1", "value": "[ \"0123\" ]"}`},
		{"GET", "/aliases/registries/1/aliases", ``},
		{"GET", "/aliases/registries/1/aliases/1", ``},
		{"PUT", "/aliases/registries/1/aliases/1", `{"key": "1", "value": "[ \"01234\" ]"}`},
		// DELETE doesn't return any content, only 204 or 404
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s %s", c.method, c.path), func(t *testing.T) {
			t.Parallel()
			helper := NewAPIHelper(t)

			input := strings.NewReader(c.input)

			r, err := http.NewRequestWithContext(helper.ctx, c.method, c.path, input)
			require.NoError(t, err)

			ent := newEntAlias("1", "1", `[ \"0123\" ]`)
			// accept any call just to make the test work
			helper.mock.EXPECT().CreateAlias(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&ent, nil)
			helper.mock.EXPECT().GetAlias(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&ent, nil)
			helper.mock.EXPECT().UpdateAlias(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&ent, nil)
			helper.mock.EXPECT().ListAliases(gomock.Any(), gomock.Any()).AnyTimes().Return([]aliasent.Alias{ent}, nil)

			helper.router.ServeHTTP(helper.rec, r)
			assert.Contains(t, helper.rec.HeaderMap.Values("Content-Type"), "application/json; charset=UTF-8;")
			assert.Contains(t, helper.rec.HeaderMap.Values("X-Content-Type-Options"), "nosniff")
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
