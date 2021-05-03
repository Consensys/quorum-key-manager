package jsonrpc

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseWriterHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.(common.WriterWrapper).Writer().(http.ResponseWriter).Header().Add("test-key", "test-value")
	assert.Equal(t, "test-value", rec.Header().Get("test-key"), "Header should have been set")
}

func TestResponseWriterWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	n, err := rw.(common.WriterWrapper).Writer().Write([]byte(`t`))
	require.NoError(t, err, "Write should not error")
	assert.Equal(t, n, 1, "Write should return correct value")

	b, err := rec.Body.ReadByte()
	require.NoError(t, err, "ReadByte should not error")
	assert.Equal(t, b, []byte(`t`)[0], "Write should have wrote correctly")
}

func TestResponseWriterWriteMsg(t *testing.T) {
	// WriteMsg with fields set
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	msg := &ResponseMsg{
		Version: "3.0",
		ID:      "39",
		Result:  true,
	}
	err := rw.WriteMsg(msg)
	require.NoError(t, err, "WriteMsg should not error")

	expectedBody := []byte(`{"jsonrpc":"3.0","result":true,"error":null,"id":"39"}`)
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "WriteMsg should write correct body")

	// WriteMsg with default values set
	rec = httptest.NewRecorder()
	rw = NewResponseWriter(rec)

	msg = &ResponseMsg{
		Version: "2.0",
		Result:  true,
	}
	err = rw.WriteMsg(msg)
	require.NoError(t, err, "WriteMsg should not error")

	expectedBody = []byte(`{"jsonrpc":"2.0","result":true,"error":null,"id":null}`)
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "WriteMsg should write correct body")
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"), "Header Content-Type should have been set")
}

func TestWriteResult(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := RWWithVersion("1.0")(RWWithID(1234)(NewResponseWriter(rec))).(ResponseWriter)

	err := WriteResult(rw, false)
	require.NoError(t, err, "WriteResult should not error")
	expectedBody := []byte(`{"jsonrpc":"1.0","result":false,"error":null,"id":1234}`)
	assert.Equal(t, expectedBody, rec.Body.Bytes()[:(rec.Body.Len()-1)], "WriteResult should write correct body")
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := RWWithVersion("")(RWWithID("abcd")(NewResponseWriter(rec)).(ResponseWriter))

	err := WriteError(rw, fmt.Errorf("test error"))
	require.NoError(t, err, "WriteError should not error")
	expectedBody := []byte(`{"jsonrpc":"2.0","result":null,"error":{"code":-32603,"message":"Internal error","data":{"message":"test error"}},"id":"abcd"}`)
