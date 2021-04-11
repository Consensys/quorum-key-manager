package nodemanager

import (
	"context"
	"fmt"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/core/manifest"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/node"
)

var NodeKind manifest.Kind = "Node"

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows to manage multiple stores
type Manager interface {
	// Load manifest
	// If any error occurs it is attached to the corresponding Message
	Load(ctx context.Context, mnfsts ...*manifest.Manifest) error

	// Node return by name
	Node(ctx context.Context, name string) (node.Node, error)

	// List stores
	List(ctx context.Context) ([]string, error)
}

type manager struct {
	mux   sync.RWMutex
	nodes map[string]*nodeBundle
}

type nodeBundle struct {
	manifest *manifest.Manifest
	node     node.Node
	err      error
}

func New() Manager {
	return &manager{
		mux:   sync.RWMutex{},
		nodes: make(map[string]*nodeBundle),
	}
}

func (m *manager) Load(ctx context.Context, mnfsts ...*manifest.Manifest) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, mnf := range mnfsts {
		if err := m.load(ctx, mnf); err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) Node(_ context.Context, name string) (node.Node, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if nodeBundle, ok := m.nodes[name]; ok {
		return nodeBundle.node, nodeBundle.err
	}

	// This piece of code is here to make sure it is possible to retrieve a default node
	for _, nodeBundle := range m.nodes {
		return nodeBundle.node, nodeBundle.err
	}

	return nil, fmt.Errorf("node not found")
}

func (m *manager) List(_ context.Context) ([]string, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	nodeNames := []string{}
	for name := range m.nodes {
		nodeNames = append(nodeNames, name)
	}

	return nodeNames, nil
}

func (m *manager) load(ctx context.Context, mnf *manifest.Manifest) error {
	logger := log.FromContext(ctx).
		WithField("kind", mnf.Kind).
		WithField("name", mnf.Name)
	logger.Infof("load manifest with specs: %v", string(mnf.Specs))

	if _, ok := m.nodes[mnf.Name]; ok {
		return fmt.Errorf("node %q already exist", mnf.Name)
	}

	switch mnf.Kind {
	case NodeKind:
		n := new(nodeBundle)
		n.manifest = mnf
		m.nodes[mnf.Name] = n

		cfg := new(node.Config)
		if err := mnf.UnmarshalSpecs(cfg); err != nil {
			n.err = err
			return err
		}

		n.node, n.err = node.New(cfg)
	default:
		return fmt.Errorf("invalid manifest kind %s", mnf.Kind)
	}

	return nil
}
