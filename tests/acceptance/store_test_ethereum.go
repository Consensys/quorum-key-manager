package acceptancetests

import (
	"fmt"
	"math/big"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/database"
	"github.com/consensys/quorum-key-manager/src/stores/entities"
	"github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ethTestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store stores.EthStore
	db    database.ETHAccounts
}

func (s *ethTestSuite) TestCreate() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()

	s.Run("should create a new Ethereum Account successfully", func() {
		id := s.newID("my-account-create")
		account, err := s.store.Create(ctx, id, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		assert.NotEmpty(s.T(), account.Address)
		assert.NotEmpty(s.T(), account.PublicKey)
		assert.NotEmpty(s.T(), account.CompressedPublicKey)
		assert.Equal(s.T(), account.KeyID, id)
		assert.Equal(s.T(), account.Tags, tags)
		assert.False(s.T(), account.Metadata.Disabled)
		assert.NotEmpty(s.T(), account.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), account.Metadata.UpdatedAt)
	})

	s.Run("should create a new Ethereum Account successfully", func() {
		id := s.newID("my-account-create")
		account, err := s.store.Create(ctx, id, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		err = s.db.Delete(ctx, account.Address.Hex())
		require.NoError(s.T(), err)
		err = s.db.Purge(ctx, account.Address.Hex())
		require.NoError(s.T(), err)

		account, err = s.store.Create(ctx, id, &entities.Attributes{
			Tags: tags,
		})
		require.NoError(s.T(), err)

		assert.NotEmpty(s.T(), account.Address)
		assert.NotEmpty(s.T(), account.PublicKey)
		assert.NotEmpty(s.T(), account.CompressedPublicKey)
		assert.Equal(s.T(), account.KeyID, id)
		assert.Equal(s.T(), account.Tags, tags)
		assert.False(s.T(), account.Metadata.Disabled)
		assert.NotEmpty(s.T(), account.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), account.Metadata.UpdatedAt)
	})
}

func (s *ethTestSuite) TestImport() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	privKey, err := crypto.GenerateKey()
	require.NoError(s.T(), err)

	s.Run("should create a new Ethereum Account successfully", func() {
		id := s.newID("my-account-import")
		account, err := s.store.Import(ctx, id, privKey.D.Bytes(), &entities.Attributes{
			Tags: tags,
		})
		if err != nil && errors.IsNotSupportedError(err) {
			return
		}
		require.NoError(s.T(), err)

		assert.Equal(s.T(), account.KeyID, id)
		assert.Equal(s.T(), crypto.PubkeyToAddress(privKey.PublicKey), account.Address)
		assert.Equal(s.T(), crypto.FromECDSAPub(&privKey.PublicKey), account.PublicKey)
		assert.Equal(s.T(), crypto.CompressPubkey(&privKey.PublicKey), account.CompressedPublicKey)
		assert.Equal(s.T(), account.Tags, tags)
		assert.False(s.T(), account.Metadata.Disabled)
		assert.NotEmpty(s.T(), account.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), account.Metadata.UpdatedAt)
	})

	s.Run("should fail with StatusConflict if we violate a constraint (same address already exists)", func() {
		id := s.newID("my-account-import-duplicate")
		account, err := s.store.Import(ctx, id, privKey.D.Bytes(), &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), account)
		assert.True(s.T(), errors.IsStatusConflictError(err) || errors.IsNotSupportedError(err))
	})

	s.Run("should fail with InvalidParameterError if private key is invalid", func() {
		id := s.newID("my-account-import-failure")
		account, err := s.store.Import(ctx, id, []byte("invalidPrivKey"), &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), account)
		assert.True(s.T(), errors.IsInvalidParameterError(err) || errors.IsInvalidFormatError(err)) // Hashicorp will return 400 and not 422
	})
}

func (s *ethTestSuite) TestGet() {
	ctx := s.env.ctx
	id := s.newID("my-account-get")
	tags := testutils.FakeTags()

	account, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.Run("should get an Ethereum Account successfully", func() {
		retrievedAccount, err := s.store.Get(ctx, account.Address)
		require.NoError(s.T(), err)

		assert.Equal(s.T(), retrievedAccount.KeyID, id)
		assert.NotEmpty(s.T(), retrievedAccount.Address)
		assert.NotEmpty(s.T(), hexutil.Encode(retrievedAccount.PublicKey))
		assert.Equal(s.T(), retrievedAccount.Tags, tags)
		assert.False(s.T(), retrievedAccount.Metadata.Disabled)
		assert.True(s.T(), retrievedAccount.Metadata.DeletedAt.IsZero())
		assert.NotEmpty(s.T(), retrievedAccount.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), retrievedAccount.Metadata.UpdatedAt)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		retrievedAccount, err := s.store.Get(ctx, ethcommon.HexToAddress("invalidAddress"))
		require.Nil(s.T(), retrievedAccount)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *ethTestSuite) TestList() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	id := s.newID("my-account-list")
	id2 := s.newID("my-account-list-2")
	id3 := s.newID("my-account-list-3")

	account1, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	account2, err := s.store.Create(ctx, id2, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	account3, err := s.store.Create(ctx, id3, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.Run("should get all account addresses", func() {
		addresses, err := s.store.List(ctx, 0, 0)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), addresses, []string{account1.Address.String(), account2.Address.String(), account3.Address.String()})
	})
	
	s.Run("should get all first account addresses", func() {
		addresses, err := s.store.List(ctx, 1, 0)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), addresses, []string{account1.Address.String()})
	})
	
	s.Run("should get last two account addresses", func() {
		addresses, err := s.store.List(ctx, 2, 1)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), addresses, []string{account2.Address.String(), account3.Address.String()})
	})
}

func (s *ethTestSuite) TestSignMessageVerify() {
	ctx := s.env.ctx
	payload := hexutil.MustDecode("0xfeaa")
	id := s.newID("my-account-sign")

	account, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), err)

	s.Run("should sign, recover an address and verify the signature successfully", func() {
		signature, err := s.store.SignMessage(ctx, account.Address, payload)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signature)

		err = s.store.VerifyMessage(ctx, account.Address, payload, signature)
		require.NoError(s.T(), err)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		signature, err := s.store.SignMessage(ctx, ethcommon.HexToAddress("invalidAddress"), payload)
		require.Empty(s.T(), signature)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *ethTestSuite) TestSignTransaction() {
	ctx := s.env.ctx
	id := s.newID("my-account-sign-tx")
	chainID := big.NewInt(1)
	tx := types.NewTransaction(
		0,
		ethcommon.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(0),
		0,
		big.NewInt(0),
		nil,
	)

	account, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), err)

	s.Run("should sign a transaction successfully", func() {
		signedRaw, err := s.store.SignTransaction(ctx, account.Address, chainID, tx)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signedRaw)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		signedRaw, err := s.store.SignTransaction(ctx, ethcommon.HexToAddress("invalidAddress"), chainID, tx)
		require.Empty(s.T(), signedRaw)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *ethTestSuite) TestSignPrivate() {
	ctx := s.env.ctx
	id := s.newID("my-account-sign-private")
	tx := quorumtypes.NewTransaction(
		0,
		ethcommon.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(0),
		0,
		big.NewInt(0),
		nil,
	)

	account, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), err)

	s.Run("should sign a transaction successfully", func() {
		signedRaw, err := s.store.SignPrivate(ctx, account.Address, tx)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signedRaw)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		signedRaw, err := s.store.SignPrivate(ctx, ethcommon.HexToAddress("invalidAddress"), tx)
		require.Empty(s.T(), signedRaw)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *ethTestSuite) TestSignEEA() {
	ctx := s.env.ctx
	id := s.newID("my-account-sign-eea")
	chainID := big.NewInt(1)
	tx := types.NewTransaction(
		0,
		ethcommon.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(0),
		0,
		big.NewInt(0),
		nil,
	)
	privateFrom := "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="
	privateFor := []string{"A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo=", "B1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="}
	privateType := formatters.PrivateTxTypeRestricted
	privateArgs := &ethereum.PrivateArgs{
		PrivateFrom: &privateFrom,
		PrivateFor:  &privateFor,
		PrivateType: &privateType,
	}

	account, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), err)

	s.Run("should sign a transaction successfully", func() {
		signedRaw, err := s.store.SignEEA(ctx, account.Address, chainID, tx, privateArgs)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signedRaw)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		signedRaw, err := s.store.SignEEA(ctx, ethcommon.HexToAddress("invalidAddress"), chainID, tx, privateArgs)
		require.Empty(s.T(), signedRaw)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *ethTestSuite) newID(name string) string {
	return fmt.Sprintf("%s-%s", name, common.RandHexString(16))
}
