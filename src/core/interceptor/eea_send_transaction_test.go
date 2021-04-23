package interceptor

import (
	"bytes"
	"io/ioutil"
	"math/big"
	"net/http"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/testutils"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	mocknode "github.com/ConsenSysQuorum/quorum-key-manager/src/node/mock"
	mockaccounts "github.com/ConsenSysQuorum/quorum-key-manager/src/store/accounts/mock"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEEASendTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, stores, nodes := newInterceptor(ctrl)
	accountsStore := mockaccounts.NewMockStore(ctrl)
	n := mocknode.NewMockNode(ctrl)
	session := mocknode.NewMockSession(ctrl)
	nodes.EXPECT().Node(gomock.Any(), "default").Return(n, nil).AnyTimes()
	n.EXPECT().Session(gomock.Any()).Return(session, nil).AnyTimes()

	transport := testutils.NewMockRoundTripper(ctrl)
	jsonrpcClient := jsonrpc.NewClient(&http.Client{Transport: transport})
	req, _ := http.NewRequest(http.MethodPost, "www.example.com", nil)
	cllr := jsonrpc.NewCaller(jsonrpc.WithVersion("2.0")(jsonrpcClient), jsonrpc.NewRequest(req))
	client := ethereum.NewClient(cllr)

	handler := NewNodeSessionMiddleware(nodes).Next(i.EEASendTransaction())
	tests := []*testHandlerCase{
		{
			desc:    "Transaction with Privacy Group ID",
			handler: handler,
			reqBody: []byte(`{"jsonrpc":"2.0","method":"test","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","nonce":"0x5", "privacyGroupId":"kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="}]}`),
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetAccountStoreByAddr(gomock.Any(), expectedFrom).Return(accountsStore, nil)

				// Get ChainID
				session.EXPECT().EthClient().Return(client)
				m := testutils.RequestMatcher(
					t,
					"www.example.com",
					[]byte(`{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":null}`),
				)
				respBody := []byte(`{"jsonrpc": "2.0","result":"0x7ce"}`)
				transport.EXPECT().RoundTrip(m).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
				}, nil)

				// Sign
				expectedTxData := &ethereum.TxData{
					Nonce:    5,
					To:       nil,
					Value:    big.NewInt(0),
					GasPrice: big.NewInt(10000000000000),
					GasLimit: 21000,
				}

				expectedPrivateArgs := (&ethereum.PrivateArgs{}).WithPrivacyGroupID("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=")
				accountsStore.EXPECT().SignEEA(gomock.Any(), big.NewInt(1998), expectedFrom, expectedTxData, expectedPrivateArgs).Return(ethcommon.FromHex("0xa6122e27"), nil)

				session.EXPECT().EthClient().Return(client)
				m = testutils.RequestMatcher(
					t,
					"www.example.com",
					[]byte(`{"jsonrpc":"2.0","method":"eea_sendRawTransaction","params":["0xa6122e27"],"id":null}`),
				)
				respBody = []byte(`{"jsonrpc": "2.0","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"}`)
				transport.EXPECT().RoundTrip(m).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
				}, nil)
			},
			expectedRespBody: []byte(`{"jsonrpc":"","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778","error":null,"id":null}`),
		},
		{
			desc:    "Transaction with privateFor",
			handler: handler,
			reqBody: []byte(`{"jsonrpc":"2.0","method":"test","params":[{"from":"0x78e6e236592597c09d5c137c2af40aecd42d12a2","gas":"0x5208","gasPrice":"0x9184e72a000","nonce":"0x5", "privateFrom":"GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=","privateFor":["KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=","eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="]}]}`),
			prepare: func() {
				expectedFrom := ethcommon.HexToAddress("0x78e6e236592597c09d5c137c2af40aecd42d12a2")
				// Get accounts
				stores.EXPECT().GetAccountStoreByAddr(gomock.Any(), expectedFrom).Return(accountsStore, nil)

				// Get ChainID
				session.EXPECT().EthClient().Return(client)
				m := testutils.RequestMatcher(
					t,
					"www.example.com",
					[]byte(`{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":null}`),
				)
				respBody := []byte(`{"jsonrpc": "2.0","result":"0x7ce"}`)
				transport.EXPECT().RoundTrip(m).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
				}, nil)

				// Sign
				expectedTxData := &ethereum.TxData{
					Nonce:    5,
					To:       nil,
					Value:    big.NewInt(0),
					GasPrice: big.NewInt(10000000000000),
					GasLimit: 21000,
				}

				expectedPrivateArgs := (&ethereum.PrivateArgs{}).WithPrivateFrom("GGilEkXLaQ9yhhtbpBT03Me9iYa7U/mWXxrJhnbl1XY=").WithPrivateFor([]string{"KkOjNLmCI6r+mICrC6l+XuEDjFEzQllaMQMpWLl4y1s=", "eLb69r4K8/9WviwlfDiZ4jf97P9czyS3DkKu0QYGLjg="})
				accountsStore.EXPECT().SignEEA(gomock.Any(), big.NewInt(1998), expectedFrom, expectedTxData, expectedPrivateArgs).Return(ethcommon.FromHex("0xa6122e27"), nil)

				session.EXPECT().EthClient().Return(client)
				m = testutils.RequestMatcher(
					t,
					"www.example.com",
					[]byte(`{"jsonrpc":"2.0","method":"eea_sendRawTransaction","params":["0xa6122e27"],"id":null}`),
				)
				respBody = []byte(`{"jsonrpc": "2.0","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778"}`)
				transport.EXPECT().RoundTrip(m).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewReader(respBody)),
				}, nil)
			},
			expectedRespBody: []byte(`{"jsonrpc":"","result":"0x6052dd2131667ef3e0a0666f2812db2defceaec91c470bb43de92268e8306778","error":null,"id":null}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
