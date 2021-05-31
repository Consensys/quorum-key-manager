package local

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/api/formatters"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	mock2 "github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/database/mock"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/entities/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/store/keys/mock"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	id               = "my-account"
	privKey          = "db337ca3295e4050586793f252e641f3b3a83739018fa4cce01a81ca920e7e1c"
	address          = "0x83a0254be47813BBff771F4562744676C4e793F0"
	pubKey           = "0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2"
	compressedPubKey = "0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f"
)

type eth1StoreTestSuite struct {
	suite.Suite
	mockKeyStore       *mock.MockStore
	mockEth1AccountsDB *mock2.MockETH1Accounts
	eth1Store          *Store
}

func TestStore(t *testing.T) {
	s := new(eth1StoreTestSuite)
	suite.Run(t, s)
}

func (s *eth1StoreTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.mockKeyStore = mock.NewMockStore(ctrl)
	s.mockEth1AccountsDB = mock2.NewMockETH1Accounts(ctrl)
	s.eth1Store = New(s.mockKeyStore, s.mockEth1AccountsDB, log.DefaultLogger())
}

func (s *eth1StoreTestSuite) TestCreate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	key := testutils.FakeKey()

	s.T().Run("should create a new ethereum account successfully", func(t *testing.T) {
		expectedAccount := &entities.ETH1Account{
			ID:                  key.ID,
			Address:             address,
			Metadata:            key.Metadata,
			PublicKey:           hexutil.MustDecode(pubKey),
			CompressedPublicKey: hexutil.MustDecode(compressedPubKey),
			Tags:                key.Tags,
		}
		s.mockKeyStore.EXPECT().Create(ctx, id, eth1KeyAlgo, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, expectedAccount).Return(nil)

		account, err := s.eth1Store.Create(ctx, id, attributes)
		assert.NoError(t, err)
		assert.Equal(t, expectedAccount, account)
	})

	s.T().Run("should fail with same error if Create key fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockKeyStore.EXPECT().Create(ctx, id, eth1KeyAlgo, attributes).Return(nil, expectedErr)

		account, err := s.eth1Store.Create(ctx, id, attributes)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})

	s.T().Run("should fail with same error if Add account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockKeyStore.EXPECT().Create(ctx, id, eth1KeyAlgo, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, gomock.Any()).Return(expectedErr)

		account, err := s.eth1Store.Create(ctx, id, attributes)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})
}

func (s *eth1StoreTestSuite) Testimport() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	key := testutils.FakeKey()
	privKeyB, _ := hex.DecodeString(privKey)

	s.T().Run("should import a new ethereum account successfully", func(t *testing.T) {
		expectedAccount := &entities.ETH1Account{
			ID:                  key.ID,
			Address:             address,
			Metadata:            key.Metadata,
			PublicKey:           hexutil.MustDecode(pubKey),
			CompressedPublicKey: hexutil.MustDecode(compressedPubKey),
			Tags:                key.Tags,
		}
		s.mockKeyStore.EXPECT().Import(ctx, id, privKeyB, eth1KeyAlgo, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, expectedAccount).Return(nil)

		account, err := s.eth1Store.Import(ctx, id, privKeyB, attributes)
		assert.NoError(t, err)
		assert.Equal(t, expectedAccount, account)
	})

	s.T().Run("should fail with same error if Create key fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockKeyStore.EXPECT().Import(ctx, id, privKeyB, eth1KeyAlgo, attributes).Return(nil, expectedErr)

		account, err := s.eth1Store.Import(ctx, id, privKeyB, attributes)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})

	s.T().Run("should fail with same error if Add account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockKeyStore.EXPECT().Import(ctx, id, privKeyB, eth1KeyAlgo, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, gomock.Any()).Return(expectedErr)

		account, err := s.eth1Store.Import(ctx, id, privKeyB, attributes)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})
}

func (s *eth1StoreTestSuite) TestGet() {
	ctx := context.Background()

	s.T().Run("should get an ethereum account successfully", func(t *testing.T) {
		fakeETH1Account := testutils.FakeETH1Account()
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeETH1Account, nil)

		account, err := s.eth1Store.Get(ctx, address)
		assert.NoError(t, err)
		assert.Equal(t, fakeETH1Account, account)
	})

	s.T().Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		account, err := s.eth1Store.Get(ctx, address)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})
}

func (s *eth1StoreTestSuite) TestGetAll() {
	ctx := context.Background()

	s.T().Run("should get all ethereum accounts successfully", func(t *testing.T) {
		expectedAccounts := []*entities.ETH1Account{testutils.FakeETH1Account(), testutils.FakeETH1Account()}
		s.mockEth1AccountsDB.EXPECT().GetAll(ctx).Return(expectedAccounts, nil)

		accounts, err := s.eth1Store.GetAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedAccounts, accounts)
	})

	s.T().Run("should fail with same error if GetAll accounts fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().GetAll(ctx).Return(nil, expectedErr)

		accounts, err := s.eth1Store.GetAll(ctx)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, accounts)
	})
}

func (s *eth1StoreTestSuite) TestList() {
	ctx := context.Background()

	s.T().Run("should list all ethereum accounts successfully", func(t *testing.T) {
		expectedAccounts := []*entities.ETH1Account{testutils.FakeETH1Account(), testutils.FakeETH1Account()}
		s.mockEth1AccountsDB.EXPECT().GetAll(ctx).Return(expectedAccounts, nil)

		addresses, err := s.eth1Store.List(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []string{expectedAccounts[0].Address, expectedAccounts[1].Address}, addresses)
	})

	s.T().Run("should fail with same error if GetAll account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().GetAll(ctx).Return(nil, expectedErr)

		accounts, err := s.eth1Store.List(ctx)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, accounts)
	})
}

func (s *eth1StoreTestSuite) TestUpdate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	key := testutils.FakeKey()
	fakeAccount := testutils.FakeETH1Account()

	s.T().Run("should update an ethereum account successfully", func(t *testing.T) {
		expectedUpdatedAccount := &entities.ETH1Account{
			ID:                  key.ID,
			Address:             address,
			Metadata:            key.Metadata,
			PublicKey:           hexutil.MustDecode(pubKey),
			CompressedPublicKey: hexutil.MustDecode(compressedPubKey),
			Tags:                key.Tags,
		}

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Update(ctx, fakeAccount.ID, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, expectedUpdatedAccount).Return(nil)

		account, err := s.eth1Store.Update(ctx, address, attributes)
		assert.NoError(t, err)
		assert.Equal(t, expectedUpdatedAccount, account)
	})

	s.T().Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		account, err := s.eth1Store.Update(ctx, address, attributes)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})

	s.T().Run("should fail with same error if Update key fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Update(ctx, fakeAccount.ID, attributes).Return(nil, expectedErr)

		account, err := s.eth1Store.Update(ctx, address, attributes)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})

	s.T().Run("should fail with same error if Add account fails", func(t *testing.T) {
		expectedUpdatedAccount := &entities.ETH1Account{
			ID:                  key.ID,
			Address:             address,
			Metadata:            key.Metadata,
			PublicKey:           hexutil.MustDecode(pubKey),
			CompressedPublicKey: hexutil.MustDecode(compressedPubKey),
			Tags:                key.Tags,
		}
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Update(ctx, fakeAccount.ID, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, expectedUpdatedAccount).Return(expectedErr)

		account, err := s.eth1Store.Update(ctx, address, attributes)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})
}

func (s *eth1StoreTestSuite) TestDelete() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()

	s.T().Run("should delete an ethereum account successfully", func(t *testing.T) {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Delete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Remove(ctx, address).Return(nil)
		s.mockEth1AccountsDB.EXPECT().AddDeleted(ctx, fakeAccount).Return(nil)

		err := s.eth1Store.Delete(ctx, address)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		err := s.eth1Store.Delete(ctx, address)
		assert.Equal(t, expectedErr, err)
	})

	s.T().Run("should fail with same error if Delete key fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Delete(ctx, fakeAccount.ID).Return(expectedErr)

		err := s.eth1Store.Delete(ctx, address)
		assert.Equal(t, expectedErr, err)
	})

	s.T().Run("should fail with same error if Remove account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Delete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Remove(ctx, address).Return(expectedErr)

		err := s.eth1Store.Delete(ctx, address)
		assert.Equal(t, expectedErr, err)
	})

	s.T().Run("should fail with same error if AddDeleted account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Delete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Remove(ctx, address).Return(nil)
		s.mockEth1AccountsDB.EXPECT().AddDeleted(ctx, fakeAccount).Return(expectedErr)

		err := s.eth1Store.Delete(ctx, address)
		assert.Equal(t, expectedErr, err)
	})
}

func (s *eth1StoreTestSuite) TestGetDeleted() {
	ctx := context.Background()

	s.T().Run("should get a deleted ethereum account successfully", func(t *testing.T) {
		fakeETH1Account := testutils.FakeETH1Account()
		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeETH1Account, nil)

		account, err := s.eth1Store.GetDeleted(ctx, address)
		assert.NoError(t, err)
		assert.Equal(t, fakeETH1Account, account)
	})

	s.T().Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(nil, expectedErr)

		account, err := s.eth1Store.GetDeleted(ctx, address)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, account)
	})
}

func (s *eth1StoreTestSuite) TestListDeleted() {
	ctx := context.Background()

	s.T().Run("should list all ethereum accounts successfully", func(t *testing.T) {
		expectedAccounts := []*entities.ETH1Account{testutils.FakeETH1Account(), testutils.FakeETH1Account()}
		s.mockEth1AccountsDB.EXPECT().GetAllDeleted(ctx).Return(expectedAccounts, nil)

		addresses, err := s.eth1Store.ListDeleted(ctx)
		assert.NoError(t, err)
		assert.Equal(t, []string{expectedAccounts[0].Address, expectedAccounts[1].Address}, addresses)
	})

	s.T().Run("should fail with same error if GetAll account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().GetAllDeleted(ctx).Return(nil, expectedErr)

		accounts, err := s.eth1Store.ListDeleted(ctx)
		assert.Equal(t, expectedErr, err)
		assert.Empty(t, accounts)
	})
}

func (s *eth1StoreTestSuite) TestUndelete() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()

	s.T().Run("should undelete an ethereum account successfully", func(t *testing.T) {
		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Undelete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, fakeAccount).Return(nil)

		err := s.eth1Store.Undelete(ctx, address)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with same error if GetDeleted account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(nil, expectedErr)

		err := s.eth1Store.Undelete(ctx, address)
		assert.Equal(t, expectedErr, err)
	})

	s.T().Run("should fail with same error if Undelete key fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Undelete(ctx, fakeAccount.ID).Return(expectedErr)

		err := s.eth1Store.Undelete(ctx, address)
		assert.Equal(t, expectedErr, err)
	})

	s.T().Run("should fail with same error if RemoveDeleted account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Undelete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(expectedErr)

		err := s.eth1Store.Undelete(ctx, address)
		assert.Equal(t, expectedErr, err)
	})

	s.T().Run("should fail with same error if Add account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Undelete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, fakeAccount).Return(expectedErr)

		err := s.eth1Store.Undelete(ctx, address)
		assert.Equal(t, expectedErr, err)
	})
}

func (s *eth1StoreTestSuite) TestDestroy() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()

	s.T().Run("should undelete an ethereum account successfully", func(t *testing.T) {
		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Destroy(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(nil)

		err := s.eth1Store.Destroy(ctx, address)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with same error if GetDeleted account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(nil, expectedErr)

		err := s.eth1Store.Destroy(ctx, address)
		assert.Equal(t, expectedErr, err)
	})

	s.T().Run("should fail with same error if Destroy key fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Destroy(ctx, fakeAccount.ID).Return(expectedErr)

		err := s.eth1Store.Destroy(ctx, address)
		assert.Equal(t, expectedErr, err)
	})

	s.T().Run("should fail with same error if RemoveDeleted account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Destroy(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(expectedErr)

		err := s.eth1Store.Destroy(ctx, address)
		assert.Equal(t, expectedErr, err)
	})
}

func (s *eth1StoreTestSuite) TestSignVerify() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()
	data := []byte("my data to sign")
	ecdsaSignature := hexutil.MustDecode("0x63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed")

	s.T().Run("should sign a payload successfully with appended V value and verify it", func(t *testing.T) {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, crypto.Keccak256(data)).Return(ecdsaSignature, nil)

		signature, err := s.eth1Store.Sign(ctx, address, data)
		assert.NoError(t, err)
		// Note this the returned signature is not the same as the ecdsaSignature! It has an appended byte
		assert.Equal(t, "0x63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed01", hexutil.Encode(signature))

		err = s.eth1Store.Verify(ctx, address, data, signature)
		assert.NoError(t, err)
	})

	s.T().Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		signature, err := s.eth1Store.Sign(ctx, address, data)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signature)
	})

	s.T().Run("should fail with same error if Sign fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, crypto.Keccak256(data)).Return(nil, expectedErr)

		signature, err := s.eth1Store.Sign(ctx, address, data)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signature)
	})
}

func (s *eth1StoreTestSuite) TestSignTransaction() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()
	chainID := big.NewInt(1)
	tx := types.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(0),
		0,
		big.NewInt(0),
		nil,
	)
	ecdsaSignature := hexutil.MustDecode("0xe276fd7524ed7af67b7f914de5be16fad6b9038009d2d78f2315351fbd48deee57a897964e80e041c674942ef4dbd860cb79a6906fb965d5e4645f5c44f7eae4")

	s.T().Run("should sign a payload successfully with appended V value", func(t *testing.T) {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, types.NewEIP155Signer(chainID).Hash(tx).Bytes()).Return(ecdsaSignature, nil)

		signedRaw, err := s.eth1Store.SignTransaction(ctx, address, chainID, tx)
		assert.NoError(t, err)
		assert.Equal(t, "0xf85d80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808025a0e276fd7524ed7af67b7f914de5be16fad6b9038009d2d78f2315351fbd48deeea057a897964e80e041c674942ef4dbd860cb79a6906fb965d5e4645f5c44f7eae4", hexutil.Encode(signedRaw))
	})

	s.T().Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignTransaction(ctx, address, chainID, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	s.T().Run("should fail with same error if Sign fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, gomock.Any()).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignTransaction(ctx, address, chainID, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})
}

func (s *eth1StoreTestSuite) TestSignPrivate() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()
	tx := quorumtypes.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(0),
		0,
		big.NewInt(0),
		nil,
	)
	ecdsaSignature := hexutil.MustDecode("0x80365b013992519479ddd83584039d66851da560dbbe67f59ab9bdcd97b6250355e93d2c8050fb413956298c10eb7b8b2c8d76f4be261e458e4987cc5fed9f01")

	s.T().Run("should sign a payload successfully with appended V value", func(t *testing.T) {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, quorumtypes.QuorumPrivateTxSigner{}.Hash(tx).Bytes()).Return(ecdsaSignature, nil)

		signedRaw, err := s.eth1Store.SignPrivate(ctx, address, tx)
		assert.NoError(t, err)
		assert.Equal(t, "0xf85d80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808026a080365b013992519479ddd83584039d66851da560dbbe67f59ab9bdcd97b62503a055e93d2c8050fb413956298c10eb7b8b2c8d76f4be261e458e4987cc5fed9f01", hexutil.Encode(signedRaw))
	})

	s.T().Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignPrivate(ctx, address, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	s.T().Run("should fail with same error if Sign fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, gomock.Any()).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignPrivate(ctx, address, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})
}

func (s *eth1StoreTestSuite) TestSignEEA() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()
	chainID := big.NewInt(1)
	tx := types.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
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
	ecdsaSignature := hexutil.MustDecode("0x6854034c21ebb5a6d4aa9a9c1462862b1e4af355383413a0dcfbba309f56ed0220c0ebc19f159ce83c24dde6f1b2d424025e45bc8b00be3e2fd4367949d4f0b3")

	s.T().Run("should sign a payload with privacyFor successfully with appended V value", func(t *testing.T) {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, hexutil.MustDecode("0x5749cc0adae7a54f9c5148a9e21719a2b472dec7b7ae7c1d68bf35e2e161f94d")).Return(ecdsaSignature, nil)

		signedRaw, err := s.eth1Store.SignEEA(ctx, address, chainID, tx, privateArgs)
		assert.NoError(t, err)
		assert.Equal(t, "0xf8cd80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808026a06854034c21ebb5a6d4aa9a9c1462862b1e4af355383413a0dcfbba309f56ed02a020c0ebc19f159ce83c24dde6f1b2d424025e45bc8b00be3e2fd4367949d4f0b3a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564", hexutil.Encode(signedRaw))
	})

	s.T().Run("should fail with same error if Get account fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignEEA(ctx, address, chainID, tx, privateArgs)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	s.T().Run("should fail with same error if Sign fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, gomock.Any()).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignEEA(ctx, address, chainID, tx, privateArgs)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})
}
