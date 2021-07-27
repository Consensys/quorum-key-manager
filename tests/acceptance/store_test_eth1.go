package acceptancetests

import (
	"encoding/hex"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/stores/store/database"
	"math/big"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/store/eth1"
	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type eth1TestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store eth1.Store
	db    database.Database // TODO: Remove when Delete and Destroy functions are implemented in all stores
}

func (s *eth1TestSuite) TearDownSuite() {
	ctx := s.env.ctx

	accounts, err := s.store.List(ctx)
	require.NoError(s.T(), err)

	s.env.logger.Info("Deleting the following accounts", "addresses", accounts)
	for _, address := range accounts {
		err = s.store.Delete(ctx, address)
		if err != nil && errors.IsNotSupportedError(err) || err != nil && errors.IsNotImplementedError(err) {
			err := s.db.ETH1Accounts().Delete(ctx, address)
			require.NoError(s.T(), err)
		}
	}

	for _, acc := range accounts {
		maxTries := MaxRetries
		for {
			err := s.store.Destroy(ctx, acc)
			if err != nil && errors.IsNotSupportedError(err) || err != nil && errors.IsNotImplementedError(err) {
				err := s.db.ETH1Accounts().Purge(ctx, acc)
				require.NoError(s.T(), err)
			}
			if err != nil && !errors.IsStatusConflictError(err) {
				break
			}
			if maxTries <= 0 {
				if err != nil {
					s.env.logger.Info("failed to destroy account", "account", acc)
				}
				break
			}

			maxTries -= 1
			waitTime := time.Second * time.Duration(MaxRetries-maxTries)
			s.env.logger.Debug("waiting for deletion to complete", "account", acc, "waitFor", waitTime.Seconds())
			time.Sleep(waitTime)
		}
	}
}

func (s *eth1TestSuite) TestCreate() {
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
		assert.Equal(s.T(), account.KeyID, id)
		assert.Equal(s.T(), account.Tags, tags)
		assert.False(s.T(), account.Metadata.Disabled)
		assert.True(s.T(), account.Metadata.DeletedAt.IsZero())
		assert.NotEmpty(s.T(), account.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), account.Metadata.UpdatedAt)
		assert.Equal(s.T(), account.Metadata.UpdatedAt, account.Metadata.CreatedAt)
	})
}

func (s *eth1TestSuite) TestImport() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	privKey, _ := hex.DecodeString(privKeyECDSA)
	id := s.newID("my-account-import")

	account, err := s.store.Import(ctx, id, privKey, &entities.Attributes{
		Tags: tags,
	})
	if err != nil && errors.IsNotSupportedError(err) {
		return
	}
	require.NoError(s.T(), err)

	s.Run("should create a new Ethereum Account successfully", func() {
		assert.Equal(s.T(), account.KeyID, id)
		assert.Equal(s.T(), "0x83a0254be47813BBff771F4562744676C4e793F0", account.Address.Hex())
		assert.Equal(s.T(), "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", hexutil.Encode(account.PublicKey))
		assert.Equal(s.T(), account.Tags, tags)
		assert.False(s.T(), account.Metadata.Disabled)
		assert.True(s.T(), account.Metadata.DeletedAt.IsZero())
		assert.NotEmpty(s.T(), account.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), account.Metadata.UpdatedAt)
		assert.Equal(s.T(), account.Metadata.UpdatedAt, account.Metadata.CreatedAt)
	})

	s.Run("should fail with StatusConflict if we violate a constraint (same address already exists)", func() {
		account, err := s.store.Import(ctx, "my-account", privKey, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), account)
		assert.True(s.T(), errors.IsStatusConflictError(err))
	})

	s.Run("should fail with InvalidParameterError if private key is invalid", func() {
		account, err := s.store.Import(ctx, "my-account", []byte("invalidPrivKey"), &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), account)
		assert.True(s.T(), errors.IsInvalidParameterError(err))
	})

	err = s.store.Delete(ctx, account.Address.Hex())
	if err != nil && errors.IsNotImplementedError(err) || err != nil && errors.IsNotSupportedError(err) {
		return
	}
	require.NoError(s.T(), err)

	err = s.store.Destroy(ctx, account.Address.Hex())
	if err != nil && errors.IsNotImplementedError(err) || err != nil && errors.IsNotSupportedError(err) {
		return
	}
	require.NoError(s.T(), err)
}

func (s *eth1TestSuite) TestGet() {
	ctx := s.env.ctx
	id := s.newID("my-account-get")
	tags := testutils.FakeTags()

	account, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.Run("should get an Ethereum Account successfully", func() {
		retrievedAccount, err := s.store.Get(ctx, account.Address.Hex())
		require.NoError(s.T(), err)

		assert.Equal(s.T(), retrievedAccount.KeyID, id)
		assert.NotEmpty(s.T(), retrievedAccount.Address)
		assert.NotEmpty(s.T(), hexutil.Encode(retrievedAccount.PublicKey))
		assert.Equal(s.T(), retrievedAccount.Tags, tags)
		assert.False(s.T(), retrievedAccount.Metadata.Disabled)
		assert.True(s.T(), retrievedAccount.Metadata.DeletedAt.IsZero())
		assert.NotEmpty(s.T(), retrievedAccount.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), retrievedAccount.Metadata.UpdatedAt)
		assert.Equal(s.T(), retrievedAccount.Metadata.UpdatedAt, retrievedAccount.Metadata.CreatedAt)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		retrievedAccount, err := s.store.Get(ctx, "invalidAccount")
		require.Nil(s.T(), retrievedAccount)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *eth1TestSuite) TestList() {
	ctx := s.env.ctx
	tags := testutils.FakeTags()
	id := s.newID("my-account-list")
	id2 := s.newID("my-account-list")

	account1, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	account2, err := s.store.Create(ctx, id2, &entities.Attributes{
		Tags: tags,
	})
	require.NoError(s.T(), err)

	s.Run("should get all account addresses", func() {
		addresses, err := s.store.List(ctx)
		require.NoError(s.T(), err)

		assert.Contains(s.T(), addresses, account1.Address.Hex())
		assert.Contains(s.T(), addresses, account2.Address.Hex())
	})
}

func (s *eth1TestSuite) TestSignVerify() {
	ctx := s.env.ctx
	payload := []byte("my data to sign")
	id := s.newID("my-account-sign")

	account, err := s.store.Create(ctx, id, &entities.Attributes{
		Tags: testutils.FakeTags(),
	})
	require.NoError(s.T(), err)

	s.Run("should sign, recover an address and verify the signature successfully", func() {
		signature, err := s.store.Sign(ctx, account.Address.Hex(), payload)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signature)

		address, err := s.store.ECRevocer(ctx, payload, signature)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), account.Address.Hex(), address)

		err = s.store.Verify(ctx, address, payload, signature)
		require.NoError(s.T(), err)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		signature, err := s.store.Sign(ctx, "invalidAccount", payload)
		require.Empty(s.T(), signature)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *eth1TestSuite) TestSignTransaction() {
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
		signedRaw, err := s.store.SignTransaction(ctx, account.Address.Hex(), chainID, tx)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signedRaw)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		signedRaw, err := s.store.SignTransaction(ctx, "invalidAccount", chainID, tx)
		require.Empty(s.T(), signedRaw)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *eth1TestSuite) TestSignPrivate() {
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
		signedRaw, err := s.store.SignPrivate(ctx, account.Address.Hex(), tx)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signedRaw)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		signedRaw, err := s.store.SignPrivate(ctx, "invalidAccount", tx)
		require.Empty(s.T(), signedRaw)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *eth1TestSuite) TestSignEEA() {
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
		signedRaw, err := s.store.SignEEA(ctx, account.Address.Hex(), chainID, tx, privateArgs)
		require.NoError(s.T(), err)
		assert.NotEmpty(s.T(), signedRaw)
	})

	s.Run("should fail with NotFoundError if account is not found", func() {
		signedRaw, err := s.store.SignEEA(ctx, "invalidAccount", chainID, tx, privateArgs)
		require.Empty(s.T(), signedRaw)
		assert.True(s.T(), errors.IsNotFoundError(err))
	})
}

func (s *eth1TestSuite) newID(name string) string {
	return fmt.Sprintf("%s-%d", name, common.RandInt(10000))
}
