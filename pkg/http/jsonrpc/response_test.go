package jsonrpc

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponse(t *testing.T) {
	// Create request
	buf := bytes.NewReader([]byte(`{"jsonrpc": "1.0", "id": "25", "result": {"test-field": "test-value", "test-field2":[1,2,3]}, "error": {"code": -32600, "message":"test error"}}`))
	resp := NewResponse(&http.Response{
		Body: ioutil.NopCloser(buf),
	})

	// Read body
	err := resp.ReadBody()
	require.NoError(t, err, "ReadBody should not error")

	assert.Equal(t, "1.0", resp.Version(), "Version should be correct")
	assert.Equal(t, -32600, resp.Error().(*ErrorMsg).Code, "Error code should be correct")
	assert.Equal(t, "test error", resp.Error().(*ErrorMsg).Message, "Error message be correct")

	// Unmarshal ID
	id := new(string)
	err = resp.UnmarshalID(id)
	require.NoError(t, err, "UnmarshalID should not error")
	assert.Equal(t, "25", *id, "ID should be correct")
	assert.Equal(t, id, resp.ID(), "ID should have been set correctly")

	type TestResult struct {
		Field string `json:"test-field"`
	}

	// Unmarshal Params
	result := new(TestResult)
	err = resp.UnmarshalResult(result)
	require.NoError(t, err, "UnmarshalResult should not error")
	assert.Equal(t, TestResult{Field: "test-value"}, *result, "Result should be correct")
	assert.Equal(t, result, resp.Result(), "Result should have been set correctly")
}
