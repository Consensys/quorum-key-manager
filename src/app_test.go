package src

import (
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/http/server"
	"github.com/consensys/quorum-key-manager/pkg/log"
	"github.com/consensys/quorum-key-manager/src/auth"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	dir := t.TempDir()
	_, err := New(&Config{
		HTTP:      server.NewDefaultConfig(),
		Auth:      &auth.Config{},
		Manifests: &manifestsmanager.Config{Path: dir},
	},
		log.NewLogger(&log.Config{}))
	require.NoError(t, err, "New must not error")
}
