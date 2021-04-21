package node

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	httpclient "github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/proxy"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
)

//go:generate mockgen -source=node.go -destination=mock/node.go -package=mock

// Node holds interface to connect to a downstream Quorum node including JSON-RPC server and private transaction manager
type Node interface {
	// ClientRPC returns client to downstream JSON-RPC
	ClientRPC() jsonrpc.Client

	// Proxy returns an JSON-RPC proxy to the downstream JSON-RPC node
	ProxyRPC() jsonrpc.Handler

	// ClientPrivTxManager returns client to downstrem private transaction manager
	ClientPrivTxManager() httpclient.Client

	// Proxy returns an HTTP proxy to the downstream Priv Tx Manager node
	ProxyPrivTxManager() http.Handler

	// Returns a session contextualized to a request
	Session(*jsonrpc.Request) (Session, error)
}

// Session holds client interface to a downstream node
// A session holds contextual data and is meant to be re-used across calls to the downstream JSON-RPC
// Typically when proxying a request which may end up in multiple calls to the downstream node
type Session interface {
	// Node returns the node that generated this session
	Node

	// CallerRPC returns a caller to downstream JSON-RPC
	CallerRPC() jsonrpc.Caller

	// EthClient returns a client to downstream JSON-RPC
	EthClient() *ethereum.Client

	// Close the session
	Close()
}

// New creates a Node
func New(cfg *Config) (Node, error) {
	cfg = cfg.Copy().SetDefault()

	n := new(node)
	var err error
	n.rpc, err = newRPCDownstream(cfg.RPC)
	if err != nil {
		return nil, err
	}

	if cfg.PrivTxManager != nil {
		n.privTMngr, err = newPrivTxMngrDownstream(cfg.PrivTxManager)
		if err != nil {
			return nil, err
		}
	}

	return n, nil
}

type node struct {
	rpc       *rpcDownstream
	privTMngr *privTxMngrDownstream
}

func (n *node) ClientRPC() jsonrpc.Client {
	return n.rpc.client
}

func (n *node) ProxyRPC() jsonrpc.Handler {
	return n.rpc.proxy
}

func (n *node) ClientPrivTxManager() httpclient.Client {
	return n.privTMngr.http.client
}

func (n *node) ProxyPrivTxManager() http.Handler {
	return n.privTMngr.http.proxy
}

// Session returns a new session
func (n *node) Session(req *jsonrpc.Request) (Session, error) {
	client := jsonrpc.WithIncrementalID(req.ID())(n.rpc.client)
	client = jsonrpc.WithVersion(req.Version())(client)

	return &session{
		node:      n,
		rpcCaller: jsonrpc.NewCaller(client, req),
	}, nil
}

type session struct {
	*node
	rpcCaller jsonrpc.Caller
	ethClient *ethereum.Client
}

func (s *session) CallerRPC() jsonrpc.Caller {
	return s.rpcCaller
}

func (s *session) EthClient() *ethereum.Client {
	return s.ethClient
}

func (s *session) Close() {
	// TODO: to be implemented in particular maybe it makes sense to recycle sessions using a sync.Pool
}

type httpDownstream struct {
	transport    http.RoundTripper
	reqPreparer  request.Preparer
	respModifier response.Modifier
	client       httpclient.Client
	proxy        http.Handler
}

func newhttpDownstream(cfg *DownstreamConfig) (*httpDownstream, error) {
	n := new(httpDownstream)
	var err error
	n.transport, err = transport.New(cfg.Transport)
	if err != nil {
		return nil, err
	}

	n.reqPreparer, err = request.Proxy(cfg.Proxy.Request)
	if err != nil {
		return nil, err
	}

	n.respModifier = response.Proxy(cfg.Proxy.Response)

	n.client, err = httpclient.New(&httpclient.Config{Timeout: cfg.ClientTimeout}, n.transport, n.reqPreparer, n.respModifier)
	if err != nil {
		return nil, err
	}

	n.proxy, err = proxy.New(nil, n.transport, n.reqPreparer, n.respModifier, nil, nil)
	if err != nil {
		return nil, err
	}

	return n, nil
}

type rpcDownstream struct {
	http   *httpDownstream
	proxy  jsonrpc.Handler
	client jsonrpc.Client
}

func newRPCDownstream(cfg *DownstreamConfig) (*rpcDownstream, error) {
	n := new(rpcDownstream)
	var err error
	n.http, err = newhttpDownstream(cfg)
	if err != nil {
		return nil, err
	}

	n.client = jsonrpc.NewClient(n.http.client)

	// Overide HTTP Proy
	n.http.proxy, err = proxy.New(nil, n.http.transport, n.http.reqPreparer, n.http.respModifier, jsonrpc.HandleProxyRoundTripError, nil)
	if err != nil {
		return nil, err
	}

	n.proxy = jsonrpc.FromHTTPHandler(n.http.proxy)

	return n, nil
}

type privTxMngrDownstream struct {
	http *httpDownstream
}

func newPrivTxMngrDownstream(cfg *DownstreamConfig) (*privTxMngrDownstream, error) {
	n := new(privTxMngrDownstream)
	var err error
	n.http, err = newhttpDownstream(cfg)
	if err != nil {
		return nil, err
	}

	return n, nil
}
