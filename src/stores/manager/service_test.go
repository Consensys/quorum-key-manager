package storemanager

import (
	"context"
	"fmt"
	"github.com/consensys/quorum-key-manager/src/infra/manifests/reader"
	"io/ioutil"
	"testing"
	"time"

	stores2 "github.com/consensys/quorum-key-manager/src/stores/connectors/stores"

	mock2 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/consensys/quorum-key-manager/src/stores/database/mock"

	"github.com/golang/mock/gomock"

	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/reader"
	"github.com/stretchr/testify/require"
)

var testManifest = []byte(`
- kind: HashicorpSecrets
  version: 0.0.1
  name: hashicorp-secrets
  specs:
    mountPoint: secret
    address: http://hashicorp:8200
    token: fakeToken
    namespace: ''
- kind: HashicorpKeys
  version: 0.0.1
  name: hashicorp-keys
  specs:
    mountPoint: quorum
    address: http://hashicorp:8200
    token: fakeToken
    namespace: ''
- kind: AKVSecrets
  version: 0.0.1
  name: akv-secrets
  allowedTenants: ['tenantOne']
  specs:
    vaultName: fakeVaultName
    tenantID: fakeTenant
    clientID: fakeClientID
    clientSecret: fakeSecret
- kind: AKVKeys
  version: 0.0.1
  name: akv-keys
  allowedTenants: ['tenantOne', 'tenantTwo']
  specs:
    vaultName: quorumkeymanager
    tenantID: fakeTenant
    clientID: fakeClientID
    clientSecret: fakeSecret
- kind: Ethereum
  version: 0.0.1
  name: eth-accounts
  specs:
    keystore: HashicorpKeys
    specs:
      mountPoint: quorum
      address: http://hashicorp:8200
      token: fakeToken
      namespace: ''
`)

// This test ensures that we do not get any panic on stores manager process
// Still this test can not ensure stores are properly created since we do not have access to dependencies
// (should be responsibility of e2e and ATs)
func TestManagerService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := testutils.NewMockLogger(ctrl)
	mockDB := mock.NewMockDatabase(ctrl)
	mockSecretDB := mock.NewMockSecrets(ctrl)
	mockKeysDB := mock.NewMockKeys(ctrl)
	mockEthDB := mock.NewMockETHAccounts(ctrl)
	mockAuthMngr := mock2.NewMockManager(ctrl)

	mockAuthMngr.EXPECT().UserPermissions(gomock.Any()).Return(types.ListPermissions()).AnyTimes()
	mockDB.EXPECT().Secrets(gomock.Any()).Return(mockSecretDB).AnyTimes()
	mockDB.EXPECT().Keys(gomock.Any()).Return(mockKeysDB).AnyTimes()
	mockDB.EXPECT().ETHAccounts(gomock.Any()).Return(mockEthDB).AnyTimes()

	dir := t.TempDir()
	err := ioutil.WriteFile(fmt.Sprintf("%v/manifest.yml", dir), testManifest, 0644)
	require.NoError(t, err, "WriteFile manifest1 must not error")

	manifests, err := reader.New(&manifestsmanager.Config{Path: dir}, mockLogger)
	require.NoError(t, err, "New on %v must not error", dir)

	err = manifests.Start(context.TODO())
	require.NoError(t, err, "Start manifests manager must not error")

	storesConnector := stores2.NewConnector(mockAuthMngr, mockDB, mockLogger)
	mngr := New(storesConnector, manifests, mockDB, mockLogger)
	err = mngr.Start(context.TODO())
	require.NoError(t, err, "Start manager manager must not error")

	// Give some time to load manifests
	time.Sleep(500 * time.Millisecond)

	stores, err := mngr.Stores().List(context.TODO(), "", &types.UserInfo{})
	require.NoError(t, err)
	assert.Contains(t, stores, "hashicorp-secrets")
	assert.Contains(t, stores, "hashicorp-keys")
	assert.Contains(t, stores, "eth-accounts")

	stores, err = mngr.Stores().List(context.TODO(), "", &types.UserInfo{Tenant: "tenantOne"})
	require.NoError(t, err)
	assert.Contains(t, stores, "akv-secrets")
	assert.Contains(t, stores, "akv-keys")

	err = manifests.Stop(context.TODO())
	require.NoError(t, err, "Stop manifests manager must not error")

	err = mngr.Stop(context.TODO())
	require.NoError(t, err, "Stop manager manager must not error")
}
