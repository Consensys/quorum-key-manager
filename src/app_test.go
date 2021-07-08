package src

import (
	postgresclient "github.com/consensys/quorum-key-manager/src/stores/infra/postgres/client"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/http/server"
	"github.com/consensys/quorum-key-manager/pkg/log/testutils"
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
		Postgres:  &postgresclient.Config{},
	}, testutils.NewMockLogger(ctrl))
	require.NoError(t, err, "New must not error")
}
