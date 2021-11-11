package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/api/middlewares"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	authentities "github.com/consensys/quorum-key-manager/src/auth/entities"
	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/consensys/quorum-key-manager/src/stores/api/types/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	secretStoreName = "SecretStore"
	secretID        = "my-secret"
)

var secretUserInfo = &authentities.UserInfo{
	Username:    "username",
	Roles:       []string{"role1", "role2"},
	Permissions: []authentities.Permission{"write:key", "read:key", "sign:key"},
}

type secretsHandlerTestSuite struct {
	suite.Suite

	ctrl        *gomock.Controller
	stores      *mock.MockStores
	secretStore *mock.MockSecretStore
	router      *mux.Router
	ctx         context.Context
}

func TestSecretsHandler(t *testing.T) {
	s := new(secretsHandlerTestSuite)
	suite.Run(t, s)
}

func (s *secretsHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.stores = mock.NewMockStores(s.ctrl)
	s.secretStore = mock.NewMockSecretStore(s.ctrl)

	s.stores.EXPECT().Secret(gomock.Any(), secretStoreName, secretUserInfo).Return(s.secretStore, nil).AnyTimes()

	s.ctx = middlewares.WithUserInfo(context.Background(), secretUserInfo)

	s.router = mux.NewRouter()
	NewStoresHandler(s.stores).Register(s.router)
}

func (s *secretsHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *secretsHandlerTestSuite) TestSet() {
	secretID := "my-secret"

	s.Run("should execute request successfully", func() {
		setSecretRequest := testutils.FakeSetSecretRequest()
		requestBytes, _ := json.Marshal(setSecretRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/SecretStore/secrets/"+secretID, bytes.NewReader(requestBytes)).WithContext(s.ctx)
		secret := testutils2.FakeSecret()

		s.secretStore.EXPECT().Set(gomock.Any(), secretID, setSecretRequest.Value, &entities.Attributes{
			Tags: setSecretRequest.Tags,
		}).Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		setSecretRequest := testutils.FakeSetSecretRequest()
		requestBytes, _ := json.Marshal(setSecretRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/SecretStore/secrets/"+secretID, bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.secretStore.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestGet() {
	s.Run("should execute request successfully with version", func() {
		version := "1"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s?version=%s", secretID, version), nil).WithContext(s.ctx)

		secret := testutils2.FakeSecret()
		s.secretStore.EXPECT().Get(gomock.Any(), secretID, version).Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request successfully without version", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s", secretID), nil).WithContext(s.ctx)

		secret := testutils2.FakeSecret()
		s.secretStore.EXPECT().Get(gomock.Any(), secretID, "").Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s", secretID), nil).WithContext(s.ctx)

		s.secretStore.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestGetDeleted() {
	s.Run("should execute request successfully with version", func() {
		version := "1"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s?version=%s&deleted=true", secretID, version), nil).WithContext(s.ctx)

		secret := testutils2.FakeSecret()
		s.secretStore.EXPECT().GetDeleted(gomock.Any(), secretID).Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request successfully without version", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s?deleted=true", secretID), nil).WithContext(s.ctx)

		secret := testutils2.FakeSecret()
		s.secretStore.EXPECT().GetDeleted(gomock.Any(), secretID).Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s?deleted=true", secretID), nil).WithContext(s.ctx)

		s.secretStore.EXPECT().GetDeleted(gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestDelete() {
	s.Run("should execute request successfully with version", func() {
		version := "1"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/SecretStore/secrets/%s?version=%s", secretID, version), nil).WithContext(s.ctx)

		s.secretStore.EXPECT().Delete(gomock.Any(), secretID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/SecretStore/secrets/%s", secretID), nil).WithContext(s.ctx)

		s.secretStore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestRestore() {
	s.Run("should execute request successfully with version", func() {
		version := "1"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/stores/SecretStore/secrets/%s/restore?version=%s", secretID, version), nil).WithContext(s.ctx)

		s.secretStore.EXPECT().Restore(gomock.Any(), secretID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/stores/SecretStore/secrets/%s/restore", secretID), nil).WithContext(s.ctx)

		s.secretStore.EXPECT().Restore(gomock.Any(), gomock.Any()).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestDestroy() {
	s.Run("should execute request successfully with version", func() {
		version := "1"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/SecretStore/secrets/%s/destroy?version=%s", secretID, version), nil).WithContext(s.ctx)

		s.secretStore.EXPECT().Destroy(gomock.Any(), secretID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/SecretStore/secrets/%s/destroy", secretID), nil).WithContext(s.ctx)

		s.secretStore.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestList() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/SecretStore/secrets", nil).WithContext(s.ctx)

		ids := []string{"secret1", "secret2"}

		s.secretStore.EXPECT().List(gomock.Any(), defaultPageSize, uint64(0)).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(http2.PageResponse{
			Data: ids,
		})
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with limit and offset successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/SecretStore/secrets?limit=1&page=5", nil).WithContext(s.ctx)

		ids := []string{"secret1", "secret2"}

		s.secretStore.EXPECT().List(gomock.Any(), uint64(1), uint64(5)).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(http2.PageResponse{
			Data: ids,
			Paging: http2.PagePagingResponse{
				Previous: "example.com?limit=1&page=4",
				Next:     "example.com?limit=1&page=6",
			},
		})
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/SecretStore/secrets", nil).WithContext(s.ctx)

		s.secretStore.EXPECT().List(gomock.Any(), defaultPageSize, uint64(0)).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestListDeleted() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/SecretStore/secrets?deleted=true", nil).WithContext(s.ctx)

		ids := []string{"secret1", "secret2"}

		s.secretStore.EXPECT().ListDeleted(gomock.Any(), defaultPageSize, uint64(0)).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(http2.PageResponse{
			Data: ids,
		})
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/SecretStore/secrets?deleted=true", nil).WithContext(s.ctx)

		s.secretStore.EXPECT().ListDeleted(gomock.Any(), defaultPageSize, uint64(0)).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}
