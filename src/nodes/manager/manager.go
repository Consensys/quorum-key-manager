package nodemanager

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/types"
	"github.com/consensys/quorum-key-manager/src/nodes/interceptor"
	"github.com/consensys/quorum-key-manager/src/nodes/node"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	storemanager "github.com/consensys/quorum-key-manager/src/stores/manager"
)

const NodeManagerID = "NodeManager"

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

	isLive bool

	logger log.Logger
}

type nodeBundle struct {
	manifest *manifest.Manifest
	node     node.Node
	err      error
	stop     func(context.Context) error
}

func New(stores storemanager.Manager, manifests manifestsmanager.Manager, logger log.Logger) *BaseManager {
	return &BaseManager{
		stores:    stores,
		manifests: manifests,
		mnfsts:    make(chan []manifestsmanager.Message),
		mux:       sync.RWMutex{},
		nodes:     make(map[string]*nodeBundle),
		logger:    logger,
	}
}

func (m *BaseManager) Start(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()
	defer func() {
		m.isLive = true
	}()

	// Subscribe to manifest of Kind node
	m.sub = m.manifests.Subscribe([]manifest.Kind{NodeKind}, m.mnfsts)

	// Start loading manifest
	go m.loadAll(ctx)

	return nil
}

func (m *BaseManager) Stop(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.isLive = false
	// Unsubscribe
	if m.sub != nil {
		_ = m.sub.Unsubscribe()
	}

	wg := &sync.WaitGroup{}
	for name, n := range m.nodes {
		wg.Add(1)
		go func(name string, n *nodeBundle) {
			err := n.stop(ctx)
			m.logger.WithError(err).Error("error closing node", "name", name)
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

	return nil, errors.NotFoundError("node not found")
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

	logger := m.logger.With("kind", mnf.Kind, "name", mnf.Name)

	if _, ok := m.nodes[mnf.Name]; ok {
		errMessage := "node already exists"
		logger.Error(errMessage)
		return errors.AlreadyExistsError(errMessage)
	}

	switch mnf.Kind {
	case NodeKind:
		n := new(nodeBundle)
		n.manifest = mnf
		m.nodes[mnf.Name] = n

		cfg := new(proxynode.Config)
		if err := mnf.UnmarshalSpecs(cfg); err != nil {
			errMessage := "invalid node specs"
			logger.WithError(err).Error(errMessage)
			n.err = errors.InvalidParameterError(errMessage)
			return n.err
		}
		cfg.SetDefault()

		// Create proxy node
		prxNode, err := proxynode.New(cfg, m.logger)
		if err != nil {
			logger.WithError(err).Error("failed to create node")
			n.err = err
			return err
		}

		// Set interceptor on proxy node
		prxNode.Handler = interceptor.New(m.stores, m.logger)

		// Start node
		err = prxNode.Start(ctx)
		if err != nil {
			logger.WithError(err).Error("error starting node")
			n.err = err
			return err
		}
		n.node = prxNode
		n.stop = prxNode.Stop
	default:
		errMessage := "invalid manifest kind"
		logger.Error(errMessage)
		return errors.InvalidParameterError(errMessage)
	}

	logger.Info("node loaded successfully")
	return nil
}

func (m *BaseManager) ID() string { return NodeManagerID }
func (m *BaseManager) CheckLiveness() error {
	if m.isLive {
		return nil
	}

	errMessage := fmt.Sprintf("service %s is not live", m.ID())
	m.logger.With("id", m.ID()).Error(errMessage)
	return errors.HealthcheckError(errMessage)
}

func (m *BaseManager) CheckReadiness() error {
	return m.Error()
}
