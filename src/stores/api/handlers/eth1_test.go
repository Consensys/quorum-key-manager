package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/api/types/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"

	mockstoremanager "github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

const (
	eth1StoreName = "Eth1Store"
	accAddress    = "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"
)

var eth1UserInfo = &types.UserInfo{
	Username: "username",
	Groups:   []string{"group1", "group2"},
}

type eth1HandlerTestSuite struct {
	suite.Suite

	ctrl         *gomock.Controller
	storeManager *mockstoremanager.MockManager
	eth1Store    *mockstoremanager.MockEth1Store
	router       *mux.Router
	ctx          context.Context
}

func TestEth1Handler(t *testing.T) {
	s := new(eth1HandlerTestSuite)
	suite.Run(t, s)
}

func (s *eth1HandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.storeManager = mockstoremanager.NewMockManager(s.ctrl)
	s.eth1Store = mockstoremanager.NewMockEth1Store(s.ctrl)
	s.ctx = authenticator.WithUserContext(context.Background(), &authenticator.UserContext{
		UserInfo: eth1UserInfo,
	})

	s.storeManager.EXPECT().GetEth1Store(gomock.Any(), eth1StoreName, eth1UserInfo).Return(s.eth1Store, nil).AnyTimes()

	s.router = mux.NewRouter()
	NewStoresHandler(s.storeManager).Register(s.router)
}

func (s *eth1HandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *eth1HandlerTestSuite) TestCreate() {
	s.Run("should execute request successfully", func() {
		createEth1AccountRequest := testutils.FakeCreateEth1AccountRequest()
		requestBytes, _ := json.Marshal(createEth1AccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/Eth1Store/eth1", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETH1Account()

		s.eth1Store.EXPECT().Create(
			gomock.Any(),
			createEth1AccountRequest.KeyID,
			&entities.Attributes{
				Tags: createEth1AccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEth1AccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request without keyID successfully", func() {
		createEth1AccountRequest := testutils.FakeCreateEth1AccountRequest()
		createEth1AccountRequest.KeyID = ""
		requestBytes, _ := json.Marshal(createEth1AccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/Eth1Store/eth1", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETH1Account()

		s.eth1Store.EXPECT().Create(
			gomock.Any(),
			gomock.Any(),
			&entities.Attributes{
				Tags: createEth1AccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEth1AccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with no request body successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/Eth1Store/eth1", nil).WithContext(s.ctx)

		acc := testutils2.FakeETH1Account()

		s.eth1Store.EXPECT().Create(gomock.Any(), gomock.Any(), &entities.Attributes{}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEth1AccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		createEth1AccountRequest := testutils.FakeCreateEth1AccountRequest()
		requestBytes, _ := json.Marshal(createEth1AccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/Eth1Store/eth1", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestImport() {
	s.Run("should execute request successfully", func() {
		importEth1AccountRequest := testutils.FakeImportEth1AccountRequest()
		requestBytes, _ := json.Marshal(importEth1AccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/Eth1Store/eth1/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETH1Account()

		s.eth1Store.EXPECT().Import(
			gomock.Any(),
			importEth1AccountRequest.KeyID,
			importEth1AccountRequest.PrivateKey,
			&entities.Attributes{
				Tags: importEth1AccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEth1AccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with without KeyID successfully", func() {
		importEth1AccountRequest := testutils.FakeImportEth1AccountRequest()
		importEth1AccountRequest.KeyID = ""
		requestBytes, _ := json.Marshal(importEth1AccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/Eth1Store/eth1/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETH1Account()

		s.eth1Store.EXPECT().Import(
			gomock.Any(),
			gomock.Any(),
			importEth1AccountRequest.PrivateKey,
			&entities.Attributes{
				Tags: importEth1AccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEth1AccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		importEth1AccountRequest := testutils.FakeImportEth1AccountRequest()
		requestBytes, _ := json.Marshal(importEth1AccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/Eth1Store/eth1/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().Import(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestUpdate() {
	s.Run("should execute request successfully", func() {
		updateEth1AccountRequest := testutils.FakeUpdateEth1AccountRequest()
		requestBytes, _ := json.Marshal(updateEth1AccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/stores/%s/eth1/%s", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETH1Account()

		s.eth1Store.EXPECT().Update(
			gomock.Any(),
			ethcommon.HexToAddress(accAddress),
			&entities.Attributes{
				Tags: updateEth1AccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEth1AccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		updateEth1AccountRequest := testutils.FakeUpdateEth1AccountRequest()
		requestBytes, _ := json.Marshal(updateEth1AccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/stores/%s/eth1/%s", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestSign() {
	s.Run("should execute request successfully", func() {
		signRequest := testutils.FakeSignHexPayloadRequest()
		requestBytes, _ := json.Marshal(signRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signature := []byte("signature")
		s.eth1Store.EXPECT().Sign(gomock.Any(), ethcommon.HexToAddress(accAddress), signRequest.Data).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signature), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signRequest := testutils.FakeSignHexPayloadRequest()
		requestBytes, _ := json.Marshal(signRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().Sign(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestSignTypedData() {
	s.Run("should execute request successfully", func() {
		signTypedDataRequest := testutils.FakeSignTypedDataRequest()
		requestBytes, _ := json.Marshal(signTypedDataRequest)
		expectedTypedData := formatters.FormatSignTypedDataRequest(signTypedDataRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign-typed-data", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signature := []byte("signature")
		s.eth1Store.EXPECT().SignTypedData(gomock.Any(), ethcommon.HexToAddress(accAddress), expectedTypedData).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signature), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signTypedDataRequest := testutils.FakeSignTypedDataRequest()
		requestBytes, _ := json.Marshal(signTypedDataRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign-typed-data", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().SignTypedData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestSignTransaction() {
	s.Run("should execute request successfully", func() {
		signTransactionRequest := testutils.FakeSignETHTransactionRequest()
		requestBytes, _ := json.Marshal(signTransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign-transaction", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signedRaw := []byte("signedRaw")
		s.eth1Store.EXPECT().SignTransaction(gomock.Any(), ethcommon.HexToAddress(accAddress), signTransactionRequest.ChainID.ToInt(), gomock.Any()).Return(signedRaw, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signedRaw), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signTransactionRequest := testutils.FakeSignETHTransactionRequest()
		requestBytes, _ := json.Marshal(signTransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign-transaction", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().SignTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestSignPrivateTransaction() {
	s.Run("should execute request successfully", func() {
		signPrivateTransactionRequest := testutils.FakeSignQuorumPrivateTransactionRequest()
		requestBytes, _ := json.Marshal(signPrivateTransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign-quorum-private-transaction", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signedRaw := []byte("signedRaw")
		s.eth1Store.EXPECT().SignPrivate(gomock.Any(), ethcommon.HexToAddress(accAddress), gomock.Any()).Return(signedRaw, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signedRaw), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signPrivateTransactionRequest := testutils.FakeSignQuorumPrivateTransactionRequest()
		requestBytes, _ := json.Marshal(signPrivateTransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign-quorum-private-transaction", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().SignPrivate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestSignEEATransaction() {
	s.Run("should execute request successfully", func() {
		signEEATransactionRequest := testutils.FakeSignEEATransactionRequest()
		requestBytes, _ := json.Marshal(signEEATransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign-eea-transaction", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signedRaw := []byte("signedRaw")
		s.eth1Store.EXPECT().SignEEA(gomock.Any(), ethcommon.HexToAddress(accAddress), signEEATransactionRequest.ChainID.ToInt(), gomock.Any(), gomock.Any()).Return(signedRaw, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signedRaw), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signEEATransactionRequest := testutils.FakeSignEEATransactionRequest()
		requestBytes, _ := json.Marshal(signEEATransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/%s/sign-eea-transaction", eth1StoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().SignEEA(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestGet() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/eth1/%s", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		acc := testutils2.FakeETH1Account()
		s.eth1Store.EXPECT().Get(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEth1AccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request to get a deleted key successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/eth1/%s?deleted=true", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		acc := testutils2.FakeETH1Account()
		s.eth1Store.EXPECT().GetDeleted(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEth1AccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/eth1/%s", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		s.eth1Store.EXPECT().Get(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestList() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/eth1", eth1StoreName), nil).WithContext(s.ctx)

		acc1 := "0xfe3b557e8fb62b89f4916b721be55ceb828dbd73"
		acc2 := "0xea674fdde714fd979de3edf0f56aa9716b898ec8"
		s.eth1Store.EXPECT().List(gomock.Any()).Return([]ethcommon.Address{
			ethcommon.HexToAddress(acc1), ethcommon.HexToAddress(acc2),
		}, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), fmt.Sprintf("[\"%s\",\"%s\"]\n", acc1, acc2), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request to get a deleted key successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/eth1?deleted=true", eth1StoreName), nil).WithContext(s.ctx)

		acc1 := "0xfe3b557e8fb62b89f4916b721be55ceb828dbd74"
		acc2 := "0xea674fdde714fd979de3edf0f56aa9716b898ec9"
		s.eth1Store.EXPECT().ListDeleted(gomock.Any()).Return([]ethcommon.Address{
			ethcommon.HexToAddress(acc1), ethcommon.HexToAddress(acc2),
		}, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), fmt.Sprintf("[\"%s\",\"%s\"]\n", acc1, acc2), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/eth1", eth1StoreName), nil).WithContext(s.ctx)

		s.eth1Store.EXPECT().List(gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestDelete() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/%s/eth1/%s", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		s.eth1Store.EXPECT().Delete(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/%s/eth1/%s", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		s.eth1Store.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestDestroy() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/%s/eth1/%s/destroy", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		s.eth1Store.EXPECT().Destroy(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/%s/eth1/%s/destroy", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		s.eth1Store.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestRestore() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/stores/%s/eth1/%s/restore", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		s.eth1Store.EXPECT().Restore(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/stores/%s/eth1/%s/restore", eth1StoreName, accAddress), nil).WithContext(s.ctx)

		s.eth1Store.EXPECT().Restore(gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestECRecover() {
	s.Run("should execute request successfully", func() {
		ecRecoverRequest := testutils.FakeECRecoverRequest()
		requestBytes, _ := json.Marshal(ecRecoverRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/ec-recover", eth1StoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().ECRecover(gomock.Any(), ecRecoverRequest.Data, ecRecoverRequest.Signature).Return(ethcommon.HexToAddress(accAddress), nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), ethcommon.HexToAddress(accAddress).String(), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		ecRecoverRequest := testutils.FakeECRecoverRequest()
		requestBytes, _ := json.Marshal(ecRecoverRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/ec-recover", eth1StoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().ECRecover(gomock.Any(), gomock.Any(), gomock.Any()).Return(ethcommon.Address{}, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestVerifySignature() {
	s.Run("should execute request successfully", func() {
		verifyRequest := testutils.FakeVerifyEth1SignatureRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/verify-signature", eth1StoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().Verify(gomock.Any(), verifyRequest.Address, verifyRequest.Data, verifyRequest.Signature).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		verifyRequest := testutils.FakeVerifyEth1SignatureRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/verify-signature", eth1StoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().Verify(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *eth1HandlerTestSuite) TestVerifyTypedDataSignature() {
	s.Run("should execute request successfully", func() {
		verifyRequest := testutils.FakeVerifyTypedDataPayloadRequest()
		requestBytes, _ := json.Marshal(verifyRequest)
		expectedTypedData := formatters.FormatSignTypedDataRequest(&verifyRequest.TypedData)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/verify-typed-data-signature", eth1StoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().VerifyTypedData(gomock.Any(), verifyRequest.Address, expectedTypedData, verifyRequest.Signature).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		verifyRequest := testutils.FakeVerifyTypedDataPayloadRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/eth1/verify-typed-data-signature", eth1StoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.eth1Store.EXPECT().VerifyTypedData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}
