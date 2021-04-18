package request

import (
	"fmt"
	"net/http"
)

// HTTPProtocol sets HTTP protocol on request
//
// Example: HTTP/1.1 or HTTP/2.
func HTTPProtocol(major, minor int) Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		req.Proto = fmt.Sprintf("HTTP/%v.%v", major, minor)
		req.ProtoMajor = major
		req.ProtoMinor = minor
		return req, nil
	})
}
