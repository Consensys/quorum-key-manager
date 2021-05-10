package proxynode

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/ethereum"
	httpclient "github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/proxy"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/tessera"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/websocket"
	gorillamux "github.com/gorilla/mux"
	gorillawebsocket "github.com/gorilla/websocket"
)

// Node is a node connected to JSON-RPC downstream and a Tessera downstream
// allowing to intercept proxy JSON-RPC
type Node struct {
	// Handler is the JSON-RPC handler
	Handler jsonrpc.Handler

	rpc        *httpDownstream
	privTxMngr *httpDownstream

	wsHandler   *websocket.Proxy
	httpHandler http.Handler
}

// New creates a Node
func New(cfg *Config) (*Node, error) {
	n := new(Node)
	var err error
	n.rpc, err = newhttpDownstream(cfg.RPC)
	if err != nil {
		return nil, err
	}

	if cfg.PrivTxManager != nil {
		n.privTxMngr, err = newhttpDownstream(cfg.PrivTxManager)
		if err != nil {
			return nil, err
		}
	}

	// Set HTTP proxy
	router := gorillamux.NewRouter()
	router.Methods(http.MethodPost).HandlerFunc(n.serveHTTP)
	n.httpHandler = router

	// Set websocket proxy
	websocketProxy := websocket.NewProxy(cfg.RPC.Proxy.WebSocket)
	websocketProxy.ReqPreparer = n.rpc.reqPreparer
	websocketProxy.RespModifier = n.rpc.respModifier
	websocketProxy.Interceptor = n.interceptWS
	websocketProxy.ErrorHandler = n.rpc.errorHandler
	n.wsHandler = websocketProxy

	return n, nil
}

func (n *Node) Start(ctx context.Context) error {
	return n.wsHandler.Start(ctx)
}

func (n *Node) Stop(ctx context.Context) error {
	return n.wsHandler.Stop(ctx)
}

func (n *Node) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if gorillawebsocket.IsWebSocketUpgrade(req) {
		// we serve websocket
		n.wsHandler.ServeHTTP(rw, req)
		return
	}

	n.httpHandler.ServeHTTP(rw, req)
}

func (n *Node) serveHTTP(rw http.ResponseWriter, req *http.Request) {
	// Create ResponseWriter
	rpcRw := jsonrpc.NewResponseWriter(rw)

	// Parse request body
	msg := new(jsonrpc.RequestMsg)
	err := json.NewDecoder(req.Body).Decode(msg)
	req.Body.Close()
	if err != nil {
		_ = jsonrpc.WriteError(rpcRw, jsonrpc.ParseError(err))
		return
	}

	// Attach session to context
	jsonrpcClient := n.newHTTPJSONRPCClient(req)
	ctx := WithSession(req.Context(), n.newSession(jsonrpcClient, msg))

	// Serve
	n.handler().ServeRPC(rpcRw, msg.WithContext(ctx))
}

func (n *Node) interceptWS(ctx context.Context, clientConn, serverConn *gorillawebsocket.Conn) (clientErrors, serverErrors <-chan error) {
	// Create a JSON-RPC client attached to the downstream server connection
	jsonrpcClient := jsonrpc.NewWebsocketClient(serverConn)
	_ = jsonrpcClient.Start(ctx)

	clientErrs := make(chan error, 1)

	// Start main loop treating client messages
	go func() {
		defer func() { _ = jsonrpcClient.Stop(ctx) }()
		for {
			// Read client message
			typ, b, err := clientConn.ReadMessage()
			if err != nil {
				clientErrs <- err
				close(clientErrs)
				return
			}

			// Create ResponseWriter
			w, err := clientConn.NextWriter(typ)
			if err != nil {
				continue
			}
			rpcRw := jsonrpc.NewResponseWriter(w)

			// Unmarshal input message
			msg := new(jsonrpc.RequestMsg)
			err = json.Unmarshal(b, msg)
			if err != nil {
				_ = jsonrpc.WriteError(rpcRw, jsonrpc.ParseError(err))
				continue
			}

			// Create and attach session to context then handle message
			sess := n.newSession(jsonrpcClient, msg)
			n.handler().ServeRPC(rpcRw, msg.WithContext(WithSession(ctx, sess)))

			// Close writer so message is sent through connection
			w.Close()
		}
	}()

	return clientErrs, jsonrpcClient.Errors()
}

func (n *Node) handler() jsonrpc.Handler {
	if n.Handler != nil {
		return n.Handler
	}
	return ProxyHandler
}

func (n *Node) newSession(jsonrpcClient jsonrpc.Client, msg *jsonrpc.RequestMsg) *session {
	return &session{
		jsonrpcClient:    jsonrpcClient,
		ethCaller:        newEthCaller(jsonrpcClient, msg),
		privTxMngrClient: n.newPrivTxMngrClient(),
	}
}

func (n *Node) newHTTPJSONRPCClient(req *http.Request) jsonrpc.Client {
	httpClient := httpclient.CombineDecorators(
		httpclient.WithModifier(n.rpc.respModifier),
		httpclient.WithRequest(req),
		httpclient.WithPreparer(n.rpc.reqPreparer),
		httpclient.WithPreparer(
			request.CombinePreparer(
				request.RemoveConnectionHeaders(),
				request.ForwardedFor(),
			),
		),
	)(n.rpc.client)
	return jsonrpc.NewHTTPClient(httpClient)
}

func newEthCaller(jsonrpcClient jsonrpc.Client, msg *jsonrpc.RequestMsg) ethereum.Caller {
	jsonrpcClient = jsonrpc.WithVersion(msg.Version)(jsonrpcClient)
	jsonrpcClient = jsonrpc.WithIncrementalID(msg.ID)(jsonrpcClient)
	return ethereum.NewCaller(jsonrpcClient)
}

func (n *Node) newPrivTxMngrClient() tessera.Client {
	if n.privTxMngr != nil {
		return tessera.NewHTTPClient(n.privTxMngr.client)
	}

	return &tessera.NotConfiguredClient{}
}

type httpDownstream struct {
	transport    http.RoundTripper
	reqPreparer  request.Preparer
	respModifier response.Modifier
	client       httpclient.Client

	errorHandler proxy.HandleRoundTripErrorFunc
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

	n.errorHandler = proxy.HandleRoundTripError

	n.client, err = httpclient.New(&httpclient.Config{Timeout: cfg.ClientTimeout}, n.transport)
	if err != nil {
		return nil, err
	}

	return n, nil
}

var ProxyHandler = jsonrpc.DefaultRWHandler(
	jsonrpc.HandlerFunc(func(rw jsonrpc.ResponseWriter, msg *jsonrpc.RequestMsg) {
		// Sen RPC request
		resp, err := SessionFromContext(msg.Context()).ClientRPC().Do(msg)
		if err != nil {
			_ = jsonrpc.WriteError(rw, err)
		} else {
			_ = rw.WriteMsg(resp)
		}
	}),
)
