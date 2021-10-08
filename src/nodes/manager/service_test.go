package nodemanager

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/consensys/quorum-key-manager/src/auth/mock"
	storesmock "github.com/consensys/quorum-key-manager/src/stores/mock"

	aliasmock "github.com/consensys/quorum-key-manager/src/aliases/mock"
	"github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthManager := mock.NewMockManager(ctrl)
	mockStoresManager := storesmock.NewMockManager(ctrl)
	mockStores := storesmock.NewMockStores(ctrl)

	mockAuthManager.EXPECT().UserPermissions(gomock.Any()).Return(types.ListPermissions()).AnyTimes()
	mockStoresManager.EXPECT().Stores().Return(mockStores).AnyTimes()

	mockAliasManager := aliasmock.NewMockService(ctrl)
	mngr := New(mockStoresManager, nil, mockAuthManager, mockAliasManager, testutils.NewMockLogger(ctrl))

	err := mngr.load(context.Background(), manifestWithTessera)
	require.NoError(t, err, "Load must not error")

	err = mngr.load(context.Background(), manifestRPCOnly)
	require.NoError(t, err, "Load must not error")

	err = mngr.load(context.Background(), manifestWithTenant)
	require.NoError(t, err, "Load must not error")

	n, err := mngr.Node(context.Background(), "node-test1", &types.UserInfo{})
	require.NoError(t, err, "Node must not error")
	require.NotNil(t, n, "Node must not be nil")

	l, err := mngr.List(context.Background(), &types.UserInfo{})
	require.NoError(t, err, "List must not error")
	require.Equal(t, []string{"node-test1", "node-test2"}, l, "List must return correct value")

	l, err = mngr.List(context.Background(), &types.UserInfo{Tenant: "tenantOne"})
	require.NoError(t, err, "List must not error")
	require.Contains(t, l, "node-test3")
}
