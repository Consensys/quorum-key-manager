package nodes

import (
	"context"
	"sync"

	"github.com/consensys/quorum-key-manager/src/aliases"
	"github.com/consensys/quorum-key-manager/src/nodes"
	"github.com/consensys/quorum-key-manager/src/nodes/entities"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	"github.com/consensys/quorum-key-manager/src/stores"

	"github.com/consensys/quorum-key-manager/src/auth"
	"github.com/consensys/quorum-key-manager/src/infra/log"
)

type Nodes struct {
	storesService stores.Stores
	roles         auth.Roles
	aliases       aliases.Aliases
	mux           sync.RWMutex
	nodes         map[string]*entities.Node
	logger        log.Logger
}

var _ nodes.Nodes = &Nodes{}

func New(storesService stores.Stores, rolesService auth.Roles, aliasesService aliases.Aliases, logger log.Logger) *Nodes {
	return &Nodes{
		storesService: storesService,
		roles:         rolesService,
		aliases:       aliasesService,
		mux:           sync.RWMutex{},
		nodes:         make(map[string]*entities.Node),
		logger:        logger,
	}
}

// TODO: Move to data layer
func (i *Nodes) createNode(_ context.Context, name string, prxNode *proxynode.Node, allowedTenants []string) {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.nodes[name] = &entities.Node{
		Name:           name,
		Node:           prxNode,
		AllowedTenants: allowedTenants,
	}
}

// TODO: Move to data layer
func (i *Nodes) getNode(_ context.Context, name string) *entities.Node {
	i.mux.RLock()
	defer i.mux.RUnlock()

	if node, ok := i.nodes[name]; ok {
		return node
	}

	return nil
}
