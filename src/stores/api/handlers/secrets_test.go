package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/api/types/testutils"
	mockstoremanager "github.com/consensys/quorum-key-manager/src/stores/manager/mock"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	mocksecrets "github.com/consensys/quorum-key-manager/src/stores/store/secrets/mock"
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

	ctrl         *gomock.Controller
	storeManager *mockstoremanager.MockManager
	secretStore  *mocksecrets.MockStore
	router       *mux.Router
}

func TestSecretsHandler(t *testing.T) {
	s := new(secretsHandlerTestSuite)
	suite.Run(t, s)
}

func (s *secretsHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.storeManager = mockstoremanager.NewMockManager(s.ctrl)
	s.secretStore = mocksecrets.NewMockStore(s.ctrl)

	s.router = mux.NewRouter()
	NewStoresHandler(s.storeManager).Register(s.router)
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
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/SecretStore/secrets/"+secretID, bytes.NewReader(requestBytes))

		secret := testutils2.FakeSecret()

		s.storeManager.EXPECT().GetSecretStore(gomock.Any(), secretStoreName).Return(s.secretStore, nil)
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
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/SecretStore/secrets/"+secretID, bytes.NewReader(requestBytes))

		s.storeManager.EXPECT().GetSecretStore(gomock.Any(), secretStoreName).Return(s.secretStore, nil)
		s.secretStore.EXPECT().Set(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestGet() {
	s.Run("should execute request successfully with version", func() {
		version := "1"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s?version=%s", secretID, version), nil)

		secret := testutils2.FakeSecret()
		s.storeManager.EXPECT().GetSecretStore(gomock.Any(), secretStoreName).Return(s.secretStore, nil)
		s.secretStore.EXPECT().Get(gomock.Any(), secretID, version).Return(secret, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatSecretResponse(secret)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request successfully without version", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s", secretID), nil)

		secret := testutils2.FakeSecret()
		s.storeManager.EXPECT().GetSecretStore(gomock.Any(), secretStoreName).Return(s.secretStore, nil)
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
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/SecretStore/secrets/%s", secretID), nil)

		s.storeManager.EXPECT().GetSecretStore(gomock.Any(), secretStoreName).Return(s.secretStore, nil)
		s.secretStore.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *secretsHandlerTestSuite) TestList() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/SecretStore/secrets", nil)

		ids := []string{"secret1", "secret2"}

		s.storeManager.EXPECT().GetSecretStore(gomock.Any(), secretStoreName).Return(s.secretStore, nil)
		s.secretStore.EXPECT().List(gomock.Any()).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(ids)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/SecretStore/secrets", nil)

		s.storeManager.EXPECT().GetSecretStore(gomock.Any(), secretStoreName).Return(s.secretStore, nil)
		s.secretStore.EXPECT().List(gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}
