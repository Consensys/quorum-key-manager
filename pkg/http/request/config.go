package request

type ProxyConfig struct {
	Addr           string            `json:"addr,omitempty"`
	PassHostHeader *bool             `json:"passHostHeader,omitempty"`
	BasicAuth      *BasicAuthConfig  `json:"basicAuth,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
}

func (cfg *ProxyConfig) Copy() *ProxyConfig {
	if cfg == nil {
		return nil
	}

	cpy := new(ProxyConfig)

	*cpy = *cfg

	if cfg.PassHostHeader != nil {
		cpy.PassHostHeader = new(bool)
		*cpy.PassHostHeader = *cfg.PassHostHeader
	}

	cpy.BasicAuth = cfg.BasicAuth.Copy()

	return cpy
}

func (cfg *ProxyConfig) SetDefault() *ProxyConfig {
	if cfg == nil {
		cfg = new(ProxyConfig)
	}

	if cfg.PassHostHeader != nil {
		cfg.PassHostHeader = new(bool)
		*cfg.PassHostHeader = true
	}

	if cfg.BasicAuth == nil {
		cfg.BasicAuth = &BasicAuthConfig{"", ""}
	}

	return cfg
}

type BasicAuthConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (cfg *BasicAuthConfig) Copy() *BasicAuthConfig {
	if cfg == nil {
		return nil
	}
	cpy := new(BasicAuthConfig)
	*cpy = *cfg
	return cpy
}
