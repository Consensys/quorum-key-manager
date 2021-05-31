package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
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
	mockkeys "github.com/ConsenSysQuorum/quorum-key-manager/src/store/keys/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

/*
const (
	eth1StoreName = "Eth1Store"
	accountID     = "my-account"
)
*/

type eth1HandlerTestSuite struct {
	suite.Suite
	keyStore *mockkeys.MockStore
	router   *mux.Router
}

func TestEth1Handler(t *testing.T) {
	s := new(eth1HandlerTestSuite)
	suite.Run(t, s)
}

func (s *eth1HandlerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	backend := mocks.NewMockBackend(ctrl)
	storeManager := mockstoremanager.NewMockStoreManager(ctrl)
	s.keyStore = mockkeys.NewMockStore(ctrl)

	backend.EXPECT().StoreManager().Return(storeManager).AnyTimes()
	storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil).AnyTimes()

	s.router = NewKeysHandler(backend)
}

func (s *eth1HandlerTestSuite) TestCreate() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		key := testutils2.FakeKey()

		s.keyStore.EXPECT().Create(
			gomock.Any(),
			createKeyRequest.ID,
			&entities.Algorithm{
				Type:          entities.KeyType(createKeyRequest.SigningAlgorithm),
				EllipticCurve: entities.Curve(createKeyRequest.Curve),
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
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 400 if curve is not supported", func(t *testing.T) {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		createKeyRequest.Curve = invalidCurve
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		createKeyRequest := testutils.FakeCreateKeyRequest()
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.keyStore.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultConnectionError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestImport() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		importKeyRequest := testutils.FakeImportKeyRequest()
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/import", bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		key := testutils2.FakeKey()

		privKey, _ := base64.URLEncoding.DecodeString(importKeyRequest.PrivateKey)
		s.keyStore.EXPECT().Import(
			gomock.Any(),
			importKeyRequest.ID,
			privKey,
			&entities.Algorithm{
				Type:          entities.KeyType(importKeyRequest.SigningAlgorithm),
				EllipticCurve: entities.Curve(importKeyRequest.Curve),
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
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	s.T().Run("should fail with 400 if curve is not supported", func(t *testing.T) {
		importKeyRequest := testutils.FakeImportKeyRequest()
		importKeyRequest.Curve = invalidCurve
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/import", bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		importKeyRequest := testutils.FakeImportKeyRequest()
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/import", bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.keyStore.EXPECT().Import(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestSign() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		signPayloadRequest := testutils.FakeSignPayloadRequest()
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/sign", keyID), bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		signature := []byte("signature")
		data, _ := base64.URLEncoding.DecodeString(signPayloadRequest.Data)
		s.keyStore.EXPECT().Sign(gomock.Any(), keyID, data).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, base64.URLEncoding.EncodeToString(signature), rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with 400 if payload is not hexadecimal", func(t *testing.T) {
		signPayloadRequest := testutils.FakeSignPayloadRequest()
		signPayloadRequest.Data = "invalidData"
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/sign", keyID), bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		signPayloadRequest := testutils.FakeSignPayloadRequest()
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/%s/sign", keyID), bytes.NewReader(requestBytes))
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.keyStore.EXPECT().Sign(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestGet() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", keyID), nil)
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		key := testutils2.FakeKey()

		s.keyStore.EXPECT().Get(gomock.Any(), keyID).Return(key, nil)

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
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.keyStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestList() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/", nil)
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

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
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.keyStore.EXPECT().List(gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestDestroy() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", keyID), nil)
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, "", rw.Body.String())
		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	s.T().Run("should execute request successfully with version", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", keyID), nil)
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, "", rw.Body.String())
		assert.Equal(t, http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/%s", keyID), nil)
		httpRequest = httpRequest.WithContext(context.WithValue(httpRequest.Context(), StoreContextID, keyStoreName))

		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusNotFound, rw.Code)
	})
}
