package nodemanager

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/auth/authorizator"

	authtypes "github.com/consensys/quorum-key-manager/src/auth/types"
	"github.com/consensys/quorum-key-manager/src/infra/log"

	"github.com/consensys/quorum-key-manager/pkg/errors"
	manifest "github.com/consensys/quorum-key-manager/src/manifests/entities"
	manifestsmanager "github.com/consensys/quorum-key-manager/src/manifests/manager"
	"github.com/consensys/quorum-key-manager/src/nodes/interceptor"
	"github.com/consensys/quorum-key-manager/src/nodes/node"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	"github.com/consensys/quorum-key-manager/src/stores"
)

const NodeManagerID = "NodeManager"

type BaseManager struct {
	stores      stores.Manager
	manifests   manifestsmanager.Manager
	authManager auth.Manager

	mux   sync.RWMutex
	nodes map[string]*nodeBundle

	isLive bool
	err    error

	logger log.Logger
}

type nodeBundle struct {
	manifest *manifest.Manifest
	node     node.Node
	err      error
	stop     func(context.Context) error
}

func New(smng stores.Manager, manifests manifestsmanager.Manager, authManager auth.Manager, logger log.Logger) *BaseManager {
	return &BaseManager{
		stores:      smng,
		manifests:   manifests,
		mux:         sync.RWMutex{},
		nodes:       make(map[string]*nodeBundle),
		authManager: authManager,
		logger:      logger,
	}
}

func (m *BaseManager) Start(ctx context.Context) error {
	messages, err := m.manifests.Load()
	if err != nil {
		return err
	}

	for _, message := range messages {
		err = m.createNodes(ctx, message.Manifest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *BaseManager) Stop(ctx context.Context) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.isLive = false

	wg := &sync.WaitGroup{}
	for name, n := range m.nodes {
		wg.Add(1)
		go func(name string, n *nodeBundle) {
			err := n.stop(ctx)
			if err != nil {
				m.logger.WithError(err).Error("node closed with errors", "name", name)
			} else {
				m.logger.Info("node closed successfully", "name", name)
			}
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
	return m.err
}

func (m *BaseManager) Node(_ context.Context, name string, userInfo *authtypes.UserInfo) (node.Node, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if nodeBundle, ok := m.nodes[name]; ok {
		permissions := m.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, m.logger)

		err := resolver.CheckAccess(nodeBundle.manifest.AllowedTenants)
		if err != nil {
			return nil, err
		}

		err = resolver.CheckPermission(&authtypes.Operation{Action: authtypes.ActionProxy, Resource: authtypes.ResourceNode})
		if err != nil {
			return nil, err
		}

		return nodeBundle.node, nodeBundle.err
	}

	return nil, errors.NotFoundError("node not found")
}

func (m *BaseManager) List(_ context.Context, userInfo *authtypes.UserInfo) ([]string, error) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	var nodeNames []string
	for name, nodeBundle := range m.nodes {
		permissions := m.authManager.UserPermissions(userInfo)
		resolver := authorizator.New(permissions, userInfo.Tenant, m.logger)

		if err := resolver.CheckAccess(nodeBundle.manifest.AllowedTenants); err != nil {
			continue
		}
		nodeNames = append(nodeNames, name)
	}

	sort.Strings(nodeNames)

	return nodeNames, nil
}

func (m *BaseManager) createNodes(ctx context.Context, mnf *manifest.Manifest) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	logger := m.logger.With("kind", mnf.Kind, "name", mnf.Name)

	if _, ok := m.nodes[mnf.Name]; ok {
		errMessage := "node already exists"
		logger.Error(errMessage)
		return errors.AlreadyExistsError(errMessage)
	}

	switch mnf.Kind {
	case manifest.Node:
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
		prxNode.Handler = interceptor.New(m.stores.Stores(), m.logger)

		// Start node
		err = prxNode.Start(ctx)
		if err != nil {
			logger.WithError(err).Error("error starting node")
			n.err = err
			return err
		}
		n.node = prxNode
		n.stop = prxNode.Stop

		logger.Info("node created successfully")
	}

	return nil
}

func (m *BaseManager) ID() string { return NodeManagerID }
func (m *BaseManager) CheckLiveness(_ context.Context) error {
	if m.isLive {
		return nil
	}

	errMessage := fmt.Sprintf("service %s is not live", m.ID())
	m.logger.With("id", m.ID()).Error(errMessage)
	return errors.HealthcheckError(errMessage)
}

func (m *BaseManager) CheckReadiness(_ context.Context) error {
	return m.Error()
}
