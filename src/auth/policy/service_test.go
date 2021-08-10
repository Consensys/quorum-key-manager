package policy

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testManifest = []byte(`
- kind: Role
  name: anonymous
  specs:
    permission:
      - proxy:nodes
      - read:nodes
- kind: Role
  name: guest
  specs:
    permission:
      - read:secret
      - read:nodes
      - proxy:nodes
- kind: Role
  name: signer
  specs:
    permission:
      - read:eth1
      - read:key
      - sign:key
      - sign:eth1
- kind: Role
  name: admin
  specs:
    permission:
      - read:eth1
      - read:key
      - sign:key
      - sign:eth1
      - create:eth1
      - create:key
`)

func TestBaseManager(t *testing.T) {
	dir := t.TempDir()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := testutils.NewMockLogger(ctrl)
	err := ioutil.WriteFile(fmt.Sprintf("%v/manifest.yml", dir), testManifest, 0644)
	require.NoError(t, err, "WriteFile manifest1 must not error")

	manifests, err := manifestsmanager.NewLocalManager(&manifestsmanager.Config{Path: dir}, logger)
	require.NoError(t, err, "NewLocalManager on %v must not error", dir)

	err = manifests.Start(context.TODO())
	require.NoError(t, err, "Start manifests manager must not error")

	mngr := New(manifests, logger)
	err = mngr.Start(context.TODO())
	require.NoError(t, err, "Start manager manager must not error")

	// Give some time to load manifests
	time.Sleep(100 * time.Millisecond)

	// Verifies that objects have been properly loaded
	guestRole, err := mngr.Role(context.TODO(), "guest")
	require.NoError(t, err)
	assert.Equal(t, "guest", guestRole.Name)
	assert.Equal(t, []types.Permission{"read:secret", "read:nodes", "proxy:nodes"}, guestRole.Permissions)
	
	otherPermission := []types.Permission{"destroy:key"}
	userInfo := &types.UserInfo{
		Roles: []string{"signer", "admin"},
		Permissions: []types.Permission{"destroy:key"},
	}
	signerRole, err := mngr.Role(context.TODO(), "signer")
	adminRole, err := mngr.Role(context.TODO(), "admin")
	permissions := mngr.UserPermissions(context.TODO(), userInfo)
	require.NoError(t, err, "Policy1A should be stored")
	assert.Equal(t, append(append(otherPermission, signerRole.Permissions...), adminRole.Permissions...), permissions)

	err = manifests.Stop(context.TODO())
	require.NoError(t, err, "Stop manifests manager must not error")

	err = mngr.Stop(context.TODO())
	require.NoError(t, err, "Stop manager manager must not error")
}
