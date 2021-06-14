package jsonrpc

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/consensysquorum/quorum-key-manager/pkg/common"
)

type ResponseWriter interface {
	WriteMsg(*ResponseMsg) error
}

func WriteResult(rw ResponseWriter, result interface{}) error {
	return rw.WriteMsg((&ResponseMsg{}).WithResult(result))
}

func WriteError(rw ResponseWriter, err error) error {
	return rw.WriteMsg((&ResponseMsg{}).WithError(err))
}

type idResponseWriter struct {
	rw ResponseWriter
	id interface{}
}

func RWWithID(id interface{}) func(ResponseWriter) ResponseWriter {
	return func(rw ResponseWriter) ResponseWriter {
		return &idResponseWriter{
			rw: rw,
			id: id,
		}
	}
}

func (rw *idResponseWriter) WriteMsg(msg *ResponseMsg) error {
	if msg.ID == nil {
		msg.ID = rw.id
	}

	return rw.rw.WriteMsg(msg)
}

func (rw *idResponseWriter) Writer() io.Writer {
	return rw.rw.(common.WriterWrapper).Writer()
}

type versionResponseWriter struct {
	rw      ResponseWriter
	version string
}

func RWWithVersion(v string) func(ResponseWriter) ResponseWriter {
	return func(rw ResponseWriter) ResponseWriter {
		if v == "" {
			v = defaultVersion
		}

		return &versionResponseWriter{
			rw:      rw,
			version: v,
		}
	}
}

func (rw *versionResponseWriter) WriteMsg(msg *ResponseMsg) error {
	if msg.Version == "" {
		msg.WithVersion(rw.version)
	}

	return rw.rw.WriteMsg(msg)
}

func (rw *versionResponseWriter) Writer() io.Writer {
	return rw.rw.(common.WriterWrapper).Writer()
}

type responseWriter struct {
	w io.Writer

	enc *json.Encoder
}

func NewResponseWriter(w io.Writer) ResponseWriter {
	return &responseWriter{
		w:   w,
		enc: json.NewEncoder(w),
	}
}

func (rw *responseWriter) WriteMsg(msg *ResponseMsg) error {
	if httpRw, ok := rw.w.(http.ResponseWriter); ok {
		httpRw.Header().Set("Content-Type", "application/json")
	}
	return rw.enc.Encode(msg)
}

func (rw *responseWriter) Writer() io.Writer {
	return rw.w
}
