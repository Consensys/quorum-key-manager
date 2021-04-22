package interceptor

import (
	nodemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/node-manager"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/core/store-manager"
)

type Interceptor struct {
	stores storemanager.Manager
	nodes  nodemanager.Manager

	nodeSession *NodeSessionMiddleware
}

func New(stores storemanager.Manager, nodes nodemanager.Manager) *Interceptor {
	return &Interceptor{
		stores:      stores,
		nodes:       nodes,
		nodeSession: NewNodeSessionMiddleware(nodes),
	}
}
