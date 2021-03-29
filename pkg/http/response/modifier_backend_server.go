package response

import "net/http"

// BackendServer set "X-Backend-Server" header to the the URL of the request
func BackendServer() Modifier {
	return ModifierFunc(func(resp *http.Response) error {
		resp.Header.Set("X-Backend-Server", resp.Request.URL.String())
		return nil
	})
}
