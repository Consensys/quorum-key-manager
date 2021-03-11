package jsonrpc

import (
	"bufio"
	"encoding/json"
	"net"
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter

	WithID(interface{}) ResponseWriter
	WithVersion(string) ResponseWriter

	WriteMsg(*ResponseMsg) error
	WriteResult(interface{}) error
	WriteError(error) error
}

type responseWriter struct {
	rw http.ResponseWriter

	id      interface{}
	version string

	enc *json.Encoder
}

func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		rw:      rw,
		version: defaultVersion,
		enc:     json.NewEncoder(rw),
	}
}

func (rw *responseWriter) Header() http.Header {
	return rw.rw.Header()
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.rw.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.rw.WriteHeader(statusCode)
}

func (rw *responseWriter) WriteMsg(msg *ResponseMsg) error {
	if msg.Version == "" {
		msg.WithVersion(rw.version)
	}

	if msg.ID == nil {
		msg.ID = rw.id
	}

	return rw.enc.Encode(msg)
}

func (rw *responseWriter) WithID(id interface{}) ResponseWriter {
	rw.id = id
	return rw
}

func (rw *responseWriter) WithVersion(version string) ResponseWriter {
	rw.version = version
	return rw
}

func (rw *responseWriter) WriteResult(result interface{}) error {
	return rw.WriteMsg((&ResponseMsg{}).WithResult(result))
}

func (rw *responseWriter) WriteError(err error) error {
	return rw.WriteMsg((&ResponseMsg{}).WithError(err))
}

func (rw *responseWriter) Flush() {
	rw.rw.(http.Flusher).Flush()
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.rw.(http.Hijacker).Hijack()
}

func (rw *responseWriter) Push(target string, opts *http.PushOptions) error {
	return rw.rw.(http.Pusher).Push(target, opts)
}

func (rw *responseWriter) CloseNotify() <-chan bool {
	return rw.rw.(http.CloseNotifier).CloseNotify() //nolint
}
