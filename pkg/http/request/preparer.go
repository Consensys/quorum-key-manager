package http

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

// Preparer is the interface that wraps alows to prepare an http.Request
//
// Prepare accepts and possibly modifies an Request (e.g., adding Headers). Implementations
// must ensure to not share or hold per-invocation state since Preparers may be shared and re-used.
type Preparer interface {
	Prepare(*http.Request) (*http.Request, error)
}

// PrepareFunc is a method that implements the Preparer interface.
type PrepareFunc func(*http.Request) (*http.Request, error)

// Prepare implements the Preparer interface on PrepareFunc.
func (f PrepareFunc) Prepare(r *http.Request) (*http.Request, error) {
	return f(r)
}

// WebSocketHeaders enforce headers to be case-insensitive

// Even if the websocket RFC says that headers should be case-insensitive,
// some servers need Sec-WebSocket-Key, Sec-WebSocket-Extensions, Sec-WebSocket-Accept,
// Sec-WebSocket-Protocol and Sec-WebSocket-Version to be case-sensitive.
// https://tools.ietf.org/html/rfc6455#page-20
func WebSocketHeaders() Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		req.Header["Sec-WebSocket-Key"] = req.Header["Sec-Websocket-Key"]
		req.Header["Sec-WebSocket-Extensions"] = req.Header["Sec-Websocket-Extensions"]
		req.Header["Sec-WebSocket-Accept"] = req.Header["Sec-Websocket-Accept"]
		req.Header["Sec-WebSocket-Protocol"] = req.Header["Sec-Websocket-Protocol"]
		req.Header["Sec-WebSocket-Version"] = req.Header["Sec-Websocket-Version"]
		delete(req.Header, "Sec-Websocket-Key")
		delete(req.Header, "Sec-Websocket-Extensions")
		delete(req.Header, "Sec-Websocket-Accept")
		delete(req.Header, "Sec-Websocket-Protocol")
		delete(req.Header, "Sec-Websocket-Version")
		return req, nil
	})
}

// CombinePreparer combines multiple preparers into a single one
func CombinePreparer(preparers ...Preparer) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		var err error
		for _, preparer := range preparers {
			req, err = preparer.Prepare(req)
			if err != nil {
				return req, err
			}
		}

		return req, nil
	})
}

// BasicAuth sets user on request if unset
// Authorization header with "Basic <username>:<password>" and attaches user
// if User is set on request URL
func BasicAuthorization(username, password string) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if req.Header.Get("Authorization") != "" {
			// If Authorization Header is already set then do not alter request
			return req, nil
		}

		u := req.URL.User
		if u == nil && username != "" && password != "" {
			// If user is not set and username/password are valid
			// then populates users
			u = url.UserPassword(username, password)
			req.URL.User = u
		}

		if u != nil {
			// If user has been set then set Authorization Header with corresponding Basic Authorization header
			username = u.Username()
			password, _ = u.Password()
			req.Header.Set("Authorization", fmt.Sprintf("Basic %v", basicAuth(username, password)))
		}

		return req, nil
	})
}

// UserAgent sets User-Agent header
func UserAgent(agent string) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if agent != "" {
			req.Header.Set("User-Agent", agent)
		} else if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Del("User-Agent")
		}

		return req, nil
	})
}

// CustomHeaders sets or deletes custom request headers
func CustomHeaders(headers map[string]string) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		// Loop through Custom request headers
		for header, value := range headers {
			switch {
			case value == "":
				req.Header.Del(header)
			default:
				req.Header.Set(header, value)
			}
		}

		return req, nil
	})
}

// Protocol sets HTTP protocol on request
//
// Example: HTTP/1.1 or HTTP/2.
func Protocol(major, minor int) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		req.Proto = fmt.Sprintf("HTTP/%v.%v", major, minor)
		req.ProtoMajor = major
		req.ProtoMinor = minor
		return req, nil
	})
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// Host set request Host
// This is useful when proxying request (to set the request host to match downstream server)

// If host is nil then it sets the Host to request URL host
func Host(host *string) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if host == nil {
			req.Host = req.URL.Host
		} else {
			req.Host = *host
		}
		return req, nil
	})
}

// ExtractURI parses request URI, updates request URL and optionnaly reset URI
func ExtractURI(reset bool) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if req.RequestURI != "" {
			u, err := url.ParseRequestURI(req.RequestURI)
			if err != nil {
				return req, err
			}

			req.URL.Path = u.Path
			req.URL.RawPath = u.RawPath
			req.URL.RawQuery = u.RawQuery
			if reset {
				req.RequestURI = ""
			}
		}

		return req, nil
	})
}

// GetBody ensures that GetBody method is set on request
func GetBody() Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if req.GetBody == nil {
			// See More https://github.com/golang/net/blob/master/http2/transport.go#L554
			if req.Body != nil {
				body, _ := ioutil.ReadAll(req.Body)
				req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				req.GetBody = func() (io.ReadCloser, error) {
					return ioutil.NopCloser(bytes.NewBuffer(body)), nil
				}
			} else {
				req.GetBody = func() (io.ReadCloser, error) {
					return nil, nil
				}
			}
		}

		return req, nil
	})
}

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
var hopByHopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}

// RemoveConnectionHeaders removes hop-by-hop headers listed in the "Connection" header
// See RFC 7230, section 6.1
func RemoveConnectionHeaders() Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		for _, f := range req.Header["Connection"] {
			for _, sf := range strings.Split(f, ",") {
				if sf = textproto.TrimString(sf); sf != "" {
					req.Header.Del(sf)
				}
			}
		}

		return req, nil
	})
}

// RemoveHopByHopHeaders remove Hop-by-hop headers from the request

// These headers are meaningful only for a single transport-level connection,
// and must not be retransmitted by proxies or cached (c.f. https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers)
func RemoveHopByHopHeaders() Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		// Remove hop-by-hop headers to the backend. Especially
		// important is "Connection" because we want a persistent
		// connection, regardless of what the client sent to us.
		for _, h := range hopByHopHeaders {
			hv := req.Header.Get(h)
			if hv == "" {
				continue
			}
			if h == "Te" && hv == "trailers" {
				// Tell backend applications that
				// care about trailer support that we support
				// trailers. (We do, but we don't go out of
				// our way to advertise that unless the
				// incoming client request thought it was
				// worth mentioning)
				continue
			}
			req.Header.Del(h)
		}

		return req, nil
	})
}

// ForwardedFor populates "X-Forwarded-For" header with the client IP address

// In case, "X-Forwarded-For" was already populated (e.g. if we are not the first proxy)
// then retains prior X-Forwarded-For information as a comma+space
// separated list and fold multiple headers into one.
func ForwardedFor() Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
			// If we aren't the first proxy retain prior
			// X-Forwarded-For information as a comma+space
			// separated list and fold multiple headers into one.
			prior, ok := req.Header["X-Forwarded-For"]
			omit := ok && prior == nil
			if len(prior) > 0 {
				clientIP = strings.Join(prior, ", ") + ", " + clientIP
			}
			if !omit {
				req.Header.Set("X-Forwarded-For", clientIP)
			}
		}

		return req, nil
	})
}
