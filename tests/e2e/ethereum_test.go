// +build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	"github.com/consensys/quorum-key-manager/src/infra/log/zap"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/api/types/testutils"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ethTestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient client.EthClient
	signAccount      *types.EthAccountResponse

	storeName string
	logger    log.Logger

	deleteQueue  *sync.WaitGroup
	destroyQueue *sync.WaitGroup
}

func TestKeyManagerEth(t *testing.T) {
	s := new(ethTestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	var cfg *tests.Config
	cfg, s.err = tests.NewConfig()
	if s.err != nil {
		t.Error(s.err)
		return
	}

	if len(cfg.SecretStores) == 0 {
		t.Error("list of secret stores cannot be empty")
		return
	}

	s.logger, s.err = zap.NewLogger(log.NewConfig(log.WarnLevel, log.TextFormat))
	if s.err != nil {
		t.Error(s.err)
		return
	}

	s.deleteQueue = &sync.WaitGroup{}
	s.destroyQueue = &sync.WaitGroup{}

	var token string
	token, s.err = generateJWT("./certificates/client.key", "*:*", "e2e|eth_test")
	if s.err != nil {
		t.Errorf("failed to generate jwt. %s", s.err)
		return
	}
	s.keyManagerClient = client.NewHTTPClient(&http.Client{
		Transport: NewTestHttpTransport(token, "", nil),
	}, &client.Config{
		URL: cfg.KeyManagerURL,
	})

	for _, storeN := range cfg.EthStores {
		s.storeName = storeN
		s.logger = s.logger.WithComponent(storeN)
		suite.Run(t, s)
	}
}

func (s *ethTestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	var err error
	s.signAccount, err = s.keyManagerClient.CreateEthAccount(s.ctx, s.storeName, testutils.FakeCreateEthAccountRequest())
	require.NoError(s.T(), err)
}

func (s *ethTestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func (s *ethTestSuite) TestCreate() {
	s.Run("should create a new account successfully", func() {
		request := testutils.FakeCreateEthAccountRequest()
		request.KeyID = "my-account-create"

		acc, err := s.keyManagerClient.CreateEthAccount(s.ctx, s.storeName, request)
		require.NoError(s.T(), err)

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

		acc, err := s.keyManagerClient.CreateEthAccount(s.ctx, s.storeName, request)
		require.NoError(s.T(), err)

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

		key, err := s.keyManagerClient.CreateEthAccount(s.ctx, "inexistentStoreName", request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), http.StatusNotFound, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestUpdate() {
	request := testutils.FakeCreateEthAccountRequest()
	request.KeyID = "my-account-create"

	acc, err := s.keyManagerClient.CreateEthAccount(s.ctx, s.storeName, request)
	require.NoError(s.T(), err)
	defer s.queueToDelete(acc)

	s.Run("should update an existing account successfully", func() {
		newTags := map[string]string{
			"tagnew": "valuenew",
		}
		acc2, err := s.keyManagerClient.UpdateEthAccount(s.ctx, s.storeName, acc.Address.Hex(), &types.UpdateEthAccountRequest{
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

		acc, err := s.keyManagerClient.ImportEthAccount(s.ctx, s.storeName, request)
		require.NoError(s.T(), err)

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

		key, err := s.keyManagerClient.ImportEthAccount(s.ctx, "inexistentStoreName", request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSign() {
	s.Run("should sign a payload successfully and verify it", func() {
		request := testutils.FakeSignMessageRequest()

		signature, err := s.keyManagerClient.SignMessage(s.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)

		hexSig, err := hexutil.Decode(signature)
		require.NoError(s.T(), err)

		err = s.keyManagerClient.VerifyMessage(s.ctx, s.storeName, &types.VerifyRequest{
			Data:      request.Message,
			Signature: hexSig,
			Address:   s.signAccount.Address,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignMessageRequest()

		signature, err := s.keyManagerClient.SignMessage(s.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSignTypedData() {
	s.Run("should sign typed data successfully and verify it", func() {
		request := testutils.FakeSignTypedDataRequest()

		signature, err := s.keyManagerClient.SignTypedData(s.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)

		hexSig, err := hexutil.Decode(signature)
		require.NoError(s.T(), err)

		err = s.keyManagerClient.VerifyTypedData(s.ctx, s.storeName, &types.VerifyTypedDataRequest{
			TypedData: *request,
			Signature: hexSig,
			Address:   s.signAccount.Address,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignTypedDataRequest()

		signature, err := s.keyManagerClient.SignTypedData(s.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSignTransaction() {
	s.Run("should sign transaction successfully", func() {
		request := testutils.FakeSignETHTransactionRequest()

		signedTx, err := s.keyManagerClient.SignTransaction(s.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)
		assert.NotNil(s.T(), signedTx)

		signer := ethtypes.NewEIP155Signer(request.ChainID.ToInt())

		tx := formatters.FormatTransaction(request)
		txData := signer.Hash(tx).Bytes()

		err = rlp.DecodeBytes(hexutil.MustDecode(signedTx), &tx)
		require.NoError(s.T(), err)
		v_, r_, s_ := tx.RawSignatureValues()
		sig := append(append(r_.Bytes(), s_.Bytes()...), v_.Bytes()...)

		err = s.keyManagerClient.Verify(s.ctx, s.storeName, &types.VerifyRequest{
			Data:      txData,
			Signature: sig,
			Address:   s.signAccount.Address,
		})
		// require.NoError(s.T(), err)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignETHTransactionRequest()

		signature, err := s.keyManagerClient.SignTransaction(s.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSignPrivateTransaction() {
	s.Run("should sign private transaction successfully", func() {
		request := testutils.FakeSignQuorumPrivateTransactionRequest()

		signature, err := s.keyManagerClient.SignQuorumPrivateTransaction(s.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignQuorumPrivateTransactionRequest()

		signature, err := s.keyManagerClient.SignQuorumPrivateTransaction(s.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestSignEEATransaction() {
	s.Run("should sign private transaction successfully", func() {
		request := testutils.FakeSignEEATransactionRequest()

		signature, err := s.keyManagerClient.SignEEATransaction(s.ctx, s.storeName, s.signAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignEEATransactionRequest()

		signature, err := s.keyManagerClient.SignEEATransaction(s.ctx, "inexistentStoreName", s.signAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestGetEthAccount() {
	s.Run("should sign private transaction successfully", func() {
		retrievedAcc, err := s.keyManagerClient.GetEthAccount(s.ctx, s.storeName, s.signAccount.Address.Hex())
		require.NoError(s.T(), err)

		assert.Equal(s.T(), s.signAccount.Address, retrievedAcc.Address)
		assert.Equal(s.T(), s.signAccount.PublicKey, retrievedAcc.PublicKey)
		assert.Equal(s.T(), s.signAccount.CompressedPublicKey, retrievedAcc.CompressedPublicKey)
		assert.Equal(s.T(), s.signAccount.KeyID, retrievedAcc.KeyID)
		assert.Equal(s.T(), s.signAccount.Tags, retrievedAcc.Tags)
		assert.Equal(s.T(), s.signAccount.Disabled, retrievedAcc.Disabled)
	})

	s.Run("should fail if account does not exist", func() {
		key, err := s.keyManagerClient.GetEthAccount(s.ctx, s.storeName, "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})

	s.Run("should parse errors successfully", func() {
		key, err := s.keyManagerClient.GetEthAccount(s.ctx, "inexistentStoreName", s.signAccount.Address.Hex())
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) TestListEthAccounts() {
	s.Run("should sign private transaction successfully", func() {
		accounts, err := s.keyManagerClient.ListEthAccounts(s.ctx, s.storeName, 999999, 0)
		require.NoError(s.T(), err)

		assert.Contains(s.T(), accounts, strings.ToLower(s.signAccount.Address.Hex()))
	})

	s.Run("should parse errors successfully", func() {
		key, err := s.keyManagerClient.ListEthAccounts(s.ctx, "inexistentStoreName", 0, 0)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *ethTestSuite) queueToDelete(accR *types.EthAccountResponse) {
	s.deleteQueue.Add(1)
	go func() {
		err := s.keyManagerClient.DeleteEthAccount(s.ctx, s.storeName, accR.Address.Hex())
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
			return s.keyManagerClient.DestroyEthAccount(s.ctx, s.storeName, accR.Address.Hex())
		}, s.T().Logf, errMsg, http.StatusConflict, MAX_RETRIES)

		if err != nil {
			s.T().Logf(errMsg)
		}
		s.destroyQueue.Done()
	}()
}
