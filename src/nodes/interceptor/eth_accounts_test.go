package interceptor

import (
	"testing"

	"github.com/consensys/quorum-key-manager/src/stores/store/entities"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEthAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, stores := newInterceptor(ctrl)
	tests := []*testHandlerCase{
		{
			desc:    "Signature",
			handler: i,
			prepare: func() {
				accts := []*entities.ETH1Account{
					{Address: ethcommon.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73")},
					{Address: ethcommon.HexToAddress("0xea674fdde714fd979de3edf0f56aa9716b898ec8")},
				}
				stores.EXPECT().ListAllAccounts(gomock.Any()).Return(accts, nil)
			},
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"eth_accounts","params":[]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":["0xfe3b557e8fb62b89f4916b721be55ceb828dbd73","0xea674fdde714fd979de3edf0f56aa9716b898ec8"],"error":null,"id":null}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
