package local

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/src/stores/api/formatters"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/store/database/mock"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	"github.com/consensys/quorum-key-manager/src/stores/store/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/store/keys/mock"
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
	s.eth1Store = New(s.mockKeyStore, s.mockEth1AccountsDB, testutils2.NewMockLogger(ctrl))
}

func (s *eth1StoreTestSuite) TestCreate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	key := testutils.FakeKey()

	s.Run("should create a new Ethereum Account successfully", func() {
		expectedAccount := &entities.ETH1Account{
			ID:                  key.ID,
			Address:             common.HexToAddress(address),
			Metadata:            key.Metadata,
			PublicKey:           hexutil.MustDecode(pubKey),
			CompressedPublicKey: hexutil.MustDecode(compressedPubKey),
			Tags:                key.Tags,
		}
		s.mockKeyStore.EXPECT().Create(ctx, id, eth1KeyAlgo, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, expectedAccount).Return(nil)

		account, err := s.eth1Store.Create(ctx, id, attributes)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedAccount, account)
	})

	s.Run("should fail with same error if Create key fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockKeyStore.EXPECT().Create(ctx, id, eth1KeyAlgo, attributes).Return(nil, expectedErr)

		account, err := s.eth1Store.Create(ctx, id, attributes)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})

	s.Run("should fail with same error if Add account fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockKeyStore.EXPECT().Create(ctx, id, eth1KeyAlgo, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, gomock.Any()).Return(expectedErr)

		account, err := s.eth1Store.Create(ctx, id, attributes)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})
}

func (s *eth1StoreTestSuite) TestImport() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	key := testutils.FakeKey()
	privKeyB, _ := hex.DecodeString(privKey)

	s.Run("should import a new Ethereum Account successfully", func() {
		expectedAccount := &entities.ETH1Account{
			ID:                  key.ID,
			Address:             common.HexToAddress(address),
			Metadata:            key.Metadata,
			PublicKey:           hexutil.MustDecode(pubKey),
			CompressedPublicKey: hexutil.MustDecode(compressedPubKey),
			Tags:                key.Tags,
		}
		s.mockKeyStore.EXPECT().Import(ctx, id, privKeyB, eth1KeyAlgo, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, expectedAccount).Return(nil)

		account, err := s.eth1Store.Import(ctx, id, privKeyB, attributes)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedAccount, account)
	})

	s.Run("should fail with same error if Create key fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockKeyStore.EXPECT().Import(ctx, id, privKeyB, eth1KeyAlgo, attributes).Return(nil, expectedErr)

		account, err := s.eth1Store.Import(ctx, id, privKeyB, attributes)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})

	s.Run("should fail with same error if Add account fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockKeyStore.EXPECT().Import(ctx, id, privKeyB, eth1KeyAlgo, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, gomock.Any()).Return(expectedErr)

		account, err := s.eth1Store.Import(ctx, id, privKeyB, attributes)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})
}

func (s *eth1StoreTestSuite) TestGet() {
	ctx := context.Background()

	s.Run("should get an Ethereum Account successfully", func() {
		fakeETH1Account := testutils.FakeETH1Account()
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeETH1Account, nil)

		account, err := s.eth1Store.Get(ctx, address)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), fakeETH1Account, account)
	})

	s.Run("should fail with same error if Get account fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		account, err := s.eth1Store.Get(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})
}

func (s *eth1StoreTestSuite) TestGetAll() {
	ctx := context.Background()

	s.Run("should get all Ethereum Accounts successfully", func() {
		expectedAccounts := []*entities.ETH1Account{testutils.FakeETH1Account(), testutils.FakeETH1Account()}
		s.mockEth1AccountsDB.EXPECT().GetAll(ctx).Return(expectedAccounts, nil)

		accounts, err := s.eth1Store.GetAll(ctx)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedAccounts, accounts)
	})

	s.Run("should fail with same error if GetAll accounts fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().GetAll(ctx).Return(nil, expectedErr)

		accounts, err := s.eth1Store.GetAll(ctx)
		assert.Equal(s.T(), expectedErr, err)
		assert.Empty(s.T(), accounts)
	})
}

func (s *eth1StoreTestSuite) TestList() {
	ctx := context.Background()

	s.Run("should list all Ethereum Accounts successfully", func() {
		expectedAccounts := []*entities.ETH1Account{testutils.FakeETH1Account(), testutils.FakeETH1Account()}
		s.mockEth1AccountsDB.EXPECT().GetAll(ctx).Return(expectedAccounts, nil)

		addresses, err := s.eth1Store.List(ctx)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), []string{expectedAccounts[0].Address.Hex(), expectedAccounts[1].Address.Hex()}, addresses)
	})

	s.Run("should fail with same error if GetAll account fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().GetAll(ctx).Return(nil, expectedErr)

		accounts, err := s.eth1Store.List(ctx)
		assert.Equal(s.T(), expectedErr, err)
		assert.Empty(s.T(), accounts)
	})
}

func (s *eth1StoreTestSuite) TestUpdate() {
	ctx := context.Background()
	attributes := testutils.FakeAttributes()
	key := testutils.FakeKey()
	fakeAccount := testutils.FakeETH1Account()

	s.Run("should update an Ethereum Account successfully", func() {
		expectedUpdatedAccount := &entities.ETH1Account{
			ID:                  key.ID,
			Address:             common.HexToAddress(address),
			Metadata:            key.Metadata,
			PublicKey:           hexutil.MustDecode(pubKey),
			CompressedPublicKey: hexutil.MustDecode(compressedPubKey),
			Tags:                key.Tags,
		}

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Update(ctx, fakeAccount.ID, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Update(ctx, expectedUpdatedAccount).Return(nil)

		account, err := s.eth1Store.Update(ctx, address, attributes)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedUpdatedAccount, account)
	})

	s.Run("should fail with same error if Get account fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		account, err := s.eth1Store.Update(ctx, address, attributes)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})

	s.Run("should fail with same error if Update key fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Update(ctx, fakeAccount.ID, attributes).Return(nil, expectedErr)

		account, err := s.eth1Store.Update(ctx, address, attributes)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})

	s.Run("should fail with same error if Add account fails", func() {
		expectedUpdatedAccount := &entities.ETH1Account{
			ID:                  key.ID,
			Address:             common.HexToAddress(address),
			Metadata:            key.Metadata,
			PublicKey:           hexutil.MustDecode(pubKey),
			CompressedPublicKey: hexutil.MustDecode(compressedPubKey),
			Tags:                key.Tags,
		}
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Update(ctx, fakeAccount.ID, attributes).Return(key, nil)
		s.mockEth1AccountsDB.EXPECT().Update(ctx, expectedUpdatedAccount).Return(expectedErr)

		account, err := s.eth1Store.Update(ctx, address, attributes)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})
}

func (s *eth1StoreTestSuite) TestDelete() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()

	s.Run("should delete an Ethereum Account successfully", func() {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Delete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Remove(ctx, address).Return(nil)
		s.mockEth1AccountsDB.EXPECT().AddDeleted(ctx, fakeAccount).Return(nil)

		err := s.eth1Store.Delete(ctx, address)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if Get account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		err := s.eth1Store.Delete(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if Delete key fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Delete(ctx, fakeAccount.ID).Return(expectedErr)

		err := s.eth1Store.Delete(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if Delete account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Delete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Remove(ctx, address).Return(expectedErr)

		err := s.eth1Store.Delete(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if AddDeleted account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Delete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Remove(ctx, address).Return(nil)
		s.mockEth1AccountsDB.EXPECT().AddDeleted(ctx, fakeAccount).Return(expectedErr)

		err := s.eth1Store.Delete(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *eth1StoreTestSuite) TestGetDeleted() {
	ctx := context.Background()

	s.Run("should get a deleted Ethereum Account successfully", func() {
		fakeETH1Account := testutils.FakeETH1Account()
		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeETH1Account, nil)

		account, err := s.eth1Store.GetDeleted(ctx, address)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), fakeETH1Account, account)
	})

	s.Run("should fail with same error if Get account fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(nil, expectedErr)

		account, err := s.eth1Store.GetDeleted(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), account)
	})
}

func (s *eth1StoreTestSuite) TestListDeleted() {
	ctx := context.Background()

	s.Run("should list all Ethereum Accounts successfully", func() {
		expectedAccounts := []*entities.ETH1Account{testutils.FakeETH1Account(), testutils.FakeETH1Account()}
		s.mockEth1AccountsDB.EXPECT().GetAllDeleted(ctx).Return(expectedAccounts, nil)

		addresses, err := s.eth1Store.ListDeleted(ctx)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), []string{expectedAccounts[0].Address.Hex(), expectedAccounts[1].Address.Hex()}, addresses)
	})

	s.Run("should fail with same error if GetAll account fails", func() {
		expectedErr := fmt.Errorf("my error")
		s.mockEth1AccountsDB.EXPECT().GetAllDeleted(ctx).Return(nil, expectedErr)

		accounts, err := s.eth1Store.ListDeleted(ctx)
		assert.Equal(s.T(), expectedErr, err)
		assert.Empty(s.T(), accounts)
	})
}

func (s *eth1StoreTestSuite) TestUndelete() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()

	s.Run("should undelete an Ethereum Account successfully", func() {
		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Undelete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, fakeAccount).Return(nil)

		err := s.eth1Store.Undelete(ctx, address)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if GetDeleted account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(nil, expectedErr)

		err := s.eth1Store.Undelete(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if Undelete key fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Undelete(ctx, fakeAccount.ID).Return(expectedErr)

		err := s.eth1Store.Undelete(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if RemoveDeleted account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Undelete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(expectedErr)

		err := s.eth1Store.Undelete(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if Add account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Undelete(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(nil)
		s.mockEth1AccountsDB.EXPECT().Add(ctx, fakeAccount).Return(expectedErr)

		err := s.eth1Store.Undelete(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *eth1StoreTestSuite) TestDestroy() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()

	s.Run("should undelete an Ethereum Account successfully", func() {
		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Destroy(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(nil)

		err := s.eth1Store.Destroy(ctx, address)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if GetDeleted account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(nil, expectedErr)

		err := s.eth1Store.Destroy(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if Destroy key fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Destroy(ctx, fakeAccount.ID).Return(expectedErr)

		err := s.eth1Store.Destroy(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})

	s.Run("should fail with same error if RemoveDeleted account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().GetDeleted(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Destroy(ctx, fakeAccount.ID).Return(nil)
		s.mockEth1AccountsDB.EXPECT().RemoveDeleted(ctx, address).Return(expectedErr)

		err := s.eth1Store.Destroy(ctx, address)
		assert.Equal(s.T(), expectedErr, err)
	})
}

func (s *eth1StoreTestSuite) TestSignVerify() {
	ctx := context.Background()
	fakeAccount := testutils.FakeETH1Account()
	data := []byte("my data to sign")
	ecdsaSignature := hexutil.MustDecode("0x63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed")

	s.Run("should sign a payload successfully with appended V value and verify it", func() {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, crypto.Keccak256(data)).Return(ecdsaSignature, nil)

		signature, err := s.eth1Store.Sign(ctx, address, data)
		assert.NoError(s.T(), err)
		// Note this the returned signature is not the same as the ecdsaSignature! It has an appended byte
		assert.Equal(s.T(), "0x63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed01", hexutil.Encode(signature))
		err = s.eth1Store.Verify(ctx, address, data, signature)
		assert.NoError(s.T(), err)
	})

	s.Run("should fail with same error if Get account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		signature, err := s.eth1Store.Sign(ctx, address, data)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), signature)
	})

	s.Run("should fail with same error if Sign fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, crypto.Keccak256(data)).Return(nil, expectedErr)

		signature, err := s.eth1Store.Sign(ctx, address, data)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), signature)
	})
}

func (s *eth1StoreTestSuite) TestSignData() {
	ctx := context.Background()
	data := crypto.Keccak256([]byte("my data to sign"))
	fakeAccount := &entities.ETH1Account{
		ID:                  "my-account",
		Address:             common.HexToAddress("0x83a0254be47813BBff771F4562744676C4e793F0"),
		PublicKey:           hexutil.MustDecode("0x04555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f0974f2265a8a7d82208f88c21a2c55663b33e5af92d919252511638e82dff8b2"),
		CompressedPublicKey: hexutil.MustDecode("0x02555214986a521f43409c1c6b236db1674332faaaf11fc42a7047ab07781ebe6f"),
	}
	recID := "01"

	s.Run("should sign payload, with no malleable signature, successfully", func() {
		R, _ := new(big.Int).SetString("63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd", 16)
		S, _ := new(big.Int).SetString("39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed", 16)
		ecdsaSignature := append(R.Bytes(), S.Bytes()...)
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, data).Return(ecdsaSignature, nil)

		expectedSignature := hexutil.Encode(ecdsaSignature) + recID
		signature, err := s.eth1Store.SignData(ctx, address, data)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), hexutil.Encode(signature), expectedSignature)
	})

	s.Run("should sign payload, with malleable signature, successfully", func() {
		R, _ := new(big.Int).SetString("63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd", 16)
		S, _ := new(big.Int).SetString("39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed", 16)
		S2 := new(big.Int).Add(S, secp256k1N)
		ecdsaSignature := append(R.Bytes(), S2.Bytes()...)
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, data).Return(ecdsaSignature, nil)

		expectedSignature := hexutil.Encode(append(R.Bytes(), S.Bytes()...)) + recID
		signature, err := s.eth1Store.SignData(ctx, address, data)
		require.NoError(s.T(), err)
		assert.Equal(s.T(), hexutil.Encode(signature), expectedSignature)
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

	s.Run("should sign a payload successfully with appended V value", func() {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, types.NewEIP155Signer(chainID).Hash(tx).Bytes()).Return(ecdsaSignature, nil)

		signedRaw, err := s.eth1Store.SignTransaction(ctx, address, chainID, tx)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "0xf85d80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808025a0e276fd7524ed7af67b7f914de5be16fad6b9038009d2d78f2315351fbd48deeea057a897964e80e041c674942ef4dbd860cb79a6906fb965d5e4645f5c44f7eae4", hexutil.Encode(signedRaw))
	})

	s.Run("should fail with same error if Get account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignTransaction(ctx, address, chainID, tx)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), signedRaw)
	})

	s.Run("should fail with same error if Sign fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, gomock.Any()).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignTransaction(ctx, address, chainID, tx)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), signedRaw)
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

	s.Run("should sign a payload successfully with appended V value", func() {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, quorumtypes.QuorumPrivateTxSigner{}.Hash(tx).Bytes()).Return(ecdsaSignature, nil)

		signedRaw, err := s.eth1Store.SignPrivate(ctx, address, tx)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "0xf85d80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808026a080365b013992519479ddd83584039d66851da560dbbe67f59ab9bdcd97b62503a055e93d2c8050fb413956298c10eb7b8b2c8d76f4be261e458e4987cc5fed9f01", hexutil.Encode(signedRaw))
	})

	s.Run("should fail with same error if Get account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignPrivate(ctx, address, tx)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), signedRaw)
	})

	s.Run("should fail with same error if Sign fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, gomock.Any()).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignPrivate(ctx, address, tx)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), signedRaw)
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

	s.Run("should sign a payload with privacyFor successfully with appended V value", func() {
		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, hexutil.MustDecode("0x5749cc0adae7a54f9c5148a9e21719a2b472dec7b7ae7c1d68bf35e2e161f94d")).Return(ecdsaSignature, nil)

		signedRaw, err := s.eth1Store.SignEEA(ctx, address, chainID, tx, privateArgs)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), "0xf8cd80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808026a06854034c21ebb5a6d4aa9a9c1462862b1e4af355383413a0dcfbba309f56ed02a020c0ebc19f159ce83c24dde6f1b2d424025e45bc8b00be3e2fd4367949d4f0b3a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564", hexutil.Encode(signedRaw))
	})

	s.Run("should fail with same error if Get account fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignEEA(ctx, address, chainID, tx, privateArgs)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), signedRaw)
	})

	s.Run("should fail with same error if Sign fails", func() {
		expectedErr := fmt.Errorf("my error")

		s.mockEth1AccountsDB.EXPECT().Get(ctx, address).Return(fakeAccount, nil)
		s.mockKeyStore.EXPECT().Sign(ctx, fakeAccount.ID, gomock.Any()).Return(nil, expectedErr)

		signedRaw, err := s.eth1Store.SignEEA(ctx, address, chainID, tx, privateArgs)
		assert.Equal(s.T(), expectedErr, err)
		assert.Nil(s.T(), signedRaw)
	})
}
