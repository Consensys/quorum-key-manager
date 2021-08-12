package interceptor

import (
	"context"
	"math/big"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	mockethereum "github.com/consensys/quorum-key-manager/pkg/ethereum/mock"
	mocktessera "github.com/consensys/quorum-key-manager/pkg/tessera/mock"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	mockaccounts "github.com/consensys/quorum-key-manager/src/stores/mock"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEthSendTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, stores := newInterceptor(ctrl)
	accountsStore := mockaccounts.NewMockEth1Store(ctrl)

	userInfo := &types.UserInfo{
		Username: "username",
		Groups:   []string{"group1", "group2"},
	}
	session := proxynode.NewMockSession(ctrl)
	ctx := proxynode.WithSession(context.TODO(), session)
	ctx = authenticator.WithUserContext(ctx, &authenticator.UserContext{
		UserInfo: userInfo,
	})

	cller := mockethereum.NewMockCaller(ctrl)
	ethCaller := mockethereum.NewMockEthCaller(ctrl)
	cller.EXPECT().Eth().Return(ethCaller).AnyTimes()
	session.EXPECT().EthCaller().Return(cller).AnyTimes()

	tesseraClient := mocktessera.NewMockClient(ctrl)
	session.EXPECT().ClientPrivTxManager().Return(tesseraClient).AnyTimes()

	tests := []*testHandlerCase{
		{
			desc:    "Public Transaction ",
			handler: i.handler,
			reqBody: []byte(`{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","data":"0x5208"}]}`),
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetEth1StoreByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)

				// Get Gas price
				ethCaller.EXPECT().GasPrice(gomock.Any()).Return(big.NewInt(1000000000), nil)

				// Get Gas Limit
				expectedCallMsg := (&ethereum.CallMsg{
					From:     &expectedFrom,
					GasPrice: big.NewInt(1000000000),
				}).WithData(ethcommon.FromHex("0x5208"))

				ethCaller.EXPECT().EstimateGas(gomock.Any(), expectedCallMsg).Return(uint64(21000), nil)

				// Get Nonce
				ethCaller.EXPECT().GetTransactionCount(gomock.Any(), expectedFrom, ethereum.PendingBlockNumber).Return(uint64(5), nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Sign
				accountsStore.EXPECT().SignTransaction(gomock.Any(), expectedFrom.Hex(), big.NewInt(1998), gomock.Any()).Return(ethcommon.FromHex("0xa6122e27"), nil)

				// SendRawTransaction
				ethCaller.EXPECT().SendRawTransaction(gomock.Any(), ethcommon.FromHex("0xa6122e27")).Return(ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)
			},
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778","error":null,"id":null}`),
		},
		{
			desc:    "Transaction private transaction",
			handler: i.handler,
			reqBody: []byte(`{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","data":"0x5208","privateFrom":"GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=","privateFor":["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="]}]}`),
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetEth1StoreByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)

				// Get Gas price
				ethCaller.EXPECT().GasPrice(gomock.Any()).Return(big.NewInt(1000000000), nil)

				// Get Gas Limit
				expectedCallMsg := (&ethereum.CallMsg{
					From:     &expectedFrom,
					GasPrice: big.NewInt(1000000000),
				}).WithData(ethcommon.FromHex("0x5208"))

				ethCaller.EXPECT().EstimateGas(gomock.Any(), expectedCallMsg).Return(uint64(21000), nil)

				// Get Nonce
				ethCaller.EXPECT().GetTransactionCount(gomock.Any(), expectedFrom, ethereum.PendingBlockNumber).Return(uint64(5), nil)

				tesseraClient.EXPECT().StoreRaw(gomock.Any(), ethcommon.FromHex("0x5208"), "GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=").Return(ethcommon.FromHex("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Sign
				accountsStore.EXPECT().SignPrivate(gomock.Any(), expectedFrom.Hex(), gomock.Any()).Return(ethcommon.FromHex("0xa6122e27"), nil)

				expectedPrivateArgs := (&ethereum.PrivateArgs{}).WithPrivateFrom("GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=").WithPrivateFor([]string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="})
				ethCaller.EXPECT().SendRawPrivateTransaction(gomock.Any(), ethcommon.FromHex("0xa6122e27"), expectedPrivateArgs).Return(ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)
			},
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778","error":null,"id":null}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
