package response

type ProxyConfig struct {
	Headers map[string][]string `json:"headers,omitempty"`
}
