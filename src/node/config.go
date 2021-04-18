package node

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

func (cfg *ProxyConfig) Copy() *ProxyConfig {
	if cfg == nil {
		return nil
	}
	return &ProxyConfig{
		Request:  cfg.Request.Copy(),
		Response: cfg.Response.Copy(),
	}
}

func (cfg *ProxyConfig) SetDefault() *ProxyConfig {
	if cfg.Request == nil {
		cfg.Request = new(request.ProxyConfig)
	}

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

func (cfg *DownstreamConfig) Copy() *DownstreamConfig {
	if cfg == nil {
		return nil
	}

	return &DownstreamConfig{
		Addr:          cfg.Addr,
		Transport:     cfg.Transport.Copy(),
		Proxy:         cfg.Proxy.Copy(),
		ClientTimeout: cfg.ClientTimeout.Copy(),
	}
}

func (cfg *DownstreamConfig) SetDefault() *DownstreamConfig {
	defaultCfg := new(httpclient.Config).SetDefault()

	if cfg.Transport == nil {
		cfg.Transport = defaultCfg.Transport
	}

	if cfg.Proxy == nil {
		cfg.Proxy = new(ProxyConfig)
	}

	if cfg.Proxy.Request == nil {
		cfg.Proxy.Request = defaultCfg.Request
	}

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
	RPC           *DownstreamConfig `json:"json-rpc,omitempty"`
	PrivTxManager *DownstreamConfig `json:"tessera,omitempty"`
}

func (cfg *Config) Copy() *Config {
	return &Config{
		RPC:           cfg.RPC.Copy(),
		PrivTxManager: cfg.PrivTxManager.Copy(),
	}
}

func (cfg *Config) SetDefault() *Config {
	if cfg == nil {
		return nil
	}

	if cfg.RPC == nil {
		cfg.RPC = new(DownstreamConfig)
	}

	cfg.RPC.SetDefault()

	if cfg.PrivTxManager != nil {
		cfg.PrivTxManager.SetDefault()
	}

	return cfg
}
