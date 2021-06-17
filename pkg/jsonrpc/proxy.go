package jsonrpc

import (
	"net/http"

	"github.com/consensysquorum/quorum-key-manager/pkg/http/proxy"
)

// HandleProxyRoundTripError allows to transform a ProxiedRoundTrip Error
func HandleProxyRoundTripError(rw http.ResponseWriter, req *http.Request, err error) {
	rpcRw, ok := rw.(ResponseWriter)
	if !ok {
		proxy.HandleRoundTripError(rw, req, err)
	}

	_ = WriteError(rpcRw, DownstreamError(err))
}
