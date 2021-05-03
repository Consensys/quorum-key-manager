package proxynode

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/tessera"
)

//go:generate mockgen -source=session.go -destination=session_mock.go -package=proxynode

// Session holds client interface to a downstream node
type Session interface {
	// ClientRPC returns a client to downstream JSON-RPC
	ClientRPC() jsonrpc.Client

	// EthClient returns a caller to downstream Ethereum JSON-RPC
	EthCaller() ethereum.Caller

	// ClientPrivTxManager returns client to downstrem private transaction manager
	ClientPrivTxManager() tessera.Client
}

type session struct {
	jsonrpcClient    jsonrpc.Client
	ethCaller        ethereum.Caller
	privTxMngrClient tessera.Client
}

func (s *session) ClientRPC() jsonrpc.Client {
	return s.jsonrpcClient
}

func (s *session) EthCaller() ethereum.Caller {
	return s.ethCaller
}

func (s *session) ClientPrivTxManager() tessera.Client {
	return s.privTxMngrClient
}
