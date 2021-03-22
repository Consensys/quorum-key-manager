package request

import (
	"net"
	"net/http"
	"strings"
)

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
