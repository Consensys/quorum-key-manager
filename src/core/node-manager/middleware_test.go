package nodemanager

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/node-manager/mock"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
	mocknode "github.com/ConsenSysQuorum/quorum-key-manager/src/node/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type requestMatcher struct {
	session node.Session
}

func (m requestMatcher) Matches(x interface{}) bool {
	req, ok := x.(*jsonrpc.Request)
	if !ok {
		return false
	}

	sess := node.SessionFromContext(req.Request().Context())

	return sess == m.session

}

func (m requestMatcher) String() string {
	return ""
}

func TestMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mngr := mock.NewMockManager(ctrl)
	mid := NewMiddleware(mngr)
	n := mocknode.NewMockNode(ctrl)
	next := jsonrpc.NewMockHandler(ctrl)

	handler := mid.Next(next)

	// Create RW
	rec := httptest.NewRecorder()
	rw := jsonrpc.NewResponseWriter(rec)

	// Create request
	body := bytes.NewReader([]byte(`{"jsonrpc": "1.0", "id": "25", "method": "testMethod", "params": {"test-field": "test-value"}}`))
	httpReq, _ := http.NewRequest(http.MethodGet, "www.test.com", body)
	req := jsonrpc.NewRequest(httpReq)
	err := req.ReadBody()
	require.NoError(t, err, "ReadBody should not error")

	mngr.EXPECT().Node(gomock.Any(), "").Return(n, nil)

	session := mocknode.NewMockSession(ctrl)
	n.EXPECT().Session(req).Return(session, nil)

	next.EXPECT().ServeRPC(rw, &requestMatcher{session})
	handler.ServeRPC(rw, req)
}
