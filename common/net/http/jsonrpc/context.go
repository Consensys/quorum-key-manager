package jsonrpc

import (
	"context"
	"net/http"
	"sync"
)

type Context struct {
	req *http.Request
	rw  http.ResponseWriter

	readMsgOnce sync.Once
	msg         *RequestMsg
	err         error
}

func fromRequest(req *http.Request, rw http.ResponseWriter) (*Context, bool) {
	hctx, ok := req.Context().Value(ctxCtxKey).(*Context)
	if ok {
		hctx.rw = rw
		hctx.req = req
	}
	return hctx, ok
}

func newContext(rw http.ResponseWriter, req *http.Request) *Context {
	hctx := new(Context)
	hctx.withRequest(req)
	hctx.rw = rw

	return hctx
}

func (hctx *Context) withRequest(req *http.Request) {
	hctx.req = req.WithContext(context.WithValue(req.Context(), ctxCtxKey, hctx))
}

func (hctx *Context) reset(rw http.ResponseWriter, req *http.Request) {
	hctx.withRequest(req)
	hctx.rw = rw
	hctx.err = nil
	hctx.msg = nil
	hctx.readMsgOnce = sync.Once{}
}

// Request returns request attach to context
func (hctx *Context) Request() *http.Request {
	return hctx.req
}

// Context returns golang context attached to underlying request
func (hctx *Context) Context() context.Context {
	return hctx.req.Context()
}

// WithContext attached a golang context to underlying request
func (hctx *Context) WithContext(ctx context.Context) *Context {
	hctx.req = hctx.req.WithContext(ctx)
	return hctx
}

// Writer returns response writer attach to context
func (hctx *Context) Writer() http.ResponseWriter {
	return hctx.rw
}

// Version returns JSON-RPC request version
func (hctx *Context) Version() string {
	return hctx.getReqMsg().Version
}

// WithMethod changes JSON-RPC request method
func (hctx *Context) WithVersion(v string) *Context {
	msg := hctx.getReqMsg()
	msg.Version = v
	return hctx
}

// Method returns JSON-RPC request method
func (hctx *Context) Method() string {
	return hctx.getReqMsg().Method
}

// WithMethod changes JSON-RPC request method
func (hctx *Context) WithMethod(method string) *Context {
	msg := hctx.getReqMsg()
	msg.Method = method
	return hctx
}

// ID returns JSON-RPC request method
func (hctx *Context) ID() []byte {
	return hctx.getReqMsg().ID
}

// WithID changes JSON-RPC request ID
func (hctx *Context) WithID(id interface{}) (*Context, error) {
	msg := hctx.getReqMsg()
	err := msg.WithID(id)
	if err == nil {
		hctx.msg = msg
	}
	return hctx, err
}

// Params unmarshals JSON-RPC request params into v
func (hctx *Context) Params(v interface{}) error {
	return hctx.getReqMsg().UnmarshalParams(v)
}

// WithParams set params
func (hctx *Context) WithParams(v interface{}) error {
	err := hctx.getReqMsg().WithParams(v)
	return err
}

// Error returns a possible error encountered while reading JSON-RPC request
//
// If Version() returns empty then it is likely that an error occured
func (hctx *Context) Error() error {
	// Force request loading
	_ = hctx.getReqMsg()
	return hctx.err
}

func (hctx *Context) getReqMsg() *RequestMsg {
	hctx.readMsgOnce.Do(func() {
		// Read request body into request message and validates it
		hctx.msg = &RequestMsg{}
		err := newServerCodec(hctx.req.Body, nil).ReadRequest(hctx.msg)
		if err == nil {
			err = hctx.msg.Validate()
		}
		hctx.err = err

		// Attach message to request context
		hctx.WithContext(WithRequestMsg(hctx.Context(), hctx.msg))
	})

	return hctx.msg
}

// Header allows to access underlying http.ResponseWriter headers
func (c *Context) Header(key, value string) http.Header {
	return c.rw.Header()
}

// WriteResult writes a successful JSON-RPC response with result
//
// If it fails at parsing v then operation is interupted
func (hctx *Context) WriteResult(v interface{}) error {
	msg := hctx.newRespMsg()

	if err := msg.WithResult(v); err != nil {
		return err
	}

	return hctx.writeResponse(msg)
}

// WriteError writes a failed JSON-RPC response with error
func (hctx *Context) WriteError(err error) error {
	msg := hctx.newRespMsg().WithError(err)
	return hctx.writeResponse(msg)
}

func (hctx *Context) newRespMsg() *ResponseMsg {
	return &ResponseMsg{
		Version: hctx.Version(),
		ID:      hctx.ID(),
	}
}

func (hctx *Context) writeResponse(msg *ResponseMsg) error {
	err := msg.Validate()
	if err != nil {
		return err
	}

	// Always respond with status 200
	hctx.rw.WriteHeader(http.StatusOK)

	return newServerCodec(nil, hctx.rw).WriteResponse(msg)
}

type ctxKey string

var (
	reqCtxKey  ctxKey = "req"
	respCtxKey ctxKey = "resp"
	ctxCtxKey  ctxKey = "ctx"
)

// WithRequestMsg attaches a RequestMsg to context
func WithRequestMsg(ctx context.Context, msg *RequestMsg) context.Context {
	return context.WithValue(ctx, reqCtxKey, msg)
}

// RequestMsgFromContext looks for a RequestMsg attached to context
func RequestMsgFromContext(ctx context.Context) *RequestMsg {
	msg, ok := ctx.Value(reqCtxKey).(*RequestMsg)
	if !ok {
		return &RequestMsg{}
	}
	return msg
}

// WithResponseMsg attaches a ResponseMsg to context
func WithResponseMsg(ctx context.Context, msg *ResponseMsg) context.Context {
	return context.WithValue(ctx, respCtxKey, msg)
}

// ResponseMsgFromContext looks for a ResponseMsg attached to context
func ResponseMsgFromContext(ctx context.Context) *ResponseMsg {
	msg, ok := ctx.Value(respCtxKey).(*ResponseMsg)
	if !ok {
		return &ResponseMsg{}
	}
	return msg
}
