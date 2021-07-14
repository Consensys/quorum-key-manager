package interceptor

import (
	"context"
	"math/big"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/ethereum"
	mockethereum "github.com/consensys/quorum-key-manager/pkg/ethereum/mock"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	mockaccounts "github.com/consensys/quorum-key-manager/src/stores/store/eth1/mock"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEEASendTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, stores := newInterceptor(ctrl)
	accountsStore := mockaccounts.NewMockStore(ctrl)

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
	eeaCaller := mockethereum.NewMockEEACaller(ctrl)
	ethCaller := mockethereum.NewMockEthCaller(ctrl)
	privCaller := mockethereum.NewMockPrivCaller(ctrl)
	cller.EXPECT().EEA().Return(eeaCaller).AnyTimes()
	cller.EXPECT().Eth().Return(ethCaller).AnyTimes()
	cller.EXPECT().Priv().Return(privCaller).AnyTimes()

	session.EXPECT().EthCaller().Return(cller).AnyTimes()

	tests := []*testHandlerCase{
		{
			desc:    "Transaction with Privacy Group ID",
			handler: i,
			reqBody: []byte(`{"jsonrpc":"2.0","method":"eea_sendTransaction","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","privacyGroupId":"kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="}],"id":"abcd"}`),
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetEth1StoreByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Get Gas price
				ethCaller.EXPECT().GasPrice(gomock.Any()).Return(big.NewInt(1000000000), nil)

				ethCaller.EXPECT().EstimateGas(gomock.Any(), gomock.Any()).Return(uint64(21000), nil)

				// Get Nonc
				privCaller.EXPECT().GetTransactionCount(gomock.Any(), expectedFrom, "kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=").Return(uint64(5), nil)

				// SignEEA
				expectedPrivateArgs := (&ethereum.PrivateArgs{PrivateType: common.ToPtr(privateTxTypeRestricted).(*string)}).WithPrivacyGroupID("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=")
				accountsStore.EXPECT().SignEEA(gomock.Any(), expectedFrom.Hex(), big.NewInt(1998), gomock.Any(), expectedPrivateArgs).Return(ethcommon.FromHex("0xa6122e27"), nil)

				// SendRawTransaction
				eeaCaller.EXPECT().SendRawTransaction(gomock.Any(), ethcommon.FromHex("0xa6122e27")).Return(ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)
			},
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778","error":null,"id":"abcd"}`),
		},
		{
			desc:    "Transaction with privateFor",
			handler: i,
			reqBody: []byte(`{"jsonrpc":"2.0","method":"eea_sendTransaction","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","privateFrom":"GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=","privateFor":["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="]}],"id":"abcd"}`),
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetEth1StoreByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Get Gas price
				ethCaller.EXPECT().GasPrice(gomock.Any()).Return(big.NewInt(1000000000), nil)

				ethCaller.EXPECT().EstimateGas(gomock.Any(), gomock.Any()).Return(uint64(21000), nil)

				// Get Nonc
				privCaller.EXPECT().GetEeaTransactionCount(gomock.Any(), expectedFrom, "GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=", []string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="}).Return(uint64(5), nil)

				// Sign
				expectedPrivateArgs := (&ethereum.PrivateArgs{PrivateType: common.ToPtr(privateTxTypeRestricted).(*string)}).WithPrivateFrom("GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=").WithPrivateFor([]string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="})
				accountsStore.EXPECT().SignEEA(gomock.Any(), expectedFrom.Hex(), big.NewInt(1998), gomock.Any(), expectedPrivateArgs).Return(ethcommon.FromHex("0xa6122e27"), nil)

				eeaCaller.EXPECT().SendRawTransaction(gomock.Any(), ethcommon.FromHex("0xa6122e27")).Return(ethcommon.HexToHash("0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"), nil)
			},
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778","error":null,"id":"abcd"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
