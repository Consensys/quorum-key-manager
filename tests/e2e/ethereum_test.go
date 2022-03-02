// +build e2e

package e2e

import (
	"fmt"
	utilstypes "github.com/consensys/quorum-key-manager/src/utils/api/types"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/api/types/testutils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ethTestSuite struct {
	suite.Suite
	err          error
	env          *Environment
	signAccount  *types.EthAccountResponse
	storeName    string
	deleteQueue  *sync.WaitGroup
	destroyQueue *sync.WaitGroup
}

func TestKeyManagerEth(t *testing.T) {
	s := new(ethTestSuite)

	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	env, err := NewEnvironment()
	require.NoError(t, err)
	s.env = env

	if len(s.env.cfg.EthStores) == 0 {
		t.Error("list of ethereum stores cannot be empty")
		return
	}

	s.deleteQueue = &sync.WaitGroup{}
	s.destroyQueue = &sync.WaitGroup{}

	for _, storeN := range s.env.cfg.EthStores {
		s.storeName = storeN
		s.env.logger = s.env.logger.WithComponent(storeN)
		suite.Run(t, s)
	}
}

func (s *ethTestSuite) SetupSuite() {
	var err error
	s.signAccount, err = s.env.client.CreateEthAccount(s.env.ctx, s.storeName, testutils.FakeCreateEthAccountRequest())
	require.NoError(s.T(), err)
}

func (s *ethTestSuite) TearDownSuite() {
	err := s.env.client.DeleteEthAccount(s.env.ctx, s.storeName, s.signAccount.Address.Hex())
	require.NoError(s.T(), err)

	time.Sleep(100 * time.Millisecond)

	_ = s.env.client.DestroyEthAccount(s.env.ctx, s.storeName, s.signAccount.Address.Hex())
}

func (s *ethTestSuite) TestCreate() {
	s.Run("should create a new account successfully", func() {
		request := testutils.FakeCreateEthAccountRequest()
		request.KeyID = "my-account-create"

		acc, err := s.env.client.CreateEthAccount(s.env.ctx, s.storeName, request)
		require.NoError(s.T(), err)
		defer s.queueToDelete(acc)

		assert.NotEmpty(s.T(), acc.Address)
		assert.NotEmpty(s.T(), acc.PublicKey)
		assert.NotEmpty(s.T(), acc.CompressedPublicKey)
		assert.Equal(s.T(), request.KeyID, acc.KeyID)
		assert.Equal(s.T(), request.Tags, acc.Tags)
		assert.False(s.T(), acc.Disabled)
		assert.NotEmpty(s.T(), acc.CreatedAt)
		assert.NotEmpty(s.T(), acc.UpdatedAt)
	})

	s.Run("should create a new account with random keyID successfully", func() {
		request := testutils.FakeCreateEthAccountRequest()
		request.KeyID = ""

		acc, err := s.env.client.CreateEthAccount(s.env.ctx, s.storeName, request)
		require.NoError(s.T(), err)
		defer s.queueToDelete(acc)

		assert.NotEmpty(s.T(), acc.Address)
		assert.NotEmpty(s.T(), acc.PublicKey)
		assert.NotEmpty(s.T(), acc.CompressedPublicKey)
		assert.NotEmpty(s.T(), acc.KeyID)
		assert.Equal(s.T(), request.Tags, acc.Tags)
		assert.False(s.T(), acc.Disabled)
		assert.NotEmpty(s.T(), acc.CreatedAt)
		assert.NotEmpty(s.T(), acc.UpdatedAt)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeCreateEthAccountRequest()

		key, err := s.env.client.CreateEthAccount(s.env.ctx, "nonExistentStoreName", request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), http.StatusNotFound, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestUpdate() {
	request := testutils.FakeCreateEthAccountRequest()
	request.KeyID = "my-account-create"

	acc, err := s.env.client.CreateEthAccount(s.env.ctx, s.storeName, request)
	require.NoError(s.T(), err)
	defer s.queueToDelete(acc)

	s.Run("should update an existing account successfully", func() {
		newTags := map[string]string{
			"tagnew": "valuenew",
		}
		acc2, err := s.env.client.UpdateEthAccount(s.env.ctx, s.storeName, acc.Address.Hex(), &types.UpdateEthAccountRequest{
			Tags: newTags,
		})

		require.NoError(s.T(), err)
		assert.Equal(s.T(), newTags, acc2.Tags)
	})
}

func (s *ethTestSuite) TestImport() {
	s.Run("should import an account successfully", func() {
		privKey, _ := crypto.GenerateKey()
		request := testutils.FakeImportEthAccountRequest()
		request.PrivateKey = privKey.D.Bytes()
		request.KeyID = "my-account-import"

		acc, err := s.env.client.ImportEthAccount(s.env.ctx, s.storeName, request)
		require.NoError(s.T(), err)
		defer s.queueToDelete(acc)

		assert.NotEmpty(s.T(), acc.Address)
		assert.NotEmpty(s.T(), acc.PublicKey)
		assert.NotEmpty(s.T(), acc.CompressedPublicKey)
		assert.Equal(s.T(), request.KeyID, acc.KeyID)
		assert.Equal(s.T(), request.Tags, acc.Tags)
		assert.False(s.T(), acc.Disabled)
		assert.NotEmpty(s.T(), acc.CreatedAt)
		assert.NotEmpty(s.T(), acc.UpdatedAt)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeImportEthAccountRequest()

		key, err := s.env.client.ImportEthAccount(s.env.ctx, "inexistentStoreName", request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSignMessage() {
	s.Run("should sign a payload successfully and verify it", func() {
		request := testutils.FakeSignMessageRequest()

		signature, err := s.env.client.SignMessage(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)

		hexSig, err := hexutil.Decode(signature)
		require.NoError(s.T(), err)

		err = s.env.client.VerifyMessage(s.env.ctx, &utilstypes.VerifyRequest{
			Data:      request.Message,
			Signature: hexSig,
			Address:   s.signAccount.Address,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignMessageRequest()

		signature, err := s.env.client.SignMessage(s.env.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSignTypedData() {
	s.Run("should sign typed data successfully and verify it", func() {
		request := testutils.FakeSignTypedDataRequest()

		signature, err := s.env.client.SignTypedData(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)

		hexSig, err := hexutil.Decode(signature)
		require.NoError(s.T(), err)

		err = s.env.client.VerifyTypedData(s.env.ctx, &utilstypes.VerifyTypedDataRequest{
			TypedData: *request,
			Signature: hexSig,
			Address:   s.signAccount.Address,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignTypedDataRequest()

		signature, err := s.env.client.SignTypedData(s.env.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSignTransaction() {
	s.Run("should sign transaction successfully (default type: DYNAMIC_FEE)", func() {
		request := testutils.FakeSignETHTransactionRequest("")

		signedTx, err := s.env.client.SignTransaction(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), signedTx)
	})

	s.Run("should sign DYNAMIC_FEE transaction successfully", func() {
		request := testutils.FakeSignETHTransactionRequest(types.DynamicFeeTxType)

		signedTx, err := s.env.client.SignTransaction(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), signedTx)
	})

	s.Run("should sign ACCESS_LIST transaction successfully", func() {
		request := testutils.FakeSignETHTransactionRequest(types.AccessListTxType)

		signedTx, err := s.env.client.SignTransaction(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), signedTx)
	})

	s.Run("should sign LEGACY transaction successfully", func() {
		request := testutils.FakeSignETHTransactionRequest(types.LegacyTxType)

		signedTx, err := s.env.client.SignTransaction(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), signedTx)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignETHTransactionRequest("")

		signature, err := s.env.client.SignTransaction(s.env.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})

	s.Run("should sign a big amount of transactions successfully", func() {
		for i := 0; i < 500; i++ {
			request := testutils.FakeSignETHTransactionRequest("")

			signedTx, err := s.env.client.SignTransaction(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
			require.NoError(s.T(), err)
			assert.NotNil(s.T(), signedTx)
		}
	})
}

func (s *ethTestSuite) TestSignPrivateTransaction() {
	s.Run("should sign private transaction successfully", func() {
		request := testutils.FakeSignQuorumPrivateTransactionRequest()

		signature, err := s.env.client.SignQuorumPrivateTransaction(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignQuorumPrivateTransactionRequest()

		signature, err := s.env.client.SignQuorumPrivateTransaction(s.env.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSignEEATransaction() {
	s.Run("should sign private transaction successfully", func() {
		request := testutils.FakeSignEEATransactionRequest()

		signature, err := s.env.client.SignEEATransaction(s.env.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignEEATransactionRequest()

		signature, err := s.env.client.SignEEATransaction(s.env.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestGetEthAccount() {
	s.Run("should sign private transaction successfully", func() {
		retrievedAcc, err := s.env.client.GetEthAccount(s.env.ctx, s.storeName, s.signAccount.Address.Hex())
		require.NoError(s.T(), err)

		assert.Equal(s.T(), s.signAccount.Address, retrievedAcc.Address)
		assert.Equal(s.T(), s.signAccount.PublicKey, retrievedAcc.PublicKey)
		assert.Equal(s.T(), s.signAccount.CompressedPublicKey, retrievedAcc.CompressedPublicKey)
		assert.Equal(s.T(), s.signAccount.KeyID, retrievedAcc.KeyID)
		assert.Equal(s.T(), s.signAccount.Tags, retrievedAcc.Tags)
		assert.Equal(s.T(), s.signAccount.Disabled, retrievedAcc.Disabled)
	})

	s.Run("should fail if account does not exist", func() {
		key, err := s.env.client.GetEthAccount(s.env.ctx, s.storeName, "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})

	s.Run("should parse errors successfully", func() {
		key, err := s.env.client.GetEthAccount(s.env.ctx, "inexistentStoreName", s.signAccount.Address.Hex())
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestListEthAccounts() {
	s.Run("should sign private transaction successfully", func() {
		accounts, err := s.env.client.ListEthAccounts(s.env.ctx, s.storeName, 999999, 0)
		require.NoError(s.T(), err)

		assert.Contains(s.T(), accounts, strings.ToLower(s.signAccount.Address.Hex()))
	})

	s.Run("should parse errors successfully", func() {
		key, err := s.env.client.ListEthAccounts(s.env.ctx, "inexistentStoreName", 0, 0)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) queueToDelete(accR *types.EthAccountResponse) {
	s.deleteQueue.Add(1)
	go func() {
		err := s.env.client.DeleteEthAccount(s.env.ctx, s.storeName, accR.Address.Hex())
		if err != nil {
			s.T().Logf("failed to delete eth account {Address: %s}", accR.Address.String())
		} else {
			s.queueToDestroy(accR)
		}
		s.deleteQueue.Done()
	}()
}

func (s *ethTestSuite) queueToDestroy(accR *types.EthAccountResponse) {
	s.destroyQueue.Add(1)
	go func() {
		errMsg := fmt.Sprintf("failed to destroy eth account {Address: %s}", accR.Address.String())
		err := retryOn(func() error {
			return s.env.client.DestroyEthAccount(s.env.ctx, s.storeName, accR.Address.Hex())
		}, s.T().Logf, errMsg, http.StatusConflict, MaxRetries)

		if err != nil {
			s.T().Logf(errMsg)
		}
		s.destroyQueue.Done()
	}()
}
