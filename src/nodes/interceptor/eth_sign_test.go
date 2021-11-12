package interceptor

import (
	"context"
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/api/http_middlewares"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"

	mockaccounts "github.com/consensys/quorum-key-manager/src/stores/mock"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEthSign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userInfo := &entities.UserInfo{
		Username:    "username",
		Roles:       []string{"role1", "role2"},
		Permissions: []entities.Permission{"write:key", "read:key", "sign:key"},
	}

	session := proxynode.NewMockSession(ctrl)
	i, stores := newInterceptor(t, ctrl)
	accountsStore := mockaccounts.NewMockEthStore(ctrl)
	ctx := proxynode.WithSession(context.TODO(), session)
	ctx = http_middlewares.WithUserInfo(ctx, userInfo)

	tests := []*testHandlerCase{
		{
			desc:    "Signature",
			handler: i.handler,
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				stores.EXPECT().EthereumByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)
				accountsStore.EXPECT().Sign(gomock.Any(), expectedFrom, ethcommon.FromHex("0x2eadbe1f")).Return(ethcommon.FromHex("0xa6122e27"), nil)
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"eth_sign","params":["0x78e6e236592597c09d5c137c2af40aecd42d12a2", "0x2eadbe1f"]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0xa6122e27","error":null,"id":null}`),
		},
		{
			desc:    "Account not found",
			handler: i.handler,
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				stores.EXPECT().EthereumByAddr(gomock.Any(), expectedFrom, userInfo).Return(nil, errors.NotFoundError("account not found"))
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"eth_sign","params":["0x78e6e236592597c09d5c137c2af40aecd42d12a2", "0x2eadbe1f"]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":null,"error":{"code":-32603,"message":"Internal error","data":{"message":"ST100: account not found"}},"id":null}`),
		},
		{
			desc:    "Error signing",
			handler: i.handler,
			ctx:     ctx,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				stores.EXPECT().EthereumByAddr(gomock.Any(), expectedFrom, userInfo).Return(accountsStore, nil)
				accountsStore.EXPECT().Sign(gomock.Any(), expectedFrom, ethcommon.FromHex("0x2eadbe1f")).Return(nil, fmt.Errorf("error signing"))
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"eth_sign","params":["0x78e6e236592597c09d5c137c2af40aecd42d12a2", "0x2eadbe1f"]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":null,"error":{"code":-32603,"message":"Internal error","data":{"message":"error signing"}},"id":null}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
