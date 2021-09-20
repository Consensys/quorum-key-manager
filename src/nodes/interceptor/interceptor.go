package interceptor

import (
	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	aliasent "github.com/consensys/quorum-key-manager/src/aliases/entities"
	"github.com/consensys/quorum-key-manager/src/infra/log"
	proxynode "github.com/consensys/quorum-key-manager/src/nodes/node/proxy"
	"github.com/consensys/quorum-key-manager/src/stores"
)

type Interceptor struct {
	stores  stores.Stores
	handler jsonrpc.Handler
	logger  log.Logger
	aliases aliasent.AliasBackend
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

	return jsonrpc.LoggedHandler(jsonrpc.DefaultRWHandler(router), i.logger)
}

func New(storesConnector stores.Stores, aliases aliasent.AliasBackend, logger log.Logger) *Interceptor {
	i := &Interceptor{
		stores:  storesConnector,
		aliases: aliases,
		logger:  logger,
	}

	i.handler = i.newHandler()

	return i
}
