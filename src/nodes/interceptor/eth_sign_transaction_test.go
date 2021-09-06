package interceptor

import (
	"context"
	"math/big"
	"testing"

	mockethereum "github.com/consensys/quorum-key-manager/pkg/ethereum/mock"
	"github.com/consensys/quorum-key-manager/src/auth/authenticator"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	mockaccounts "github.com/consensys/quorum-key-manager/src/stores/mock"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEthSignTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, stores := newInterceptor(ctrl)
	accountsStore := mockaccounts.NewMockEthStore(ctrl)

	session := proxynode.NewMockSession(ctrl)
	userInfo := &types.UserInfo{
		Username:    "username",
		Roles:       []string{"role1", "role2"},
		Permissions: []types.Permission{"write:key", "read:key", "sign:key"},
	}
	ctx := proxynode.WithSession(context.TODO(), session)
	ctx = authenticator.WithUserContext(ctx, &authenticator.UserContext{
		UserInfo: userInfo,
	})

	cller := mockethereum.NewMockCaller(ctrl)
	eeaCaller := mockethereum.NewMockEEACaller(ctrl)
	ethCaller := mockethereum.NewMockEthCaller(ctrl)
	cller.EXPECT().EEA().Return(eeaCaller).AnyTimes()
	cller.EXPECT().Eth().Return(ethCaller).AnyTimes()

	session.EXPECT().EthCaller().Return(cller).AnyTimes()

	tests := []*testHandlerCase{
		{
			desc:    "Public transaction",
			handler: i,
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetEthStoreByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Sign
				accountsStore.EXPECT().SignTransaction(gomock.Any(), expectedFrom, big.NewInt(1998), gomock.Any()).Return(ethcommon.FromHex("0xa6122e27"), nil)
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"eth_signTransaction","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9172a000","nonce":"0x5","data":"0x5208","value":"0x1"}]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0xa6122e27","error":null,"id":null}`),
		},
		{
			desc:    "Private transaction",
			handler: i,
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetEthStoreByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Sign
				accountsStore.EXPECT().SignPrivate(gomock.Any(), expectedFrom, gomock.Any()).Return(ethcommon.FromHex("0xa6122e27"), nil)
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"eth_signTransaction","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","nonce":"0x5","data":"0x5208","value":"0x1","privateFor":["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="]}]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0xa6122e27","error":null,"id":null}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
