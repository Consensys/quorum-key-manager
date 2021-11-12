package interceptor

import (
	"context"
	"math/big"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/api/http_middlewares"

	aliasmock "github.com/consensys/quorum-key-manager/src/aliases/mock"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	mockaccounts "github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	mockethereum "github.com/consensys/quorum-key-manager/pkg/ethereum/mock"
	mocktessera "github.com/consensys/quorum-key-manager/pkg/tessera/mock"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEthSendTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	session := proxynode.NewMockSession(ctrl)
	caller := mockethereum.NewMockCaller(ctrl)
	ethCaller := mockethereum.NewMockEthCaller(ctrl)
	tesseraClient := mocktessera.NewMockClient(ctrl)
	accountsStore := mockaccounts.NewMockEthStore(ctrl)
	stores := mockaccounts.NewMockStores(ctrl)
	aliases := aliasmock.NewMockService(ctrl)

	hexFrom := "0x78e6e236592597c09d5c137c2af40aecd42d12a2"
	from := ethcommon.HexToAddress(hexFrom)
	privateFrom := "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="
	userInfo := &entities.UserInfo{
		Username:    "username",
		Roles:       []string{"role1", "role2"},
		Permissions: []entities.Permission{"write:key", "read:key", "sign:key"},
	}
	ctx := proxynode.WithSession(context.TODO(), session)
	ctx = http_middlewares.WithUserInfo(ctx, userInfo)
	gasPrice := big.NewInt(38)
	chainID := big.NewInt(1)
	value := big.NewInt(45)

	caller.EXPECT().Eth().Return(ethCaller).AnyTimes()
	session.EXPECT().EthCaller().Return(caller).AnyTimes()
	session.EXPECT().ClientPrivTxManager().Return(tesseraClient).AnyTimes()
	stores.EXPECT().EthereumByAddr(gomock.Any(), from, userInfo).Return(accountsStore, nil).AnyTimes()

	i, err := New(stores, aliases, testutils.NewMockLogger(ctrl))
	require.NoError(t, err)

	t.Run("should send a private tx successfully", func(t *testing.T) {
		privateFor := []string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}

		privateArgs := (&ethereum.PrivateArgs{}).
			WithPrivateFrom(privateFrom).
			WithPrivateFor(privateFor)
		msg := &ethereum.SendTxMsg{
			From:        from,
			PrivateArgs: *privateArgs,
		}
		expectedEstimateGasCall := &ethereum.CallMsg{
			From:     &msg.From,
			To:       msg.To,
			Value:    msg.Value,
			Data:     msg.Data,
			GasPrice: gasPrice,
		}
		expectedData := new([]byte)
		expectedSignedTx := []byte("mysignature")
		expectedHash := ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778")

		ethCaller.EXPECT().GasPrice(ctx).Return(gasPrice, nil)
		ethCaller.EXPECT().EstimateGas(ctx, expectedEstimateGasCall).Return(uint64(21000), nil)
		ethCaller.EXPECT().GetTransactionCount(ctx, msg.From, ethereum.PendingBlockNumber).Return(uint64(0), nil)
		tesseraClient.EXPECT().StoreRaw(ctx, *expectedData, *msg.PrivateFrom).Return(ethcommon.FromHex("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)
		ethCaller.EXPECT().ChainID(gomock.Any()).Return(chainID, nil)
		accountsStore.EXPECT().SignPrivate(ctx, msg.From, gomock.Any()).Return(expectedSignedTx, nil)
		ethCaller.EXPECT().SendRawPrivateTransaction(ctx, expectedSignedTx, privateArgs).Return(expectedHash, nil)
		aliases.EXPECT().ReplaceSimpleAlias(gomock.Any(), privateFrom).Return(privateFrom, nil)
		aliases.EXPECT().ReplaceAliases(gomock.Any(), privateFor).Return(privateFor, nil)

		hash, err := i.ethSendTransaction(ctx, msg)
		require.NoError(t, err)

		assert.Equal(t, hash.Hex(), expectedHash.Hex())
	})

	t.Run("should send an alias private tx successfully", func(t *testing.T) {
		privateFor := []string{"{{JPM:Group-A}}", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}
		privateForExp := []string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}

		privateArgs := (&ethereum.PrivateArgs{}).
			WithPrivateFrom(privateFrom).
			WithPrivateFor(privateFor)
		msg := &ethereum.SendTxMsg{
			From:        from,
			PrivateArgs: *privateArgs,
		}
		expectedEstimateGasCall := &ethereum.CallMsg{
			From:     &msg.From,
			To:       msg.To,
			Value:    msg.Value,
			Data:     msg.Data,
			GasPrice: gasPrice,
		}
		expectedData := new([]byte)
		expectedSignedTx := []byte("mysignature")
		expectedHash := ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778")

		ethCaller.EXPECT().GasPrice(ctx).Return(gasPrice, nil)
		ethCaller.EXPECT().EstimateGas(ctx, expectedEstimateGasCall).Return(uint64(21000), nil)
		ethCaller.EXPECT().GetTransactionCount(ctx, msg.From, ethereum.PendingBlockNumber).Return(uint64(0), nil)
		tesseraClient.EXPECT().StoreRaw(ctx, *expectedData, *msg.PrivateFrom).Return(ethcommon.FromHex("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)
		ethCaller.EXPECT().ChainID(gomock.Any()).Return(chainID, nil)
		accountsStore.EXPECT().SignPrivate(ctx, msg.From, gomock.Any()).Return(expectedSignedTx, nil)
		ethCaller.EXPECT().SendRawPrivateTransaction(ctx, expectedSignedTx, privateArgs).Return(expectedHash, nil)
		aliases.EXPECT().ReplaceSimpleAlias(gomock.Any(), privateFrom).Return(privateFrom, nil)
		aliases.EXPECT().ReplaceAliases(gomock.Any(), privateFor).Return(privateForExp, nil)

		hash, err := i.ethSendTransaction(ctx, msg)
		require.NoError(t, err)

		assert.Equal(t, hash.Hex(), expectedHash.Hex())
	})

	t.Run("should send an alias privacy group id successfully", func(t *testing.T) {
		privateForExp := []string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}

		privacyGroupIDAlias := "{{JPM:Group-A}}"

		privateArgs := &ethereum.PrivateArgs{
			PrivateFrom:    &privateFrom,
			PrivacyGroupID: &privacyGroupIDAlias,
		}

		privateArgsExp := &ethereum.PrivateArgs{
			PrivateFrom: &privateFrom,
			PrivateFor:  &privateForExp,
		}

		msg := &ethereum.SendTxMsg{
			From:        from,
			PrivateArgs: *privateArgs,
		}
		expectedEstimateGasCall := &ethereum.CallMsg{
			From:     &msg.From,
			To:       msg.To,
			Value:    msg.Value,
			Data:     msg.Data,
			GasPrice: gasPrice,
		}
		expectedData := new([]byte)
		expectedSignedTx := []byte("mysignature")
		expectedHash := ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778")

		ethCaller.EXPECT().GasPrice(ctx).Return(gasPrice, nil)
		ethCaller.EXPECT().EstimateGas(ctx, expectedEstimateGasCall).Return(uint64(21000), nil)
		ethCaller.EXPECT().GetTransactionCount(ctx, msg.From, ethereum.PendingBlockNumber).Return(uint64(0), nil)
		tesseraClient.EXPECT().StoreRaw(ctx, *expectedData, *msg.PrivateFrom).Return(ethcommon.FromHex("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)
		ethCaller.EXPECT().ChainID(gomock.Any()).Return(chainID, nil)
		accountsStore.EXPECT().SignPrivate(ctx, msg.From, gomock.Any()).Return(expectedSignedTx, nil)
		ethCaller.EXPECT().SendRawPrivateTransaction(gomock.Any(), expectedSignedTx, privateArgsExp).Return(expectedHash, nil)
		aliases.EXPECT().ReplaceAliases(gomock.Any(), []string{*privateArgs.PrivacyGroupID}).Return(privateForExp, nil)
		aliases.EXPECT().ReplaceSimpleAlias(gomock.Any(), privateFrom).Return(privateFrom, nil)

		hash, err := i.ethSendTransaction(ctx, msg)
		require.NoError(t, err)

		assert.Equal(t, hash.Hex(), expectedHash.Hex())
	})

	t.Run("should send a legacy tx successfully", func(t *testing.T) {
		msg := &ethereum.SendTxMsg{
			From:     from,
			GasPrice: gasPrice,
		}
		expectedEstimateGasCall := &ethereum.CallMsg{
			From:     &msg.From,
			To:       msg.To,
			Value:    msg.Value,
			Data:     msg.Data,
			GasPrice: gasPrice,
		}
		expectedSignedTx := []byte("mysignature")
		expectedHash := ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778")

		ethCaller.EXPECT().EstimateGas(ctx, expectedEstimateGasCall).Return(uint64(21000), nil)
		ethCaller.EXPECT().GetTransactionCount(ctx, msg.From, ethereum.PendingBlockNumber).Return(uint64(0), nil)
		ethCaller.EXPECT().ChainID(gomock.Any()).Return(chainID, nil)
		accountsStore.EXPECT().SignTransaction(ctx, msg.From, chainID, gomock.Any()).Return(expectedSignedTx, nil)
		ethCaller.EXPECT().SendRawTransaction(ctx, expectedSignedTx).Return(expectedHash, nil)

		hash, err := i.ethSendTransaction(ctx, msg)
		require.NoError(t, err)

		assert.Equal(t, hash.Hex(), expectedHash.Hex())
	})

	t.Run("should send a dynamic fee tx successfully", func(t *testing.T) {
		msg := &ethereum.SendTxMsg{
			From:  from,
			Value: value,
		}
		expectedEstimateGasCall := &ethereum.CallMsg{
			From:       &msg.From,
			To:         msg.To,
			Value:      value,
			Data:       msg.Data,
			GasFeeCap:  new(big.Int).Add(gasPrice, big.NewInt(0)),
			GasTipCap:  msg.GasTipCap,
			AccessList: msg.AccessList,
		}
		expectedSignedTx := []byte("mysignature")
		expectedHash := ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778")

		ethCaller.EXPECT().BaseFeePerGas(ctx, ethereum.LatestBlockNumber).Return(gasPrice, nil)
		ethCaller.EXPECT().EstimateGas(ctx, expectedEstimateGasCall).Return(uint64(21000), nil)
		ethCaller.EXPECT().GetTransactionCount(ctx, msg.From, ethereum.PendingBlockNumber).Return(uint64(0), nil)
		ethCaller.EXPECT().ChainID(gomock.Any()).Return(chainID, nil)
		accountsStore.EXPECT().SignTransaction(ctx, msg.From, chainID, gomock.Any()).Return(expectedSignedTx, nil)
		ethCaller.EXPECT().SendRawTransaction(ctx, expectedSignedTx).Return(expectedHash, nil)

		hash, err := i.ethSendTransaction(ctx, msg)
		require.NoError(t, err)

		assert.Equal(t, hash.Hex(), expectedHash.Hex())
	})

	t.Run("should revert to legacy tx if baseFeePerGas is nil", func(t *testing.T) {
		msg := &ethereum.SendTxMsg{
			From:  from,
			Value: value,
		}
		expectedEstimateGasCall := &ethereum.CallMsg{
			From:     &msg.From,
			To:       msg.To,
			Value:    msg.Value,
			Data:     msg.Data,
			GasPrice: gasPrice,
		}
		expectedSignedTx := []byte("mysignature")
		expectedHash := ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778")

		ethCaller.EXPECT().BaseFeePerGas(ctx, ethereum.LatestBlockNumber).Return(nil, nil)
		ethCaller.EXPECT().GasPrice(ctx).Return(gasPrice, nil)
		ethCaller.EXPECT().EstimateGas(ctx, expectedEstimateGasCall).Return(uint64(21000), nil)
		ethCaller.EXPECT().GetTransactionCount(ctx, msg.From, ethereum.PendingBlockNumber).Return(uint64(0), nil)
		ethCaller.EXPECT().ChainID(gomock.Any()).Return(chainID, nil)
		accountsStore.EXPECT().SignTransaction(ctx, msg.From, chainID, gomock.Any()).Return(expectedSignedTx, nil)
		ethCaller.EXPECT().SendRawTransaction(ctx, expectedSignedTx).Return(expectedHash, nil)

		hash, err := i.ethSendTransaction(ctx, msg)
		require.NoError(t, err)

		assert.Equal(t, hash.Hex(), expectedHash.Hex())
	})
}
