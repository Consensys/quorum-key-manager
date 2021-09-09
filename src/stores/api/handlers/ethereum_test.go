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
	http2 "github.com/consensys/quorum-key-manager/src/infra/http"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	apiTypes "github.com/consensys/quorum-key-manager/src/stores/api/types"
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
	ethStoreName = "EthStores"
	accAddress   = "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4"
)

var ethUserInfo = &types.UserInfo{
	Username:    "username",
	Roles:       []string{"role1", "role2"},
	Permissions: []types.Permission{"write:key", "read:key", "sign:key"},
}

var defaultPageSize = uint64(100)

type ethHandlerTestSuite struct {
	suite.Suite

	ctrl         *gomock.Controller
	storeManager *mockstoremanager.MockManager
	ethStore     *mockstoremanager.MockEthStore
	router       *mux.Router
	ctx          context.Context
}

func TestEthHandler(t *testing.T) {
	s := new(ethHandlerTestSuite)
	suite.Run(t, s)
}

func (s *ethHandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.storeManager = mockstoremanager.NewMockManager(s.ctrl)
	s.ethStore = mockstoremanager.NewMockEthStore(s.ctrl)
	s.ctx = authenticator.WithUserContext(context.Background(), &authenticator.UserContext{
		UserInfo: ethUserInfo,
	})

	s.storeManager.EXPECT().GetEthStore(gomock.Any(), ethStoreName, ethUserInfo).Return(s.ethStore, nil).AnyTimes()

	s.router = mux.NewRouter()
	NewStoresHandler(s.storeManager).Register(s.router)
}

func (s *ethHandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ethHandlerTestSuite) TestCreate() {
	s.Run("should execute request successfully", func() {
		createEthAccountRequest := testutils.FakeCreateEthAccountRequest()
		requestBytes, _ := json.Marshal(createEthAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/EthStores/ethereum", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETHAccount()

		s.ethStore.EXPECT().Create(
			gomock.Any(),
			createEthAccountRequest.KeyID,
			&entities.Attributes{
				Tags: createEthAccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEthAccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request without keyID successfully", func() {
		createEthAccountRequest := testutils.FakeCreateEthAccountRequest()
		createEthAccountRequest.KeyID = ""
		requestBytes, _ := json.Marshal(createEthAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/EthStores/ethereum", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETHAccount()

		s.ethStore.EXPECT().Create(
			gomock.Any(),
			gomock.Any(),
			&entities.Attributes{
				Tags: createEthAccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEthAccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with no request body successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/EthStores/ethereum", nil).WithContext(s.ctx)

		acc := testutils2.FakeETHAccount()

		s.ethStore.EXPECT().Create(gomock.Any(), gomock.Any(), &entities.Attributes{}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEthAccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		createEthAccountRequest := testutils.FakeCreateEthAccountRequest()
		requestBytes, _ := json.Marshal(createEthAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/EthStores/ethereum", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestImport() {
	s.Run("should execute request successfully", func() {
		importEthAccountRequest := testutils.FakeImportEthAccountRequest()
		requestBytes, _ := json.Marshal(importEthAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/EthStores/ethereum/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETHAccount()

		s.ethStore.EXPECT().Import(
			gomock.Any(),
			importEthAccountRequest.KeyID,
			importEthAccountRequest.PrivateKey,
			&entities.Attributes{
				Tags: importEthAccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEthAccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with without KeyID successfully", func() {
		importEthAccountRequest := testutils.FakeImportEthAccountRequest()
		importEthAccountRequest.KeyID = ""
		requestBytes, _ := json.Marshal(importEthAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/EthStores/ethereum/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETHAccount()

		s.ethStore.EXPECT().Import(
			gomock.Any(),
			gomock.Any(),
			importEthAccountRequest.PrivateKey,
			&entities.Attributes{
				Tags: importEthAccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEthAccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		importEthAccountRequest := testutils.FakeImportEthAccountRequest()
		requestBytes, _ := json.Marshal(importEthAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, "/stores/EthStores/ethereum/import", bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().Import(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestUpdate() {
	s.Run("should execute request successfully", func() {
		updateEthAccountRequest := testutils.FakeUpdateEthAccountRequest()
		requestBytes, _ := json.Marshal(updateEthAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/stores/%s/ethereum/%s", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		acc := testutils2.FakeETHAccount()

		s.ethStore.EXPECT().Update(
			gomock.Any(),
			ethcommon.HexToAddress(accAddress),
			&entities.Attributes{
				Tags: updateEthAccountRequest.Tags,
			}).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEthAccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		updateEthAccountRequest := testutils.FakeUpdateEthAccountRequest()
		requestBytes, _ := json.Marshal(updateEthAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/stores/%s/ethereum/%s", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestSignTypedData() {
	s.Run("should execute request successfully", func() {
		signTypedDataRequest := testutils.FakeSignTypedDataRequest()
		requestBytes, _ := json.Marshal(signTypedDataRequest)
		expectedTypedData := formatters.FormatSignTypedDataRequest(signTypedDataRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-typed-data", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signature := []byte("signature")
		s.ethStore.EXPECT().SignTypedData(gomock.Any(), ethcommon.HexToAddress(accAddress), expectedTypedData).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signature), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signTypedDataRequest := testutils.FakeSignTypedDataRequest()
		requestBytes, _ := json.Marshal(signTypedDataRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-typed-data", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().SignTypedData(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestSignMessage() {
	s.Run("should execute request successfully", func() {
		signEIP191Request := testutils.FakeSignMessageRequest()
		requestBytes, _ := json.Marshal(signEIP191Request)

		expectedSignature := "0xb91467e570a6466aa9e9876cbcd013baba02900b8979d43fe208a4a4f339f5fd6007e74cd82e037b800186422fc2da167c747ef045e5d18a5f5d4300f8e1a0291c"

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-message", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signature, _ := hexutil.Decode(expectedSignature)
		s.ethStore.EXPECT().SignMessage(gomock.Any(), ethcommon.HexToAddress(accAddress), gomock.Any()).Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		res := &apiTypes.SignMessageRequest{}

		_ = json.Unmarshal(rw.Body.Bytes(), res)

		assert.Equal(s.T(), hexutil.Encode(signature), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestSignTransaction() {
	s.Run("should execute request successfully", func() {
		signTransactionRequest := testutils.FakeSignETHTransactionRequest()
		requestBytes, _ := json.Marshal(signTransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signedRaw := []byte("signedRaw")
		s.ethStore.EXPECT().SignTransaction(gomock.Any(), ethcommon.HexToAddress(accAddress), signTransactionRequest.ChainID.ToInt(), gomock.Any()).Return(signedRaw, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signedRaw), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signTransactionRequest := testutils.FakeSignETHTransactionRequest()
		requestBytes, _ := json.Marshal(signTransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-transaction", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().SignTransaction(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestSignPrivateTransaction() {
	s.Run("should execute request successfully", func() {
		signPrivateTransactionRequest := testutils.FakeSignQuorumPrivateTransactionRequest()
		requestBytes, _ := json.Marshal(signPrivateTransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-quorum-private-transaction", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signedRaw := []byte("signedRaw")
		s.ethStore.EXPECT().SignPrivate(gomock.Any(), ethcommon.HexToAddress(accAddress), gomock.Any()).Return(signedRaw, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signedRaw), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signPrivateTransactionRequest := testutils.FakeSignQuorumPrivateTransactionRequest()
		requestBytes, _ := json.Marshal(signPrivateTransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-quorum-private-transaction", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().SignPrivate(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestSignEEATransaction() {
	s.Run("should execute request successfully", func() {
		signEEATransactionRequest := testutils.FakeSignEEATransactionRequest()
		requestBytes, _ := json.Marshal(signEEATransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-eea-transaction", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		signedRaw := []byte("signedRaw")
		s.ethStore.EXPECT().SignEEA(gomock.Any(), ethcommon.HexToAddress(accAddress), signEEATransactionRequest.ChainID.ToInt(), gomock.Any(), gomock.Any()).Return(signedRaw, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), hexutil.Encode(signedRaw), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		signEEATransactionRequest := testutils.FakeSignEEATransactionRequest()
		requestBytes, _ := json.Marshal(signEEATransactionRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/%s/sign-eea-transaction", ethStoreName, accAddress), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().SignEEA(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestGet() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/ethereum/%s", ethStoreName, accAddress), nil).WithContext(s.ctx)

		acc := testutils2.FakeETHAccount()
		s.ethStore.EXPECT().Get(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEthAccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request to get a deleted key successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/ethereum/%s?deleted=true", ethStoreName, accAddress), nil).WithContext(s.ctx)

		acc := testutils2.FakeETHAccount()
		s.ethStore.EXPECT().GetDeleted(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(acc, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatEthAccResponse(acc)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/ethereum/%s", ethStoreName, accAddress), nil).WithContext(s.ctx)

		s.ethStore.EXPECT().Get(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestList() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/ethereum", ethStoreName), nil).WithContext(s.ctx)

		acc1 := "0xfe3b557e8fb62b89f4916b721be55ceb828dbd73"
		acc2 := "0xea674fdde714fd979de3edf0f56aa9716b898ec8"
		s.ethStore.EXPECT().List(gomock.Any(), defaultPageSize, uint64(0)).Return([]ethcommon.Address{
			ethcommon.HexToAddress(acc1), ethcommon.HexToAddress(acc2),
		}, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(http2.PageResponse{
			Data: []string{acc1, acc2},
		})
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request with limit and offset successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/ethereum?limit=10&page=2", ethStoreName), nil).WithContext(s.ctx)

		acc1 := "0xfe3b557e8fb62b89f4916b721be55ceb828dbd74"
		acc2 := "0xea674fdde714fd979de3edf0f56aa9716b898ec9"
		s.ethStore.EXPECT().List(gomock.Any(), uint64(10), uint64(20)).Return([]ethcommon.Address{
			ethcommon.HexToAddress(acc1), ethcommon.HexToAddress(acc2),
		}, nil)

		s.router.ServeHTTP(rw, httpRequest)

		expectedBody, _ := json.Marshal(http2.PageResponse{
			Data: []string{acc1, acc2},
			Paging: http2.PagePagingResponse{
				Previous: "example.com?limit=10&page=1",
			},
		})
		assert.Equal(s.T(), string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	s.Run("should execute request to get a deleted key successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/ethereum?deleted=true", ethStoreName), nil).WithContext(s.ctx)

		acc1 := "0xfe3b557e8fb62b89f4916b721be55ceb828dbd75"
		acc2 := "0xea674fdde714fd979de3edf0f56aa9716b898ed9"
		s.ethStore.EXPECT().ListDeleted(gomock.Any(), defaultPageSize, uint64(0)).Return([]ethcommon.Address{
			ethcommon.HexToAddress(acc1), ethcommon.HexToAddress(acc2),
		}, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), fmt.Sprintf("{\"data\":[\"%s\",\"%s\"],\"paging\":{}}\n", acc1, acc2), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/stores/%s/ethereum", ethStoreName), nil).WithContext(s.ctx)

		s.ethStore.EXPECT().List(gomock.Any(), defaultPageSize, uint64(0)).Return(nil, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestDelete() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/%s/ethereum/%s", ethStoreName, accAddress), nil).WithContext(s.ctx)

		s.ethStore.EXPECT().Delete(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/%s/ethereum/%s", ethStoreName, accAddress), nil).WithContext(s.ctx)

		s.ethStore.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestDestroy() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/%s/ethereum/%s/destroy", ethStoreName, accAddress), nil).WithContext(s.ctx)

		s.ethStore.EXPECT().Destroy(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/stores/%s/ethereum/%s/destroy", ethStoreName, accAddress), nil).WithContext(s.ctx)

		s.ethStore.EXPECT().Destroy(gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestRestore() {
	s.Run("should execute request successfully", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/stores/%s/ethereum/%s/restore", ethStoreName, accAddress), nil).WithContext(s.ctx)

		s.ethStore.EXPECT().Restore(gomock.Any(), ethcommon.HexToAddress(accAddress)).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/stores/%s/ethereum/%s/restore", ethStoreName, accAddress), nil).WithContext(s.ctx)

		s.ethStore.EXPECT().Restore(gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestECRecover() {
	s.Run("should execute request successfully", func() {
		ecRecoverRequest := testutils.FakeECRecoverRequest()
		requestBytes, _ := json.Marshal(ecRecoverRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/ec-recover", ethStoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().ECRecover(gomock.Any(), ecRecoverRequest.Data, ecRecoverRequest.Signature).Return(ethcommon.HexToAddress(accAddress), nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), ethcommon.HexToAddress(accAddress).String(), rw.Body.String())
		assert.Equal(s.T(), http.StatusOK, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		ecRecoverRequest := testutils.FakeECRecoverRequest()
		requestBytes, _ := json.Marshal(ecRecoverRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/ec-recover", ethStoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().ECRecover(gomock.Any(), gomock.Any(), gomock.Any()).Return(ethcommon.Address{}, errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}

func (s *ethHandlerTestSuite) TestVerifyTypedDataSignature() {
	s.Run("should execute request successfully", func() {
		verifyRequest := testutils.FakeVerifyTypedDataPayloadRequest()
		requestBytes, _ := json.Marshal(verifyRequest)
		expectedTypedData := formatters.FormatSignTypedDataRequest(&verifyRequest.TypedData)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/verify-typed-data", ethStoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().VerifyTypedData(gomock.Any(), verifyRequest.Address, expectedTypedData, verifyRequest.Signature).Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(s.T(), "", rw.Body.String())
		assert.Equal(s.T(), http.StatusNoContent, rw.Code)
	})

	// Sufficient test to check that the mapping to HTTP errors is working. All other status code tests are done in integration tests
	s.Run("should fail with correct error code if use case fails", func() {
		verifyRequest := testutils.FakeVerifyTypedDataPayloadRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/stores/%s/ethereum/verify-typed-data", ethStoreName), bytes.NewReader(requestBytes)).WithContext(s.ctx)

		s.ethStore.EXPECT().VerifyTypedData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.HashicorpVaultError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(s.T(), http.StatusFailedDependency, rw.Code)
	})
}
