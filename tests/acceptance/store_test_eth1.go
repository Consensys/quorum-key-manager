package acceptancetests

import (
	"encoding/hex"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/mocks"
	"github.com/consensys/quorum-key-manager/src/stores/manager/local"
	"github.com/consensys/quorum-key-manager/src/stores/store/database/postgres"
	"github.com/golang/mock/gomock"
	"math/big"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/store/eth1"
	eth1local "github.com/consensys/quorum-key-manager/src/stores/store/eth1/local"
	hashicorpkey "github.com/consensys/quorum-key-manager/src/stores/store/keys/hashicorp"
	quorumtypes "github.com/consensys/quorum/core/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type eth1TestSuite struct {
	suite.Suite
	env   *IntegrationEnvironment
	store eth1.Store
}

func (s *eth1TestSuite) TearDownSuite() {
	ctx := s.env.ctx

	accounts, err := s.store.List(ctx)
	require.NoError(s.T(), err)

	s.env.logger.Info("Deleting the following accounts", "addresses", accounts)
	for _, address := range accounts {
		err = s.store.Delete(ctx, address)
		if err != nil && errors.IsNotSupportedError(err) || err != nil && errors.IsNotImplementedError(err) {
			return
		}
	}

	for _, acc := range accounts {
		maxTries := MaxRetries
		for {
			err := s.store.Destroy(ctx, acc)
			if err != nil && errors.IsNotSupportedError(err) || err != nil && errors.IsNotImplementedError(err) {
				return
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

func (s *eth1TestSuite) TestInit() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	logger := s.env.logger
	ctx := s.env.ctx
	attr := &entities.Attributes{
		Tags: testutils.FakeTags(),
	}
	algo := &entities.Algorithm{
		Type:          entities.Ecdsa,
		EllipticCurve: entities.Secp256k1,
	}

	keyStore := hashicorpkey.New(s.env.hashicorpClient, HashicorpKeyMountPoint, logger)
	db := postgres.New(logger, mocks.NewMockClient(ctrl))

	key1, err := keyStore.Create(ctx, "init-key-1", algo, attr)
	require.NoError(s.T(), err)

	key2, err := keyStore.Create(ctx, "init-key-2", algo, attr)
	require.NoError(s.T(), err)

	_, err = keyStore.Create(ctx, "init-key-eddsa", &entities.Algorithm{
		Type:          entities.Eddsa,
		EllipticCurve: entities.Bn254,
	}, attr)
	require.NoError(s.T(), err)

	err = local.InitDB(ctx, keyStore, db.ETH1Accounts())
	require.NoError(s.T(), err)

	ethStore := eth1local.New(keyStore, db.ETH1Accounts(), logger)

	s.Run("should load ETH1 keys", func() {
		pubKey1, _ := crypto.UnmarshalPubkey(key1.PublicKey)
		account1, err := ethStore.Get(ctx, crypto.PubkeyToAddress(*pubKey1).Hex())
		require.NoError(s.T(), err)
		assert.Equal(s.T(), account1.ID, key1.ID)

		pubKey2, _ := crypto.UnmarshalPubkey(key2.PublicKey)
		account2, err := ethStore.Get(ctx, crypto.PubkeyToAddress(*pubKey2).Hex())
		require.NoError(s.T(), err)
		assert.Equal(s.T(), account2.ID, key2.ID)
	})
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

		assert.Equal(s.T(), account.ID, id)
		assert.NotEmpty(s.T(), account.Address)
		assert.NotEmpty(s.T(), account.PublicKey)
		assert.NotEmpty(s.T(), account.CompressedPublicKey)
		assert.Equal(s.T(), account.Tags, tags)
		assert.False(s.T(), account.Metadata.Disabled)
		assert.True(s.T(), account.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), account.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), account.Metadata.ExpireAt.IsZero())
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
		assert.Equal(s.T(), account.ID, id)
		assert.Equal(s.T(), "0x83a0254be47813BBff771F4562744676C4e793F0", account.Address.Hex())
		assert.Equal(s.T(), "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2", hexutil.Encode(account.PublicKey))
		assert.Equal(s.T(), "0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f", hexutil.Encode(account.CompressedPublicKey))
		assert.Equal(s.T(), account.Tags, tags)
		assert.False(s.T(), account.Metadata.Disabled)
		assert.True(s.T(), account.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), account.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), account.Metadata.ExpireAt.IsZero())
		assert.NotEmpty(s.T(), account.Metadata.CreatedAt)
		assert.NotEmpty(s.T(), account.Metadata.UpdatedAt)
		assert.Equal(s.T(), account.Metadata.UpdatedAt, account.Metadata.CreatedAt)
	})

	s.Run("should fail with AlreadyExistsError if the account already exists (same address)", func() {
		account, err := s.store.Import(ctx, "my-account", privKey, &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), account)
		assert.True(s.T(), errors.IsAlreadyExistsError(err))
	})

	s.Run("should fail with InvalidParameterError if private key is invalid", func() {
		account, err := s.store.Import(ctx, "my-account", []byte("invalidPrivKey"), &entities.Attributes{
			Tags: tags,
		})

		require.Nil(s.T(), account)
		assert.True(s.T(), errors.IsInvalidParameterError(err))
	})
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

		assert.Equal(s.T(), retrievedAccount.ID, id)
		assert.NotEmpty(s.T(), retrievedAccount.Address)
		assert.NotEmpty(s.T(), hexutil.Encode(retrievedAccount.PublicKey))
		assert.NotEmpty(s.T(), hexutil.Encode(retrievedAccount.CompressedPublicKey))
		assert.Equal(s.T(), retrievedAccount.Tags, tags)
		assert.False(s.T(), retrievedAccount.Metadata.Disabled)
		assert.True(s.T(), retrievedAccount.Metadata.DestroyedAt.IsZero())
		assert.True(s.T(), retrievedAccount.Metadata.DeletedAt.IsZero())
		assert.True(s.T(), retrievedAccount.Metadata.ExpireAt.IsZero())
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

func (s *eth1TestSuite) TestSignDataVerify() {
	ctx := s.env.ctx
	id := s.newID("my-account-sign-data")
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

	s.Run("should sign a transaction data successfully", func() {
		signer := types.NewEIP155Signer(chainID)
		txData := signer.Hash(tx).Bytes()
		signature, err := s.store.SignData(ctx, account.Address.Hex(), txData)
		require.NoError(s.T(), err)
		signedTx, err := tx.WithSignature(signer, signature)
		require.NoError(s.T(), err)
		sender, err := types.Sender(signer, signedTx)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), account.Address.Hex(), sender.Hex())
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
