package jsonrpc

import (
	"fmt"
	"net/http"
	"reflect"
	"sync/atomic"

	httpclient "github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/client"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/request"
	"github.com/ConsenSysQuorum/quorum-key-manager/pkg/http/response"
)

var defaultVersion = "2.0"

//go:generate mockgen -source=client.go -destination=mock/client.go -package=mock

// Client is an jsonrpc HTTPClient interface
type Client interface {
	// Do sends an jsonrpc request and returns an jsonrpc response
	Do(*RequestMsg) (*ResponseMsg, error)
}

// HTTPClient is a connector to a jsonrpc server
type HTTPClient struct {
	client httpclient.Client
}

// NewClient creates a new jsonrpc HTTPClient from an HTTP HTTPClient
func NewHTTPClient(c httpclient.Client) *HTTPClient {
	return &HTTPClient{
		client: c,
	}
}

// Do sends an jsonrpc request over the underlying HTTP client and returns a jsonrpc response
func (c *HTTPClient) Do(reqMsg *RequestMsg) (*ResponseMsg, error) {
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

type incrementalIDlient struct {
	client Client

	baseID    string
	idCounter uint32
}

// WithIncrementalID wraps a HTTPClient with an ID counter an increases it each time a new request comes out
func WithIncrementalID(id interface{}) func(Client) Client {
	return func(c Client) Client {
		idC := &incrementalIDlient{
			client: c,
		}

		if id != nil {
			idC.baseID = fmt.Sprintf("%v.", id)
		}

		return idC
	}
}

func (c *incrementalIDlient) nextID() string {
	return fmt.Sprintf("%v%v", c.baseID, atomic.AddUint32(&c.idCounter, 1))
}

func (c *incrementalIDlient) Do(msg *RequestMsg) (*ResponseMsg, error) {
	if msg.ID == nil {
		msg.WithID(c.nextID())
	}

	return c.client.Do(msg)
}

type versionClient struct {
	client Client

	version string
}

// WithVersion wraps a HTTPClient to set version each time a new request comes out
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

type validateIDClient struct {
	client Client
}

func ValidateID(client Client) Client {
	return &validateIDClient{client: client}
}

func (c *validateIDClient) Do(msg *RequestMsg) (*ResponseMsg, error) {
	respMsg, err := c.client.Do(msg)
	if err != nil {
		return respMsg, err
	}

	if msg.ID != nil {
		var respIDVal = reflect.New(reflect.TypeOf(msg.ID))
		err = respMsg.UnmarshalID(respIDVal.Interface())
		if err != nil {
			return respMsg, InvalidDownstreamResponse(err)
		}

		if respIDVal.Elem().Interface() != msg.ID {
			fmt.Printf("Piou %T %T\n", respIDVal.Elem().Interface(), msg.ID)
			return respMsg, InvalidDownstreamResponse(fmt.Errorf("response id does not match request id"))
		}
	}

	return respMsg, nil
}
