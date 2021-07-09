package nodemanager

import (
	"context"
	"encoding/json"
	"testing"

	testutils2 "github.com/consensys/quorum-key-manager/src/infra/log/testutils"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
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

func TestManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mngr := New(nil, nil, testutils2.NewMockLogger(ctrl))

	err := mngr.load(context.Background(), manifestWithTessera)
	require.NoError(t, err, "Load must not error")

	err = mngr.load(context.Background(), manifestRPCOnly)
	require.NoError(t, err, "Load must not error")

	n, err := mngr.Node(context.Background(), "node-test1")
	require.NoError(t, err, "Node must not error")
	require.NotNil(t, n, "Node must not be nil")

	l, err := mngr.List(context.Background())
	require.NoError(t, err, "List must not error")
	require.Equal(t, []string{"node-test1", "node-test2"}, l, "List must return correct value")

	n, err = mngr.Node(context.Background(), "")
	require.NoError(t, err, "Default node must not error")
	require.NotNil(t, n, "Default node must not be nil")
}
