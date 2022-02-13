package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/src/aliases/api/types"
	"github.com/consensys/quorum-key-manager/src/aliases/api/types/testutils"
	"github.com/consensys/quorum-key-manager/src/aliases/mock"
	authapi "github.com/consensys/quorum-key-manager/src/auth/api/http"
	authentities "github.com/consensys/quorum-key-manager/src/auth/entities"
	testutils2 "github.com/consensys/quorum-key-manager/src/entities/testutils"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var reqUserInfo = &authentities.UserInfo{
	Username:    "username",
	Roles:       []string{"role1", "role2"},
	Permissions: []authentities.Permission{"*:*"},
}

type registriesHandlerTestSuite struct {
	suite.Suite

	ctrl       *gomock.Controller
	router     *mux.Router
	registries *mock.MockRegistries
	ctx        context.Context
}

func TestRegistryHandler(t *testing.T) {
	s := new(registriesHandlerTestSuite)
	suite.Run(t, s)
}

func (s *registriesHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.registries = mock.NewMockRegistries(s.ctrl)

	s.ctx = authapi.WithUserInfo(context.Background(), reqUserInfo)

	s.router = mux.NewRouter()
	NewRegistryHandler(s.registries).Register(s.router)
}

func (s *registriesHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *registriesHandlerTestSuite) TestSet() {
	registry := testutils2.FakeAliasRegistry()

	s.Run("should execute request body successfully", func() {
		registryReq := testutils.FakeCreateRegistryRequest()
		requestBytes, _ := json.Marshal(registryReq)
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/registries/"+registry.Name, bytes.NewReader(requestBytes)).
			WithContext(s.ctx)

		s.registries.EXPECT().Create(gomock.Any(), registry.Name, registryReq.AllowedTenants, reqUserInfo).Return(registry, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := types.NewRegistryResponse(registry)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with no request body successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/registries/"+registry.Name, nil).WithContext(s.ctx)

		s.registries.EXPECT().Create(gomock.Any(), registry.Name, nil, reqUserInfo).Return(registry, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := types.NewRegistryResponse(registry)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})
}
