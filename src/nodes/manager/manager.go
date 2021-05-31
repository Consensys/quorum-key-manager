package nodemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	manifestsmanager "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/manager"
	manifest "github.com/ConsenSysQuorum/quorum-key-manager/src/manifests/types"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/nodes/interceptor"
	"github.com/ConsenSysQuorum/quorum-key-manager/src/nodes/node"
	proxynode "github.com/ConsenSysQuorum/quorum-key-manager/src/nodes/node/proxy"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/stores/manager"
)

var NodeKind manifest.Kind = "Node"

//go:generate mockgen -source=manager.go -destination=mock/manager.go -package=mock

// Manager allows to manage multiple stores
type Manager interface {
	// Node return by name
	Node(ctx context.Context, name string) (node.Node, error)

	// List stores
	List(ctx context.Context) ([]string, error)
}

type BaseManager struct {
	stores    storemanager.Manager
	manifests manifestsmanager.Manager

	mux   sync.RWMutex
	nodes map[string]*nodeBundle

	sub    manifestsmanager.Subscription
	mnfsts chan []manifestsmanager.Message
}

type nodeBundle struct {
	manifest *manifest.Manifest
	node     node.Node
	err      error
	stop     func(context.Context) error
}

func New(stores storemanager.Manager, manifests manifestsmanager.Manager) *BaseManager {
	return &BaseManager{
		stores:    stores,
		manifests: manifests,
		mnfsts:    make(chan []manifestsmanager.Message),
		mux:       sync.RWMutex{},
		nodes:     make(map[string]*nodeBundle),
	}
}

func (m *BaseManager) Start(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	// Subscribe to manifest of Kind node
	sub, err := m.manifests.Subscribe([]manifest.Kind{NodeKind}, m.mnfsts)
	if err != nil {
		return err
	}
	m.sub = sub

	// Start loading manifest
	go m.loadAll(ctx)

	return nil
}

func (m *BaseManager) Stop(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	// Unsubscribe
	if m.sub != nil {
		_ = m.sub.Unsubscribe()
	}

	wg := &sync.WaitGroup{}
	for name, n := range m.nodes {
		wg.Add(1)
		go func(name string, n *nodeBundle) {
			err := n.stop(ctx)
			log.FromContext(ctx).WithError(err).WithField("name", name).Errorf("error closing node")
			wg.Done()
		}(name, n)
	}
	wg.Wait()

	return nil
}

func (m *BaseManager) Close() error {
	return nil
}

func (m *BaseManager) Error() error {
	return nil
}

func (m *BaseManager) loadAll(ctx context.Context) {
	for mnfsts := range m.mnfsts {
		for _, mnf := range mnfsts {
			_ = m.load(ctx, mnf.Manifest)
		}
	}
}

func (m *BaseManager) Node(_ context.Context, name string) (node.Node, error) {
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

func (m *BaseManager) List(_ context.Context) ([]string, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	nodeNames := []string{}
	for name := range m.nodes {
		nodeNames = append(nodeNames, name)
	}

	sort.Strings(nodeNames)

	return nodeNames, nil
}

func (m *BaseManager) load(ctx context.Context, mnf *manifest.Manifest) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	logger := log.FromContext(ctx).
		WithField("kind", mnf.Kind).
		WithField("name", mnf.Name)

	if _, ok := m.nodes[mnf.Name]; ok {
		err := fmt.Errorf("node %q already exist", mnf.Name)
		logger.WithError(err).Errorf("error loading node manifest")
		return err
	}

	switch mnf.Kind {
	case NodeKind:
		n := new(nodeBundle)
		n.manifest = mnf
		m.nodes[mnf.Name] = n

		cfg := new(proxynode.Config)
		if err := mnf.UnmarshalSpecs(cfg); err != nil {
			err = fmt.Errorf("invalid node specs: %v", err)
			logger.WithError(err).Errorf("error loading node manifest")
			n.err = err
			return err
		}
		cfg.SetDefault()

		b, _ := json.Marshal(cfg)
		logger.Infof("creating node with config %v", string(b))

		// Create proxy node
		prxNode, err := proxynode.New(cfg)
		if err != nil {
			logger.WithError(err).Errorf("error creating node")
			n.err = err
			return err
		}

		// Set interceptor on proxy node
		prxNode.Handler = interceptor.New(m.stores)

		// Start node
		err = prxNode.Start(ctx)
		if err != nil {
			logger.WithError(err).Errorf("error starting node")
			n.err = err
			return err
		}
		n.node = prxNode
		n.stop = prxNode.Stop
	default:
		err := fmt.Errorf("invalid manifest kind %s", mnf.Kind)
		logger.WithError(err).Errorf("error starting node")
		return err
	}

	return nil
}
