package response

import "net/http"

// Headers set "X-Backend-Server" header to the the URL of the request
func Headers(headers map[string]string) Modifier {
	return ModifierFunc(func(resp *http.Response) error {
		for header, value := range headers {
			switch {
			case value == "":
				resp.Header.Del(header)
			default:
				resp.Header.Set(header, value)
			}
		}
		return nil
	})
}
