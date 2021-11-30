package interceptor

import (
	"context"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/api/http"

	"github.com/consensys/quorum-key-manager/src/auth/entities"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
)

func TestEthAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInfo := &entities.UserInfo{
		Username:    "username",
		Roles:       []string{"role1", "role2"},
		Permissions: []entities.Permission{"write:key", "read:key", "sign:key"},
	}

	session := proxynode.NewMockSession(ctrl)
	ctx := proxynode.WithSession(context.TODO(), session)
	ctx = http.WithUserInfo(ctx, userInfo)

	i, stores, _ := newInterceptor(t, ctrl)
	tests := []*testHandlerCase{
		{
			desc:    "Signature",
			handler: i,
			ctx:     ctx,
			prepare: func() {
				accts := []ethcommon.Address{
					ethcommon.HexToAddress("0xfe3b557e8fb62b89f4916b721be55ceb828dbd73"),
					ethcommon.HexToAddress("0xea674fdde714fd979de3edf0f56aa9716b898ec8"),
				}
				stores.EXPECT().ListAllAccounts(gomock.Any(), userInfo).Return(accts, nil)
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
