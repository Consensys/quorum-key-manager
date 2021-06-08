package storemanager

import (
	"context"
	"fmt"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var testManifest = []byte(`
- kind: HashicorpSecrets
  version: 0.0.1
  name: hashicorp-secrets
  specs:
    mountPoint: secret
    address: http://hashicorp:8200
    tokenPath: /vault/token/.root
    namespace: ''
- kind: HashicorpKeys
  version: 0.0.1
  name: hashicorp-keys
  specs:
    mountPoint: orchestrate
    address: http://hashicorp:8200
    tokenPath: /vault/token/.root
    namespace: ''
- kind: AKVSecrets
  version: 0.0.1
  name: akv-secrets
  specs:
    vaultName: quorumkeymanager
    tenantID: 17255fb0-373b-4a1a-bd47-d211ab86df81
    clientID: 8c925036-dd6f-4a1e-a315-5e6fab4f2f09
    clientSecret: Cp1BSu50gx-._Q6UJQsSc2oQE-2b.cF.2y
- kind: AKVKeys
  version: 0.0.1
  name: akv-keys
  specs:
    vaultName: quorumkeymanager
    tenantID: 17255fb0-373b-4a1a-bd47-d211ab86df81
    clientID: 8c925036-dd6f-4a1e-a315-5e6fab4f2f09
    clientSecret: Cp1BSu50gx-._Q6UJQsSc2oQE-2b.cF.2y
- kind: Eth1Account
  version: 0.0.1
  name: eth1-accounts
  specs:
    keystore: HashicorpKeys
    specs:
      mountPoint: orchestrate
      address: http://hashicorp:8200
      tokenPath: /vault/token/.root
      namespace: ''
- kind: Node
  name: quorum-node
  version: 0.0.0
  specs:
    rpc:
      addr: http://quorum1:8545
    tessera:
      addr: http://tessera1:9080
- kind: Node
  name: besu-node
  version: 0.0.0
  specs:
    rpc:
      addr: http://validator1:8545
`)

// This test ensure that we do not get any panic on stores manager process
// Still this test can not ensure stores are properly created since we do not have access to dependencies
// (should be responsibility of e2e and ATs)
func TestBaseManager(t *testing.T) {
	dir := t.TempDir()
	err := ioutil.WriteFile(fmt.Sprintf("%v/manifest.yml", dir), testManifest, 0644)
	require.NoError(t, err, "WriteFile manifest1 must not error")

	manifests, err := manager.NewLocalManager(&manager.Config{Path: dir})
	require.NoError(t, err, "NewLocalManager on %v must not error", dir)

	err = manifests.Start(context.TODO())
	require.NoError(t, err, "Start manifests manager must not error")

	mngr := New(manifests)
	err = mngr.Start(context.TODO())
	require.NoError(t, err, "Start manager manager must not error")

	// Give some time to load manifests
	time.Sleep(100 * time.Millisecond)

	err = manifests.Stop(context.TODO())
	require.NoError(t, err, "Stop manifests manager must not error")

	err = mngr.Stop(context.TODO())
	require.NoError(t, err, "Stop manager manager must not error")
}
