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
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type utilsHandlerTestSuite struct {
	suite.Suite

	ctrl      *gomock.Controller
	utilities *mock.MockUtils
	router    *mux.Router
}

func TestUtilsHandler(t *testing.T) {
	s := new(utilsHandlerTestSuite)
	suite.Run(t, s)
}

func (s *utilsHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	manager := mock.NewMockManager(s.ctrl)
	s.utilities = mock.NewMockUtils(s.ctrl)

	manager.EXPECT().Utilities().Return(s.utilities, nil).AnyTimes()

	s.router = mux.NewRouter()
	NewStoresHandler(manager).Register(s.router)
}

func (s *utilsHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *utilsHandlerTestSuite) TestECRecover() {
	s.Run("should execute request successfully", func() {
		ecRecoverRequest := testutils.FakeECRecoverRequest()
		requestBytes, _ := json.Marshal(ecRecoverRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/ethereum/ec-recover", bytes.NewReader(requestBytes))

		s.utilities.EXPECT().ECRecover(ecRecoverRequest.Data, ecRecoverRequest.Signature).Return(ethcommon.HexToAddress(accAddress), nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), ethcommon.HexToAddress(accAddress).String(), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		ecRecoverRequest := testutils.FakeECRecoverRequest()
		requestBytes, _ := json.Marshal(ecRecoverRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/ethereum/ec-recover", bytes.NewReader(requestBytes))

		s.utilities.EXPECT().ECRecover(gomock.Any(), gomock.Any()).Return(ethcommon.Address{}, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *utilsHandlerTestSuite) TestVerifyMessage() {
	s.Run("should execute request successfully", func() {
		verifyRequest := testutils.FakeVerifyRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/verify-message", ethStoreName), bytes.NewReader(requestBytes))

		s.utilities.EXPECT().VerifyMessage(verifyRequest.Address, verifyRequest.Data, verifyRequest.Signature).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		verifyRequest := testutils.FakeVerifyRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/verify-message", ethStoreName), bytes.NewReader(requestBytes))

		s.utilities.EXPECT().VerifyMessage(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *utilsHandlerTestSuite) TestVerifyTypedData() {
	s.Run("should execute request successfully", func() {
		verifyRequest := testutils.FakeVerifyTypedDataPayloadRequest()
		requestBytes, _ := json.Marshal(verifyRequest)
		expectedTypedData := formatters.FormatSignTypedDataRequest(&verifyRequest.TypedData)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/ethereum/verify-typed-data", bytes.NewReader(requestBytes))

		s.utilities.EXPECT().VerifyTypedData(verifyRequest.Address, expectedTypedData, verifyRequest.Signature).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		verifyRequest := testutils.FakeVerifyTypedDataPayloadRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "ethereum/verify-typed-data", bytes.NewReader(requestBytes))

		s.utilities.EXPECT().VerifyTypedData(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}
