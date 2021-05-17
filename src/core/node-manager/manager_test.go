package nodemanager

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
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
	mngr := New(nil)

	mnfsts := []*manifest.Manifest{
		manifestWithTessera,
		manifestRPCOnly,
	}

	err := mngr.Load(context.Background(), mnfsts...)
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
