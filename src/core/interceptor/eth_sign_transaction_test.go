package interceptor

import (
	"math/big"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	mockethereum "github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum/mock"
	mocknode "github.com/ConsenSysQuorum/quorum-key-manager/src/node/mock"
	mockaccounts "github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts/mock"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEthSignTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, stores, nodes := newInterceptor(ctrl)
	accountsStore := mockaccounts.NewMockStore(ctrl)

	n := mocknode.NewMockNode(ctrl)
	session := mocknode.NewMockSession(ctrl)
	nodes.EXPECT().Node(gomock.Any(), "default").Return(n, nil).AnyTimes()
	n.EXPECT().Session(gomock.Any()).Return(session, nil).AnyTimes()

	cller := mockethereum.NewMockCaller(ctrl)
	eeaCaller := mockethereum.NewMockEEACaller(ctrl)
	ethCaller := mockethereum.NewMockEthCaller(ctrl)
	cller.EXPECT().EEA().Return(eeaCaller).AnyTimes()
	cller.EXPECT().Eth().Return(ethCaller).AnyTimes()

	session.EXPECT().EthCaller().Return(cller).AnyTimes()

	handler := NewNodeSessionMiddleware(nodes).Next(i.EthSignTransaction())
	tests := []*testHandlerCase{
		{
			desc:    "Public transaction",
			handler: handler,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetAccountStoreByAddr(gomock.Any(), expectedFrom).Return(accountsStore, nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Sign
				expectedTxData := &ethereum.TxData{
					Nonce:    5,
					To:       nil,
					Value:    big.NewInt(0),
					GasPrice: big.NewInt(10000000000000),
					GasLimit: 21000,
				}
				accountsStore.EXPECT().SignEIP155(gomock.Any(), big.NewInt(1998), expectedFrom, expectedTxData).Return(ethcommon.FromHex("0xa6122e27"), nil)
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"test","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","nonce":"0x5"}]}`),
			expectedRespBody: []byte(`{"jsonrpc":"","result":"0xa6122e27","error":null,"id":null}`),
		},
		{
			desc:    "Private transaction",
			handler: handler,
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetAccountStoreByAddr(gomock.Any(), expectedFrom).Return(accountsStore, nil)

				// Get ChainID
				ethCaller.EXPECT().ChainID(gomock.Any()).Return(big.NewInt(1998), nil)

				// Sign
				expectedTxData := &ethereum.TxData{
					Nonce:    5,
					To:       nil,
					Value:    big.NewInt(0),
					GasPrice: big.NewInt(10000000000000),
					GasLimit: 21000,
				}
				accountsStore.EXPECT().SignPrivate(gomock.Any(), expectedFrom, expectedTxData).Return(ethcommon.FromHex("0xa6122e27"), nil)
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"test","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","nonce":"0x5","privateFor":["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="]}]}`),
			expectedRespBody: []byte(`{"jsonrpc":"","result":"0xa6122e27","error":null,"id":null}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
