package tls_test

import (
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/tls"
	"github.com/consensys/quorum-key-manager/pkg/tls/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	cfg, err := tls.NewConfig(testutils.ValidOpt)
	require.NoError(t, err, "Creating config should not error")
	require.NotNil(t, cfg, "Config should not be nil")

	assert.Len(t, cfg.Certificates, 1, "Certificates should have been set")
	assert.NotNil(t, cfg.ClientCAs, "ClientCA Pool should have been set")
}
