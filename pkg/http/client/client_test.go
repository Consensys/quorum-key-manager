package httpclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	client, err := New(nil, nil, nil, nil)
	require.NoError(t, err, "New must not error")
	assert.Implements(t, (*Client)(nil), client, "Client should match Client interface")
}
