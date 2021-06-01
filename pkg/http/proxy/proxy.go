package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/oxtoacart/bpool"
)

// New creates a new HTTP proxy
func New(
	cfg *Config,
	trnsprt http.RoundTripper,
	preparer request.Preparer,
	modifier response.Modifier,
	errorHandler HandleRoundTripErrorFunc,
	pool httputil.BufferPool,
) (*httputil.ReverseProxy, error) {
	var err error
	if trnsprt == nil {
		trnsprt, err = transport.New(cfg.Transport)
		if err != nil {
			return nil, err
		}
	}

	if preparer == nil {
		preparer, err = request.Proxy(cfg.Request)
		if err != nil {
			return nil, err
		}
	}

	if modifier == nil {
		modifier = response.Proxy(cfg.Response)
	}

	if pool == nil {
		pool = bpool.NewBytePool(32, 1024)
	}

	if errorHandler == nil {
		errorHandler = HandleRoundTripError
	}

	return &httputil.ReverseProxy{
		Director: func(outReq *http.Request) {
			u := outReq.URL
			log.FromContext(outReq.Context()).
				Errorf("proxy call - outReq.URL: %s - u.EscapedPath(): %s, outReq.RequestURI: %s", outReq.URL, u.EscapedPath(), outReq.RequestURI)
			newReq, _ := preparer.Prepare(outReq)
			*outReq = *newReq
		},
		Transport:      trnsprt,
		FlushInterval:  cfg.FlushInterval.Duration,
		ModifyResponse: modifier.Modify,
		BufferPool:     pool,
		ErrorHandler:   errorHandler,
	}, nil
}
