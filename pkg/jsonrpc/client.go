package jsonrpc

import (
	"fmt"
	"reflect"
	"sync/atomic"
)

var defaultVersion = "2.0"

//go:generate mockgen -source=client.go -destination=mock/client.go -package=mock

// Client is an jsonrpc HTTPClient interface
type Client interface {
	// Do sends an jsonrpc request and returns an jsonrpc response
	Do(*RequestMsg) (*ResponseMsg, error)
}

type incrementalIDClient struct {
	client Client

	baseID    string
	idCounter uint32
}

// WithIncrementalID wraps a HTTPClient with an ID counter an increases it each time a new request comes out
func WithIncrementalID(id interface{}) func(Client) Client {
	return func(c Client) Client {
		idC := &incrementalIDClient{
			client: c,
		}

		if id != nil {
			idC.baseID = fmt.Sprintf("%v.", id)
		}

		return idC
	}
}

func (c *incrementalIDClient) nextID() string {
	return fmt.Sprintf("%v%v", c.baseID, atomic.AddUint32(&c.idCounter, 1))
}

func (c *incrementalIDClient) Do(msg *RequestMsg) (*ResponseMsg, error) {
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
			return respMsg, InvalidDownstreamResponse(fmt.Errorf("response id does not match request id"))
		}
	}

	return respMsg, nil
}
