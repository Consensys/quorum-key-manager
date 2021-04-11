package jsonrpc

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequest(t *testing.T) {
	// Create request
	body := bytes.NewReader([]byte(`{"jsonrpc": "1.0", "id": "25", "method": "testMethod", "params": {"test-field": "test-value"}}`))
	httpReq, _ := http.NewRequest(http.MethodGet, "www.test.com", body)
	req := NewRequest(httpReq)

	// Read body
	err := req.ReadBody()
	require.NoError(t, err, "ReadBody should not error")

	assert.Equal(t, "1.0", req.Version(), "Version should be correct")
	assert.Equal(t, "testMethod", req.Method(), "Method should be correct")

	// Unmarshal ID
	id := new(string)
	err = req.UnmarshalID(id)
	require.NoError(t, err, "UnmarshalID should not error")
	assert.Equal(t, "25", *id, "ID should be correct")
	assert.Equal(t, id, req.ID(), "ID should have been set correctly")

	type TestParams struct {
		Field string `json:"test-field"`
	}

	// Unmarshal Params
	params := new(TestParams)
	err = req.UnmarshalParams(params)
	require.NoError(t, err, "UnmarshalParams should not error")
	assert.Equal(t, TestParams{Field: "test-value"}, *params, "Params should be correct")
	assert.Equal(t, params, req.Params(), "Params should have been set correctly")

	// Set fields
	req.WithVersion("2.0")
	req.WithMethod("testMethod2")
	req.WithID("abcd")
	req.WithParams(true)

	// Write body
	err = req.WriteBody()
	require.NoError(t, err, "WriteBody should not error")

	b := make([]byte, req.req.ContentLength-1)
	_, err = io.ReadFull(req.req.Body, b)
	require.NoError(t, err, "Read body should not error")

	expectedBody := []byte(`{"jsonrpc":"2.0","method":"testMethod2","params":true,"id":"abcd"}`)
	assert.Equal(t, expectedBody, b, "Body should match")
}
