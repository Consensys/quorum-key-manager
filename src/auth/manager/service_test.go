package manager

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	"github.com/consensys/quorum-key-manager/src/infra/manifests/mock"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaseManager(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := testutils.NewMockLogger(ctrl)
	mockManifestReader := mock.NewMockReader(ctrl)

	mngr := New(mockManifestReader, mockLogger)

	t.Run("should start service successfully by loading roles", func(t *testing.T) {
		testManifests := []*manifest.Manifest{
			{
				Kind:  "Role",
				Name:  "anonymous",
				Specs: json.RawMessage(`{"permission": ["proxy:nodes"]}`),
			},
			{
				Kind:  "Role",
				Name:  "guest",
				Specs: json.RawMessage(`{"permission": ["read:secrets","proxy:nodes"]}`),
			},
			{
				Kind:  "Role",
				Name:  "signer",
				Specs: json.RawMessage(`{"permission": ["read:ethereum","read:keys","sign:keys","sign:ethereum"]}`),
			}, {
				Kind:  "Role",
				Name:  "admin",
				Specs: json.RawMessage(`{"permission": ["read:ethereum","read:keys","sign:keys","sign:ethereum","create:ethereum","create:keys"]}`),
			},
		}

		mockManifestReader.EXPECT().Load().Return(testManifests, nil)

		err := mngr.Start(ctx)
		require.NoError(t, err)

		// Verifies that objects have been properly loaded
		guestRole, err := mngr.Role("guest")
		require.NoError(t, err)
		assert.Equal(t, "guest", guestRole.Name)
		assert.Equal(t, []entities.Permission{"read:secrets", "proxy:nodes"}, guestRole.Permissions)

		otherPermission := []entities.Permission{"destroy:keys"}
		userInfo := &entities.UserInfo{
			Roles:       []string{"signer", "admin"},
			Permissions: []entities.Permission{"destroy:keys"},
		}
		signerRole, err := mngr.Role("signer")
		require.NoError(t, err)
		adminRole, err := mngr.Role("admin")
		require.NoError(t, err)

		permissions := mngr.UserPermissions(userInfo)
		assert.Equal(t, append(append(otherPermission, signerRole.Permissions...), adminRole.Permissions...), permissions)
	})

	t.Run("should fail with ConfigError if manifest fails to be loaded", func(t *testing.T) {
		mockManifestReader.EXPECT().Load().Return(nil, fmt.Errorf("error"))

		err := mngr.Start(ctx)
		assert.True(t, errors.IsConfigError(err))
	})
}
