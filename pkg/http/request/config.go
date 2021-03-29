package request

type ProxyConfig struct {
	PassHostHeader *bool             `json:"passHostHeader,omitempty"`
	BasicAuth      *BasicAuthConfig  `json:"basicAuth,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
}

func (cfg *ProxyConfig) SetDefault() {
	if cfg.PassHostHeader != nil {
		cfg.PassHostHeader = new(bool)
		*cfg.PassHostHeader = true
	}

	if cfg.BasicAuth == nil {
		cfg.BasicAuth = &BasicAuthConfig{"", ""}
	}
}

type BasicAuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}
