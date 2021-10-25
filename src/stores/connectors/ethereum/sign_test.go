package eth

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	common2 "github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	mock3 "github.com/consensys/quorum-key-manager/src/auth/mock"
	authtypes "github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mock2 "github.com/consensys/quorum-key-manager/src/stores/database/mock"
	testutils2 "github.com/consensys/quorum-key-manager/src/stores/entities/testutils"
	"github.com/consensys/quorum-key-manager/src/stores/mock"
	quorumtypes "github.com/consensys/quorum/core/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignMessage(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := hexutil.MustDecode("0xfeaa")
	expectedData := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%v", 2, string(data))
	expectedErr := fmt.Errorf("my error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	t.Run("should sign successfully", func(t *testing.T) {
		acc := testutils2.FakeETHAccount()
		ecdsaSignature := hexutil.MustDecode("0xe276fd7524ed7af67b7f914de5be16fad6b9038009d2d78f2315351fbd48deee57a897964e80e041c674942ef4dbd860cb79a6906fb965d5e4645f5c44f7eae4")
		acc.PublicKey = hexutil.MustDecode("0x04e2e7621c0c08e43905648be731a482e8eb3d3186023335812f52130e4a18dd729b22d88fbf0f22b8fa4390267ef0c54367dc638a25b38ea74290bdb9f79ff917")

		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(gomock.Any(), acc.KeyID, crypto.Keccak256([]byte(expectedData)), ethAlgo).Return(ecdsaSignature, nil)

		expectedSignature := hexutil.Encode(ecdsaSignature) + "1b"
		signature, err := connector.SignMessage(ctx, acc.Address, data)

		require.NoError(t, err)
		assert.Equal(t, hexutil.Encode(signature), expectedSignature)
	})

	t.Run("should sign and convert malleable signature successfully", func(t *testing.T) {
		acc := testutils2.FakeETHAccount()
		R, _ := new(big.Int).SetString("63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd", 16)
		S, _ := new(big.Int).SetString("39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed", 16)
		S2 := new(big.Int).Add(S, secp256k1N)
		ecdsaSignatureMalleable := append(R.Bytes(), S2.Bytes()...)
		acc.PublicKey = hexutil.MustDecode("0x048e12a3333368c1c4c7c8fe9d1fd172444f6de328cd1650b76836928be9e71e952f61d75d87133f881bf4d25e02c3147aebf69abf7587c39cb5d177693ce1ca8f")

		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(gomock.Any(), acc.KeyID, crypto.Keccak256([]byte(expectedData)), ethAlgo).Return(ecdsaSignatureMalleable, nil)

		expectedSignature := hexutil.Encode(append(R.Bytes(), S.Bytes()...)) + "1c"
		signature, err := connector.SignMessage(ctx, acc.Address, data)

		require.NoError(t, err)
		assert.Equal(t, hexutil.Encode(signature), expectedSignature)
	})

	t.Run("should fail to sign if address is not recoverable", func(t *testing.T) {
		R, _ := new(big.Int).SetString("63341e2c837449de3735b6f4402b154aa0a118d02e45a2b311fba39c444025dd", 16)
		S, _ := new(big.Int).SetString("39db7699cb3d8a5caf7728a87e778c2cdccc4085cf2a346e37c1823dec5ce2ed", 16)
		ecdsaSignature := append(R.Bytes(), S.Bytes()...)
		acc := testutils2.FakeETHAccount()
		acc.PublicKey = hexutil.MustDecode("0x148a6e95f1f0f5d1b0aa4cc16a4b9d8bcfc666a538eb49af436e92285673a56830a57bf228fa5e4fff9445ed51b7923153519b316c4d71bea83911cae1c5952a91")

		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(gomock.Any(), acc.KeyID, crypto.Keccak256([]byte(expectedData)), ethAlgo).Return(ecdsaSignature, nil)

		_, err := connector.SignMessage(ctx, acc.Address, data)

		require.Error(t, err)
		assert.True(t, errors.IsCryptoOperationError(err))
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		acc := testutils2.FakeETHAccount()

		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(expectedErr)

		_, err := connector.SignMessage(ctx, acc.Address, data)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to sign if db fails", func(t *testing.T) {
		acc := testutils2.FakeETHAccount()

		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(nil, expectedErr)

		_, err := connector.SignMessage(ctx, acc.Address, data)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})

	t.Run("should fail to sign if store fails", func(t *testing.T) {
		acc := testutils2.FakeETHAccount()

		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(gomock.Any(), acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(gomock.Any(), acc.KeyID, crypto.Keccak256([]byte(expectedData)), ethAlgo).Return(nil, expectedErr)

		_, err := connector.SignMessage(ctx, acc.Address, data)

		assert.Error(t, err)
		assert.Equal(t, err, expectedErr)
	})
}

func TestSignTransaction(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("my error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	acc := testutils2.FakeETHAccount()
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

	t.Run("should sign a payload successfully with appended V value", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(ctx, acc.KeyID, types.NewEIP155Signer(chainID).Hash(tx).Bytes(), ethAlgo).Return(ecdsaSignature, nil)

		signedRaw, err := connector.SignTransaction(ctx, acc.Address, chainID, tx)
		assert.NoError(t, err)
		assert.Equal(t, "0xf85d80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808025a0e276fd7524ed7af67b7f914de5be16fad6b9038009d2d78f2315351fbd48deeea057a897964e80e041c674942ef4dbd860cb79a6906fb965d5e4645f5c44f7eae4", hexutil.Encode(signedRaw))
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(expectedErr)

		signedRaw, err := connector.SignTransaction(ctx, acc.Address, chainID, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	t.Run("should fail with same error if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(nil, expectedErr)

		signedRaw, err := connector.SignTransaction(ctx, acc.Address, chainID, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	t.Run("should fail with same error if store fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(ctx, acc.KeyID, gomock.Any(), ethAlgo).Return(nil, expectedErr)

		signedRaw, err := connector.SignTransaction(ctx, acc.Address, chainID, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})
}

func TestSignPrivate(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("my error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	acc := testutils2.FakeETHAccount()
	tx := quorumtypes.NewTransaction(
		0,
		common.HexToAddress("0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18"),
		big.NewInt(0),
		0,
		big.NewInt(0),
		nil,
	)
	ecdsaSignature := hexutil.MustDecode("0x80365b013992519479ddd83584039d66851da560dbbe67f59ab9bdcd97b6250355e93d2c8050fb413956298c10eb7b8b2c8d76f4be261e458e4987cc5fed9f01")

	t.Run("should sign a payload successfully with appended V value", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(ctx, acc.KeyID, quorumtypes.QuorumPrivateTxSigner{}.Hash(tx).Bytes(), ethAlgo).Return(ecdsaSignature, nil)

		signedRaw, err := connector.SignPrivate(ctx, acc.Address, tx)
		assert.NoError(t, err)
		assert.Equal(t, "0xf85d80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808026a080365b013992519479ddd83584039d66851da560dbbe67f59ab9bdcd97b62503a055e93d2c8050fb413956298c10eb7b8b2c8d76f4be261e458e4987cc5fed9f01", hexutil.Encode(signedRaw))
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(expectedErr)

		signedRaw, err := connector.SignPrivate(ctx, acc.Address, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	t.Run("should fail with same error if db fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(nil, expectedErr)

		signedRaw, err := connector.SignPrivate(ctx, acc.Address, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	t.Run("should fail with same error if store fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(ctx, acc.KeyID, gomock.Any(), ethAlgo).Return(nil, expectedErr)

		signedRaw, err := connector.SignPrivate(ctx, acc.Address, tx)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})
}

func TestSignEEA(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedErr := fmt.Errorf("my error")

	store := mock.NewMockKeyStore(ctrl)
	db := mock2.NewMockETHAccounts(ctrl)
	logger := testutils.NewMockLogger(ctrl)
	auth := mock3.NewMockAuthorizator(ctrl)

	connector := NewConnector(store, db, auth, logger)

	acc := testutils2.FakeETHAccount()
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
	privateArgs := &ethereum.PrivateArgs{
		PrivateFrom: &privateFrom,
		PrivateFor:  &privateFor,
		PrivateType: common2.ToPtr(ethereum.PrivateTypeRestricted).(*ethereum.PrivateType),
	}
	ecdsaSignature := hexutil.MustDecode("0x6854034c21ebb5a6d4aa9a9c1462862b1e4af355383413a0dcfbba309f56ed0220c0ebc19f159ce83c24dde6f1b2d424025e45bc8b00be3e2fd4367949d4f0b3")

	t.Run("should sign a payload with privacyFor successfully with appended V value", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(ctx, acc.KeyID,
			hexutil.MustDecode("0x5749cc0adae7a54f9c5148a9e21719a2b472dec7b7ae7c1d68bf35e2e161f94d"),
			ethAlgo).Return(ecdsaSignature, nil)

		signedRaw, err := connector.SignEEA(ctx, acc.Address, chainID, tx, privateArgs)
		assert.NoError(t, err)
		assert.Equal(t, "0xf8cd80808094905b88eff8bda1543d4d6f4aa05afef143d27e18808026a06854034c21ebb5a6d4aa9a9c1462862b1e4af355383413a0dcfbba309f56ed02a020c0ebc19f159ce83c24dde6f1b2d424025e45bc8b00be3e2fd4367949d4f0b3a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486af842a0035695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486aa0075695b4cc4b0941e60551d7a19cf30603db5bfc23e5ac43a56f57f25f75486a8a72657374726963746564", hexutil.Encode(signedRaw))
	})

	t.Run("should fail with same error if authorization fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(expectedErr)

		signedRaw, err := connector.SignEEA(ctx, acc.Address, chainID, tx, privateArgs)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	t.Run("should fail with same error if Get account fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(nil, expectedErr)

		signedRaw, err := connector.SignEEA(ctx, acc.Address, chainID, tx, privateArgs)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})

	t.Run("should fail with same error if Sign fails", func(t *testing.T) {
		auth.EXPECT().CheckPermission(&authtypes.Operation{Action: authtypes.ActionSign, Resource: authtypes.ResourceEthAccount}).Return(nil)
		db.EXPECT().Get(ctx, acc.Address.Hex()).Return(acc, nil)
		store.EXPECT().Sign(ctx, acc.KeyID, gomock.Any(), ethAlgo).Return(nil, expectedErr)

		signedRaw, err := connector.SignEEA(ctx, acc.Address, chainID, tx, privateArgs)
		assert.Equal(t, expectedErr, err)
		assert.Nil(t, signedRaw)
	})
}
