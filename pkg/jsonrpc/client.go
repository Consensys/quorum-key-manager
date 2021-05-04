package jsonrpc

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	httpclient "github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
)

var defaultVersion = "2.0"

type ClientConfig struct {
	Version string             `json:"version,omitempty"`
	HTTP    *httpclient.Config `json:"http,omitempty"`
}

func (cfg *ClientConfig) SetDefault() *ClientConfig {
	if cfg.HTTP == nil {
		cfg.HTTP = new(httpclient.Config)
	}

	cfg.HTTP.SetDefault()

	if cfg.Version == "" {
		cfg.Version = defaultVersion
	}

	return cfg
}

// Client is an jsonrpc HTTPClient interface
type Client interface {
	// Do sends an jsonrpc request and returns an jsonrpc response
	Do(*RequestMsg) (*ResponseMsg, error)
}

// client is a connector to a jsonrpc server
type client struct {
	client httpclient.Client
}

// NewClient creates a new jsonrpc client from an HTTP client
func NewClient(c httpclient.Client) Client {
	return &client{
		client: c,
	}
}

// Do sends an jsonrpc request and returns an jsonrpc response
func (c *client) Do(reqMsg *RequestMsg) (*ResponseMsg, error) {
	err := reqMsg.Validate()
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(reqMsg.Context(), http.MethodPost, "", nil)

	// write request body
	err = request.WriteJSON(req, reqMsg)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, DownstreamError(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, InvalidDownstreamHTTPStatuError(resp.StatusCode)
	}

	// Create response and reads body
	respMsg := new(ResponseMsg)
	err = response.ReadJSON(resp, respMsg)
	if err != nil {
		return nil, InvalidDownstreamResponse(err)
	}

	err = respMsg.Validate()
	if err != nil {
		return nil, InvalidDownstreamResponse(err)
	}

	return respMsg, nil
}

type idClient struct {
	client Client

	baseID    string
	idCounter uint32
}

// WithID wraps a client with an ID counter an increases it each time a new request comes out
func WithID(id interface{}) func(Client) Client {
	return func(c Client) Client {
		idC := &idClient{
			client: c,
		}

		if id != nil {
			idC.baseID = fmt.Sprintf("%v.", id)
		}

		return idC
	}
}

func (c *idClient) nextID() string {
	return fmt.Sprintf("%v%v", c.baseID, atomic.AddUint32(&c.idCounter, 1))
}

func (c *idClient) Do(msg *RequestMsg) (*ResponseMsg, error) {
	msg.WithID(c.nextID())
	return c.client.Do(msg)
}

type versionClient struct {
	client Client

	version string
}

// WithVersion wraps a client to set version each time a new request comes out
func WithVersion(version string) func(Client) Client {
	return func(c Client) Client {
		if version == "" {
			version = defaultVersion
		}
		return &versionClient{
			client:  c,
			version: version,
		}
	}
}

func (c *versionClient) Do(msg *RequestMsg) (*ResponseMsg, error) {
	if msg.Version == "" {
		msg.WithVersion(c.version)
	}

	return c.client.Do(msg)
}

// Caller is an interface for a JSON-RPC caller
type Caller interface {
	Call(ctx context.Context, method string, params interface{}) (*ResponseMsg, error)
}
