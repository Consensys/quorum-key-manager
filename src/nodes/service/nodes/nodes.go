package nodes

import (
	"context"
	"sync"

	"github.com/consensys/quorum-key-manager/pkg/errors"
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
	aliases       aliases.Service
	mux           sync.RWMutex
	nodes         map[string]*entities.Node
	logger        log.Logger
}

var _ nodes.Nodes = &Nodes{}

func New(storesService stores.Stores, rolesService auth.Roles, aliasesService aliases.Service, logger log.Logger) *Nodes {
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
func (i *Nodes) getNode(_ context.Context, name string) (*entities.Node, error) {
	i.mux.RLock()
	defer i.mux.RUnlock()

	if node, ok := i.nodes[name]; ok {
		return node, nil
	}

	errMessage := "node was not found"
	i.logger.Error(errMessage, "name", name)
	return nil, errors.NotFoundError(errMessage)
}
