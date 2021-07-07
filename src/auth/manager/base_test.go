package manager

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/consensys/quorum-key-manager/pkg/log/testutils"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testManifest = []byte(`
- kind: Group
  version: 0.0.1
  name: test-group1
  specs:
    policies:
      - test-policy1A
      - test-policy1B
- kind: Group
  version: 0.0.1
  name: test-group2
  specs:
    policies:
      - test-policy2A
- kind: Policy
  version: 0.0.1
  name: test-policy1A
  specs:
    statements:
      - name: DoAction1
        effect: Allow
        actions:
          - Action1
        resource:
          - /path/to/resource
      - name: DoAction23
        effect: Allow
        actions:
          - Action3
          - Action2
        resource:
          - /path/to/resource
- kind: Policy
  version: 0.0.1
  name: test-policy2A
  specs:
    statements:
      - name: DoAction3
        effect: Allow
        actions:
          - Action3
        resource:
          - /path/to/resource
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
	group1, err := mngr.Group(context.TODO(), "test-group1")
	require.NoError(t, err, "Group1 should be stored")
	assert.Equal(t, "test-group1", group1.Name, "Group1 should have correct name")
	assert.Equal(t, []string{"test-policy1A", "test-policy1B"}, group1.Policies, "Group1 should have correct policies")

	group2, err := mngr.Group(context.TODO(), "test-group2")
	require.NoError(t, err, "Group2 should be stored")
	assert.Equal(t, "test-group2", group2.Name, "Group2 should have correct name")
	assert.Equal(t, []string{"test-policy2A"}, group2.Policies, "Group2 should have correct policies")

	policy1A, err := mngr.Policy(context.TODO(), "test-policy1A")
	require.NoError(t, err, "Policy1A should be stored")
	assert.Equal(t, "test-policy1A", policy1A.Name, "Policy1A should have correct name")
	assert.Len(t, policy1A.Statements, 2, "Policy1A should have correct statements")

	policy2A, err := mngr.Policy(context.TODO(), "test-policy2A")
	require.NoError(t, err, "policy2A should be stored")
	assert.Equal(t, "test-policy2A", policy2A.Name, "Policy2A should have correct name")
	assert.Len(t, policy2A.Statements, 1, "Policy2A should have correct statements")

	err = manifests.Stop(context.TODO())
	require.NoError(t, err, "Stop manifests manager must not error")

	err = mngr.Stop(context.TODO())
	require.NoError(t, err, "Stop manager manager must not error")
}
