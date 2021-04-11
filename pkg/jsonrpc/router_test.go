package jsonrpc

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	value int
}

func (h *mockHandler) ServeRPC(rw ResponseWriter, req *Request) {
	_ = rw.WithVersion(req.Version()).WithID(req.ID()).WriteResult(h.value)
}

func TestRouter(t *testing.T) {
	handlerDefault := &mockHandler{value: 1}
	handlerVersionDefault := &mockHandler{value: 2}
	handlerEthDefault := &mockHandler{value: 3}
	handlerEthSendTransaction := &mockHandler{value: 4}
	handlerEthSendRawTransaction := &mockHandler{value: 5}

	router := NewRouter().DefaultHandler(handlerDefault)
	v2Router := router.Version("2.0").Subrouter().DefaultHandler(handlerVersionDefault)
	ethRouter := v2Router.MethodPrefix("eth_").Subrouter().DefaultHandler(handlerEthDefault)
	ethRouter.Handle("eth_sendTransaction", handlerEthSendTransaction)
	ethRouter.Handle("eth_sendRawTransaction", handlerEthSendRawTransaction)

	tests := []struct {
		desc                   string
		req                    *Request
		shouldMatch            bool
		expectedMatchedHandler *mockHandler
	}{
		{
			desc:                   "invalid version",
			req:                    NewRequest(nil).WithVersion("1.0").WithMethod("testMethod"),
			shouldMatch:            true,
			expectedMatchedHandler: handlerDefault,
		},
		{
			desc:                   "valid version invalid prefix",
			req:                    NewRequest(nil).WithVersion("2.0").WithMethod("unknown_testMethod"),
			shouldMatch:            true,
			expectedMatchedHandler: handlerVersionDefault,
		},
		{
			desc:                   "valid version, valid prefix invalid method",
			req:                    NewRequest(nil).WithVersion("2.0").WithMethod("eth_unknown"),
			shouldMatch:            true,
			expectedMatchedHandler: handlerEthDefault,
		},
		{
			desc:                   "valid version, valid prefix, valid method eth_sendTransaction",
			req:                    NewRequest(nil).WithVersion("2.0").WithMethod("eth_sendTransaction"),
			shouldMatch:            true,
			expectedMatchedHandler: handlerEthSendTransaction,
		},
		{
			desc:                   "valid version, valid prefix, valid method eth_sendRawTransaction",
			req:                    NewRequest(nil).WithVersion("2.0").WithMethod("eth_sendRawTransaction"),
			shouldMatch:            true,
			expectedMatchedHandler: handlerEthSendRawTransaction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var match RouteMatch
			if tt.shouldMatch {
				require.True(t, router.Match(tt.req, &match), "Should match")
			} else {
				require.False(t, router.Match(tt.req, &match), "Should not match")
			}

			assert.Equal(t, tt.expectedMatchedHandler, match.Handler, "Matched handler should be correct")
		})
	}
}

func TestRouterServeRPC(t *testing.T) {
	router := NewRouter()
	router.Handle("known", &mockHandler{value: 1})

	// Request matching router
	req := NewRequest(nil).
		WithVersion("3.0").
		WithMethod("known").
		WithID("abcd")

	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	router.ServeRPC(rw, req)

	expectedBody := []byte(`{"jsonrpc":"3.0","result":1,"error":null,"id":"abcd"}`)
	assert.Equal(t, http.StatusOK, rec.Code, "Code should be correct")
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Correct body should have been written")

	// Request not matching router
	req = NewRequest(nil).
		WithVersion("3.0").
		WithMethod("unknown").
		WithID("abcd")

	rec = httptest.NewRecorder()
	rw = NewResponseWriter(rec)

	router.ServeRPC(rw, req)

	expectedBody = []byte(`{"jsonrpc":"2.0","result":null,"error":{"code":0,"message":"not supported","data":null},"id":null}`)
	assert.Equal(t, http.StatusOK, rec.Code, "Code should be correct")
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Correct body should have been written")
}
