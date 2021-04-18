package request

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// Body ensures that GetBody method is set on request
func Body() Preparer {
	return PrepareFunc(func(req *http.Request) (*http.Request, error) {
		if req.GetBody != nil {
			return req, nil
		}

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

		return req, nil
	})
}
