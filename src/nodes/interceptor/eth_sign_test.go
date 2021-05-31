package interceptor

import (
	"fmt"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/errors"
	mockaccounts "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/store/eth1/mock"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEthSign(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, stores := newInterceptor(ctrl)
	accountsStore := mockaccounts.NewMockStore(ctrl)

	tests := []*testHandlerCase{
		{
			desc:    "Signature",
			handler: i.handler,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				stores.EXPECT().GetEth1StoreByAddr(gomock.Any(), expectedFrom).Return(accountsStore, nil)
				accountsStore.EXPECT().Sign(gomock.Any(), expectedFrom.Hex(), ethcommon.FromHex("0x2eadbe1f")).Return(ethcommon.FromHex("0xa6122e27"), nil)
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"eth_sign","params":["0x78e6e236592597c09d5c137c2af40aecd42d12a2", "0x2eadbe1f"]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":"0xa6122e27","error":null,"id":null}`),
		},
		{
			desc:    "Account not found",
			handler: i.handler,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				stores.EXPECT().GetEth1StoreByAddr(gomock.Any(), expectedFrom).Return(nil, errors.NotFoundError("account not found"))
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"eth_sign","params":["0x78e6e236592597c09d5c137c2af40aecd42d12a2", "0x2eadbe1f"]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":null,"error":{"code":-32603,"message":"Internal error","data":{"message":"account not found"}},"id":null}`),
		},
		{
			desc:    "Error signing",
			handler: i.handler,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				stores.EXPECT().GetEth1StoreByAddr(gomock.Any(), expectedFrom).Return(accountsStore, nil)
				accountsStore.EXPECT().Sign(gomock.Any(), expectedFrom.Hex(), ethcommon.FromHex("0x2eadbe1f")).Return(nil, fmt.Errorf("error signing"))
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
