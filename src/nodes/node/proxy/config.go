package proxynode

import (
	httpclient "github.com/consensys/quorum-key-manager/pkg/http/client"
	"github.com/consensys/quorum-key-manager/pkg/http/request"
	"github.com/consensys/quorum-key-manager/pkg/http/response"
	"github.com/consensys/quorum-key-manager/pkg/http/transport"
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/consensys/quorum-key-manager/pkg/websocket"
)

type ProxyConfig struct {
	Request   *request.ProxyConfig   `json:"request,omitempty" yaml:"request,omitempty"`
	Response  *response.ProxyConfig  `json:"response,omitempty" yaml:"response,omitempty"`
	WebSocket *websocket.ProxyConfig `json:"websocket,omitempty" yaml:"websocket,omitempty"`
}

func (cfg *ProxyConfig) SetDefault() *ProxyConfig {
	if cfg.Request == nil {
		cfg.Request = new(request.ProxyConfig)
	}
	cfg.Request.SetDefault()

	if cfg.Response == nil {
		cfg.Response = new(response.ProxyConfig)
	}

	if cfg.WebSocket == nil {
		cfg.WebSocket = new(websocket.ProxyConfig)
	}

	cfg.WebSocket.SetDefault()

	return cfg
}

type DownstreamConfig struct {
	Addr          string            `json:"addr" yaml:"addr" validate:"required" example:"http://geth:8545"`
	Transport     *transport.Config `json:"transport,omitempty" yaml:"transport,omitempty"`
	Proxy         *ProxyConfig      `json:"proxy,omitempty" yaml:"proxy,omitempty"`
	ClientTimeout *json.Duration    `json:"clientTimeout,omitempty" yaml:"client_timeout,omitempty"`
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
		cfg.Proxy.Request = new(request.ProxyConfig)
	}

	if cfg.Proxy.Request.Addr == "" {
		cfg.Proxy.Request.Addr = cfg.Addr
	}

	cfg.Proxy.Request.SetDefault()

	if cfg.Proxy.Response == nil {
		cfg.Proxy.Response = new(response.ProxyConfig)
	}
	cfg.Proxy.Response.SetDefault()

	if cfg.ClientTimeout == nil {
		cfg.ClientTimeout = defaultCfg.Timeout
	}

	if cfg.Proxy.WebSocket == nil {
		cfg.Proxy.WebSocket = new(websocket.ProxyConfig)
	}
	cfg.Proxy.WebSocket.SetDefault()

	return cfg
}

// Config is the cfg format for a Hashicorp Vault secret store
type Config struct {
	RPC           *DownstreamConfig `json:"rpc,omitempty" yaml:"rpc,omitempty"`
	PrivTxManager *DownstreamConfig `json:"tessera,omitempty" yaml:"tessera,omitempty"`
}

func (cfg *Config) SetDefault() *Config {
	if cfg.RPC == nil {
		cfg.RPC = new(DownstreamConfig)
	}
	cfg.RPC.SetDefault()

	if cfg.PrivTxManager != nil {
		cfg.PrivTxManager.SetDefault()
	}

	return cfg
}
