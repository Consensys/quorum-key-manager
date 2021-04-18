package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testHandler(rw ResponseWriter, req *Request) {
	_ = rw.WriteResult(req.msg.raw.Params)
}

func TestToHTTPHandler(t *testing.T) {
	httpHandler := ToHTTPHandler(HandlerFunc(testHandler))

	// Fresh request
	body := bytes.NewReader([]byte(`{"jsonrpc": "1.0", "id": "abcd", "method": "testMethod", "params": {"test-field": "test-value"}}`))
	req, _ := http.NewRequest(http.MethodPost, "www.test.com", body)

	rec := httptest.NewRecorder()
	httpHandler.ServeHTTP(rec, req)

	expectedBody := []byte(`{"jsonrpc":"1.0","result":{"test-field":"test-value"},"error":null,"id":"abcd"}`)
	assert.Equal(t, http.StatusOK, rec.Code, "Code should be correct")
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Correct body should have been written")

	// Request with attached context
	rpcReq := NewRequest(nil).
		WithVersion("3.0").
		WithMethod("testMethod").
		WithID("abcd")
	rpcReq.msg.raw = &jsonReqMsg{}
	rpcReq.msg.raw.Params = new(json.RawMessage)
	*rpcReq.msg.raw.Params = json.RawMessage(`{"test-field":"test-value"}`)

	ctx := WithRequest(context.Background(), rpcReq)
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, "www.test.com", body)

	rec = httptest.NewRecorder()

	httpHandler.ServeHTTP(rec, req)

	expectedBody = []byte(`{"jsonrpc":"3.0","result":{"test-field":"test-value"},"error":null,"id":"abcd"}`)
	assert.Equal(t, http.StatusOK, rec.Code, "Code should be correct")
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Correct body should have been written")

	// Already wrapped responsewriter
	rpcReq = NewRequest(nil).
		WithVersion("3.0").
		WithMethod("testMethod").
		WithID("abcd")
	rpcReq.msg.raw = &jsonReqMsg{}
	rpcReq.msg.raw.Params = new(json.RawMessage)
	*rpcReq.msg.raw.Params = json.RawMessage(`{"test-field":"test-value"}`)

	ctx = WithRequest(context.Background(), rpcReq)
	req, _ = http.NewRequestWithContext(ctx, http.MethodPost, "www.test.com", body)

	rec = httptest.NewRecorder()
	rw := NewResponseWriter(rec).WithVersion("2.1").WithID(1234)

	httpHandler.ServeHTTP(rw, req)

	expectedBody = []byte(`{"jsonrpc":"2.1","result":{"test-field":"test-value"},"error":null,"id":1234}`)
	assert.Equal(t, http.StatusOK, rec.Code, "Code should be correct")
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Correct body should have been written")
}

func testHTTPHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)

	b := make([]byte, req.ContentLength)
	_, _ = io.ReadFull(req.Body, b)
	_, _ = rw.Write(b)
}

func TestFromHTTPHandler(t *testing.T) {
	rpcHandler := FromHTTPHandler(http.HandlerFunc(testHTTPHandler))

	// Fresh request
	req, _ := http.NewRequest(http.MethodPost, "www.test.com", nil)
	rpcReq := NewRequest(req).WithVersion("3.0").WithID("abcd").WithMethod("testMethod").WithParams([]int{1, 2, 3})

	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rpcHandler.ServeRPC(rw, rpcReq)

	expectedBody := []byte(`{"jsonrpc":"3.0","method":"testMethod","params":[1,2,3],"id":"abcd"}`)
	assert.Equal(t, http.StatusOK, rec.Code, "Code should be correct")
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "Correct body should have been written")

}
