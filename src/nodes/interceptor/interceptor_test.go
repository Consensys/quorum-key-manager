package interceptor

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	aliasmock "github.com/consensys/quorum-key-manager/src/aliases/mock"
	placeholdermock "github.com/consensys/quorum-key-manager/src/aliases/placeholder/mock"
	mockstoremanager "github.com/consensys/quorum-key-manager/src/stores/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newInterceptor(t *testing.T, ctrl *gomock.Controller) (*Interceptor, *mockstoremanager.MockStores) {
	stores := mockstoremanager.NewMockStores(ctrl)
	aliases := aliasmock.NewMockAliasBackend(ctrl)
	aliasParser := placeholdermock.NewMockAliasParser(ctrl)
	i, err := New(stores, aliases, aliasParser, testutils.NewMockLogger(ctrl))
	require.NoError(t, err)
	return i, stores
}

type testHandlerCase struct {
	desc string

	// JSON body of the response
	ctx context.Context

	prepare func()
	handler jsonrpc.Handler

	reqBody          []byte
	expectedRespBody []byte
}

func assertHandlerScenario(t *testing.T, tt *testHandlerCase) {
	if tt.prepare != nil {
		tt.prepare()
	}

	rec := httptest.NewRecorder()
	rw := jsonrpc.NewResponseWriter(rec)

	msg := new(jsonrpc.RequestMsg)
	err := json.Unmarshal(tt.reqBody, msg)
	require.NoError(t, err, "Unmarshal must not error")

	tt.handler.ServeRPC(rw, msg.WithContext(tt.ctx))

	assert.Equal(t, http.StatusOK, rec.Code, "Response code should be correct")
	assert.Equal(t, tt.expectedRespBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Response body should be correct")
}

func TestPersonal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	i, _ := newInterceptor(t, ctrl)
	tests := []*testHandlerCase{
		{
			desc:             "Personal",
			handler:          i,
			reqBody:          []byte(`{"jsonrpc":"2.0","method":"personal_accounts","params":[]}`),
			expectedRespBody: []byte(`{"jsonrpc":"2.0","result":null,"error":{"code":-32601,"message":"Method not found","data":null},"id":null}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assertHandlerScenario(t, tt)
		})
	}
}
