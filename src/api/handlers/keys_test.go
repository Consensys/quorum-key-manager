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
	mocks2 "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/mocks"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	testutils2 "github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	mocks3 "github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/mocks"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	keyStoreName            = "KeyStore"
	invalidSigningAlgorithm = "invalidSigningAlgorithm"
	invalidCurve            = "invalidCurve"
	keyID                   = "my-key"
)

type keysHandlerTestSuite struct {
	suite.Suite
	keyStore *mocks3.MockKeyStore
	router   *mux.Router
}

func TestKeysHandler(t *testing.T) {
	s := new(keysHandlerTestSuite)
	suite.Run(t, s)
}

func (s *keysHandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	backend := mocks.NewMockBackend(ctrl)
	storeManager := mocks2.NewMockStoreManager(ctrl)
	s.keyStore = mocks3.NewMockKeyStore(ctrl)

	backend.EXPECT().StoreManager().Return(storeManager).AnyTimes()
	storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil).AnyTimes()

	s.router = NewKeysHandler(backend)
}

func (s *keysHandlerTestSuite) TestCreate() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		key := testutils2.FakeKey()

		s.keyStore.EXPECT().Create(
			gomock.Any(),
			createKeyRequest.ID,
			&entities.Algorithm{
				Type:          createKeyRequest.SigningAlgorithm,
				EllipticCurve: createKeyRequest.Curve,
			},
			&entities.Attributes{
				Tags: createKeyRequest.Tags,
			}).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with 400 if signing algorithm is not supported", func(t *testing.T) {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		createKeyRequest.SigningAlgorithm = invalidSigningAlgorithm
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 400 if curve is not supported", func(t *testing.T) {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		createKeyRequest.Curve = invalidCurve
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.keyStore.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultConnectionError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusFailedDependency, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestImport() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		importKeyRequest := testutils.FakeImportKeyRequest()
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/import", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		key := testutils2.FakeKey()

		s.keyStore.EXPECT().Import(
			gomock.Any(),
			importKeyRequest.ID,
			importKeyRequest.PrivateKey,
			&entities.Algorithm{
				Type:          importKeyRequest.SigningAlgorithm,
				EllipticCurve: importKeyRequest.Curve,
			},
			&entities.Attributes{
				Tags: importKeyRequest.Tags,
			}).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with 400 if signing algorithm is not supported", func(t *testing.T) {
		importKeyRequest := testutils.FakeImportKeyRequest()
		importKeyRequest.SigningAlgorithm = invalidSigningAlgorithm
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/import", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 400 if curve is not supported", func(t *testing.T) {
		importKeyRequest := testutils.FakeImportKeyRequest()
		importKeyRequest.Curve = invalidCurve
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/import", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		importKeyRequest := testutils.FakeImportKeyRequest()
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/import", bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.keyStore.EXPECT().Import(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestSign() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		signPayloadRequest := testutils.FakeSignPayloadRequest()
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/sign", keyID), bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		signature := "0xsignature"

		s.keyStore.EXPECT().Sign(
			gomock.Any(),
			keyID,
			signPayloadRequest.Data,
			signPayloadRequest.Version,
		).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, signature, rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with 400 if payload is not hexadecimal", func(t *testing.T) {
		signPayloadRequest := testutils.FakeSignPayloadRequest()
		signPayloadRequest.Data = "invalidData"
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/sign", keyID), bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		signPayloadRequest := testutils.FakeSignPayloadRequest()
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/sign", keyID), bytes.NewReader(requestBytes))
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.keyStore.EXPECT().Sign(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestGet() {
	s.T().Run("should execute request successfully with version", func(t *testing.T) {
		version := "1"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s?version=%s", keyID, version), nil)
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		key := testutils2.FakeKey()

		s.keyStore.EXPECT().Get(gomock.Any(), keyID, version).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should execute request successfully without version", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", keyID), nil)
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		key := testutils2.FakeKey()

		s.keyStore.EXPECT().Get(gomock.Any(), keyID, "").Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", keyID), nil)
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.keyStore.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestList() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/", nil)
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		ids := []string{"key1", "key2"}

		s.keyStore.EXPECT().List(gomock.Any()).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(ids)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/", nil)
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.keyStore.EXPECT().List(gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestDestroy() {
	version := "1"

	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", keyID), nil)
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID, "").Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, "", rw.Body.String())
		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	s.T().Run("should execute request successfully with version", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s?version=%s", keyID, version), nil)
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID, version).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, "", rw.Body.String())
		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", keyID), nil)
		httpRequest.Header.Set(StoreIDHeader, keyStoreName)

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID, "").Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}
