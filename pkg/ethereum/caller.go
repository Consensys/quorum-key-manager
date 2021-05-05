package ethereum

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
)

//go:generate mockgen -source=caller.go -destination=mock/caller.go -package=mock

type Caller interface {
	Eth() EthCaller
	EEA() EEACaller
	Priv() PrivCaller
}

// Caller implement methods to interface with the JSON-RPC API of an ethereum caller
// It is some kind of Web3 interface
type caller struct {
	eth  *ethCaller
	eea  *eeaCaller
	priv *privCaller
}

// NewCaller creates a caller from a jsonrpc.Client
func NewCaller(client jsonrpc.Client) Caller {
	return &caller{
		eth:  &ethCaller{client},
		eea:  &eeaCaller{client},
		priv: &privCaller{client},
	}
}

// Eth return eth namespace caller
func (c *caller) Eth() EthCaller { // nolint
	return c.eth
}

// EEA return eea namespace caller
func (c *caller) EEA() EEACaller { // nolint
	return c.eea
}

// Priv return priv namespace caller
func (c *caller) Priv() PrivCaller { // nolint
	return c.priv
}
