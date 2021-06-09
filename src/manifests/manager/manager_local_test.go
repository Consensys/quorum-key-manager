package manager

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	manifest1 = []byte(`
kind: KindA
name: test-1.1
specs:
  field: value
---
kind: KindB
name: test-1.2
specs:
  field: value
`)
	manifest2 = []byte(`
kind: KindB
name: test-2.1
specs:
  field: value
---
kind: KindC
name: test-2.2
specs:
  field: value
`)
)

func assertMessage(t *testing.T, expected []Message, msgs chan []Message) {
	select {
	case msg := <-msgs:
		assert.Equal(t, expected, msg, "Messages should match")
	case <-time.After(20 * time.Millisecond):
		assert.Equal(t, expected, nil, "No message")
	}
}

func TestLocalManager(t *testing.T) {
	dir := t.TempDir()
	err := ioutil.WriteFile(fmt.Sprintf("%v/manifest1.yml", dir), manifest1, 0644)
	require.NoError(t, err, "WriteFile manifest1 must not error")

	err = ioutil.WriteFile(fmt.Sprintf("%v/manifest2.yml", dir), manifest2, 0644)
	require.NoError(t, err, "WriteFile manifest2 must not error")

	mngr, err := NewLocalManager(&Config{Path: dir})
	require.NoError(t, err, "NewLocalManager on %v must not error", dir)

	chanAB := make(chan []Message)
	subAB, err := mngr.Subscribe([]manifest.Kind{"KindA", "KindB"}, chanAB)
	require.NoError(t, err, "Subscribe AB must not error")
	defer func() { _ = subAB.Unsubscribe() }()

	err = mngr.Start(context.TODO())
	require.NoError(t, err, "Start must not error")

	chanBC := make(chan []Message)
	subBC, err := mngr.Subscribe([]manifest.Kind{"KindC", "KindB"}, chanBC)
	require.NoError(t, err, "Subscribe BC must not error")
	defer func() { _ = subBC.Unsubscribe() }()

	chanAll := make(chan []Message)
	subAll, err := mngr.Subscribe(nil, chanAll)
	require.NoError(t, err, "Subscribe All must not error")
	defer func() { _ = subAll.Unsubscribe() }()

	chanNone := make(chan []Message)
	subNone, err := mngr.Subscribe([]manifest.Kind{}, chanNone)
	require.NoError(t, err, "Subscribe None must not error")
	defer func() { _ = subNone.Unsubscribe() }()

	assertMessage(t, []Message{
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindA",
				Name:  "test-1.1",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindB",
				Name:  "test-1.2",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindB",
				Name:  "test-2.1",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
	}, chanAB)

	assertMessage(t, []Message{
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindB",
				Name:  "test-1.2",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindB",
				Name:  "test-2.1",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindC",
				Name:  "test-2.2",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
	}, chanBC)

	assertMessage(t, []Message{
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindA",
				Name:  "test-1.1",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindB",
				Name:  "test-1.2",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindB",
				Name:  "test-2.1",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
		Message{
			Loader: "LocalManager",
			Manifest: &manifest.Manifest{
				Kind:  "KindC",
				Name:  "test-2.2",
				Specs: map[interface{}]interface{}{"field": "value"},
			},
			Action: CreateAction,
		},
	}, chanAll)

	assertMessage(t, nil, chanNone)

	err = mngr.Stop(context.TODO())
	require.NoError(t, err, "Stop must not error")
}
