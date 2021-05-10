package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/formatters"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/types/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/mocks"
	mockstoremanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/mock"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	testutils2 "github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	mocksecrets "github.com/ConsenSysQuorum/quorum-key-manager/src/store/secrets/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	secretStoreName = "SecretStore"
	secretID        = "my-secret"
)

type secretsHandlerTestSuite struct {
	suite.Suite
<<<<<<< HEAD
	secretStore *mocksecretstore.MockStore
=======
	secretStore *mocksecrets.MockStore
>>>>>>> refactor@store(mock): rename mock package
	router      *mux.Router
}

func TestSecretsHandler(t *testing.T) {
	s := new(secretsHandlerTestSuite)
	suite.Run(t, s)
}

func (s *secretsHandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	backend := mocks.NewMockBackend(ctrl)
	storeManager := mockstoremanager.NewMockStoreManager(ctrl)
<<<<<<< HEAD
	s.secretStore = mocksecretstore.NewMockStore(ctrl)
=======
	s.secretStore = mocksecrets.NewMockStore(ctrl)
>>>>>>> refactor@store(mock): rename mock package

	backend.EXPECT().StoreManager().Return(storeManager).AnyTimes()
	storeManager.EXPECT().GetSecretStore(gomock.Any(), secretStoreName).Return(s.secretStore, nil).AnyTimes()

	s.router = NewSecretsHandler(backend)
}

func (s *secretsHandlerTestSuite) TestSet() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		setSecretRequest := testutils.FakeSetSecretRequest()
		requestBytes, _ := json.Marshal(setSecretRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, secretStoreName)

		secret := testutils2.FakeSecret()

		s.secretStore.EXPECT().Set(gomock.Any(), setSecretRequest.ID, setSecretRequest.Value, &entities.Attributes{
			Tags: setSecretRequest.Tags,
		}).Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		setSecretRequest := testutils.FakeSetSecretRequest()
		requestBytes, _ := json.Marshal(setSecretRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, secretStoreName)

		s.secretStore.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultConnectionError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusFailedDependency, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestGet() {
	s.T().Run("should execute request successfully with version", func(t *testing.T) {
		version := "1"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s?version=%s", secretID, version), nil)
		httpRequest.Header.Set(StoreIDHeader, secretStoreName)

		secret := testutils2.FakeSecret()

		s.secretStore.EXPECT().Get(gomock.Any(), secretID, version).Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should execute request successfully without version", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", secretID), nil)
		httpRequest.Header.Set(StoreIDHeader, secretStoreName)

		secret := testutils2.FakeSecret()

		s.secretStore.EXPECT().Get(gomock.Any(), secretID, "").Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", secretID), nil)
		httpRequest.Header.Set(StoreIDHeader, secretStoreName)

		s.secretStore.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestList() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/", nil)
		httpRequest.Header.Set(StoreIDHeader, secretStoreName)

		ids := []string{"secret1", "secret2"}

		s.secretStore.EXPECT().List(gomock.Any()).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(ids)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/", nil)
		httpRequest.Header.Set(StoreIDHeader, secretStoreName)

		s.secretStore.EXPECT().List(gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}
