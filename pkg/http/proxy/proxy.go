package proxy

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/log"
	"github.com/oxtoacart/bpool"
)

// StatusClientClosedRequest non-standard HTTP status code for client disconnection
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection
const StatusClientClosedRequestText = "Client Closed Connection"

// New creates a new HTTP proxy
func New(cfg *Config, trnsprt http.RoundTripper, pool httputil.BufferPool) (*httputil.ReverseProxy, error) {
	cfg.SetDefault()

	var err error
	if trnsprt == nil {
		trnsprt, err = transport.New(cfg.Transport)
		if err != nil {
			return nil, err
		}
	}

	if pool == nil {
		pool = bpool.NewBytePool(32, 1024)
	}

	return &httputil.ReverseProxy{
		Director: func(outReq *http.Request) {
			fmt.Println("Proxying Piou")

			u := outReq.URL
			if outReq.RequestURI != "" {
				parsedURL, err := url.ParseRequestURI(outReq.RequestURI)
				if err == nil {
					u = parsedURL
				}
			}

			outReq.URL.Path = u.Path
			outReq.URL.RawPath = u.RawPath
			outReq.URL.RawQuery = u.RawQuery
			outReq.RequestURI = "" // Outgoing request should not have RequestURI

			outReq.Proto = "HTTP/1.1"
			outReq.ProtoMajor = 1
			outReq.ProtoMinor = 1

			if _, ok := outReq.Header["User-Agent"]; !ok {
				outReq.Header.Set("User-Agent", "")
			}

			if cfg.PassHostHeader != nil && !*cfg.PassHostHeader {
				outReq.Host = outReq.URL.Host
			}

			// Even if the websocket RFC says that headers should be case-insensitive,
			// some servers need Sec-WebSocket-Key, Sec-WebSocket-Extensions, Sec-WebSocket-Accept,
			// Sec-WebSocket-Protocol and Sec-WebSocket-Version to be case-sensitive.
			// https://tools.ietf.org/html/rfc6455#page-20
			outReq.Header["Sec-WebSocket-Key"] = outReq.Header["Sec-Websocket-Key"]
			outReq.Header["Sec-WebSocket-Extensions"] = outReq.Header["Sec-Websocket-Extensions"]
			outReq.Header["Sec-WebSocket-Accept"] = outReq.Header["Sec-Websocket-Accept"]
			outReq.Header["Sec-WebSocket-Protocol"] = outReq.Header["Sec-Websocket-Protocol"]
			outReq.Header["Sec-WebSocket-Version"] = outReq.Header["Sec-Websocket-Version"]
			delete(outReq.Header, "Sec-Websocket-Key")
			delete(outReq.Header, "Sec-Websocket-Extensions")
			delete(outReq.Header, "Sec-Websocket-Accept")
			delete(outReq.Header, "Sec-Websocket-Protocol")
			delete(outReq.Header, "Sec-Websocket-Version")

			// It allows to proxy servers that uses authentication through URL (e.g. https://user:password@example.com)
			// In particular it allows to support nodes on Kaleido
			if u := outReq.URL.User; u != nil && outReq.Header.Get("Authorization") == "" {
				username := u.Username()
				password, _ := u.Password()
				outReq.Header.Set("Authorization", fmt.Sprintf("Basic %v", basicAuth(username, password)))
			}

			// See More https://github.com/golang/net/blob/master/http2/transport.go#L554
			if outReq.Body != nil {
				body, _ := ioutil.ReadAll(outReq.Body)
				outReq.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				outReq.GetBody = func() (io.ReadCloser, error) {
					return ioutil.NopCloser(bytes.NewBuffer(body)), nil
				}
			} else {
				outReq.GetBody = func() (io.ReadCloser, error) {
					return nil, nil
				}
			}
		},
		Transport:     trnsprt,
		FlushInterval: cfg.FlushInterval.Duration,
		ModifyResponse: func(resp *http.Response) error {
			resp.Header.Set("X-Backend-Server", resp.Request.URL.String())
			if resp.StatusCode >= 300 {
				body, _ := ioutil.ReadAll(resp.Body)
				resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				log.FromContext(resp.Request.Context()).
					Debugf("'%d %s' caused by: %q", resp.StatusCode, statusText(resp.StatusCode), string(body))
			}

			return nil
		},
		BufferPool: pool,
		ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
			logger := log.FromContext(req.Context())
			fmt.Println("Error Handling Piou")
			statusCode := http.StatusInternalServerError

			switch {
			case errors.Is(err, io.EOF):
				statusCode = http.StatusBadGateway
			case errors.Is(err, context.Canceled):
				statusCode = StatusClientClosedRequest
			default:
				var netErr net.Error
				if errors.As(err, &netErr) {
					if netErr.Timeout() {
						statusCode = http.StatusGatewayTimeout
					} else {
						statusCode = http.StatusBadGateway
					}
				}
			}

			logger.Debugf("'%d %s' caused by: %v", statusCode, statusText(statusCode), err)
			rw.WriteHeader(statusCode)
			_, werr := rw.Write([]byte(statusText(statusCode)))
			if werr != nil {
				logger.Debugf("Error while writing status code", werr)
			}
		},
	}, nil
}

func statusText(statusCode int) string {
	if statusCode == StatusClientClosedRequest {
		return StatusClientClosedRequestText
	}
	return http.StatusText(statusCode)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
