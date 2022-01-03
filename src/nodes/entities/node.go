package entities

import (
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
)

type Node struct {
	Name           string
	Node           *proxynode.Node
	AllowedTenants []string
}
