package src

import (
	"testing"

	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/consensys/quorum-key-manager/src/infra/postgres/client"

	"github.com/consensys/quorum-key-manager/pkg/http/server"
	"github.com/consensys/quorum-key-manager/src/auth"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dir := t.TempDir()
	_, err := New(&Config{
		HTTP:      server.NewDefaultConfig(),
		Auth:      &auth.Config{},
		Manifests: &manifestsmanager.Config{Path: dir},
		Postgres:  &client.Config{},
	}, testutils.NewMockLogger(ctrl))
	require.NoError(t, err, "New must not error")
}
