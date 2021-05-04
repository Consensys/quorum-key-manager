package proxynode

import (
	httpclient "github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/transport"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/json"
)

// ProxyConfig
type ProxyConfig struct {
	Request  *request.ProxyConfig  `json:"request,omitempty"`
	Response *response.ProxyConfig `json:"response,omitempty"`
}

func (cfg *ProxyConfig) SetDefault() *ProxyConfig {
	if cfg.Request == nil {
		cfg.Request = new(request.ProxyConfig)
	}
	cfg.Request.SetDefault()

	if cfg.Response == nil {
		cfg.Response = new(response.ProxyConfig)
	}

	return cfg
}

// DownstreamConfig
type DownstreamConfig struct {
	Addr          string            `json:"addr,omitempty"`
	Transport     *transport.Config `json:"transport,omitempty"`
	Proxy         *ProxyConfig      `json:"proxy,omitempty"`
	ClientTimeout *json.Duration    `json:"clientTimeout,omitempty"`
}

func (cfg *DownstreamConfig) SetDefault() *DownstreamConfig {
	defaultCfg := new(httpclient.Config)
	defaultCfg.SetDefault()

	if cfg.Transport == nil {
		cfg.Transport = defaultCfg.Transport
	}
	cfg.Transport.SetDefault()

	if cfg.Proxy == nil {
		cfg.Proxy = new(ProxyConfig)
	}
	cfg.Proxy.SetDefault()

	if cfg.Proxy.Request == nil {
		cfg.Proxy.Request = defaultCfg.Request
	}
	cfg.Proxy.Request.SetDefault()

	if cfg.Proxy.Request.Addr == "" {
		cfg.Proxy.Request.Addr = cfg.Addr
	}

	if cfg.Proxy.Response == nil {
		cfg.Proxy.Response = defaultCfg.Response
	}

	if cfg.ClientTimeout == nil {
		cfg.ClientTimeout = defaultCfg.Timeout
	}

	return cfg
}

// Config is the cfg format for an Hashicorp Vault secret store
type Config struct {
	RPC           *DownstreamConfig `json:"rpc,omitempty"`
	PrivTxManager *DownstreamConfig `json:"tessera,omitempty"`
}

func (cfg *Config) SetDefault() *Config {
	if cfg.RPC == nil {
		cfg.RPC = new(DownstreamConfig)
	}
	cfg.RPC.SetDefault()

	if cfg.RPC.Addr == "" {
		cfg.RPC.Addr = "localhost:8545"
	}

	if cfg.PrivTxManager != nil {
		cfg.PrivTxManager.SetDefault()
	}

	return cfg
}
