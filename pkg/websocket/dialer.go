package websocket

import (
	"github.com/consensys/quorum-key-manager/pkg/common"
	"github.com/consensys/quorum-key-manager/pkg/json"
	"github.com/gorilla/websocket"
)

type DialerConfig struct {
	HandshakeTimeout  *json.Duration `json:"handshakeTimeout,omitempty"`
	ReadBufferSize    *int           `json:"readBufferSize,omitempty"`
	WriteBufferSize   *int           `json:"writeBufferSize,omitempty"`
	EnableCompression *bool          `json:"enableCompression,omitempty"`
}

func (cfg *DialerConfig) SetDefault() *DialerConfig {
	if cfg.HandshakeTimeout == nil {
		cfg.HandshakeTimeout = &json.Duration{Duration: 0}
	}

	if cfg.ReadBufferSize == nil {
		cfg.ReadBufferSize = common.ToPtr(int(1024)).(*int)
	}

	if cfg.WriteBufferSize == nil {
		cfg.WriteBufferSize = common.ToPtr(int(1024)).(*int)
	}

	if cfg.EnableCompression == nil {
		cfg.EnableCompression = common.ToPtr(false).(*bool)
	}

	return cfg
}

func NewDialer(cfg *DialerConfig) *websocket.Dialer {
	return &websocket.Dialer{
		HandshakeTimeout:  cfg.HandshakeTimeout.Duration,
		ReadBufferSize:    *cfg.ReadBufferSize,
		WriteBufferSize:   *cfg.WriteBufferSize,
		EnableCompression: *cfg.EnableCompression,
	}
}
