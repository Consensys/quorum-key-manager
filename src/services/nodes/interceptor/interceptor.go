package interceptor

import (
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/jsonrpc"
	storemanager "github.com/ConsenSysQuorum/quorum-key-manager/src/services/stores/manager"
	proxynode "github.com/ConsenSysQuorum/quorum-key-manager/src/services/nodes/node/proxy"
)

type Interceptor struct {
	stores storemanager.StoreManager

	handler jsonrpc.Handler
}

func (i *Interceptor) ServeRPC(rw jsonrpc.ResponseWriter, msg *jsonrpc.RequestMsg) {
	i.handler.ServeRPC(rw, msg)
}

func (i *Interceptor) newHandler() jsonrpc.Handler {
	// Only JSON-RPC v2 is supported
	router := jsonrpc.NewRouter().DefaultHandler(jsonrpc.NotSupportedVersionHandler())
	v2Router := router.Version("2.0").Subrouter().DefaultHandler(proxynode.ProxyHandler)

	// Set JSON-RPC interceptors
	v2Router.Method("eth_accounts").Handle(i.EthAccounts())
	v2Router.Method("eth_sendTransaction").Handle(i.EthSendTransaction())
	v2Router.Method("eth_sign").Handle(i.EthSign())
	v2Router.Method("eth_signTransaction").Handle(i.EthSignTransaction())
	v2Router.Method("eea_sendTransaction").Handle(i.EEASendTransaction())

	// Silence JSON-RPC personal
	v2Router.MethodPrefix("personal_").Handle(jsonrpc.MethodNotFoundHandler())

	return jsonrpc.LoggedHandler(jsonrpc.DefaultRWHandler(router))
}

func New(stores storemanager.StoreManager) *Interceptor {
	i := &Interceptor{
		stores: stores,
	}
	i.handler = i.newHandler()
	return i
}
