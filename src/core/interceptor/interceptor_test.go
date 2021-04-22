package interceptor

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	mocknodemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/node-manager/mock"
	mockstoremanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newInterceptor(ctrl *gomock.Controller) (i *Interceptor, stores *mockstoremanager.MockManager, nodes *mocknodemanager.MockManager) {
	stores = mockstoremanager.NewMockManager(ctrl)
	nodes = mocknodemanager.NewMockManager(ctrl)

	return New(stores, nodes), stores, nodes
}

type testHandlerCase struct {
	desc string

	// JSON body of the response
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

	httpReq, _ := http.NewRequest(http.MethodPost, "www.example.com", bytes.NewReader(tt.reqBody))
	req := jsonrpc.NewRequest(httpReq)
	err := req.ReadBody()
	require.NoError(t, err, "ReadBody must not error")

	tt.handler.ServeRPC(rw, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Response code should be correct")
	assert.Equal(t, tt.expectedRespBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Response body should be correct")
}
