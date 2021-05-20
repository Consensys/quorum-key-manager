package acceptancetests

import (
	"encoding/hex"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/store/eth1"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"math/big"
	"testing"
)

// TODO: Destroy secrets when done with the tests to avoid conflicts between tests

type eth1TestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store eth1.Store
}

func (s *eth1TestSuite) TestCreate() {
	ctx := s.env.ctx
	id := "my-account-create"
	tags := testutils.FakeTags()

	s.T().Run("should create a new ethereum account successfully", func(t *testing.T) {
		account, err := s.store.Create(ctx, id, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)
		assert.Equal(t, account.ID, id)
		assert.NotEmpty(t, account.Address)
		assert.NotEmpty(t, account.PublicKey)
		assert.NotEmpty(t, account.CompressedPublicKey)
		assert.Equal(t, account.Tags, tags)
		assert.Equal(t, "1", account.Metadata.Version)
		assert.False(t, account.Metadata.Disabled)
		assert.True(t, account.Metadata.DestroyedAt.IsZero())
		assert.True(t, account.Metadata.DeletedAt.IsZero())
		assert.True(t, account.Metadata.ExpireAt.IsZero())
		assert.NotEmpty(t, account.Metadata.CreatedAt)
		assert.NotEmpty(t, account.Metadata.UpdatedAt)
		assert.Equal(t, account.Metadata.UpdatedAt, account.Metadata.CreatedAt)
	})
}

func (s *eth1TestSuite) TestImport() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()

	s.T().Run("should create a new ethereum account successfully", func(t *testing.T) {
		id := fmt.Sprintf("my-account-import-%d", common.RandInt(1000))
		privKey, _ := hex.DecodeString(privKeyECDSA)

		account, err := s.store.Import(ctx, id, privKey, &entities.Attributes{
			Tags: tags,
		})

		require.NoError(t, err)
		assert.Equal(t, account.ID, id)
		assert.Equal(t, "0x83a0254be47813BBff771F4562744676C4e793F0", account.Address)
		assert.Equal(t, "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", hexutil.Encode(account.PublicKey))
		assert.Equal(t, "0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f", hexutil.Encode(account.CompressedPublicKey))
		assert.Equal(t, account.Tags, tags)
		assert.Equal(t, "1", account.Metadata.Version)
		assert.False(t, account.Metadata.Disabled)
		assert.True(t, account.Metadata.DestroyedAt.IsZero())
		assert.True(t, account.Metadata.DeletedAt.IsZero())
		assert.True(t, account.Metadata.ExpireAt.IsZero())
		assert.NotEmpty(t, account.Metadata.CreatedAt)
		assert.NotEmpty(t, account.Metadata.UpdatedAt)
		assert.Equal(t, account.Metadata.UpdatedAt, account.Metadata.CreatedAt)
	})

	s.T().Run("should fail with InvalidParameterError if private key is invalid", func(t *testing.T) {
		id := "my-account"

		account, err := s.store.Import(ctx, id, []byte("invalidPrivKey"), &entities.Attributes{
			Tags: tags,
		})

		require.Nil(t, account)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}

func (s *eth1TestSuite) TestGet() {
	ctx := s.env.ctx
	id := fmt.Sprintf("my-account-get-%d", common.RandInt(1000))
	tags := testutils.FakeTags()
	privKey, _ := hex.DecodeString(privKeyECDSA)

	account, err := s.store.Import(ctx, id, privKey, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.T().Run("should get an ethereum account successfully", func(t *testing.T) {
		retrievedAccount, err := s.store.Get(ctx, account.Address)
		require.NoError(t, err)

		assert.Equal(t, retrievedAccount.ID, id)
		assert.Equal(t, "0x83a0254be47813BBff771F4562744676C4e793F0", retrievedAccount.Address)
		assert.Equal(t, "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", hexutil.Encode(retrievedAccount.PublicKey))
		assert.Equal(t, "0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f", hexutil.Encode(retrievedAccount.CompressedPublicKey))
		assert.Equal(t, retrievedAccount.Tags, tags)
		assert.Equal(t, "1", retrievedAccount.Metadata.Version)
		assert.False(t, retrievedAccount.Metadata.Disabled)
		assert.True(t, retrievedAccount.Metadata.DestroyedAt.IsZero())
		assert.True(t, retrievedAccount.Metadata.DeletedAt.IsZero())
		assert.True(t, retrievedAccount.Metadata.ExpireAt.IsZero())
		assert.NotEmpty(t, retrievedAccount.Metadata.CreatedAt)
		assert.NotEmpty(t, retrievedAccount.Metadata.UpdatedAt)
		assert.Equal(t, retrievedAccount.Metadata.UpdatedAt, retrievedAccount.Metadata.CreatedAt)
	})

	s.T().Run("should fail with NotFoundError if account is not found", func(t *testing.T) {
		retrievedAccount, err := s.store.Get(ctx, "invalidAccount")
		require.Nil(t, retrievedAccount)
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *eth1TestSuite) TestList() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	id := fmt.Sprintf("my-account-list-%s", common.RandString(5))
	id2 := fmt.Sprintf("my-account-list-%s", common.RandString(5))

	account1, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	account2, err := s.store.Create(ctx, id2, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.T().Run("should get all account addresses", func(t *testing.T) {
		addresses, err := s.store.List(ctx)
		require.NoError(t, err)

		assert.Contains(t, addresses, account1.Address)
		assert.Contains(t, addresses, account2.Address)
	})
}

func (s *eth1TestSuite) TestSign() {
	ctx := s.env.ctx
	payload := []byte("my data to sign")
	id := fmt.Sprintf("my-account-sign-%d", common.RandInt(1000))
	privKey, _ := hex.DecodeString(privKeyECDSA)

	account, err := s.store.Import(ctx, id, privKey, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), err)

	s.T().Run("should sign a payload successfully", func(t *testing.T) {
		signature, err := s.store.Sign(ctx, account.Address, payload)
		require.NoError(t, err)
		assert.Equal(t, "0x63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed01", hexutil.Encode(signature))

		verified, err := verifySignature(signature, payload, privKey)
		require.NoError(t, err)
		assert.True(t, verified)
	})

	s.T().Run("should fail with NotFoundError if account is not found", func(t *testing.T) {
		signature, err := s.store.Sign(ctx, "invalidAccount", payload)
		require.Empty(t, signature)
		assert.True(t, errors.IsNotFoundError(err))
	})
}

func (s *eth1TestSuite) TestSignTransaction() {
	ctx := s.env.ctx
	id := fmt.Sprintf("my-account-sign-tx-%d", common.RandInt(1000))
	chainID := big.NewInt(1)
	to := ethcommon.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18")
	txData := &ethereum.TxData{
		Nonce:    0,
		To:       &to,
		Value:    big.NewInt(0),
		GasPrice: big.NewInt(0),
		GasLimit: 0,
		Data:     nil,
	}
	privKey, _ := hex.DecodeString(privKeyECDSA)

	account, err := s.store.Import(ctx, id, privKey, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), err)

	s.T().Run("should sign a transaction successfully", func(t *testing.T) {
		signature, err := s.store.SignTransaction(ctx, account.Address, chainID, txData)
		require.NoError(t, err)
		assert.Equal(t, "YzQeLIN0Sd43Nbb0QCsVSqChGNAuRaKzEfujnERAJd0523aZyz2KXK93KKh-d4ws3MxAhc8qNG43wYI97Fzi7Q==", signature)
	})

	s.T().Run("should fail with NotFoundError if account is not found", func(t *testing.T) {
		signature, err := s.store.SignTransaction(ctx, "invalidAccount", chainID, txData)
		require.Empty(t, signature)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should fail with InvalidParameterError if tx data is invalid", func(t *testing.T) {
		signature, err := s.store.SignTransaction(ctx, account.Address, chainID, txData)
		require.Empty(t, signature)
		assert.True(t, errors.IsInvalidParameterError(err))
	})
}
