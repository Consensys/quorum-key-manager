package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

// Request wraps an http.Request with JSON-RPC capabilities
type Request struct {
	req *http.Request

	setMsgOnce sync.Once
	msg        *RequestMsg
	err        error
}

// NewRequest creates a new Request
func NewRequest(req *http.Request) *Request {
	jsonRPCReq := new(Request)
	if req != nil {
		jsonRPCReq = jsonRPCReq.WithRequest(req)
	}

	return jsonRPCReq
}

// Request returns attached http.Request
func (req *Request) Request() *http.Request {
	return req.req
}

// WithRequest attaches http.Request
func (req *Request) WithRequest(r *http.Request) *Request {
	req.req = r.WithContext(WithRequest(r.Context(), req))
	return req
}

func (req *Request) Context() context.Context {
	return req.Request().Context()
}

func (req *Request) WithContext(ctx context.Context) *Request {
	return req.WithRequest(req.Request().WithContext(ctx))
}

// Clone request
func (req *Request) Clone(ctx context.Context) *Request {
	newReq := new(Request)
	if req.msg != nil {
		newReq.msg = req.msg.Copy()
	}

	newReq.err = req.err
	newReq.req = req.req.Clone(WithRequest(ctx, newReq))

	return newReq
}

// WriteBody prepares underlying http.Request body with JSON-RPC message
func (req *Request) WriteBody() error {
	err := req.getMsg().Validate()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(req.msg)
	if err != nil {
		return err
	}

	// Set request body with buffer
	req.req.ContentLength = int64(buf.Len())
	req.req.Body = ioutil.NopCloser(buf)
	snapshot := *buf
	req.req.GetBody = func() (io.ReadCloser, error) {
		r := snapshot
		return ioutil.NopCloser(&r), nil
	}

	return nil
}

// ReadBody reads underlying http.Request body into a JON-RPC message
func (req *Request) ReadBody() error {
	req.setMsgOnce.Do(func() {
		if req.msg == nil {
			req.msg = new(RequestMsg)
		}

		// Read request body into request message and validates it
		err := json.NewDecoder(req.req.Body).Decode(req.msg)
		req.req.Body.Close()
		if err != nil {
			req.err = err
			return
		}

		err = req.msg.Validate()
		if err != nil {
			req.err = err
			return
		}
	})

	return req.err
}

func (req *Request) getMsg() *RequestMsg {
	req.setMsgOnce.Do(func() {
		if req.msg == nil {
			req.msg = new(RequestMsg)
		}
	})
	return req.msg
}

// Version returns JSON-RPC request version
func (req *Request) Version() string {
	return req.getMsg().Version
}

// WithVersion changes JSON-RPC request version
func (req *Request) WithVersion(v string) *Request {
	req.getMsg().WithVersion(v)
	return req
}

// Method returns JSON-RPC request method
func (req *Request) Method() string {
	return req.getMsg().Method
}

// WithMethod changes JSON-RPC request method
func (req *Request) WithMethod(method string) *Request {
	req.getMsg().WithMethod(method)
	return req
}

// ID returns JSON-RPC request method
func (req *Request) ID() interface{} {
	return req.getMsg().ID
}

// WithID changes JSON-RPC request ID
func (req *Request) WithID(id interface{}) *Request {
	req.getMsg().WithID(id)
	return req
}

// UnmarshalID unmarshals JSON-RPC request id into v
func (req *Request) UnmarshalID(v interface{}) error {
	return req.getMsg().UnmarshalID(v)
}

// ID returns JSON-RPC request parameters
func (req *Request) Params() interface{} {
	return req.getMsg().Params
}

// WithParams set params
func (req *Request) WithParams(v interface{}) *Request {
	req.getMsg().WithParams(v)
	return req
}

// UnmarshalParams unmarshals JSON-RPC request params into v
func (req *Request) UnmarshalParams(v interface{}) error {
	return req.getMsg().UnmarshalParams(v)
}

// Error returns a possible error encountered while reading JSON-RPC request
func (req *Request) Error() error {
	// call getMsg so if ReadBody has been called in parallel
	// then it will block until ReadBody completes and return the error
	_ = req.getMsg()
	return req.err
}
