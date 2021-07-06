// +build e2e

package e2e

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/api/types/testutils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type eth1TestSuite struct {
	suite.Suite
	err              error
	ctx              context.Context
	keyManagerClient *client.HTTPClient
	cfg              *tests.Config
	mainAccount      *types.Eth1AccountResponse
}

func (s *eth1TestSuite) SetupSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}

	var err error

	s.keyManagerClient = client.NewHTTPClient(&http.Client{}, &client.Config{
		URL: s.cfg.KeyManagerURL,
	})

	s.mainAccount, err = s.keyManagerClient.CreateEth1Account(s.ctx, s.cfg.Eth1Store, testutils.FakeCreateEth1AccountRequest())
	require.NoError(s.T(), err)
}

func (s *eth1TestSuite) TearDownSuite() {
	if s.err != nil {
		s.T().Error(s.err)
	}
}

func TestKeyManagerEth1(t *testing.T) {
	s := new(eth1TestSuite)

	s.ctx = context.Background()
	sig := common.NewSignalListener(func(signal os.Signal) {
		s.err = fmt.Errorf("interrupt signal was caught")
		t.FailNow()
	})
	defer sig.Close()

	s.cfg, s.err = tests.NewConfig()
	suite.Run(t, s)
}

func (s *eth1TestSuite) TestCreate() {
	s.Run("should create a new account successfully", func() {
		request := testutils.FakeCreateEth1AccountRequest()
		request.KeyID = "my-account-create"

		acc, err := s.keyManagerClient.CreateEth1Account(s.ctx, s.cfg.Eth1Store, request)
		require.NoError(s.T(), err)

		assert.NotEmpty(s.T(), acc.Address)
		assert.NotEmpty(s.T(), acc.PublicKey)
		assert.NotEmpty(s.T(), acc.CompressedPublicKey)
		assert.Equal(s.T(), request.KeyID, acc.KeyID)
		assert.Equal(s.T(), request.Tags, acc.Tags)
		assert.False(s.T(), acc.Disabled)
		assert.NotEmpty(s.T(), acc.CreatedAt)
		assert.NotEmpty(s.T(), acc.UpdatedAt)
		assert.True(s.T(), acc.ExpireAt.IsZero())
		assert.True(s.T(), acc.DeletedAt.IsZero())
		assert.True(s.T(), acc.DestroyedAt.IsZero())
	})
	
	s.Run("should create a new account with random keyID successfully", func() {
		request := testutils.FakeCreateEth1AccountRequest()
		request.KeyID = ""

		acc, err := s.keyManagerClient.CreateEth1Account(s.ctx, s.cfg.Eth1Store, request)
		require.NoError(s.T(), err)

		assert.NotEmpty(s.T(), acc.Address)
		assert.NotEmpty(s.T(), acc.PublicKey)
		assert.NotEmpty(s.T(), acc.CompressedPublicKey)
		assert.NotEmpty(s.T(), acc.KeyID)
		assert.Equal(s.T(), request.Tags, acc.Tags)
		assert.False(s.T(), acc.Disabled)
		assert.NotEmpty(s.T(), acc.CreatedAt)
		assert.NotEmpty(s.T(), acc.UpdatedAt)
		assert.True(s.T(), acc.ExpireAt.IsZero())
		assert.True(s.T(), acc.DeletedAt.IsZero())
		assert.True(s.T(), acc.DestroyedAt.IsZero())
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeCreateEth1AccountRequest()

		key, err := s.keyManagerClient.CreateEth1Account(s.ctx, "inexistentStoreName", request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *eth1TestSuite) TestUpdate() {
	request := testutils.FakeCreateEth1AccountRequest()
	request.KeyID = "my-account-create"

	acc, err := s.keyManagerClient.CreateEth1Account(s.ctx, s.cfg.Eth1Store, request)
	require.NoError(s.T(), err)

	s.Run("should update an existing account successfully", func() {
		newTags := map[string]string{
			"tagnew": "valuenew",
		}
		acc2, err := s.keyManagerClient.UpdateEth1Account(s.ctx, s.cfg.Eth1Store, acc.Address.String(), &types.UpdateEth1AccountRequest{
			Tags: newTags,
		})

		require.NoError(s.T(), err)
		assert.Equal(s.T(), newTags, acc2.Tags)
	})
}

func (s *eth1TestSuite) TestImport() {
	s.Run("should import an account successfully", func() {
		privKey, _ := crypto.GenerateKey()
		request := testutils.FakeImportEth1AccountRequest()
		request.PrivateKey = privKey.D.Bytes()
		request.KeyID = "my-account-import"

		acc, err := s.keyManagerClient.ImportEth1Account(s.ctx, s.cfg.Eth1Store, request)
		require.NoError(s.T(), err)

		assert.NotEmpty(s.T(), acc.Address)
		assert.NotEmpty(s.T(), acc.PublicKey)
		assert.NotEmpty(s.T(), acc.CompressedPublicKey)
		assert.Equal(s.T(), request.KeyID, acc.KeyID)
		assert.Equal(s.T(), request.Tags, acc.Tags)
		assert.False(s.T(), acc.Disabled)
		assert.NotEmpty(s.T(), acc.CreatedAt)
		assert.NotEmpty(s.T(), acc.UpdatedAt)
		assert.True(s.T(), acc.ExpireAt.IsZero())
		assert.True(s.T(), acc.DeletedAt.IsZero())
		assert.True(s.T(), acc.DestroyedAt.IsZero())
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeImportEth1AccountRequest()

		key, err := s.keyManagerClient.ImportEth1Account(s.ctx, "inexistentStoreName", request)
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *eth1TestSuite) TestSign() {
	s.Run("should sign a payload successfully and verify it", func() {
		request := testutils.FakeSignHexPayloadRequest()

		signature, err := s.keyManagerClient.SignEth1(s.ctx, s.cfg.Eth1Store, s.mainAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)

		err = s.keyManagerClient.VerifyEth1Signature(s.ctx, s.cfg.Eth1Store, &types.VerifyEth1SignatureRequest{
			Data:      request.Data,
			Signature: hexutil.MustDecode(signature),
			Address:   s.mainAccount.Address,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignHexPayloadRequest()

		signature, err := s.keyManagerClient.SignEth1(s.ctx, "inexistentStoreName", s.mainAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *eth1TestSuite) TestSignTypedData() {
	s.Run("should sign typed data successfully and verify it", func() {
		request := testutils.FakeSignTypedDataRequest()

		signature, err := s.keyManagerClient.SignTypedData(s.ctx, s.cfg.Eth1Store, s.mainAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)

		err = s.keyManagerClient.VerifyTypedDataSignature(s.ctx, s.cfg.Eth1Store, &types.VerifyTypedDataRequest{
			TypedData: *request,
			Signature: hexutil.MustDecode(signature),
			Address:   s.mainAccount.Address,
		})
		require.NoError(s.T(), err)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignTypedDataRequest()

		signature, err := s.keyManagerClient.SignTypedData(s.ctx, "inexistentStoreName", s.mainAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *eth1TestSuite) TestSignTransaction() {
	s.Run("should sign transaction successfully", func() {
		request := testutils.FakeSignETHTransactionRequest()

		signature, err := s.keyManagerClient.SignTransaction(s.ctx, s.cfg.Eth1Store, s.mainAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignETHTransactionRequest()

		signature, err := s.keyManagerClient.SignTransaction(s.ctx, "inexistentStoreName", s.mainAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *eth1TestSuite) TestSignPrivateTransaction() {
	s.Run("should sign private transaction successfully", func() {
		request := testutils.FakeSignQuorumPrivateTransactionRequest()

		signature, err := s.keyManagerClient.SignQuorumPrivateTransaction(s.ctx, s.cfg.Eth1Store, s.mainAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignQuorumPrivateTransactionRequest()

		signature, err := s.keyManagerClient.SignQuorumPrivateTransaction(s.ctx, "inexistentStoreName", s.mainAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *eth1TestSuite) TestSignEEATransaction() {
	s.Run("should sign private transaction successfully", func() {
		request := testutils.FakeSignEEATransactionRequest()

		signature, err := s.keyManagerClient.SignEEATransaction(s.ctx, s.cfg.Eth1Store, s.mainAccount.Address.Hex(), request)
		require.NoError(s.T(), err)

		assert.NotNil(s.T(), signature)
	})

	s.Run("should parse errors successfully", func() {
		request := testutils.FakeSignEEATransactionRequest()

		signature, err := s.keyManagerClient.SignEEATransaction(s.ctx, "inexistentStoreName", s.mainAccount.Address.Hex(), request)
		require.Empty(s.T(), signature)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *eth1TestSuite) TestGet() {
	s.Run("should sign private transaction successfully", func() {
		retrievedAcc, err := s.keyManagerClient.GetEth1Account(s.ctx, s.cfg.Eth1Store, s.mainAccount.Address.Hex())
		require.NoError(s.T(), err)

		assert.Equal(s.T(), s.mainAccount.Address, retrievedAcc.Address)
		assert.Equal(s.T(), s.mainAccount.PublicKey, retrievedAcc.PublicKey)
		assert.Equal(s.T(), s.mainAccount.CompressedPublicKey, retrievedAcc.CompressedPublicKey)
		assert.Equal(s.T(), s.mainAccount.KeyID, retrievedAcc.KeyID)
		assert.Equal(s.T(), s.mainAccount.Tags, retrievedAcc.Tags)
		assert.Equal(s.T(), s.mainAccount.Disabled, retrievedAcc.Disabled)
		assert.Equal(s.T(), s.mainAccount.CreatedAt, retrievedAcc.CreatedAt)
		assert.Equal(s.T(), s.mainAccount.UpdatedAt, retrievedAcc.UpdatedAt)
		assert.Equal(s.T(), s.mainAccount.ExpireAt, retrievedAcc.ExpireAt)
		assert.Equal(s.T(), s.mainAccount.DeletedAt, retrievedAcc.DeletedAt)
		assert.Equal(s.T(), s.mainAccount.DestroyedAt, retrievedAcc.DestroyedAt)
	})

	s.Run("should fail if account does not exist", func() {
		key, err := s.keyManagerClient.GetEth1Account(s.ctx, s.cfg.Eth1Store, "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})

	s.Run("should parse errors successfully", func() {
		key, err := s.keyManagerClient.GetEth1Account(s.ctx, "inexistentStoreName", s.mainAccount.Address.Hex())
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}

func (s *eth1TestSuite) TestList() {
	s.Run("should sign private transaction successfully", func() {
		accounts, err := s.keyManagerClient.ListEth1Accounts(s.ctx, s.cfg.Eth1Store)
		require.NoError(s.T(), err)

		assert.Contains(s.T(), accounts, s.mainAccount.Address.Hex())
	})

	s.Run("should parse errors successfully", func() {
		key, err := s.keyManagerClient.ListEth1Accounts(s.ctx, "inexistentStoreName")
		require.Nil(s.T(), key)

		httpError := err.(*client.ResponseError)
		assert.Equal(s.T(), 404, httpError.StatusCode)
	})
}
