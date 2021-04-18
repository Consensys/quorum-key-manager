package response

type ProxyConfig struct {
	Headers map[string]string `json:"headers,omitempty"`
}

func (cfg *ProxyConfig) Copy() *ProxyConfig {
	if cfg == nil {
		return nil
	}

	cpy := new(ProxyConfig)
	*cpy = *cfg
	return cpy
}
