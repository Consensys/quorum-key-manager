package nodes

import (
	"context"
	"github.com/consensys/quorum-key-manager/pkg/errors"
	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/nodes"
	"github.com/consensys/quorum-key-manager/src/nodes/entities"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	"github.com/consensys/quorum-key-manager/src/stores"
	"sync"

	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type Interactor struct {
	storesService stores.Stores
	roles         auth.Roles
	aliasParser   aliases.Parser
	mux           sync.RWMutex
	nodes         map[string]*entities.Node
	logger        log.Logger
}

var _ nodes.Nodes = &Interactor{}

func New(storesService stores.Stores, rolesService auth.Roles, aliasParser aliases.Parser, logger log.Logger) *Interactor {
	return &Interactor{
		storesService: storesService,
		roles:         rolesService,
		aliasParser:   aliasParser,
		mux:           sync.RWMutex{},
		nodes:         make(map[string]*entities.Node),
		logger:        logger,
	}
}

// TODO: Move to data layer
func (i *Interactor) createNode(_ context.Context, name string, prxNode *proxynode.Node, allowedTenants []string) {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.nodes[name] = &entities.Node{
		Name:           name,
		Node:           prxNode,
		AllowedTenants: allowedTenants,
	}
}

// TODO: Move to data layer
func (i *Interactor) getNode(_ context.Context, name string) (*entities.Node, error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	if node, ok := i.nodes[name]; ok {
		return node, nil
	}

	errMessage := "node was not found"
	i.logger.Error(errMessage, "name", name)
	return nil, errors.NotFoundError(errMessage)
}
