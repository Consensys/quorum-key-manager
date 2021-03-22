package request

import (
	"net/http"
)

// Headers sets or deletes custom request headers
func Headers(headers map[string]string) Preparer {
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
