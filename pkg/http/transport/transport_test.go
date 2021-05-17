package transport

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransport(t *testing.T) {
	cfg := (&Config{
		EnableHTTP2: true,
		EnableH2C:   true,
	}).SetDefault()

	_, err := New(cfg)

	require.NoError(t, err, "New should not error")
}
