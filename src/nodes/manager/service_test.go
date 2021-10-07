package nodemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifest "github.com/consensys/quorum-key-manager/src/infra/manifests/entities"
	manifestmock "github.com/consensys/quorum-key-manager/src/infra/manifests/mock"

	"github.com/consensys/quorum-key-manager/src/auth/mock"
	storesmock "github.com/consensys/quorum-key-manager/src/stores/mock"

	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"
)

var manifestWithTessera = &manifest.Manifest{
	Kind:    "Node",
	Version: "v1alpha",
	Name:    "node-test1",
	Tags: map[string]string{
		"key1": "value1",
		"key2": "value2",
	},
	Specs: json.RawMessage(`
{
	"rpc": {
		"addr": "www.test-rpc.com",
		"transport": {
			"idleConnTimeout": "15s"
		},
		"proxy": {
			"request": {
				"headers": {
					"CUSTOM-HEADER": ["test"]
				}
			}
		}
	},
	"tessera": {
		"addr": "www.test-tessera.com",
		"transport": {
			"dialer": {
				"timeout": "30s"
			}
		}
	}
}`),
}

var manifestRPCOnly = &manifest.Manifest{
	Kind:    "Node",
	Version: "v1alpha",
	Name:    "node-test2",
	Tags: map[string]string{
		"key1": "value1",
		"key2": "value2",
	},
	Specs: json.RawMessage(`
{
	"rpc": {
		"addr": "www.test-rpc.com",
		"transport": {
			"idleConnTimeout": "15s"
		},
		"request": {
			"headers": {
				"CUSTOMER-HEADER": ["test"]
			}
		}
	}
}
`),
}

var manifestWithTenant = &manifest.Manifest{
	Kind:           "Node",
	Version:        "v2alpha",
	Name:           "node-test3",
	AllowedTenants: []string{"tenantOne"},
	Specs: json.RawMessage(`
{
	"rpc": {
		"addr": "www.test-rpc.com",
		"transport": {
			"idleConnTimeout": "15s"
		},
		"request": {
			"headers": {
				"CUSTOMER-HEADER": ["test"]
			}
		}
	}
}
`),
}

func TestManager(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthManager := mock.NewMockManager(ctrl)
	mockStoresManager := storesmock.NewMockManager(ctrl)
	mockStores := storesmock.NewMockStores(ctrl)
	mockManifestReader := manifestmock.NewMockReader(ctrl)

	mngr := New(mockStoresManager, mockManifestReader, mockAuthManager, testutils.NewMockLogger(ctrl))

	t.Run("should start service successfully loading nodes from mnf", func(t *testing.T) {
		mockManifestReader.EXPECT().Load().Return([]*manifest.Manifest{manifestWithTessera, manifestRPCOnly, manifestWithTenant}, nil)
		mockAuthManager.EXPECT().UserPermissions(gomock.Any()).Return(types.ListPermissions()).AnyTimes()
		mockStoresManager.EXPECT().Stores().Return(mockStores).AnyTimes()

		err := mngr.Start(ctx)
		require.NoError(t, err)

		n, err := mngr.Node(ctx, "node-test1", &types.UserInfo{})
		require.NoError(t, err)
		require.NotNil(t, n)

		l, err := mngr.List(ctx, &types.UserInfo{})
		require.NoError(t, err)
		require.Equal(t, []string{"node-test1", "node-test2"}, l)

		l, err = mngr.List(ctx, &types.UserInfo{Tenant: "tenantOne"})
		require.NoError(t, err)
		require.Contains(t, l, "node-test3")
	})

	t.Run("should fail with ConfigError if mnf fails to load", func(t *testing.T) {
		mockManifestReader.EXPECT().Load().Return(nil, fmt.Errorf("error"))

		err := mngr.Start(ctx)
		require.True(t, errors.IsConfigError(err))
	})
}
