package httpclient

import (
	"net/http"

	"github.com/consensys/quorum-key-manager/pkg/http/request"
	"github.com/consensys/quorum-key-manager/pkg/http/response"
)

// Decorator decorates a Client
type Decorator func(Client) Client

func WithPreparer(preparer request.Preparer) Decorator {
	return func(c Client) Client {
		return &preparedClient{
			client:   c,
			preparer: preparer,
		}
	}
}

type preparedClient struct {
	client Client

	preparer request.Preparer
}

func (c *preparedClient) Do(req *http.Request) (*http.Response, error) {
	preparedReq, err := c.preparer.Prepare(req)
	if err != nil {
		return nil, err
	}

	return c.client.Do(preparedReq)
}

func (c *preparedClient) CloseIdleConnections() {
	c.client.CloseIdleConnections()
}

type modifiedClient struct {
	client Client

	modifier response.Modifier
}

func WithModifier(modifier response.Modifier) Decorator {
	return func(c Client) Client {
		return &modifiedClient{
			client:   c,
			modifier: modifier,
		}
	}
}

func (c *modifiedClient) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	err = c.modifier.Modify(resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *modifiedClient) CloseIdleConnections() {
	c.client.CloseIdleConnections()
}

func WithRequest(req *http.Request) Decorator {
	return func(c Client) Client {
		return WithPreparer(request.Request(req))(c)
	}
}

func CombineDecorators(decorators ...Decorator) Decorator {
	return func(c Client) Client {
		for i := len(decorators); i > 0; i-- {
			c = decorators[i-1](c)
		}
		return c
	}
}
