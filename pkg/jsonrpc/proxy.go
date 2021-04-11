package jsonrpc

import (
	"net/http"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/proxy"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
)

// HandleProxyRoundTripError allows to transform a ProxiedRoundTrip Error
func HandleProxyRoundTripError(rw http.ResponseWriter, req *http.Request, err error) {
	rpcRw, ok := rw.(ResponseWriter)
	if !ok {
		proxy.HandleRoundTripError(rw, req, err)
	}

	statusCode := proxy.StatusCodeFromRoundTripError(err)
	statusText := proxy.StatusText(statusCode)

	logger := log.FromContext(req.Context())
	logger.Debugf("'%d %s' caused by: %v", statusCode, statusText, err)

	errMsg := &ErrorMsg{
		Code:    -32000,
		Message: statusText,
	}

	werr := rpcRw.WriteError(errMsg)
	if werr != nil {
		logger.Debugf("Error while writing error message", werr)
	}
}
