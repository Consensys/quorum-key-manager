package storemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	mock3 "github.com/consensys/quorum-key-manager/src/infra/manifests/mock"
	stores2 "github.com/consensys/quorum-key-manager/src/stores/connectors/stores"

	mock2 "github.com/consensys/quorum-key-manager/src/auth/mock"
	"github.com/consensys/quorum-key-manager/src/auth/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"
	"github.com/stretchr/testify/assert"

	"github.com/consensys/quorum-key-manager/src/stores/database/mock"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
)

func TestManagerService(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := testutils.NewMockLogger(ctrl)
	mockDB := mock.NewMockDatabase(ctrl)
	mockSecretDB := mock.NewMockSecrets(ctrl)
	mockKeysDB := mock.NewMockKeys(ctrl)
	mockEthDB := mock.NewMockETHAccounts(ctrl)
	mockAuthMngr := mock2.NewMockManager(ctrl)
	mockManifestReader := mock3.NewMockReader(ctrl)

	mockAuthMngr.EXPECT().UserPermissions(gomock.Any()).Return(entities.ListPermissions()).AnyTimes()
	mockDB.EXPECT().Secrets(gomock.Any()).Return(mockSecretDB).AnyTimes()
	mockDB.EXPECT().Keys(gomock.Any()).Return(mockKeysDB).AnyTimes()
	mockDB.EXPECT().ETHAccounts(gomock.Any()).Return(mockEthDB).AnyTimes()

	storesConnector := stores2.NewConnector(mockAuthMngr, mockDB, mockLogger)
	mngr := New(storesConnector, mockManifestReader, mockDB, mockLogger)

	t.Run("should start stores service successfully by loading stores", func(t *testing.T) {
		testManifests := []*manifest.Manifest{
			{
				Kind:  "HashicorpSecrets",
				Name:  "hashicorp-secrets",
				Specs: json.RawMessage(`{"mountPoint": "secret", "address":"http://hashicorp:8200", "token": "fakeToken", "namespace": ""}`),
			},
			{
				Kind:  "HashicorpKeys",
				Name:  "hashicorp-keys",
				Specs: json.RawMessage(`{"mountPoint": "quorum", "address":"http://hashicorp:8200", "token": "fakeToken", "namespace": ""}`),
			},
			{
				Kind:           "AKVSecrets",
				Name:           "akv-secrets",
				AllowedTenants: []string{"tenantOne"},
				Specs:          json.RawMessage(`{"vaultName": "fakeVaultName", "tenantID":"fakeTenant", "clientID": "fakeClientID", "clientSecret": "fakeSecret"}`),
			},
			{
				Kind:           "AKVKeys",
				Name:           "akv-keys",
				AllowedTenants: []string{"tenantOne", "tenantTwo"},
				Specs:          json.RawMessage(`{"vaultName": "quorumkeymanager", "tenantID":"fakeTenant", "clientID": "fakeClientID", "clientSecret": "fakeSecret"}`),
			},
		}

		mockManifestReader.EXPECT().Load().Return(testManifests, nil)

		err := mngr.Start(ctx)
		require.NoError(t, err)

		stores, err := mngr.Stores().List(context.TODO(), "", &entities.UserInfo{})
		require.NoError(t, err)
		assert.Contains(t, stores, "hashicorp-secrets")
		assert.Contains(t, stores, "hashicorp-keys")

		stores, err = mngr.Stores().List(context.TODO(), "", &entities.UserInfo{Tenant: "tenantOne"})
		require.NoError(t, err)
		assert.Contains(t, stores, "akv-secrets")
		assert.Contains(t, stores, "akv-keys")

	})

	t.Run("should fail with ConfigError if manifest fails to be loaded", func(t *testing.T) {
		mockManifestReader.EXPECT().Load().Return(nil, fmt.Errorf("error"))

		err := mngr.Start(ctx)
		assert.True(t, errors.IsConfigError(err))
	})
}
