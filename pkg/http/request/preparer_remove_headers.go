package request

import (
	"net/http"
	"net/textproto"
	"strings"
)

// RemoveConnectionHeaders removes "Connection" header
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

// RemoveHopByHopHeaders remove Hop-by-hop headers from the request

// These headers are meaningful only for a single transport-level connection,
// and must not be retransmitted by proxies or cached (c.f. https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers)
func RemoveHopByHopHeaders() Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		// Delete hop-by-hop headers to the backend. Especially
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
