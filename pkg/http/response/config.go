package response

type ProxyConfig struct {
	Headers map[string][]string `json:"headers,omitempty"`
}

func (cfg *ProxyConfig) SetDefault() *ProxyConfig {
	return cfg
}
