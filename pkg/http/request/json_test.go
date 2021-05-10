package request

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestMsg struct {
	Field []string `json:"field"`
}

func TestWriteJSON(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "", nil)

	msg := &TestMsg{
		Field: []string{"Hello", "World"},
	}
	err := WriteJSON(req, msg)
	require.NoError(t, err, "WriteJSON must not error")
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"), "Content-Type Header must be correct")

	b, err := ioutil.ReadAll(req.Body)
	require.NoError(t, err, "ReadAll must not error")
	expectedBody := []byte(`{"field":["Hello","World"]}`)
	assert.Equal(t, expectedBody, b, "Body should be correct")
}
