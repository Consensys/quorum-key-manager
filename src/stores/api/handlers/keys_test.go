package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	formatters2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/formatters"
	testutils3 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/api/types/testutils"
	mock2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager/mock"
	entities2 "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/keys/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
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

	ctrl         *gomock.Controller
	storeManager *mock2.MockManager
	keyStore     *mock.MockStore
	router       *mux.Router
}

func TestKeysHandler(t *testing.T) {
	s := new(keysHandlerTestSuite)
	suite.Run(t, s)
}

func (s *keysHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.storeManager = mock2.NewMockManager(s.ctrl)
	s.keyStore = mock.NewMockStore(s.ctrl)

	s.router = mux.NewRouter()
	NewStoresHandler(s.storeManager).Register(s.router)
}

func (s *keysHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *keysHandlerTestSuite) TestCreate() {
	s.Run("should execute request successfully", func() {
		createKeyRequest := testutils3.FakeCreateKeyRequest()
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys", bytes.NewReader(requestBytes))

		key := testutils.FakeKey()

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().Create(
			gomock.Any(),
			createKeyRequest.ID,
			&entities2.Algorithm{
				Type:          entities2.KeyType(createKeyRequest.SigningAlgorithm),
				EllipticCurve: entities2.Curve(createKeyRequest.Curve),
			},
			&entities2.Attributes{
				Tags: createKeyRequest.Tags,
			}).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters2.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should fail with 400 if signing algorithm is not supported", func() {
		createKeyRequest := testutils3.FakeCreateKeyRequest()
		createKeyRequest.SigningAlgorithm = invalidSigningAlgorithm
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys", bytes.NewReader(requestBytes))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	s.Run("should fail with 400 if curve is not supported", func() {
		createKeyRequest := testutils3.FakeCreateKeyRequest()
		createKeyRequest.Curve = invalidCurve
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys", bytes.NewReader(requestBytes))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		createKeyRequest := testutils3.FakeCreateKeyRequest()
		requestBytes, _ := json.Marshal(createKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys", bytes.NewReader(requestBytes))

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultConnectionError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestImport() {
	s.Run("should execute request successfully", func() {
		importKeyRequest := testutils3.FakeImportKeyRequest()
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/import", bytes.NewReader(requestBytes))

		key := testutils.FakeKey()
		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		privKey, _ := base64.URLEncoding.DecodeString(importKeyRequest.PrivateKey)
		s.keyStore.EXPECT().Import(
			gomock.Any(),
			importKeyRequest.ID,
			privKey,
			&entities2.Algorithm{
				Type:          entities2.KeyType(importKeyRequest.SigningAlgorithm),
				EllipticCurve: entities2.Curve(importKeyRequest.Curve),
			},
			&entities2.Attributes{
				Tags: importKeyRequest.Tags,
			}).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters2.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should fail with 400 if signing algorithm is not supported", func() {
		importKeyRequest := testutils3.FakeImportKeyRequest()
		importKeyRequest.SigningAlgorithm = invalidSigningAlgorithm
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/import", bytes.NewReader(requestBytes))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	s.Run("should fail with 400 if curve is not supported", func() {
		importKeyRequest := testutils3.FakeImportKeyRequest()
		importKeyRequest.Curve = invalidCurve
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/import", bytes.NewReader(requestBytes))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		importKeyRequest := testutils3.FakeImportKeyRequest()
		requestBytes, _ := json.Marshal(importKeyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/KeyStore/keys/import", bytes.NewReader(requestBytes))

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().Import(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestSign() {
	s.Run("should execute request successfully", func() {
		signPayloadRequest := testutils3.FakeSignPayloadRequest()
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/KeyStore/keys/%s/sign", keyID), bytes.NewReader(requestBytes))

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)

		signature := []byte("signature")
		data, _ := base64.URLEncoding.DecodeString(signPayloadRequest.Data)
		s.keyStore.EXPECT().Sign(gomock.Any(), keyID, data).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), base64.URLEncoding.EncodeToString(signature), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should fail with 400 if payload is not hexadecimal", func() {
		signPayloadRequest := testutils3.FakeSignPayloadRequest()
		signPayloadRequest.Data = "invalidData"
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/KeyStore/keys/%s/sign", keyID), bytes.NewReader(requestBytes))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusBadRequest, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signPayloadRequest := testutils3.FakeSignPayloadRequest()
		requestBytes, _ := json.Marshal(signPayloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/KeyStore/keys/%s/sign", keyID), bytes.NewReader(requestBytes))

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().Sign(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestGet() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil)

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		key := testutils.FakeKey()
		s.keyStore.EXPECT().Get(gomock.Any(), keyID).Return(key, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters2.FormatKeyResponse(key)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil)

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestList() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/KeyStore/keys", nil)

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)

		ids := []string{"key1", "key2"}
		s.keyStore.EXPECT().List(gomock.Any()).Return(ids, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(ids)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, "/stores/KeyStore/keys", nil)

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().List(gomock.Any()).Return(nil, errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}

func (s *keysHandlerTestSuite) TestDestroy() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil)

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	s.Run("should execute request successfully with version", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil)
		httpRequest = httpRequest.WithContext(WithStoreName(httpRequest.Context(), keyStoreName))

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/KeyStore/keys/%s", keyID), nil)

		s.storeManager.EXPECT().GetKeyStore(gomock.Any(), keyStoreName).Return(s.keyStore, nil)
		s.keyStore.EXPECT().Destroy(gomock.Any(), keyID).Return(errors.NotFoundError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusNotFound, rw.Code)
	})
}
