package jsonrpc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Response wraps an http.Response with JSON-RPC capabilities
type Response struct {
	resp *http.Response

	setMsgOnce sync.Once
	msg        *ResponseMsg
	err        error
}

// NewResponse creates a new response object
func NewResponse(resp *http.Response) *Response {
	return &Response{
		resp: resp,
	}
}

// Response returns encapsulated response
func (resp *Response) Response() *http.Response {
	return resp.resp
}

// ReadBody reads underlying http.Response body into a JON-RPC message
func (resp *Response) ReadBody() error {
	resp.setMsgOnce.Do(func() {
		if resp.msg == nil {
			resp.msg = new(ResponseMsg)
		}

		if resp.resp.StatusCode < 200 || resp.resp.StatusCode >= 300 {
			resp.err = fmt.Errorf("invalid http response: %v (code=%v)", http.StatusText(resp.resp.StatusCode), resp.resp.StatusCode)
			return
		}

		defer resp.resp.Body.Close()

		// Read response body into response message and validates it
		err := json.NewDecoder(resp.resp.Body).Decode(resp.msg)
		if err != nil {
			resp.err = err
			return
		}

		err = resp.msg.Validate()
		if err != nil {
			resp.err = err
			return
		}
	})

	return resp.err
}

func (resp *Response) getMsg() *ResponseMsg {
	resp.setMsgOnce.Do(func() {
		if resp.msg == nil {
			resp.msg = new(ResponseMsg)
		}
	})
	return resp.msg
}

// Version returns JSON-RPC response version
func (resp *Response) Version() string {
	return resp.getMsg().Version
}

// WithVersion changes JSON-RPC response version
func (resp *Response) WithVersion(v string) *Response {
	resp.getMsg().WithVersion(v)
	return resp
}

// ID returns JSON-RPC response method
func (resp *Response) ID() interface{} {
	return resp.getMsg().ID
}

// WithID changes JSON-RPC response ID
func (resp *Response) WithID(id interface{}) *Response {
	resp.getMsg().WithID(id)
	return resp
}

// UnmarshalID unmarshals JSON-RPC response id into v
func (resp *Response) UnmarshalID(v interface{}) error {
	return resp.getMsg().UnmarshalID(v)
}

// Result returns JSON-RPC response Result
func (resp *Response) Result() interface{} {
	return resp.getMsg().Result
}

// WithResult set result
func (resp *Response) WithResult(v interface{}) *Response {
	resp.getMsg().WithResult(v)
	return resp
}

// UnmarshalResult unmarshals JSON-RPC response params into v
func (resp *Response) UnmarshalResult(v interface{}) error {
	return resp.getMsg().UnmarshalResult(v)
}

// WithError attached error
func (resp *Response) WithError(err error) *Response {
	resp.getMsg().WithError(err)
	return resp
}

// Error returns a possible error encountered while reading JSON-RPC response or JSON-RPC error
func (resp *Response) Error() error {
	// Force response loading
	msg := resp.getMsg()
	if resp.err != nil {
		return resp.err
	}

	return msg.Error
}
