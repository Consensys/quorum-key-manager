package node

import (
	"net/http"
)

//go:generate mockgen -source=node.go -destination=mock/node.go -package=mock

// Node holds interface to connect to a downstream Quorum node including JSON-RPC server and private transaction manager
type Node interface {
	// Proxy returns an JSON-RPC proxy to the downstream JSON-RPC node
	http.Handler
}
